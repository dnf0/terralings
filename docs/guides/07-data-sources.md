# Chapter 07: Data Sources & State Querying

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Read-Only Data Sources, Filesystem Queries, External Providers, and Pre/Postconditions
-   :material-play-circle: **Interactive Challenges** &bull; 4 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Read-Only Queries

Data sources allow configurations to fetch information defined outside of the immediate Terraform workspace (e.g. existing subnets, AMIs, local files, or remote state). Unlike `resource` blocks, data sources are read-only and never create, mutate, or destroy remote infrastructure.

```text
    ┌──────────────────────────────┐
    │     Data Source Block        │ (Filter / Search Query)
    └──────────────┬───────────────┘
                   │
                   ▼ (Provider Read API)
    ┌──────────────────────────────┐
    │  External System / Cloud API │
    └──────────────┬───────────────┘
                   │
                   ▼ (Data Export)
    ┌──────────────────────────────┐
    │  data.local_file.spec.content│ ──► [ Consumed by Downstream Resources ]
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
# Read existing local file data
data "local_file" "schema" {
  filename = "${path.module}/schema.json"
}

# Resource enforcing preconditions on data source output
resource "terraform_data" "validator" {
  input = data.local_file.schema.content

  lifecycle {
    precondition {
      condition     = length(data.local_file.schema.content) > 0
      error_message = "The schema file must not be empty."
    }
  }
}
```

---

## 3. Production Best Practices

1. **Use Preconditions and Postconditions**: Guard data lookups with `lifecycle.precondition` and `lifecycle.postcondition` blocks to catch external API drift or empty search query results early.
2. **Avoid Non-Deterministic Data Sources in Plan**: Be cautious with data sources that evaluate during the apply phase (when arguments depend on uncreated resources), as this defers validation to apply time.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `datasources01` | Local File Data Sources | `plan` | Query and consume local filesystem contents via `data "local_file"`. |
| `datasources02` | Archive Zip Data Generation | `plan` | Dynamically package directories into zip payloads for deployment. |
| `datasources03` | External Script Data Integration | `plan` | Query structured JSON data from external shell scripts. |
| `datasources04` | Data Preconditions & Postconditions | `plan` | Enforce assertions on fetched data before plan execution. |
