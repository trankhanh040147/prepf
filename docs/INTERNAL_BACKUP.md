# Internal Directory Backup Documentation

This document records the structure of `internal/` before cloning packages from crush.

## Directory Structure

```
internal/
├── cli/          # CLI commands (prepf-specific)
├── config/       # Config system (prepf-specific)
├── gemini/       # Gemini AI client (prepf-specific)
├── mock/         # Mock interview module (prepf-specific, currently empty)
├── storage/      # Storage layer (prepf-specific)
├── ui/           # prepf-specific UI components
└── util/         # prepf utilities
    └── stringutil/
```

## prepf-Specific Packages (To Restore)

1. **cli/** - CLI commands
   - `root.go`, `version.go`, `config.go`
   - `constants.go`, `root_test.go`, `noninteractive_test.go`

2. **config/** - Config system
   - `config.go`, `config_test.go`, `constants.go`

3. **gemini/** - Gemini AI client
   - `client.go`, `client_config.go`, `client_chat.go`, `constants.go`

4. **mock/** - Mock interview module (empty directory)

5. **storage/** - Storage layer
   - `profile.go`, `profile_test.go`

6. **ui/** - prepf-specific UI components
   - `model.go`, `view.go`, `update.go`, `update_*.go`
   - `keys.go`, `constants.go`, `messages.go`
   - `renderer.go`, `viewport.go`, `spinner.go`
   - `search.go`, `help.go`, `yank.go`
   - `base_test.go`, `piped_test.go`

7. **util/** - prepf utilities
   - `stringutil/fuzzy.go`

## Notes

- All packages use module path: `github.com/trankhanh040147/prepf`
- The `mock/` directory exists but is currently empty
- `util/stringutil/` should be merged with crush's `stringext/` package

