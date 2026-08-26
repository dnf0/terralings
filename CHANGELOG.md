# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-08-26

### Added
- **Interactive Watch Mode**: Continuous filesystem watcher (`watcher.RunWatch`) using `fsnotify` with event debouncing, live UI feedback, and automatic exercise advancement.
- **Dual Engine Auto-Detection**: Binary discovery system (`detector.DetectBinary`) supporting OpenTofu (`tofu`) and Terraform (`terraform`) with `--bin` and `TERRALINGS_BIN` overrides.
- **Execution & Isolation Runner**: Sandboxed workspace runner (`runner.Runner`) supporting `validate`, `plan`, and `test` execution modes with shared provider plugin caching for fast sub-100ms evaluation.
- **Comprehensive 13-Chapter Curriculum**: 56 progressive exercises and reference solutions covering:
  - Chapter 01: HCL Foundations & Core Primitives (`primitives01`–`primitives06`)
  - Chapter 02: Input Variables, Types & Validations (`variables01`–`variables05`)
  - Chapter 03: Outputs, Locals & Expressions (`outputs01`, `locals01`, `expr01`, `expr02`)
  - Chapter 04: Built-in Functions & Collections (`func01`–`func05`)
  - Chapter 05: Meta-Arguments & Resource Scaling (`meta01`–`meta05`)
  - Chapter 06: Dynamic Blocks & Advanced HCL (`dynamic01`–`dynamic04`)
  - Chapter 07: Data Sources & State Querying (`data01`–`data04`)
  - Chapter 08: Modular Infrastructure Architecture (`module01`–`module05`)
  - Chapter 09: State Refactoring & Surgery (`state01`–`state04`)
  - Chapter 10: Native Unit & Integration Testing (`test01`–`test04`)
  - Chapter 11: Production Infrastructure Patterns (`pattern01`–`pattern04`)
  - Chapter 12: OpenTofu Advanced Features (`tofu01`–`tofu03`)
  - Chapter 13: Governance & Security Policies (`gov01`–`gov03`)
- **CLI Commands**: Rich Cobra CLI (`watch`, `run`, `hint`, `list`, `verify`, `version`) with progressive hints and chapter progress tracking.
- **Terminal UI**: Stylized terminal formatting (`ui`) powered by Lip Gloss, including exercise banners, progress bars, and bordered hint cards.
- **End-to-End Test Suite**: Complete E2E tests (`test/e2e_test.go`) validating all 56 exercises, 56 reference solutions, and CLI subprocess workflows.
- **CI/CD Automation**: GitHub Actions workflow testing matrix across Ubuntu and macOS on both OpenTofu and Terraform engines.
