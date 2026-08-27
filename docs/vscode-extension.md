# VS Code Companion Extension

The **Terralings Companion Extension for Visual Studio Code** provides a seamless, integrated learning experience directly within your editor.

---

## Key Features

- **Activity Bar Curriculum Explorer (`terralings-exercises`)**: Visual sidebar tree displaying all 13 chapters and 56 exercises with live passing (`✓`), in-progress, or failing status icons.
- **Embedded Language Server Protocol (LSP) Client**: Automatically connects to `terralings lsp` to provide real-time compilation diagnostics, error squiggles, and hover hints directly in `.tf` and `.tftest.hcl` files.
- **Interactive 5-Step Walkthrough**: Onboarding guide integrated into VS Code's Welcome experience (`Welcome to Terralings`).
- **Integrated Terminal Runners**: One-click actions to launch `terralings watch`, `terralings tui`, `terralings doctor`, or verify the active exercise.
- **Status Bar Progress Indicator**: Real-time completion counter and progress badge in the VS Code status bar.
- **Inline Context Menus**: Quick actions in the editor title bar and sidebar context menu to run, hint, or reset exercises.

---

## Installation

### Method 1: Install from VSIX Package (Recommended)

You can install the packaged extension directly from the repository build or downloaded `.vsix` file:

```bash
# Build the VSIX locally from the repo root
cd extensions/vscode
npm install
npm run build
npx @vscode/vsce package -o ../../dist/terralings-vscode.vsix

# Install via VS Code CLI
code --install-extension ../../dist/terralings-vscode.vsix
```

### Method 2: Manual Installation via VS Code UI

1. Open Visual Studio Code.
2. Open the **Extensions** view (`Cmd+Shift+X` on macOS, `Ctrl+Shift+X` on Linux/Windows).
3. Click the **Views and More Actions** menu (`...`) at the top-right of the Extensions pane.
4. Select **Install from VSIX...**.
5. Navigate to and select `terralings-vscode.vsix`.
6. Reload or restart VS Code when prompted.

---

## Extension Configuration Settings

Configure extension behaviors in VS Code Settings (`Cmd+,` or `Ctrl+,`), searching for `terralings`:

| Setting Key | Type | Default | Description |
|---|:---:|:---:|---|
| `terralings.binaryPath` | `string` | `"terralings"` | Path to the `terralings` CLI executable. If `terralings` is in your system `$PATH`, the default value works automatically. |
| `terralings.enableLsp` | `boolean` | `true` | Enable or disable the built-in Terralings Language Server client for diagnostics and hover tooltips. |
| `terralings.autoOpenWalkthrough` | `boolean` | `true` | Automatically show the Welcome walkthrough when opening an exercise workspace for the first time. |

### Example `settings.json`

```json
{
  "terralings.binaryPath": "/usr/local/bin/terralings",
  "terralings.enableLsp": true,
  "terralings.autoOpenWalkthrough": true
}
```

---

## Command Palette Reference

All extension commands can be triggered via the VS Code Command Palette (`Cmd+Shift+P` on macOS, `Ctrl+Shift+P` on Linux/Windows):

| Command Identifier | Command Palette Title | Description |
|---|---|---|
| `terralings.openExercise` | `Terralings: Open Exercise` | Prompt for and open a specific curriculum exercise file |
| `terralings.runCurrent` | `Terralings: Verify Current Exercise` | Run instant evaluation against the active editor file |
| `terralings.watch` | `Terralings: Start Watch Mode` | Launch `terralings watch` in an integrated terminal |
| `terralings.tui` | `Terralings: Open Terminal Dashboard (TUI)` | Open `terralings tui` in an integrated terminal |
| `terralings.hint` | `Terralings: Show Hint` | Display progressive hints for the currently active exercise |
| `terralings.reset` | `Terralings: Reset Exercise Template` | Reset current exercise back to its starting template |
| `terralings.doctor` | `Terralings: Run Environment Diagnostics (Doctor)` | Run pre-flight health checks in terminal |
| `terralings.tour` | `Terralings: Open Guided Tour` | Start interactive onboarding tour in terminal |
| `terralings.refreshTree` | `Terralings: Refresh Exercises` | Force refresh of the curriculum sidebar tree and status badges |

---

## Sidebar & Editor Integration

### Activity Bar Sidebar Tree

Click the Terralings icon on the VS Code Activity Bar to reveal the **Curriculum & Exercises** view:

- **Chapters**: Collapsible chapters showing overall chapter progress.
- **Exercises**: Click any exercise to open the file in the editor.
- **Inline Actions**: Hover over any exercise to see quick action buttons:
  - `$(play)` Verify exercise
  - `$(lightbulb)` Show hint
  - `$(discard)` Reset to template

### Editor Title Actions

When viewing any `.tf` or `.tftest.hcl` file in an exercise folder, action buttons appear in the upper-right corner of the editor title bar:
- Click **▶** to verify the active file.
- Click **💡** to view hints for the active file.

---

## Troubleshooting

### Issue: "Terralings binary not found in PATH"

**Symptom**: The extension warns that it cannot find the `terralings` executable.

**Resolution**:
1. Confirm where `terralings` is installed by running `which terralings` in your terminal.
2. Open VS Code Settings (`Cmd+,` / `Ctrl+,`).
3. Set `terralings.binaryPath` to the absolute path (e.g., `/usr/local/bin/terralings` or `/Users/username/go/bin/terralings`).

### Issue: LSP Diagnostics Not Appearing

**Symptom**: Syntax errors or hover hints are not displayed in `.tf` files.

**Resolution**:
1. Ensure `terralings.enableLsp` is set to `true`.
2. Ensure you have an IaC engine (`tofu` or `terraform`) installed and working (`terralings doctor`).
3. Open the VS Code Output panel (`Cmd+Shift+U` / `Ctrl+Shift+U`) and select **Terralings Language Server** from the dropdown to view LSP connection logs.
