# gojira

Small cli helper for adding/updating work logs in Jira / Tempo.
Based on [urfave/cli](https://github.com/urfave/cli/), [charmbracelet/huh](https://github.com/charmbracelet/huh)
and [rivo/tview](https://github.com/rivo/tview)

## Features

`gojira`

Argument-less call will try to detect jira issue from git branch name. If detected it will automatically
perform `gojira log ISSUE`, otherwise will display issues from `gojira issues` command.

`gojira issues`

Displays last 5 jira issues that were recently updated and are assigned to you. Select one to update worklog with
additional time spent.

`gojira worklogs`

Displays today's work logs - select one to edit existing work log.

`gojira log ISSUE [TIME_SPENT]`

Adds or updates existing `ISSUE` work log with given `TIME_SPENT`. Is `TIME_SPENT` is not provided you will be prompted
for it.

- `ISSUE` could be straight Issue Key like `TICKET-999`, jira url
  like `https://instance.atlassian.net/browse/TICKET-999`
  or any other string containing single issue key. Uppercase is important!
- `TIME_SPENT` accepts jira format like `1h30m / 2h 20m`.

## Installation

[Check releases page](https://github.com/jzyinq/gojira/releases)
or clone repository and run `make install`.

## Configuration

`gojira` needs a couple of env variables right now that you have to configure:

- Export below values in your .bashrc / .zshrc / .profile file:

```
export GOJIRA_JIRA_INSTANCE_URL="https://<INSTANCE>.atlassian.net"
export GOJIRA_JIRA_LOGIN="your@email.com"
export GOJIRA_JIRA_TOKEN="generate at https://id.atlassian.com/manage-profile/security/api-tokens"
export GOJIRA_TEMPO_TOKEN="generate at https://<INSTANCE>.atlassian.net/plugins/servlet/ac/io.tempo.jira/tempo-app#!/configuration/api-integration"
```  

- Now we need to fetch one last env variable using previously saved values:

`export GOJIRA_JIRA_ACCOUNT_ID=` - fetch it using this curl:

```bash 
curl --request GET \
  --url "$GOJIRA_JIRA_INSTANCE_URL/rest/api/3/user/bulk/migration?username=$GOJIRA_JIRA_LOGIN" \
  --header "Authorization: Basic $(echo -n $GOJIRA_JIRA_LOGIN:$GOJIRA_JIRA_TOKEN | base64)"
```

If you receive `unknown` as `accountId`, just get it from your jira profile url instead:

```
https://<INSTANCE>.atlassian.net/jira/people/<ACCOUNT_ID>
```

Just remember to urldecode it. Save it and you should ready to go!

## [Changelog](./CHANGELOG.md)

## Todo list

- [ ] delete worklog through simple cli version for today
- [ ] ticket status change prompt after logging time
- [ ] tests
- [ ] unify workLogs and worklogsIssues structs - use one for both
  - Reduce jira/tempo spaghetti and unnecessary structs and functions
- [ ] godtools cli semantics update
  - `gojira log -i TICKET` -> `gojira log -i TICKET`
  - `gojira log -i TICKET -t 1h30m`
  - `gojira` -> `gojira recent`
  - `gojira` -> `gojira --help`
- [ ] trigger ui updates after worklog change more efficiently
- [x] cli version does not update worklogs if they exist already
  - [x] fetch worklogs from current day and propose them for selection
- [x] Add worklogs from ui
- [x] gojira worklog delete option
- [x] recent jira task list for easy time logging
- [x] delete worklogs
- [x] error handling
- [x] call for worklogs for whole ~week~ month instead of day
- [x] show calendar with colorized dates
    - fix colors git
    - [x] show by colors if worklog is incomplete/full/overhours for date-
- [x] accept full jira url in `gojira log` and scrap issue key from it
- [x] prompt validation
- [x] while logging time check if worklog exists - if yes, append logged time (config.UpdateExistingWorkLog)
- [x] cli help arguments & better handling
- [x] more details on worklog list - goroutine details fetching?
- [x] interactive edit worklog prompt
- [x] detect git branch name (jira ticket)
- [x] display todays logged working hours
- [x] NewWorklog view - add input for date and date period optionally
- [x] Remove app.ui.flex from the picture
- [x] Hour summary to present day without counting worklogs from the future
- [x] While deleting freshly set worklog, fetch it's data from jira to delete it properly - currently there is:
  ```
  The worklog has either been deleted or you no longer have permission to view it
  ```