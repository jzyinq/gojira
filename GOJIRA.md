## Goal

Separate Text UI (TUI) code from application logic by:
- Renaming `gojira/` to `pkg/`
- Creating two clear packages: `pkg/app` (domain + services) and `pkg/tui` (all `tview` UI)
- Avoiding circular dependencies and reducing coupling to enable future non-TUI UIs (CLI forms, daemon, API)

## Current state overview

The project keeps all code in `gojira/` with tight coupling via a global `app` variable that mixes UI state and application state.

### TUI/UI code (rivo/tview, visuals, input handling)
- `gojira/ui.go` (root layout and `UserInteface`)
- `gojira/calendar.go`
- `gojira/dayview.go`
- `gojira/summary.go`
- `gojira/loaderview.go`
- `gojira/errorview.go`
- `gojira/custommodal.go`
- `gojira/prompt.go` (huh console forms – still UI)

### Application logic (domain + services)
- `gojira/jira.go` (Jira client, domain `Issue`)
- `gojira/tempo.go` (Tempo client)
- `gojira/worklog.go` (domain `Worklog`, aggregate, use-cases)
- `gojira/http.go` (HTTP helper)
- `gojira/hoilidays.go` (holidays fetch + helpers)
- `gojira/config.go` (env config – global)
- `gojira/utils.go` (mixed UI + domain helpers)
- `gojira/cli.go` (urfave CLI commands – currently mixes UI and domain and even defines a domain method)
- `gojira/gojira.go` (bootstrapping, global `app`, CLI app construction)

### Examples of problematic coupling

1) Global mutable `app` contains both UI and domain state:

```11:18:/home/jzy/projects/private/gojira/gojira/gojira.go
type gojira struct {
    cli            *cli.App
    ui             *UserInteface
    time           *time.Time
    holidays       *Holidays
    workLogs       Worklogs
    workLogsIssues WorklogsIssues
}
```

```77:77:/home/jzy/projects/private/gojira/gojira/gojira.go
var app gojira
```

2) A domain operation is implemented inside `cli.go` and mutates global UI/domain state:

```274:300:/home/jzy/projects/private/gojira/gojira/cli.go
func (issue Issue) LogWork(logTime *time.Time, timeSpent string) error {
    logrus.Infof("Logging %s of time to ticket %s at %s", timeSpent, issue.Key, logTime)
    todayWorklog, err := app.workLogs.LogsOnDate(logTime)
    if err != nil {
        return err
    }
    if Config.UpdateExistingWorklog {
        for index, workLog := range todayWorklog {
            if strconv.Itoa(workLog.Issue.Id) == issue.Id {
                timeSpentSum := FormatTimeSpent(TimeSpentToSeconds(timeSpent) + workLog.TimeSpentSeconds)
                err := todayWorklog[index].Update(timeSpentSum)
                if err != nil {
                    return err
                }
                return nil
            }
        }
    }
    worklog, err := NewWorklog(issue.GetIdAsInt(), logTime, timeSpent)
    if err != nil {
        return err
    }
    // add this workload to global object
    app.workLogs.logs = append(app.workLogs.logs, &worklog)
    app.workLogsIssues.issues = append(app.workLogsIssues.issues, WorklogIssue{Issue: issue, Worklog: &worklog})
    return nil
}
```

3) UI concerns leak into helpers (`tcell` colors) used by domain/TUI:

```45:55:/home/jzy/projects/private/gojira/gojira/utils.go
func GetTimeSpentColor(timeSpentInSeconds int, hours int) tcell.Color {
    switch {
    case timeSpentInSeconds < hours*60*60 && timeSpentInSeconds > 0:
        return tcell.ColorOrange
    case timeSpentInSeconds == hours*60*60:
        return tcell.ColorGreen
    case timeSpentInSeconds > hours*60*60:
        return tcell.ColorBlue
    default:
        return tcell.ColorWhite
    }
}
```

4) TUI directly reaches into services and global state instead of a boundary:

```195:219:/home/jzy/projects/private/gojira/gojira/dayview.go
func (d *DayView) loadLatest() {
    d.latestIssuesStatus.SetText("Latest issues").SetDynamicColors(true)
    issues, err := NewJiraClient().GetLatestIssues()
    if err != nil {
        app.ui.errorView.ShowError(err.Error(), nil)
        return
    }
    d.latestIssuesList.Clear()
    // ...
    d.latestIssuesList.SetSelectedFunc(func(row, column int) {
        NewAddWorklogForm(d, issues.Issues, row)
    })
}
```

## Risks and code smells to address

- Global mutable `app` couples UI and domain state; difficult to test and parallelize
- Domain operations inside CLI/TUI modules (e.g., `Issue.LogWork` in `cli.go`)
- UI-specific types (`tcell`) referenced from helpers used by domain
- Services (Jira/Tempo) constructed ad-hoc (`NewJiraClient()`) inside UI
- Data races: `NewWorklogIssues` appends to shared slices from goroutines without synchronization
- `worklogs.Delete` directly manipulates global slices with pointer aliasing (commented as buggy)
- `Config` is global and `GetEnv` calls `os.Exit(1)` (hard to test)
- `ResolveIssueKey` depends on `urfave/cli.Context` in a utility function
- Hardcoded paths (log `/tmp/gojira.log`) and magic numbers (8h for day)
- Misspelling `UserInteface` (rename in refactor)

These will block clean separation unless addressed.

## Target architecture

Layered with clear package boundaries:

- pkg/app – pure domain + services
  - Domain models: `Issue`, `Worklog`, aggregates (`Worklogs`, `WorklogsIssues`), `Holidays`
  - Use-cases/services via interfaces:
    - `JiraService` (GetIssue, GetIssuesByJQL, GetLatestIssues, Create/Update/DeleteWorklog, GetIssuesByKeys)
    - `TempoService` (GetWorklogs, UpdateWorklog, DeleteWorklog)
    - `HolidayService` (Fetch by country)
    - `Clock` (Now, used for ranges)
    - Optional `Logger`
  - Implementations: HTTP-based clients inside app (or `pkg/app/jirahttp`, `pkg/app/tempohttp`)
  - Helpers: time range utils, parse/format (no UI types), FindIssueKeyInString
  - App state: `type State` holding current time, cached worklogs/issues/holidays with methods to refresh

- pkg/tui – all `tview` UI
  - Widgets: calendar, day view, summary, loader, error, modal
  - Styling helpers (color mapping), key handling, async loading guards
  - A single entry point: `func Run(state *app.State, services Services) error`
    - `Services` is a small interface the TUI needs to call (bridging to app)

- pkg/cli (optional now or later)
  - Console forms (huh), spinners, and urfave CLI commands that do not require TUI
  - A `New()` returning `*cli.App` wired to app services

- cmd/gojira – main
  - Bootstrap config, logger, services, state; run CLI app and/or TUI

Critical rule: `pkg/app` must not import `pkg/tui` or UI libs. `pkg/tui` imports `pkg/app` types only.

## Boundary contracts (sketch)

- App-facing services used by UI/CLI:
  - `type Services interface {
      FetchWorklogs(from, to time.Time) (app.Worklogs, error)
      FetchLatestIssues() ([]app.Issue, error)
      LogWork(issue app.Issue, date time.Time, timeSpent string) (app.Worklog, error)
      UpdateWorklog(w *app.Worklog, timeSpent string) error
      DeleteWorklog(w *app.Worklog) error
      SearchIssues(jql string, max int) ([]app.Issue, error)
    }`
  - The concrete implementation adapts to `JiraService` and `TempoService`.

- `pkg/tui` owns all `tview` state and never touches globals.

## Migration plan (incremental, low-risk)

Phase 1 – Create new packages and move pure domain/util code
- Create `pkg/app`
  - Move: `jira.go`, `tempo.go`, `worklog.go`, `hoilidays.go` (rename to `holidays.go`), `http.go`
  - Split `utils.go`:
    - Keep domain-safe helpers: `FormatTimeSpent`, `CalculateTimeSpent`, `FindIssueKeyInString`, `MonthRange`, `DayRange`
    - Move UI helpers to `pkg/tui/style.go`: `GetTimeSpentColor`, `GetTimeSpentColorTag`
    - Move CLI-specific helper to `pkg/cli/resolve.go`: `ResolveIssueKey` (remove `*cli.Context` arg; make it accept a string or args)
- Replace global `Config` with a constructor: `app.LoadConfig() (Config, error)` and pass `Config` into service constructors
- Extract `Issue.LogWork` out of UI/CLI: make it a use-case in app that returns a `Worklog` and let callers update their own caches. Remove direct mutation of global slices from inside domain methods.

Phase 2 – Introduce app `State` and services interface
- Define `type State` in `pkg/app/state.go` with: `Time time.Time`, `Worklogs`, `WorklogsIssues`, `Holidays`
- Add methods: `LoadWorklogs(from, to)`, `EnsureHolidays(countryCode)`; handle synchronization internally (mutex) to fix data races
- Define `type JiraService`, `type TempoService` interfaces and implement HTTP versions in `pkg/app/jirahttp` and `pkg/app/tempohttp` (or keep in `pkg/app` but behind interfaces)

Phase 3 – Move all TUI files to `pkg/tui` and adapt
- Move: `ui.go`, `calendar.go`, `dayview.go`, `summary.go`, `loaderview.go`, `errorview.go`, `custommodal.go`
- Change references from global `app` to injected `state` and `services`
- Replace calls like `NewJiraClient()` with `services.FetchLatestIssues()` etc.
- Replace direct access to `app.workLogs`, `app.workLogsIssues` with state methods
- Keep UI-only concurrency constructs (e.g., `loadingWorklogs` channel) within `pkg/tui`

Phase 4 – Isolate console forms and CLI
- Move `prompt.go` into `pkg/cli/prompt` (still UI but console, not TUI)
- Update `cli.go` to:
  - Not define domain methods
  - Use app services and console forms only
  - Add a dedicated `tui` command that calls `tui.Run(state, services)`; `worklogs` can delegate to that

Phase 5 – Rename folder and update imports
- `git mv gojira pkg`
- Update imports across the repo:
  - TUI: `import "gojira/pkg/app"`
  - CLI: `import ("gojira/pkg/app"; "gojira/pkg/tui")`
  - main: `import "gojira/pkg/cli/cmd"` (or keep building the `*cli.App` in main for now)
- Fix `main.go` to bootstrap from `pkg/app` and run the CLI

Phase 6 – Cleanups and safety
- Fix misspelling: `UserInteface` -> `UserInterface`
- Remove hardcoded `/tmp/gojira.log` path; make it configurable or disable by default
- Replace `os.Exit(1)` in `GetEnv` with returned errors; handle in main
- Remove magic numbers; define constants (e.g., default workday hours)
- Add synchronization in `State` where slices are mutated concurrently
- Replace brittle slice deletions with ID-based maps or safe helpers

## Detailed file mapping and edits

- pkg/app
  - `config/config.go`: `LoadConfig()` returns `Config`; remove global `Config` var
  - `http/httpclient.go`: `SendHttpRequest` (internal)
  - `jira/client.go`: behind `JiraService`
  - `tempo/client.go`: behind `TempoService`
  - `domain/types.go`: `Issue`, `Worklog`, aggregates
  - `domain/timeutil.go`: `MonthRange`, `DayRange`, `FormatTimeSpent`
  - `domain/parse.go`: `FindIssueKeyInString`, `TimeSpentToSeconds`
  - `state.go`: `State` with mutex, caches, helpers

- pkg/tui
  - `ui.go`, `calendar.go`, `dayview.go`, `summary.go`, `loaderview.go`, `errorview.go`, `custommodal.go`
  - `style.go`: `GetTimeSpentColor`, `GetTimeSpentColorTag`
  - exports: `Run(state *app.State, services Services) error`

- pkg/cli
  - `prompt/prompt.go`: `IssueWorklogForm`, `InputTimeSpentForm`, `SelectActionForm`
  - `cmd/cli.go`: commands using services; one command to launch TUI
  - `resolve.go`: `ResolveIssueKey(args []string, branch string) string`

## Avoiding circular dependencies

- `pkg/app` must not import `pkg/tui` or UI libs; only stdlib + optional logging interface
- `pkg/tui` imports `pkg/app` types and `Services` interface
- `pkg/cli` imports `pkg/app` and optionally `pkg/tui`
- `cmd` imports `pkg/cli` (or builds CLI directly) and `pkg/app`

## Concurrency and data safety changes

- Introduce a mutex in `app.State` protecting `Worklogs` and `WorklogsIssues`
- Ensure updates (append/delete) happen under lock
- Return copies/snapshots to TUI to render

## Step-by-step execution checklist

1) Introduce `pkg/app` with copies of domain files; update imports locally in those files
2) Extract `Issue.LogWork` from `cli.go` into `pkg/app` use-case and remove global mutations
3) Split `utils.go` into domain-safe and UI-specific helpers; move color functions to `pkg/tui`
4) Create `pkg/tui` and move all tview files; change references to use injected `state`/`services`
5) Create `pkg/cli/prompt` and move console forms; update `cli.go` to only orchestrate and call services/forms; add `tui` command
6) Replace `GetEnv`/`Config` globals with `LoadConfig()` and pass config to services
7) Introduce `app.State` and move caching there; add locks; update call sites
8) `git mv gojira pkg` and fix all import paths
9) Build and run; fix compile errors
10) Manually test: TUI navigation, add/update/delete worklogs, issues search, log work via CLI-only flow

## Backout plan

This plan is incremental. If any phase introduces instability, revert to the previous commit. The final folder rename should be done after code compiles in the new structure to minimize diff blast radius.

## Notes on untracked `gojira/tui/`

There is an untracked `gojira/tui/` folder mirroring some UI files. Prefer moving authoritative UI files to `pkg/tui` and discard duplicates unless they contain newer changes. Diff them before moving.


