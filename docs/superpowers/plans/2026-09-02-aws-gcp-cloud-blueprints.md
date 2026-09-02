# AWS & GCP Production Cloud Blueprints Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Choose an execution mode:
> 1. `superpowers:subagent-driven-development` (recommended for multi-agent reviews, backed by `SKILL.state` / `.agent-state/state.json`)
> 2. `agent-rules:stateful-execution` (SKILL.state) (recommended for deterministic single-agent linear execution)
> 3. `superpowers:executing-plans` (batch execution with manual checkpoints)
> Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expand the Terralings interactive curriculum with Chapters 14 (AWS) and 15 (GCP), introducing 12 production cloud architecture exercises and comprehensive 6-section reference guides.

**Architecture:** Client-side WebAssembly evaluation in Pyodide with zero external cloud credentials, comprehensive MkDocs reference guides with failure triage trees and field tables, and deep bidirectional playground linking.

**Tech Stack:** HCL/OpenTofu, Pyodide/WebAssembly, MkDocs Material, JavaScript, JSON.

## Global Constraints

- Zero cloud credentials required (all exercises run locally in Wasm).
- 6-section guide format strictly matching Chapters 01–13 and Kubelings standard.
- Strict Markdown link and navigation validation via `mkdocs build --strict`.
- Conventional commit messages (`docs: ...`, `feat: ...`).

---

### Task 1: Author Chapter 14 Reference Guide (AWS Cloud Blueprints)

**Files:**
- Create: `docs/guides/14-aws-architecture.md`

**Interfaces:**
- Produces: 6-section reference guide covering `aws01` through `aws06` with field tables and triage tree.

- [ ] **Step 1: Write `docs/guides/14-aws-architecture.md`** with complete 6-section template:
  - Header cards (`aws01`–`aws06`).
  - Section 1: AWS Core Plane & VPC/Compute/Serverless Architecture ASCII diagram.
  - Section 2: Annotated Production HCL Anatomy & Key Field Schema Reference table.
  - Section 3: Real-World Architectural Patterns (ALB + ASG compute & Serverless API Gateway + Lambda).
  - Section 4: Production Hardening & Operational Governance.
  - Section 5: Failure Modes & Diagnostic Triage Tree (`??? failure "..."`).
  - Section 6: Interactive Practice Matrix for `aws01` through `aws06`.
- [ ] **Step 2: Verify markdown formatting with `mkdocs build --strict`**
- [ ] **Step 3: Commit Chapter 14 Guide**
  ```bash
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git add docs/guides/14-aws-architecture.md
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "docs: add chapter 14 AWS cloud architecture reference guide"
  ```

---

### Task 2: Author Chapter 15 Reference Guide (GCP Cloud Blueprints)

**Files:**
- Create: `docs/guides/15-gcp-architecture.md`

**Interfaces:**
- Produces: 6-section reference guide covering `gcp01` through `gcp06` with field tables and triage tree.

- [ ] **Step 1: Write `docs/guides/15-gcp-architecture.md`** with complete 6-section template:
  - Header cards (`gcp01`–`gcp06`).
  - Section 1: GCP Resource Hierarchy, VPC & Cloud Run Architecture ASCII diagram.
  - Section 2: Annotated Production HCL Anatomy & Key Field Schema Reference table.
  - Section 3: Real-World Architectural Patterns (Managed Instance Groups & Cloud Run v2 with Eventarc).
  - Section 4: Production Hardening & Operational Governance.
  - Section 5: Failure Modes & Diagnostic Triage Tree (`??? failure "..."`).
  - Section 6: Interactive Practice Matrix for `gcp01` through `gcp06`.
- [ ] **Step 2: Verify markdown formatting with `mkdocs build --strict`**
- [ ] **Step 3: Commit Chapter 15 Guide**
  ```bash
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git add docs/guides/15-gcp-architecture.md
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "docs: add chapter 15 GCP cloud architecture reference guide"
  ```

---

### Task 3: Expand Playground Bundle & Guide Routing Engine

**Files:**
- Modify: `docs/assets/playground/playground-bundle.json`
- Modify: `docs/assets/playground/playground.js:15-35`

**Interfaces:**
- Consumes: Exercise IDs `aws01`–`aws06` and `gcp01`–`gcp06`.
- Produces: Active Wasm exercise data, prompt instructions, initial code, solutions, and chapter guide routing.

- [ ] **Step 1: Add Chapter 14 (`aws01`–`aws06`) and Chapter 15 (`gcp01`–`gcp06`) exercises into `playground-bundle.json`**
- [ ] **Step 2: Update `CHAPTER_GUIDES` map in `playground.js`** to route chapters `14` and `15` to their corresponding guides.
- [ ] **Step 3: Commit bundle and router changes**
  ```bash
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git add docs/assets/playground/
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "feat(playground): add AWS and GCP exercises to Wasm bundle and router"
  ```

---

### Task 4: Update Syllabus Matrix and Site Navigation

**Files:**
- Modify: `docs/syllabus.md`
- Modify: `mkdocs.yml`

**Interfaces:**
- Consumes: Chapters 14 & 15 guides and 12 exercise IDs.
- Produces: Integrated 15-chapter curriculum overview and updated documentation sidebar.

- [ ] **Step 1: Update `docs/syllabus.md`** to add Chapters 14 & 15 tables with direct `playground/index.html?exercise=<id>` solve links.
- [ ] **Step 2: Update `mkdocs.yml` navigation** to include Chapters 14 & 15 under a "Cloud Architecture & Blueprints" section.
- [ ] **Step 3: Commit syllabus and navigation updates**
  ```bash
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git add docs/syllabus.md mkdocs.yml
  PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "docs: update syllabus and mkdocs navigation with AWS and GCP tracks"
  ```

---

### Task 5: Strict Verification & Release Build

**Files:**
- Verify: Full documentation suite and interactive playground.

- [ ] **Step 1: Run `mkdocs build --strict`** to ensure zero link or syntax warnings.
- [ ] **Step 2: Push changes to remote repository**
  ```bash
  git push origin main
  ```
