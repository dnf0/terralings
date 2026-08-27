# Editor LSP & Autocompletion

Terralings features a built-in Language Server Protocol (LSP) daemon to supercharge your editor experience.

---

### Language Server Features

- **Inline Compiler Diagnostics:** Error squiggles pointing directly to the offending HCL lines.
- **Progressive Hint Hover & Code Actions:** Request hints and explanations right inside the editor.
- **Multi-Editor Support:** Seamless integration with VS Code, Neovim (`nvim-lspconfig`), and Helix (`languages.toml`).

---

### Environment Diagnostics (Doctor)

Ensure all prerequisites (OpenTofu / Terraform binary, git, and path resolution) are properly configured:

```bash
terralings doctor
```

---

### Ready to Learn?

Start your Terralings journey now!

- [Run Environment Diagnostics](command:terralings.doctor)
- [Start Watch Mode](command:terralings.watch)
- [Open Guided Tour](command:terralings.tour)
