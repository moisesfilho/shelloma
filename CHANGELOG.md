# Changelog

All notable changes to Shelloma are documented in this file.

## [1.4.0] - 2026-08-04

### Added

- **Command Learning**: new `shelloma learn <cmd>` command (alias `aprender`) that captures the help output of any CLI command available on the machine and saves it locally under `~/.config/shelloma/learned/<command>.json`.
- **Learned-Help Prompt Injection**: when a request mentions a learned command, its help reference is automatically injected into the Ollama system prompt, improving command accuracy.
- `CHANGELOG.md` documenting project changes.

### Fixed

- **Direct-command mode output overwriting**: when running `shelloma "instruction"` (non-interactive), the terminal app's fixed footer positioning (`ClearInputArea`, `DrawInputSeparator`, `DrawLegendAtBottom`, `MoveToInputLine`, `MoveToContentStart`) emitted absolute cursor escapes that cleared/overwrote the returned information. These decorations are now disabled in direct mode via a new `ui.IsTerminalApp` flag.
- **Direct-command mode is now single-shot**: after executing the suggested command, Shelloma exits immediately instead of showing the "Anything else?" (`Prompt:`) continuation loop. The interactive terminal app (run without arguments) keeps the continuous loop.
- **Direct-command mode legend**: the options legend is now displayed right below the "Choose an option" prompt when there is no fixed footer.

## [1.3.0] - 2026-07-31

### Added

- `shelloma config docs`: compact configuration schema & commands reference.
- **Natural language configuration**: requests like "change model to llama3" are parsed by the local AI into the matching `shelloma config set` command.
- Desktop mode (`--desktop`), continuous command loop, rich keyboard shortcuts, Snap packages and Bash completion.
- CI/CD pipeline documentation and new Makefile targets (icons, build-all, deb, tar, appimage, flatpak, install-user).
- Go requirement updated from 1.20+ to 1.23+.

## [1.2.2] - 2026-07-25

- Updated Windows build process.

## [1.2.1] - 2026-07-24

- Fixed a bug in console input reading.
