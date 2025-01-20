# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.11.0] - 2025-01-20
### Changed
- Tempo API version from 3 to 4 due to incoming shutdown of the old version

## [0.10.2] - 2025-01-08
### Fixed
- hardcoded `2024` while fetching holidays for given year

## [0.10.1] - 2024-08-11
### Fixed
- uppercase `L` was used instead of lowercase version for latest issues

## [0.10.0] - 2024-07-26
### Added
- `l` shortcut in `gojira worklogs` to reload latest issues

## [0.9.0] - 2024-06-03
### Added
- Jump by a month on calendar using shift + left/right arrow key

## [0.8.1] - 2024-05-08
### Fixed
- Calendar control stopped working after introducing delete on worklog list

## [0.8.0] - 2024-05-08
### Added
- `gojira worklogs` search bar - looks for a text or issue key if passed uppercased like `ISSUE-123`
- `Enter` now submits worklog straight from time spent input, `Delete` removes it
- `Delete` also works on selected worklog on the list - no confirmation required though so watch out

### Changed
- `SetFocusFunc`/`SetBlurFunc` now handles decorated windows instead of original mess

## [0.7.0] - 2024-05-05
### Added
- Fetch national holidays from [date.nager.at](https://date.nager.at) if LC_TIME is present in environment. Holidays will be marked on calendar and excluded from month summary.

## [0.6.0] - 2024-03-30
### Changed
- Replace [manifoldco/promptui](https://github.com/charmbracelet/huh) with [charmbracelet/huh](https://github.com/charmbracelet/huh) due to lack of maintainer
- code cleanup pass, simplified structs, improved messaging, aligned variable naming
- Include already logged issues in recent issues list - allow editing existing time
- Add issues with worklogs for current day while launching `gojira issues`

## [0.5.4] - 2024-03-29
### Fixed
- MonthRange function returned first day of next month which causes invalid summaries

## [0.5.3] - 2024-03-27
### Fixed
- calendar controls not working while focus is on latest issues view

### Changed
- Set focus on time spent field while adding new worklog

## [0.5.2] - 2024-03-26
### Changed
- More detailed loader while adding worklogs in a batch

### Fixed
- Version number is now properly synced with the latest release

## [0.5.1] - 2024-03-25
### Fixed
- Summary not updating with worklog changes

## [0.5.0] - 2024-03-25
### Added
- Summary now shows time diff for worklogs

### Fixed
- Loader flickering while moving faster through calendar

## [0.4.0] - 2024-03-23
### Added
- `gojira worklogs` now have a calendar which tracks month of currently selected date 
  - days are colored depending on time logged
    - `white` are without any logs 
    - `yellow` are for incomplete logs 8h hours is considered as full day 
    - `purple` shows days with exceeded worklogs (> 8 hours) 
    - `grey` is for weekends by default 
  - calendar also shows currently log time against estimated work hours for whole month
- error modal for nicer error display in `gojira worklogs`
- loader modal for handling http requests in `gojira worklogs`

### Changed
- refactor pass which cleans up a bit ui functions
- app time is based on the UTC instead of local time - it's a tentative fix for near midnight scenarios
- extracted Jira & Tempo API calls to separate packages
- use latest go (1.22)

### Fixed
- Mostly `gojira worklogs` fixes: 
  - UI is now based on grid instead of flex - should be more responsive 
  - Rearrange UI elements for better readability

## [0.3.1] - 2022-03-09
### Changed
- Replace multiple methods for fetching issue key with single `ResolveIssueKey` method

## [0.3.0] - 2022-02-21
### Added
- `gojira view` command that opens up issue in default browser - works with issue key as arg and fetched from git branch name

### Changed
- default behavior of argument-less `gojira` call if git branch is detected, now it allows to select do you want to log worklog or view issue in browser

## [0.2.2] - 2021-05-04
### Fixed
- time spent input now properly handles lack of whitespace between time parts - `1h30m`

## [0.2.1] - 2021-05-04
### Fixed
- `gojira worklogs` now properly reports overall time spent after editing work log

## [0.2.0] - 2021-05-04
### Added
- `gojira log` improvements:
  - command verifies entered issue key by asking jira api about details
  - details are displayed along logging time or before time spent prompt
  - command accepts as TICKET also jira url like `https://instance.atlassian.net/browse/TICKET-999`
    basically any string containing something that looks like jira ticket will be accepted

## [0.1.0] - 2021-05-03
### Added
- Initial release of gojira

[Unreleased]: https://github.com/jzyinq/gojira/compare/0.11.0...master
[0.11.0]: https://github.com/jzyinq/gojira/compare/0.10.2...0.11.0
[0.10.2]: https://github.com/jzyinq/gojira/compare/0.10.1...0.10.2
[0.10.1]: https://github.com/jzyinq/gojira/compare/0.10.0...0.10.1
[0.10.0]: https://github.com/jzyinq/gojira/compare/0.9.0...0.10.0
[0.9.0]: https://github.com/jzyinq/gojira/compare/0.8.1...0.9.0
[0.8.1]: https://github.com/jzyinq/gojira/compare/0.8.0...0.8.1
[0.8.0]: https://github.com/jzyinq/gojira/compare/0.7.0...0.8.0
[0.7.0]: https://github.com/jzyinq/gojira/compare/0.6.0...0.7.0
[0.6.0]: https://github.com/jzyinq/gojira/compare/0.5.4...0.6.0
[0.5.4]: https://github.com/jzyinq/gojira/compare/0.5.3...0.5.4
[0.5.3]: https://github.com/jzyinq/gojira/compare/0.5.2...0.5.3
[0.5.2]: https://github.com/jzyinq/gojira/compare/0.5.1...0.5.2
[0.5.1]: https://github.com/jzyinq/gojira/compare/0.5.0...0.5.1
[0.5.0]: https://github.com/jzyinq/gojira/compare/0.4.0...0.5.0
[0.4.0]: https://github.com/jzyinq/gojira/compare/0.3.1...0.4.0
[0.3.1]: https://github.com/jzyinq/gojira/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/jzyinq/gojira/compare/0.2.2...0.3.0
[0.2.2]: https://github.com/jzyinq/gojira/compare/0.2.1...0.2.2
[0.2.1]: https://github.com/jzyinq/gojira/compare/0.2.0...0.2.1
[0.2.0]: https://github.com/jzyinq/gojira/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/jzyinq/gojira/tree/0.1.0
