package githubapp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
        "os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"

	"github.com/yourname/guardrail-saas/internal/engine"
	"github.com/yourname/guardrail-saas/internal/storage"
)

var dbInstance = mustInitDB()

func mustInitDB() *sql.DB {
	path := os.Getenv("DATABASE_PATH")
	if path == "" {
		path = "./guardrail.db"
	}
	db, err := storage.Init(path)
	if err != nil {
		// In production you’d log and exit; for now fail fast
		panic(fmt.Errorf("db init failed: %w", err))
	}
	return db
}

type PREvent struct {
	Action string `json:"action"`
	Number int    `json:"number"`
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
	Repository struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	} `json:"repository"`
	PullRequest struct {
		Head struct {
			SHA string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		http.Error(w, "server misconfigured: missing GITHUB_WEBHOOK_SECRET", http.StatusInternalServerError)
		return
	}

	body, _ := io.ReadAll(r.Body)

	if !verify(secret, r.Header.Get("X-Hub-Signature-256"), body) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var ev PREvent
	if err := json.Unmarshal(body, &ev); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	if ev.Action != "opened" && ev.Action != "synchronize" {
		w.WriteHeader(http.StatusOK)
		return
	}

	client, err := installationClient(ev.Installation.ID)
	if err != nil {
		http.Error(w, "installation auth failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Replace with real repo checkout + full engine scan.
	// For now, keep the placeholder deterministic value.
	violationsCount := 2
	score, tier := engine.Score(violationsCount, "prod")

	_ = storage.Save(dbInstance, ev.Repository.Owner.Login, ev.Repository.Name, ev.Number, score, tier)

	ctx := context.Background()
	conclusion := "success"
	if score > 65 {
		conclusion = "failure"
	}

	check := github.CreateCheckRunOptions{
		Name:       "Guardrail Deployment Risk",
		HeadSHA:    ev.PullRequest.Head.SHA,
		Status:     github.String("completed"),
		Conclusion: github.String(conclusion),
		Output: &github.CheckRunOutput{
			Title:   github.String(fmt.Sprintf("Risk Score: %.0f (%s)", score, tier)),
			Summary: github.String("Automated deployment risk analysis (MVP)."),
		},
	}

	_, _, _ = client.Checks.CreateCheckRun(ctx, ev.Repository.Owner.Login, ev.Repository.Name, check)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func verify(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func installationClient(installationID int64) (*github.Client, error) {
	appIDStr := os.Getenv("GITHUB_APP_ID")
	keyPath := os.Getenv("GITHUB_PRIVATE_KEY_PATH")
	if appIDStr == "" || keyPath == "" {
		return nil, fmt.Errorf("missing GITHUB_APP_ID or GITHUB_PRIVATE_KEY_PATH")
	}

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid GITHUB_APP_ID: %w", err)
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse private key pem: %w", err)
	}

	j := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": appID,
	})

	signedJWT, err := j.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign app jwt: %w", err)
	}

	// Client authenticated as the GitHub App
	appTS := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signedJWT})
	appClient := github.NewClient(oauth2.NewClient(context.Background(), appTS))

	// Exchange for installation token
	instToken, _, err := appClient.Apps.CreateInstallationToken(context.Background(), installationID, nil)
	if err != nil {
		return nil, fmt.Errorf("create installation token: %w", err)
	}

	instTS := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: instToken.GetToken()})
	return github.NewClient(oauth2.NewClient(context.Background(), instTS)), nil
}
