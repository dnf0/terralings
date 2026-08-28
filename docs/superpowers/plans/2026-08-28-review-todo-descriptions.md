# Review and Enhance All Exercise TODO Descriptions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Conduct a comprehensive review of all TODO descriptions across all 56 exercises in the 13 chapters of Terralings, upgrading every TODO to clearly and informatively communicate **WHAT** is being done and **WHY** it is being done from a Terraform/OpenTofu engine and architectural perspective.

**Architecture:** Review and update each chapter's exercise files systematically. Every TODO is structured with explicit `TODO (What): ...` and `TODO (Why): ...` commentary (or equivalent comprehensive multi-line guidance), providing rich educational value while preserving the exact exercise verification mechanics and deterministic failure states.

**Tech Stack:** HCL (Terraform 1.6+ / OpenTofu 1.6+), Go 1.23+ test suite, TypeScript VS Code extension tests.

---

### Chapter Breakdown & Tasks

- **Task 1: Chapter 01 (`01_primitives`) & Chapter 02 (`02_variables`) Review & Upgrade (11 exercises)**
  - `exercises/01_primitives/primitives01.tf` .. `primitives06.tf`
  - `exercises/02_variables/variables01.tf` .. `variables05.tf`
- **Task 2: Chapter 03 (`03_outputs_locals`) & Chapter 04 (`04_functions`) Review & Upgrade (9 exercises)**
  - `exercises/03_outputs_locals/outputs01.tf`, `locals01.tf`, `expr01.tf`, `expr02.tf`
  - `exercises/04_functions/func01.tf` .. `func05.tf`
- **Task 3: Chapter 05 (`05_meta_arguments`) & Chapter 06 (`06_dynamic_blocks`) Review & Upgrade (9 exercises)**
  - `exercises/05_meta_arguments/meta01.tf` .. `meta05.tf`
  - `exercises/06_dynamic_blocks/dynamic01.tf` .. `dynamic04.tf`
- **Task 4: Chapter 07 (`07_data_sources`) & Chapter 08 (`08_modules`) Review & Upgrade (9 exercises)**
  - `exercises/07_data_sources/data01.tf` .. `data04.tf`
  - `exercises/08_modules/module01` .. `module05`
- **Task 5: Chapter 09 (`09_state_refactoring`) & Chapter 10 (`10_testing`) Review & Upgrade (8 exercises)**
  - `exercises/09_state_refactoring/state01.tf` .. `state04.tf`
  - `exercises/10_testing/test01` .. `test04`
- **Task 6: Chapter 11 (`11_patterns`), Chapter 12 (`12_opentofu`) & Chapter 13 (`13_governance`) Review & Upgrade (10 exercises)**
  - `exercises/11_patterns/pattern01.tf` .. `pattern04.tf`
  - `exercises/12_opentofu/tofu01.tf` .. `tofu03.tf`
  - `exercises/13_governance/gov01.tf` .. `gov03.tf`
- **Task 7: Comprehensive Verification & Test Suite Execution**
  - Run `make all`, `make docs-build`, `git diff` review, and `roborev` inspection.

---

### Task 1: Chapter 01 (`01_primitives`) & Chapter 02 (`02_variables`)
**Files:**
- Modify: `exercises/01_primitives/primitives01.tf`
- Modify: `exercises/01_primitives/primitives02.tf`
- Modify: `exercises/01_primitives/primitives03.tf`
- Modify: `exercises/01_primitives/primitives04.tf`
- Modify: `exercises/01_primitives/primitives05.tf`
- Modify: `exercises/01_primitives/primitives06.tf`
- Modify: `exercises/02_variables/variables01.tf`
- Modify: `exercises/02_variables/variables02.tf`
- Modify: `exercises/02_variables/variables03.tf`
- Modify: `exercises/02_variables/variables04.tf`
- Modify: `exercises/02_variables/variables05.tf`

- [ ] **Step 1: Upgrade Chapter 01 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 2: Upgrade Chapter 02 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 3: Run `make test` to verify Chapter 01 and 02 exercises fail initially as expected**
- [ ] **Step 4: Commit Chapter 01 & 02 TODO improvements**

---

### Task 2: Chapter 03 (`03_outputs_locals`) & Chapter 04 (`04_functions`)
**Files:**
- Modify: `exercises/03_outputs_locals/outputs01.tf`
- Modify: `exercises/03_outputs_locals/locals01.tf`
- Modify: `exercises/03_outputs_locals/expr01.tf`
- Modify: `exercises/03_outputs_locals/expr02.tf`
- Modify: `exercises/04_functions/func01.tf`
- Modify: `exercises/04_functions/func02.tf`
- Modify: `exercises/04_functions/func03.tf`
- Modify: `exercises/04_functions/func04.tf`
- Modify: `exercises/04_functions/func05.tf`

- [ ] **Step 1: Upgrade Chapter 03 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 2: Upgrade Chapter 04 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 3: Run `make test` to verify Chapter 03 and 04 exercises fail initially as expected**
- [ ] **Step 4: Commit Chapter 03 & 04 TODO improvements**

---

### Task 3: Chapter 05 (`05_meta_arguments`) & Chapter 06 (`06_dynamic_blocks`)
**Files:**
- Modify: `exercises/05_meta_arguments/meta01.tf`
- Modify: `exercises/05_meta_arguments/meta02.tf`
- Modify: `exercises/05_meta_arguments/meta03.tf`
- Modify: `exercises/05_meta_arguments/meta04.tf`
- Modify: `exercises/05_meta_arguments/meta05.tf`
- Modify: `exercises/06_dynamic_blocks/dynamic01.tf`
- Modify: `exercises/06_dynamic_blocks/dynamic02.tf`
- Modify: `exercises/06_dynamic_blocks/dynamic03.tf`
- Modify: `exercises/06_dynamic_blocks/dynamic04.tf`

- [ ] **Step 1: Upgrade Chapter 05 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 2: Upgrade Chapter 06 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 3: Run `make test` to verify Chapter 05 and 06 exercises fail initially as expected**
- [ ] **Step 4: Commit Chapter 05 & 06 TODO improvements**

---

### Task 4: Chapter 07 (`07_data_sources`) & Chapter 08 (`08_modules`)
**Files:**
- Modify: `exercises/07_data_sources/data01.tf`
- Modify: `exercises/07_data_sources/data02.tf`
- Modify: `exercises/07_data_sources/data03.tf`
- Modify: `exercises/07_data_sources/data04.tf`
- Modify: `exercises/08_modules/module01/main.tf`
- Modify: `exercises/08_modules/module02/main.tf`
- Modify: `exercises/08_modules/module03/main.tf`
- Modify: `exercises/08_modules/module04/main.tf`
- Modify: `exercises/08_modules/module05/main.tf`

- [ ] **Step 1: Upgrade Chapter 07 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 2: Upgrade Chapter 08 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 3: Run `make test` to verify Chapter 07 and 08 exercises fail initially as expected**
- [ ] **Step 4: Commit Chapter 07 & 08 TODO improvements**

---

### Task 5: Chapter 09 (`09_state_refactoring`) & Chapter 10 (`10_testing`)
**Files:**
- Modify: `exercises/09_state_refactoring/state01.tf`
- Modify: `exercises/09_state_refactoring/state02.tf`
- Modify: `exercises/09_state_refactoring/state03.tf`
- Modify: `exercises/09_state_refactoring/state04.tf`
- Modify: `exercises/10_testing/test01/naming.tftest.hcl`
- Modify: `exercises/10_testing/test02/apply.tftest.hcl`
- Modify: `exercises/10_testing/test03/mocking.tftest.hcl`
- Modify: `exercises/10_testing/test04/failures.tftest.hcl`

- [ ] **Step 1: Upgrade Chapter 09 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 2: Upgrade Chapter 10 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 3: Run `make test` to verify Chapter 09 and 10 exercises fail initially as expected**
- [ ] **Step 4: Commit Chapter 09 & 10 TODO improvements**

---

### Task 6: Chapter 11 (`11_patterns`), Chapter 12 (`12_opentofu`) & Chapter 13 (`13_governance`)
**Files:**
- Modify: `exercises/11_patterns/pattern01.tf`
- Modify: `exercises/11_patterns/pattern02.tf`
- Modify: `exercises/11_patterns/pattern03.tf`
- Modify: `exercises/11_patterns/pattern04.tf`
- Modify: `exercises/12_opentofu/tofu01.tf`
- Modify: `exercises/12_opentofu/tofu02.tf`
- Modify: `exercises/12_opentofu/tofu03.tf`
- Modify: `exercises/13_governance/gov01.tf`
- Modify: `exercises/13_governance/gov02.tf`
- Modify: `exercises/13_governance/gov03.tf`

- [ ] **Step 1: Upgrade Chapter 11 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 2: Upgrade Chapter 12 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 3: Upgrade Chapter 13 exercises with informative WHAT & WHY TODOs**
- [ ] **Step 4: Run `make test` to verify Chapter 11, 12, and 13 exercises fail initially as expected**
- [ ] **Step 5: Commit Chapter 11, 12 & 13 TODO improvements**

---

### Task 7: Comprehensive Verification & Handoff
- [ ] **Step 1: Run `make all` across entire repository**
- [ ] **Step 2: Run `make docs-build`**
- [ ] **Step 3: Run `roborev status` and `roborev show HEAD`**
- [ ] **Step 4: Present completion summary and PR options to user**
