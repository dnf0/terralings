# Terralings 🌍

**An interactive, client-side WebAssembly learning platform and comprehensive reference manual for Terraform & OpenTofu.**

[![Playground](https://img.shields.io/badge/Playground-⚡%20Launch%20Interactive%20IDE-blueviolet)](playground/index.html)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Curriculum](https://img.shields.io/badge/Curriculum-13%20Chapters%20%7C%2056%20Exercises-brightgreen)](syllabus.md)

---

## ⚡ The Modern Way to Master Infrastructure as Code

Terralings combines a **zero-install, 100% client-side WebAssembly interactive playground** with **13 comprehensive architectural reference guides** spanning the entire Terraform & OpenTofu ecosystem.

<div class="grid cards" markdown>

-   :material-play-circle-outline: **Zero-Install Web IDE**
    ---
    Run Monaco Editor, Pyodide WebAssembly, and real-time HCL validation 100% in your browser. No cloud credentials, no local binary installs, and zero backend servers required.
    
    [**Launch Playground →**](playground/index.html){ .md-button .md-button--primary }

-   :material-book-open-page-variant-outline: **13-Chapter Reference Manual**
    ---
    Deep architectural documentation, annotated HCL specs, production best practices, and state surgery workflows for modern IaC.

    [**Explore Reference Guides →**](#-comprehensive-13-chapter-reference-guides){ .md-button }

</div>

---

## 📚 Comprehensive 13-Chapter Reference Guides

Explore in-depth architectural guides and launch linked practice exercises directly into the playground:

<div class="grid cards" markdown>

-   ### Core HCL & Data Flow
    ---
    - [**01. HCL Foundations & Core Primitives**](guides/01-primitives.md) &bull; Blocks, provider constraints, resources, heredocs, lifecycle
    - [**02. Input Variables, Types & Validations**](guides/02-variables.md) &bull; Primitives, collections, objects, custom validation rules
    - [**03. Outputs, Locals & Expressions**](guides/03-outputs-locals.md) &bull; Outputs, sensitive masking, locals, ternaries, splat expressions
    - [**04. Built-in Functions & Collections**](guides/04-functions.md) &bull; String ops, collection math, encoding, files, `try`/`can`

-   ### Dynamic Logic & Modules
    ---
    - [**05. Meta-Arguments & Scaling**](guides/05-meta-arguments.md) &bull; `count`, `for_each`, `depends_on`, lifecycle rules, replacement triggers
    - [**06. Dynamic Blocks & Advanced HCL**](guides/06-dynamic-blocks.md) &bull; `dynamic` block generators, custom iterators, nested loops
    - [**07. Data Sources & State Querying**](guides/07-data-sources.md) &bull; Local files, archive zips, external scripts, preconditions
    - [**08. Modular Infrastructure Architecture**](guides/08-modules.md) &bull; Child modules, multi-instance calling, provider aliases

-   ### State Surgery, Testing & Governance
    ---
    - [**09. State Management & Refactoring Surgery**](guides/09-state-management.md) &bull; Declarative `moved` blocks, `import` blocks, replacements
    - [**10. Native Unit & Integration Testing**](guides/10-testing.md) &bull; `.tftest.hcl`, `run` blocks, mocks, `expect_failures`
    - [**11. Production Patterns & Anti-Patterns**](guides/11-production-patterns.md) &bull; Environment mapping, feature flags, tagging factories
    - [**12. OpenTofu Innovations & Enterprise**](guides/12-opentofu.md) &bull; State encryption at rest, key providers, open registries
    - [**13. Architecture Governance & Standards**](guides/13-governance.md) &bull; Root encapsulation, policy scoping, ephemeral isolation

</div>

---

## 💡 How the Playground Works

The Terralings web playground runs entirely on client-side WebAssembly technology:

```text
┌────────────────────────────────────────────────────────────────────────┐
│                          Web Browser Tab                               │
│  ┌───────────────────────┐             ┌─────────────────────────────┐ │
│  │     Monaco Editor     │   HCL Code  │  Web Worker (Pyodide Wasm)  │ │
│  │  (VS Code in Browser) ├────────────►│  • Python 3.12 Engine       │ │
│  └───────────────────────┘             │  • 13 Chapter Validators    │ │
│                                        │  • Progressive Hint Engine  │ │
│  ┌───────────────────────┐             └──────────────┬──────────────┘ │
│  │   Interactive xterm   │◄───────────────────────────┘                │
│  │    Terminal Output    │         Instant Test & Validation Result    │
│  └───────────────────────┘         (< 1ms in WebAssembly)              │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Navigation

<div class="grid cards" markdown>

-   :material-laptop: **Interactive Playground**
    ---
    Dive straight into the browser IDE with real-time feedback.
    
    [Open Playground →](playground/index.html)

-   :material-format-list-numbered: **Curriculum Syllabus**
    ---
    View the complete breakdown of all 56 exercises and concepts.
    
    [Explore Syllabus →](syllabus.md)

-   :material-github: **Contribute on GitHub**
    ---
    Add new exercises, improve validators, or submit guides.
    
    [Contributing Guide →](contributing.md)

</div>
