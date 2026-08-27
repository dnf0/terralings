# Anatomy of an Exercise

Terralings structures its curriculum into 13 progressive chapters covering every major facet of Terraform and OpenTofu.

---

### Exercise Structure

All exercises reside in the `exercises/` directory:

```text
exercises/
├── 01_primitives/
│   ├── primitives01.tf
│   └── primitives02.tf
├── 02_resources/
├── 03_variables/
...
└── 13_advanced/
```

Each exercise file contains:
- `# TODO:` comments detailing the objective and syntax requirements.
- Deliberate syntax errors, missing configurations, or failing test assertions.

---

### Solving Exercises

1. Open the exercise in your editor.
2. Read the instructions and the compiler/linter error messages.
3. Edit the code to fix the issues.
4. Save the file. When valid, Terralings marks it complete!

If you ever want to see the canonical approach, check the corresponding file in `solutions/`.

---

### Quick Actions

- [Refresh Curriculum Tree](command:terralings.refreshTree)
- [Verify Current File](command:terralings.runCurrent)
