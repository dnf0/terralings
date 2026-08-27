# Getting Started with Terralings

This guide covers everything you need to install, configure, and begin solving exercises with Terralings in under five minutes.

---

## Prerequisites

Before installing Terralings, make sure you have an IaC engine installed on your system.

### 1. Infrastructure-as-Code Engine (Required)

Terralings requires either **OpenTofu** or **Terraform** available in your `$PATH`:

- **OpenTofu** &ge; 1.6.0 *(Recommended)*:
  ```bash
  # macOS (Homebrew)
  brew install opentofu

  # Linux (Debian / Ubuntu)
  sudo apt-get install -y apt-transport-https ca-certificates curl gnupg
  sudo install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://get.opentofu.org/opentofu.gpg | sudo tee /etc/apt/keyrings/opentofu.gpg >/dev/null
  echo "deb [signed-by=/etc/apt/keyrings/opentofu.gpg] https://packages.opentofu.org/opentofu/tofu/any/ any main" | sudo tee /etc/apt/sources.list.d/opentofu.list
  sudo apt-get update && sudo apt-get install -y tofu
  ```
- **Terraform** &ge; 1.5.0:
  ```bash
  # macOS (Homebrew)
  brew tap hashicorp/tap
  brew install hashicorp/tap/terraform
  ```

Verify your installation:
```bash
tofu version
# or
terraform version
```

### 2. Go Toolchain (Optional)

If you plan to build Terralings from source or contribute to development, install **Go** &ge; 1.22:
```bash
go version
```

---

## Installation Methods

Choose your preferred installation method:

=== "Option 1: 1-Line Installer Script (macOS / Linux)"

    The fastest way to install the latest pre-compiled binary:

    ```bash
    curl -sSL https://raw.githubusercontent.com/dnf0/terralings/main/install.sh | bash
    ```

    The script automatically detects your OS (`darwin`, `linux`), architecture (`amd64`, `arm64`), downloads the latest release from GitHub, and installs the binary into `~/.local/bin` (or `/usr/local/bin`).

=== "Option 2: Go Toolchain"

    If you have Go 1.22+ installed, compile and install directly into `$GOPATH/bin`:

    ```bash
    go install github.com/dnf0/terralings/cmd/terralings@latest
    ```

    Make sure `$GOPATH/bin` or `~/go/bin` is in your `$PATH`:
    ```bash
    export PATH="$HOME/go/bin:$PATH"
    ```

=== "Option 3: GitHub Releases"

    Download pre-compiled binaries for macOS (Apple Silicon & Intel), Linux (`amd64`, `arm64`), or Windows from the [GitHub Releases](https://github.com/dnf0/terralings/releases) page:

    ```bash
    # Example for macOS Apple Silicon (ARM64)
    curl -sSL -o terralings.tar.gz https://github.com/dnf0/terralings/releases/latest/download/terralings_darwin_arm64.tar.gz
    tar -xzf terralings.tar.gz
    sudo mv terralings /usr/local/bin/
    ```

=== "Option 4: Build from Source"

    Clone the repository and build using the provided Makefile:

    ```bash
    git clone https://github.com/dnf0/terralings.git
    cd terralings
    make build
    # Binary is placed at ./bin/terralings
    sudo mv ./bin/terralings /usr/local/bin/
    ```

---

## Step-by-Step First Run

Follow these five steps to verify your installation and solve your first exercise.

### Step 1: Pre-Flight Diagnostic Probe (`terralings doctor`)

Run `terralings doctor` to confirm your IaC engine and workspace permissions:

```bash
terralings doctor
```

Output:
```text
🩺 Terralings Doctor Diagnostic Report
────────────────────────────────────────────────────────────

 ✓ IaC Engine Binary
   Found opentofu at /usr/local/bin/tofu (OpenTofu v1.8.0)

 ✓ Provider Plugin Cache
   Plugin cache directory ready at /Users/username/.terralings/plugin-cache

 ✓ Progress Persistence Store
   State store ready at .terralings/state.json

────────────────────────────────────────────────────────────
Diagnostic summary: 3 checks passed, 0 warnings, 0 failures.
Your environment is fully ready for Terralings!
```

### Step 2: Guided Onboarding Tour (`terralings tour`)

Take the interactive 2-minute guided walkthrough to learn the terminal interface, keyboard shortcuts, and exercise mechanics:

```bash
terralings tour
```

You can navigate through each step interactively or specify `--step <1-5>` directly.

### Step 3: Scaffold Exercise Workspace (`terralings init`)

Extract the complete embedded 56-exercise curriculum into a new or existing working directory:

```bash
# In your desired learning directory:
mkdir my-terralings-workspace && cd my-terralings-workspace
terralings init
```

This extracts all 13 chapters into the `exercises/` folder:
```
exercises/
├── 01_primitives/
│   ├── primitives01.tf
│   ├── primitives02.tf
│   └── ...
├── 02_variables/
├── ...
└── 13_governance/
```

### Step 4: Launch Continuous Watch Mode (`terralings watch`)

Start the continuous file watcher:

```bash
terralings watch
```

Terralings begins monitoring `exercises/` using `fsnotify`. It immediately evaluates the first exercise and displays compilation output:

```text
================================================================================
  Exercise: primitives01 (exercises/01_primitives/primitives01.tf)
  Chapter:  01_primitives - HCL Foundations & Core Primitives
  Mode:     validate
================================================================================

[FAIL] Verification failed with diagnostics:

Error: Missing required provider configuration
  on exercises/01_primitives/primitives01.tf line 12:
  12:   required_providers {
  13:   }

Press 'h' for a hint, 'n' to skip, 'r' to re-run, 'q' to quit.
```

### Step 5: Solving Your First Exercise (`primitives01`)

Open `exercises/01_primitives/primitives01.tf` in your code editor:

```hcl
// I AM NOT DONE
// ============================================================================
// Exercise: primitives01
// Chapter:  01_primitives - HCL Foundations & Core Primitives
// ============================================================================

terraform {
  # TODO: Specify required_version >= 1.6.0
  # TODO: Declare required_providers with local provider from hashicorp/local
  required_providers {
    # Add local provider requirement here
  }
}
```

Fix the configuration by completing the `terraform` block:

```hcl
terraform {
  required_version = ">= 1.6.0"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4"
    }
  }
}
```

Save the file. Within 30 milliseconds, the watcher automatically re-evaluates `primitives01`, reports success, and advances to `primitives02`!

---

## Shell Autocompletions

Terralings provides rich autocompletion scripts for Bash, Zsh, Fish, and PowerShell, including interactive exercise name completions.

=== "Zsh"

    Add to your `~/.zshrc`:
    ```bash
    # If not already present in ~/.zshrc:
    autoload -Uz compinit && compinit

    source <(terralings completion zsh)
    ```

=== "Bash"

    Add to your `~/.bashrc`:
    ```bash
    source <(terralings completion bash)
    ```

=== "Fish"

    Generate the fish completion configuration:
    ```bash
    terralings completion fish > ~/.config/fish/completions/terralings.fish
    ```

=== "PowerShell"

    Add to your PowerShell `$PROFILE`:
    ```powershell
    terralings completion powershell | Out-String | Invoke-Expression
    ```

---

## Editor Setup

Terralings provides two complementary ways to integrate with your code editor:

### Option A: VS Code Companion Extension

If you use **Visual Studio Code**, install the official companion extension:
- **Activity Bar View**: Browse the 13-chapter curriculum tree with live passing/failing status icons.
- **Embedded LSP Client**: Automatic in-editor syntax diagnostics, hover hints, and error squiggles.
- **Interactive Walkthroughs**: Built-in 5-step interactive onboarding guides.
- **One-Click Runners**: Launch `watch`, `tui`, `doctor`, and `hint` with dedicated editor buttons.

See the [VS Code Extension Guide](vscode-extension.md) for full installation and configuration instructions.

### Option B: Native Language Server (`terralings lsp`)

For **Neovim**, **Helix**, or any editor supporting the Language Server Protocol (LSP), connect directly to `terralings lsp`:

=== "Neovim (`nvim-lspconfig`)"

    ```lua
    local lspconfig = require('lspconfig')
    local configs = require('lspconfig.configs')

    if not configs.terralings then
      configs.terralings = {
        default_config = {
          cmd = { 'terralings', 'lsp' },
          filetypes = { 'terraform', 'hcl' },
          root_dir = lspconfig.util.root_pattern('exercises', '.terralings'),
          settings = {},
        },
      }
    end

    lspconfig.terralings.setup{}
    ```

=== "Helix (`languages.toml`)"

    ```toml
    [[language]]
    name = "hcl"
    scope = "source.hcl"
    language-servers = ["terralings-lsp"]

    [language-server.terralings-lsp]
    command = "terralings"
    args = ["lsp"]
    ```

---

## Next Steps

Now that your environment is configured:
- Read the **[Curriculum Syllabus](syllabus.md)** to see all 13 chapters and 56 exercises.
- Explore all CLI commands in the **[CLI Reference](cli-reference.md)**.
- Dive deep into the **[Onboarding Guide](onboarding-guide.md)** for learning tips and advanced workflows.
