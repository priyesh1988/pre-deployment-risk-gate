
# Guardrail — Deployment Risk Intelligence for Kubernetes

![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Language](https://img.shields.io/badge/language-Go-blue)
![Status](https://img.shields.io/badge/status-early--stage-orange)

Guardrail is a **Deployment Risk Intelligence platform** built from the ground up to help engineering teams prevent unsafe Kubernetes releases before they reach production.

It is not a policy engine.
It is not a vulnerability scanner.
It is not cloud posture management.

It is a **release decision layer** for Kubernetes.

---

## 🚨 The Problem

Most Kubernetes security tools:
- List violations
- Generate compliance reports
- Enforce policies at runtime

But they don’t answer the most important CI/CD question:

> Is this deployment safe enough to merge and ship?

Engineering teams in regulated environments need:
- Quantified release risk
- Automated merge gating
- Audit-ready evidence
- Org-level risk visibility

Guardrail provides exactly that.

---

## 🔍 Current Market Landscape

Here’s how existing tools approach the space:

| Product | Primary Focus | Strength | Gap |
|----------|---------------|----------|------|
| Snyk IaC | IaC scanning | Finds misconfigurations & vulnerabilities | No unified deployment risk score |
| Wiz | Cloud security platform | Deep cloud visibility | Not PR-native release gating |
| Datadog KSPM | Runtime posture monitoring | Strong observability | Not a release decision engine |
| Kubescape | Kubernetes compliance scanning | CIS/NSA framework checks | No normalized risk scoring |
| OPA / Gatekeeper | Policy enforcement | Flexible policy engine | Requires custom policies, no risk model |

---

## 🧠 How Guardrail Is Different

Guardrail shifts from **finding issues** to **quantifying release risk**.

Instead of 20 disconnected findings, Guardrail produces:

- 🎯 Risk Score (0–100)
- 📊 Risk Tier (LOW, MEDIUM, HIGH, CRITICAL)
- 🔍 Top Risk Drivers
- 🔐 Merge Gate Decision
- 🧾 Compliance Evidence Artifact

This enables automated release governance.

---

## 🚀 Core Capabilities

### 1️⃣ PR-Native Risk Analysis

On every pull request:

- Analyze Kubernetes manifests
- Calculate deployment risk score
- Post GitHub Check Run status
- Block merge if threshold exceeded

Example:

```
Deployment Risk: 74 / 100 (HIGH)

Top Drivers:
- Missing resource limits
- Privileged container
- Implicit :latest image

Status: ❌ Blocked
```

---

### 2️⃣ Compliance Evidence Export

Each scan generates a structured artifact:

```json
{
  "repository": "org/service",
  "environment": "prod",
  "score": 72,
  "tier": "HIGH",
  "violations": [...],
  "timestamp": "2026-03-03T14:32:00Z"
}
```

Designed for:
- SOC2
- PCI
- HIPAA
- Internal audit reviews

---

### 3️⃣ Org-Level Risk Tracking

Aggregate risk metrics across repositories:

```
GET /api/org/acme/risk
```

Response:

```json
{
  "average_risk": 47.3
}
```

Enables leadership-level visibility.

---

### 4️⃣ Multi-Tenant SaaS Architecture

Each organization has:

- Isolated risk policies
- Independent thresholds
- Role-based access
- Historical risk tracking
- Billing integration

---

## ⚙️ Architecture Overview

GitHub Webhook  
↓  
Guardrail Engine  
↓  
Risk Score + Tier  
↓  
Check Run Status  
↓  
Compliance Artifact  
↓  
Multi-Tenant Database  
↓  
Org Risk API  

---

## 📦 Local Development

Build:

```bash
go test ./...
go build -o guardrail ./cmd/guardrail
```

Scan:

```bash
./guardrail -path ./k8s -env prod -threshold 65
```

---

## 🔧 Configuration

Create `.guardrail.yaml`:

```yaml
environment: prod
threshold: 65

rules:
  IMG_LATEST_TAG:
    enabled: true
  MISSING_RESOURCES:
    enabled: true
  PRIVILEGED_CONTAINER:
    enabled: true
```

---

## 🎯 Design Principles

- Deterministic rule engine
- Quantitative risk modeling
- CI-native feedback
- Minimal developer friction
- Enterprise-grade auditability

---

## 🛣 Roadmap

- Full rule engine expansion
- Policy packs for regulated industries
- Risk trend dashboards
- Waiver + approval workflow
- Enterprise SSO integration
- Stripe billing for SaaS plans

---

## 🔐 Why It Matters

Kubernetes is powerful — but YAML misconfigurations cause real outages.

Guardrail prevents unsafe deployments before they reach production, giving teams confidence to ship faster without sacrificing compliance or reliability.

---

Built from scratch.  
Designed for regulated teams.  
Focused on release intelligence.
