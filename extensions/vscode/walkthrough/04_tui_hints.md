# TUI Dashboard & Hints

Terralings includes an interactive terminal dashboard and progressive multi-level hints.

---

### Terminal Dashboard (TUI)

Launch the full-screen interactive dashboard with:

```bash
terralings tui
```

Features:
- Split-pane live view with exercise instructions and real-time execution outputs.
- Keyboard navigation through chapters and exercises.
- Progress metrics, completion percentages, and attempt statistics.

---

### Multi-Level Hints

Stuck on a tricky syntax quirk or provider block? Terralings provides progressive hints without spoiling the entire solution:

- **Level 1:** High-level conceptual guidance.
- **Level 2:** Concrete HCL pointers and syntax hints.
- **Level 3:** Complete solution breakdown.

Use `terralings hint <exercise>` or the VS Code command to view hints.

---

### Resetting an Exercise

Want to start an exercise from scratch? You can reset any exercise file back to its default template:

```bash
terralings reset <exercise>
```

---

### Quick Actions

- [Open Terminal Dashboard (TUI)](command:terralings.tui)
- [Show Hint for Current Exercise](command:terralings.hint)
- [Reset Exercise Template](command:terralings.reset)
