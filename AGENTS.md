# AGENTS.md

Rules for AI coding agents (and people) who work on this repo. Read it before you change anything.

This is wlecki's fork of [biomassa/godoist](https://github.com/biomassa/godoist), a Todoist TUI and CLI in Go. `origin` is the fork, `upstream` is biomassa. Open PRs against the fork. Never open issues, PRs or comments on upstream unless the owner (wlecki) asks for it.

## Layout

- `cmd/godoist`: the CLI (cobra) and the start of the TUI.
- `internal/todoist`: the API client (Todoist API v1), the sync replica and the offline queue.
- `internal/tui`: the Bubble Tea v2 UI. One file per job: `model.go` (Model, Update), `keys.go`, `input.go`, `actions.go`, `sync.go`, `nav.go`, `rows.go`, `view.go`, and one file per dialog or feature.
- `internal/cache`, `internal/state`, `internal/config`: files on disk.
- `scripts/release.sh` and `install.sh`: releases. Don't touch them unless asked.

## Checks

Run all of these before you commit. CI runs the same:

```sh
gofmt -l .        # must print nothing
go vet ./...
go build ./...
go test -race -count=1 ./...
```

## Rules

**Scope**
- One change per PR. A refactor doesn't change behaviour, and a fix or feature doesn't move code around.
- Do what was asked. No drive-by renames, reformatting or "improvements" in other files.
- Don't add a dependency without asking first.
- Don't add abstractions for one caller: no interfaces, options structs or helpers that nothing else needs yet.
- Don't split or reformat existing files unless asked. Some are already big (`view.go`, `noteedit.go`); put a new feature in a new file instead of growing them.

**Code**
- Every new top-level declaration gets a doc comment. Write comments in plain, short English, like the ones that are already there. Say why, not what the next line does.
- Don't leave commented-out code, TODOs without a reason, or debug prints.
- Error strings start lower case. The one exception is the name "Todoist".
- Task and comment writes go through `taskWrite` (`internal/todoist/queue.go`), so they work offline. Don't call the REST API for them directly.
- Files with user data (config, cache, state, queue) are written with `0o600` and their directories with `0o700`.
- The API token lives only in the config file. Never log or print it, and never write it anywhere else. The cache uses a fingerprint of it.

**Tests**
- New behaviour and bug fixes come with a test. A bug fix test fails without the fix.
- Tests never call the real Todoist API. Use `httptest` (see `internal/todoist/client_test.go`) and `t.TempDir()` with `t.Setenv` for files.
- Never skip, delete or weaken a test to get green.
- A new test has a short comment that says what must hold, like the existing ones.
- TUI tests drive the model with the `send(t, m, msg...)` helper. No real terminal.

**Go and dependencies**
- The module path stays `github.com/biomassa/godoist`. Don't rename imports to the fork.
- This is Bubble Tea v2: imports are `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2` and `charm.land/bubbles/v2`, not the v1 `github.com/charmbracelet/...` paths. Copy from existing code.
- Don't change the `go` line in `go.mod`. If imports change, run `go mod tidy`.

**Docs**
- A change that users can see gets a line under `## [Unreleased]` in `CHANGELOG.md`.
- New keys go in the help screen (`helpSections` in `internal/tui/view.go`) and in the README.
- Don't bump the version, tag, or move `[Unreleased]` notes into a release. `scripts/release.sh` does that.

**Commits and PRs**
- Commit titles are plain sentences in the imperative, like "Add a CI workflow for gofmt, vet, build and tests". No Conventional Commits prefixes.
- A PR description says what changed, why, and what was checked. Keep it short.
- Work on a branch. No direct pushes to `main` and no force pushes to shared branches.
