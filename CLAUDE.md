# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build       # builds binary to dist/gojira
make install     # builds and installs to /usr/local/bin
make tests       # runs tests with coverage
make fixer       # runs golangci-lint static analysis
```

Run a single test:
```bash
go test ./gojira/... -run TestName
```

Debug logs are written to `/tmp/gojira.log` at runtime.

## Architecture

`gojira` is a TUI/CLI tool for logging work to Jira + Tempo. Entry point is `main.go` → `gojira.Run()`.

**Global state:** A single `app gojira` package-level variable (in `gojira/gojira.go`) holds all runtime state: API clients, UI references, current time, worklogs, and a `sync.RWMutex`. Everything mutates this singleton.

**Two modes of operation:**
- **CLI mode** (non-interactive): `log`, `issues`, `view`, `config` commands — use `charmbracelet/huh` spinners and forms, return to terminal after completing.
- **TUI mode** (interactive): `worklogs` command — launches a `rivo/tview` full-screen app with a grid layout: day view (left), summary (top-right), calendar (bottom-right).

**API layer:** Two separate clients with interfaces defined in `interfaces.go`:
- `JiraClient` (`jira.go`) — Jira REST API v2/v3 for issues and worklog CRUD
- `TempoClient` (`tempo.go`) — Tempo REST API v4 for fetching/updating/deleting worklogs

Both are initialized after config loads in the `cli.Before` hook. Mocks live in `mocks.go` and use `testify/mock`.

**Key data flow:**
1. Worklogs fetched from Tempo (`GetWorklogs`) into `app.workLogs` (`Worklogs` struct)
2. Issue details fetched from Jira and joined into `app.workLogsIssues` (`WorklogsIssues` struct)  
3. Mutations (create/update/delete) call the appropriate API then update both slices via `app.mu.Lock()`

**Config** (`config.go`): Loaded from env vars (`GOJIRA_JIRA_INSTANCE_URL`, `GOJIRA_JIRA_LOGIN`, `GOJIRA_JIRA_TOKEN`, `GOJIRA_JIRA_ACCOUNT_ID`, `GOJIRA_TEMPO_TOKEN`). `UpdateExistingWorklog: true` by default — when logging to a ticket that already has a worklog today, time is summed rather than creating a new entry.

**TUI components** (all in `gojira/` package): `Calendar`, `Summary`, `DayView`, `LoaderView`, `ErrorView` — each implements a `tview` primitive. The `UserInterface` struct in `ui.go` assembles them into the grid.

## Linting

Line length limit is 120 chars, function length limit is 65 lines / 40 statements. Use `//nolint:lll` for unavoidable long lines (e.g., in `cli.go` for the config help text).
