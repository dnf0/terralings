# Design Specification: Documentation Architecture Diagrams to Native Mermaid

**Date:** 2026-09-02  
**Status:** Approved / In Progress  
**Objective:** Convert all ASCII/text architecture plots in the Terralings documentation hub (`docs/index.md`) and all 15 reference guides (`docs/guides/*.md`) into rich, responsive, interactive Mermaid diagrams.

---

## 1. Background & Motivation

The Terralings documentation previously used monospaced ASCII-box diagrams to illustrate execution engines, variable precedence ladders, state surgery graphs, and multi-cloud architectural topologies. While functional, ASCII diagrams:
1. Do not adapt dynamically to the Material for MkDocs dark and light palettes.
2. Are difficult to read and scale on mobile/tablet viewports.
3. Lack visual distinction (colors, node shapes, directional flow styling) for complex cloud architectures.

By enabling `pymdownx.superfences` with the Mermaid custom fence in `mkdocs.yml`, all diagrams will render as crisp, vector-based SVG graphics with automatic theme synchronization.

---

## 2. Configuration (`mkdocs.yml`)

Configure `pymdownx.superfences` in `mkdocs.yml`:

```yaml
markdown_extensions:
  - pymdownx.superfences:
      custom_fences:
        - name: mermaid
          class: mermaid
          format: !!python/name:pymdownx.superfences.fence_code_format
```

---

## 3. Diagram Inventory & Mermaid Specifications

### 3.1 Homepage (`docs/index.md`)
- **Theme**: Client-Side WebAssembly Architecture
- **Type**: `flowchart LR`
- **Nodes**:
  - `Browser["Web Browser Tab"]`
  - `Monaco["Monaco Editor (VS Code in Browser)"]`
  - `Worker["Web Worker (Pyodide Wasm Engine)"]`
  - `Validators["15 Chapter Validators & Hint Engine"]`
  - `Terminal["Interactive xterm.js Terminal Output"]`

### 3.2 Chapter 01: HCL Foundations (`docs/guides/01-primitives.md`)
- **Theme**: HCL Syntax Parsing & DAG Construction
- **Type**: `flowchart TD`
- **Nodes**: `HCL Source Code` -> `Lexer & AST Parser` -> `Provider Schema Handshake` -> `Resource DAG Builder` -> `Topological Execution Walk`

### 3.3 Chapter 02: Input Variables (`docs/guides/02-variables.md`)
- **Theme**: Variable Precedence Hierarchy & Validation Pipeline
- **Type**: `flowchart TD`
- **Nodes**:
  - Priority 1: CLI `-var` & `-var-file`
  - Priority 2: `*.auto.tfvars`
  - Priority 3: `terraform.tfvars`
  - Priority 4: Environment Variables (`TF_VAR_*`)
  - Priority 5: HCL `default` fallback
  - Gate: Custom `validation { condition = ... }`

### 3.4 Chapter 03: Outputs & Locals (`docs/guides/03-outputs-locals.md`)
- **Theme**: Locals DAG Evaluation & Sensitive Masking
- **Type**: `flowchart LR`
- **Nodes**: `Resource State/Attributes` -> `Locals Transformation Block` -> `Sensitive Masking Filter (sensitive = true)` -> `CLI Output / State Persistence`

### 3.5 Chapter 04: Functions & Safe Evaluation (`docs/guides/04-functions.md`)
- **Theme**: Built-in Functions & `can`/`try` Guard Clauses
- **Type**: `flowchart TD`
- **Nodes**: `Raw Input Expression` -> `Function Invocation` -> `Evaluation Success / Error Catching` -> `try() Fallback Handler`

### 3.6 Chapter 05: Meta-Arguments & Scaling (`docs/guides/05-meta-arguments.md`)
- **Theme**: Resource Expansion Engine & Lifecycle Hooks
- **Type**: `flowchart TD`
- **Nodes**: `Resource Block` -> `Branching: count (List Index) vs for_each (Map Key)` -> `Lifecycle Evaluator (create_before_destroy, prevent_destroy)` -> `State Instance Nodes`

### 3.7 Chapter 06: Dynamic Blocks (`docs/guides/06-dynamic-blocks.md`)
- **Theme**: Dynamic Block AST Generation Engine
- **Type**: `flowchart TD`
- **Nodes**: `Complex Collection (List/Map)` -> `dynamic "block_name" Loop` -> `Custom Iterator (iterator.value)` -> `Injected Nested Block AST`

### 3.8 Chapter 07: Data Sources (`docs/guides/07-data-sources.md`)
- **Theme**: Data Source Query Lifecycle & Assertion Gates
- **Type**: `flowchart TD`
- **Nodes**: `Config Compilation` -> `External / Provider Read Phase` -> `Precondition & Postcondition Gate` -> `Downstream Resource Reference Graph`

### 3.9 Chapter 08: Modules (`docs/guides/08-modules.md`)
- **Theme**: Module Hierarchy & Provider Alias Forwarding
- **Type**: `flowchart TD`
- **Nodes**: `Root Module (Default Provider)` -> `Provider Alias Mapping ({ aws = aws.us_west })` -> `Child Module Instance` -> `Encapsulated Resources`

### 3.10 Chapter 09: State Management (`docs/guides/09-state-management.md`)
- **Theme**: Declarative State Surgery (`moved` & `import`)
- **Type**: `flowchart TD`
- **Nodes**: `Prior State Schema` -> `Declarative moved/import Directives` -> `In-Memory State AST Transformation` -> `Zero-Destruction Plan Execution`

### 3.11 Chapter 10: Native Testing (`docs/guides/10-testing.md`)
- **Theme**: `.tftest.hcl` Execution Lifecycle
- **Type**: `flowchart TD`
- **Nodes**: `Test Harness Init` -> `Mock Provider Setup` -> `Sequential run Blocks (command = plan | apply)` -> `Assertion Matrix Evaluation` -> `expect_failures Error Matching`

### 3.12 Chapter 11: Production Patterns (`docs/guides/11-production-patterns.md`)
- **Theme**: Environment Mapping & Tagging Factory
- **Type**: `flowchart LR`
- **Nodes**: `Workspace / Environment Target` -> `Environment Matrix Lookup` -> `Tagging Factory (merge defaults + custom)` -> `Standardized Cloud Resources`

### 3.13 Chapter 12: OpenTofu Innovations (`docs/guides/12-opentofu.md`)
- **Theme**: State Encryption at Rest Engine
- **Type**: `flowchart LR`
- **Nodes**: `Raw State AST` -> `AEAD Encryption Pipeline (AES-GCM)` -> `Key Provider (AWS KMS / GCP Cloud KMS / Passphrase)` -> `Encrypted State at Rest`

### 3.14 Chapter 13: Governance (`docs/guides/13-governance.md`)
- **Theme**: Root & Policy Encapsulation Architecture (ADR-0005)
- **Type**: `flowchart TD`
- **Nodes**: `Root Module Boundary` -> `Service Module Enclosure` -> `Integrated Security Policy / IAM Role (ADR-0005 Owner)` -> `Resource Target`

### 3.15 Chapter 14: AWS Production Architecture (`docs/guides/14-aws-architecture.md`)
- **Theme**: Multi-AZ Resilient AWS Cloud Blueprint
- **Type**: `flowchart TD`
- **Nodes**:
  - `Internet` -> `Internet Gateway`
  - `Public Subnets`: `Application Load Balancer (ALB)`, `NAT Gateways (Multi-AZ)`
  - `Private App Subnets`: `Auto Scaling Groups (ASG) / Launch Templates`, `Serverless Lambda Functions`
  - `Async Messaging`: `SQS FIFO Queues`, `SNS Event Topics`
  - `Data Tier`: `Encrypted S3 Buckets`, `DynamoDB (PAY_PER_REQUEST)`

### 3.16 Chapter 15: GCP Production Architecture (`docs/guides/15-gcp-architecture.md`)
- **Theme**: High-Availability GCP Enterprise Blueprint
- **Type**: `flowchart TD`
- **Nodes**:
  - `Internet` -> `External HTTPS Load Balancer & Cloud Armor`
  - `Custom Regional VPC`: `Regional Subnets with Private Google Access`
  - `Compute Tier`: `Regional Managed Instance Groups (MIG)`, `Cloud Run v2 Container Microservices`
  - `Eventing`: `Cloud Pub/Sub Topics with Dead-Letter Queues (DLQ)`
  - `Security & Data`: `Workload Identity Federation`, `Cloud Storage (Uniform Access)`

---

## 4. Verification & Testing Strategy
- Run `mkdocs build --strict` to verify syntax and rendering for all 16 diagrams.
- Run `uv run pytest -q` and `go test ./test` to confirm no regressions.
- Verify visually in dark and light modes via browser.
