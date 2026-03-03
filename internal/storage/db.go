package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Init(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		org TEXT,
		repo TEXT,
		pr INTEGER,
		score REAL,
		tier TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)
	return db, err
}

func Save(db *sql.DB, org, repo string, pr int, score float64, tier string) error {
	_, err := db.Exec("INSERT INTO scans(org,repo,pr,score,tier) VALUES (?,?,?,?,?)",
		org, repo, pr, score, tier)
	return err
}

func OrgRisk(db *sql.DB, org string) (float64, error) {
	row := db.QueryRow("SELECT AVG(score) FROM scans WHERE org=?", org)
	var avg float64
	return avg, row.Scan(&avg)
}
