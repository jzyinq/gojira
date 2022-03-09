# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/jzyinq/gojira/compare/0.3.1...master
[0.3.1]: https://github.com/jzyinq/gojira/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/jzyinq/gojira/compare/0.2.2...0.3.0
[0.2.2]: https://github.com/jzyinq/gojira/compare/0.2.1...0.2.2
[0.2.1]: https://github.com/jzyinq/gojira/compare/0.2.0...0.2.1
[0.2.0]: https://github.com/jzyinq/gojira/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/jzyinq/gojira/tree/0.1.0