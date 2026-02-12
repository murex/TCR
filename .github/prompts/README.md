# GitHub Copilot Prompt Files

This directory contains custom prompt files (slash commands) for GitHub Copilot to simplify common tasks in the TCR project.

## Available Prompts

### `/upgrade-golang-version`

Upgrades the Go version across the entire project from 1.n to 1.n+1.

**Usage:**
```
/upgrade-golang-version from=1.25 to=1.26
```

**What it does:**
- Updates all `go.mod` files (src, tcr-doc, examples)
- Updates all `go.work` workspace files
- Updates GitHub workflow files (.github/workflows/)
- Updates golangci-lint configuration
- Updates Bazel MODULE files
- Updates documentation (AGENTS.md, dev-doc/build.md)
- Rebuilds golangci-lint with the new Go version
- Runs full build verification (`make prepare`)
- Verifies all Go examples build and test successfully

**Prerequisites:**
- Go installed on your system
- golangci-lint installed
- Make installed
- Bazel installed (for go-bazel example)

## How to Use Prompt Files

1. **In Chat View**: Type `/` followed by the prompt name (e.g., `/upgrade-golang-version`)
2. **Command Palette**: Run `Chat: Run Prompt` and select a prompt
3. **From Editor**: Open the `.prompt.md` file and click the play button

## Creating New Prompts

To add a new prompt file:

1. Create a new `.prompt.md` file in this directory
2. Add YAML frontmatter with metadata (description, agent, tools)
3. Write the prompt instructions in Markdown format
4. Test the prompt using the editor play button

See [VS Code Prompt Files Documentation](https://code.visualstudio.com/docs/copilot/customization/prompt-files) for detailed guidance.

## Best Practices

- Use clear, descriptive names for prompts
- Include input parameter hints in the frontmatter
- Provide step-by-step instructions
- Reference project-specific documentation
- Include quality gates and verification steps
- Document common issues and solutions
- Test prompts before committing

## Related Resources

- [AGENTS.md](../../AGENTS.md) - Comprehensive project documentation
- [.github/copilot-instructions.md](../copilot-instructions.md) - Project-wide Copilot instructions