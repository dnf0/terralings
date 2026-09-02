# Update Repository README Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Choose an execution mode:
> 1. `superpowers:subagent-driven-development` (recommended for multi-agent reviews, backed by `SKILL.state` / `.agent-state/state.json`)
> 2. `agent-rules:stateful-execution` (SKILL.state) (recommended for deterministic single-agent linear execution)
> 3. `superpowers:executing-plans` (batch execution with manual checkpoints)
> Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Modernize `README.md` to accurately document the 15-chapter / 68-exercise curriculum, WebAssembly interactive browser playground, updated CLI diagnostics/analytics mockups, and complete exercise syllabus.

**Architecture:** Update `README.md` in-place following `makeareadme.com` and project guidelines. Maintain all existing badges, links, and quickstart flows while expanding the curriculum matrix to include Chapters 14 (AWS Architecture) and 15 (GCP Architecture), and adding a dedicated section for the WebAssembly browser playground.

**Tech Stack:** Markdown, Material for MkDocs, Go 1.22+, OpenTofu, Terraform, GitHub Actions.

## Global Constraints
- Target 15 Chapters and 68 Exercises everywhere in the document.
- Follow `makeareadme.com` standard structure and maintain existing shield badges and repository links.
- Ensure all exercise links in the curriculum matrix link correctly to documentation and playground where applicable.
- All code blocks, mockups, and tables must be valid GitHub-Flavored Markdown.

---

### Task 1: Update Header, Badges & WebAssembly Playground Showcase

**Files:**
- Modify: `README.md:1-120`

- [ ] **Step 1: Update badges and overview copy**
  - Add WebAssembly Playground badge linking to `https://dnf0.github.io/terralings/playground/`.
  - Update overview description from 56 to 68 exercises.
- [ ] **Step 2: Add WebAssembly Playground Showcase section**
  - Introduce the browser IDE (Monaco Editor, Pyodide Wasm, xterm.js).
  - Link directly to the live playground.
- [ ] **Step 3: Update Architecture ASCII diagram**
  - Update `(13 Chapters / 56 Exercises)` to `(15 Chapters / 68 Exercises)` and include WebAssembly browser client.

---

### Task 2: Update CLI Output Mockups & Chapter Counts

**Files:**
- Modify: `README.md:121-395`

- [ ] **Step 1: Update `terralings doctor` mockup**
  - Update `(56 configuration files found)` to `(68 configuration files found)`.
- [ ] **Step 2: Update `terralings list` description**
  - Update to `15 curriculum chapters and 68 exercises`.
- [ ] **Step 3: Update `terralings verify` mockup**
  - Update progress bar and counter from `56/56 (100.0%)` to `68/68 (100.0%)`.
- [ ] **Step 4: Update `terralings stats` mockup**
  - Update progress counter from `34/56` to `48/68`.
- [ ] **Step 5: Update VS Code Extension section**
  - Update chapter/exercise counts to 15 chapters and 68 exercises.

---

### Task 3: Expand Curriculum Matrix with Chapters 14 and 15

**Files:**
- Modify: `README.md:396-500`

- [ ] **Step 1: Update Curriculum Matrix summary line**
  - Update to `15 structured chapters containing 68 exercises`.
- [ ] **Step 2: Append Chapter 14 (AWS Infrastructure & Production Blueprints)**
  - `aws01`: Multi-AZ VPC Networking (`plan`)
  - `aws02`: Resilient Compute & Load Balancing (`plan`)
  - `aws03`: Serverless Microservice Pipeline (`plan`)
  - `aws04`: Event-Driven Async Decoupling (`plan`)
  - `aws05`: Zero-Trust IAM & Security Hardening (`plan`)
  - `aws06`: Storage & Data Tier Architecture (`plan`)
- [ ] **Step 3: Append Chapter 15 (Google Cloud (GCP) Architecture Blueprints)**
  - `gcp01`: Custom VPC Networking & Firewall Rules (`plan`)
  - `gcp02`: Managed Instance Groups & Load Balancing (`plan`)
  - `gcp03`: Serverless Cloud Run Services (`plan`)
  - `gcp04`: Pub/Sub Event Pipelines (`plan`)
  - `gcp05`: Workload Identity & IAM Federation (`plan`)
  - `gcp06`: Resilient Storage & Cloud Databases (`plan`)

---

### Task 4: Verification & Final Polish

**Files:**
- Verify: `README.md`
- Run: `mkdocs build --strict`
- Run: `uv run pytest -q`

- [ ] **Step 1: Verify markdown rendering and links**
- [ ] **Step 2: Run mkdocs build --strict**
- [ ] **Step 3: Run pytest test suite**
