# Design Spec: AWS & GCP Production Cloud Blueprints Curriculum

**Author:** Antigravity  
**Date:** 2026-09-02  
**Status:** Approved  
**Target:** Terralings WebAssembly Playground & MkDocs Documentation Guides

---

## 1. Executive Summary

This specification expands the Terralings interactive curriculum from 13 chapters (56 exercises) to **15 chapters (68 exercises)** by introducing two dedicated enterprise cloud tracks:
- **Chapter 14: AWS Infrastructure & Production Blueprints** (6 exercises)
- **Chapter 15: Google Cloud (GCP) Architecture Blueprints** (6 exercises)

All exercises execute entirely client-side in WebAssembly with zero external cloud credentials, using in-memory HCL AST parsing, semantic schema validation, and bidirectional linking to comprehensive 6-section reference guides.

---

## 2. Curriculum Scope & Exercise Inventory

### Chapter 14: AWS Infrastructure & Production Blueprints (`docs/guides/14-aws-architecture.md`)
- **`aws01` (Multi-AZ VPC Networking)**: Configure `aws_vpc`, public and private `aws_subnet` across availability zones, `aws_internet_gateway`, and explicit `aws_route_table_association`.
- **`aws02` (Resilient Compute & Load Balancing)**: Author `aws_launch_template`, `aws_autoscaling_group`, and `aws_lb_target_group` with health check rules.
- **`aws03` (Serverless Microservice Pipeline)**: Deploy `aws_lambda_function`, `aws_apigatewayv2_api` (HTTP API), `aws_apigatewayv2_integration`, and IAM execution role.
- **`aws04` (Event-Driven Async Decoupling)**: Build an asynchronous event pipeline with `aws_sqs_queue` (FIFO), `aws_sns_topic`, `aws_sns_topic_subscription`, and dead-letter queues.
- **`aws05` (Zero-Trust IAM & Security Hardening)**: Define `aws_iam_role` with strict `assume_role_policy` documents, least-privilege `aws_iam_policy`, and egress-filtered `aws_security_group`.
- **`aws06` (Storage & Data Tier Architecture)**: Provision secure `aws_s3_bucket` (SSE-KMS, Versioning, Public Access Block) and `aws_dynamodb_table` with `PAY_PER_REQUEST` and Global Secondary Indexes.

### Chapter 15: Google Cloud (GCP) Architecture Blueprints (`docs/guides/15-gcp-architecture.md`)
- **`gcp01` (Custom VPC Networking & Firewall Rules)**: Configure `google_compute_network` with `auto_create_subnetworks = false`, regional `google_compute_subnetwork`, and target-tagged `google_compute_firewall`.
- **`gcp02` (Managed Instance Groups & Load Balancing)**: Author `google_compute_instance_template` and `google_compute_region_instance_group_manager` with rolling update policies.
- **`gcp03` (Serverless Cloud Run Services)**: Build `google_cloud_run_v2_service` with container concurrency limits, traffic splits, and `google_cloud_run_service_iam_member`.
- **`gcp04` (Pub/Sub Event Pipelines)**: Construct an event distribution bus using `google_pubsub_topic`, `google_pubsub_subscription` with push/pull configurations, and `dead_letter_policy`.
- **`gcp05` (Workload Identity & IAM Federation)**: Configure `google_service_account`, scoped `google_project_iam_member`, and Workload Identity binding for zero-secret CI/CD.
- **`gcp06` (Resilient Cloud Storage & Databases)**: Provision `google_storage_bucket` with uniform bucket-level access, retention lifecycle rules, and `google_sql_database_instance`.

---

## 3. WebAssembly Client-Side Validation Engine

- **Engine Execution**: Runs inside the Pyodide Web Worker (`playground-worker.js`).
- **Validation Rules**: Uses Python-based AST analysis and schema inspection to verify:
  1. Correct resource block types, labels, and arguments.
  2. Attribute cross-references and dependency wiring (e.g. `subnet_ids = [aws_subnet.public_a.id, aws_subnet.public_b.id]`).
  3. Security attributes (e.g., SSE-KMS encryption, private subnets, egress filtering).
- **Diagnostics**: Returns structured error messages with contextual tips pointing to the relevant chapter guide.

---

## 4. Documentation & Navigation Deliverables

1. **Chapter 14 Reference Guide** ([`docs/guides/14-aws-architecture.md`](file:///Users/danielfisher/repos/terralings/docs/guides/14-aws-architecture.md)):
   - 6-section format: Header card, DAG data-flow diagram, annotated HCL & schema table, 2 real-world patterns, production hardening, failure triage tree, interactive practice matrix.
2. **Chapter 15 Reference Guide** ([`docs/guides/15-gcp-architecture.md`](file:///Users/danielfisher/repos/terralings/docs/guides/15-gcp-architecture.md)):
   - 6-section format with complete GCP resource field references and diagnostic triage steps.
3. **Curriculum Syllabus** ([`docs/syllabus.md`](file:///Users/danielfisher/repos/terralings/docs/syllabus.md)):
   - Updated overview table (15 chapters, 68 exercises) with solve links for all 68 exercises.
4. **Site Navigation** ([`mkdocs.yml`](file:///Users/danielfisher/repos/terralings/mkdocs.yml)):
   - Adds "Cloud Architecture & Blueprints" domain cluster to the navigation tree.
5. **Playground Bundle & Linking** ([`docs/assets/playground/playground-bundle.json`](file:///Users/danielfisher/repos/terralings/docs/assets/playground/playground-bundle.json) & [`docs/assets/playground/playground.js`](file:///Users/danielfisher/repos/terralings/docs/assets/playground/playground.js)):
   - Extends bundle with 12 new exercises and links chapter IDs 14 & 15 to their guides.

---

## 5. Spec Self-Review Checklist
- [x] **Placeholder Scan**: No TBDs or incomplete sections.
- [x] **Internal Consistency**: Exercise IDs match guide matrices and bundle schema.
- [x] **Scope Check**: Tightly focused on Chapters 14 & 15 without altering existing Core HCL mechanics.
- [x] **Ambiguity Check**: Zero-credential client-side Wasm evaluation explicitly specified.
