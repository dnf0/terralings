# Continuous Watch Mode

The fastest way to work through Terralings is using **Watch Mode**.

---

### The Edit-Save-Verify Loop

When Watch Mode is running, Terralings monitors your files for changes:

1. You edit and save an `.tf` exercise file.
2. Terralings instantly re-evaluates the exercise in an isolated runner sandbox.
3. If it fails, the error output is displayed immediately.
4. When it succeeds, you are prompted to advance to the next exercise.

```text
  [Enter / n] Next exercise  |  [p] Previous  |  [r] Rerun  |  [q] Quit
```

---

### Running On Demand

You can also run or verify exercises individually:

- **Run current exercise:** `terralings run <exercise>`
- **Verify entire curriculum:** `terralings verify`

---

### Quick Actions

- [Start Watch Mode](command:terralings.watch)
- [Verify Current Exercise](command:terralings.runCurrent)
