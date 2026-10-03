# StudyCTL

A local-first Linux CLI application for tracking study and work time.

## Overview

`studyctl` is a terminal-based productivity tracker designed for Linux.
It records study and work sessions locally, calculates
daily/weekly/monthly productivity, tracks daily goals, and provides a
historical record of how time was spent.

The application is intentionally local-first:

-   No cloud service is required.
-   No account is required.
-   Data is stored locally in SQLite.
-   The application can run completely offline.
-   The final application is distributed as a single Go binary.

## Goals

The application is designed to answer four questions:

1.  What am I doing right now?
2.  How much time did I spend studying and working?
3.  What exactly did I spend that time on?
4.  Am I meeting my productivity targets?

## Technology Stack

  Component          Technology
  ------------------ ----------------------------------------------
  Language           Go
  CLI framework      Cobra
  Database           SQLite
  SQLite driver      `modernc.org/sqlite`
  Terminal UI        Standard terminal output / optional Lipgloss
  Operating system   Linux
  Storage            Local filesystem

## Requirements

Install:

-   Go 1.22 or newer
-   Git
-   Linux

Verify Go:

``` bash
go version
```

## Installation from Source

Clone the project:

``` bash
git clone https://github.com/yourname/studyctl.git
cd studyctl
```

Initialize dependencies:

``` bash
go mod tidy
```

Build:

``` bash
go build -o studyctl .
```

Run:

``` bash
./studyctl --help
```

Install system-wide:

``` bash
sudo mv studyctl /usr/local/bin/
```

Verify:

``` bash
studyctl version
```

## Project Structure

``` text
studyctl/
├── cmd/
│   ├── root.go
│   ├── start.go
│   ├── stop.go
│   ├── status.go
│   ├── log.go
│   ├── today.go
│   ├── week.go
│   ├── month.go
│   ├── history.go
│   ├── delete.go
│   ├── goal.go
│   ├── version.go
│   ├── completion.go
│   └── helpers.go
│
├── internal/
│   ├── database/
│   │   ├── database.go
│   │   └── goals.go
│   │
│   ├── session/
│   │   ├── session.go
│   │   └── report.go
│   │
│   └── config/
│       └── config.go
│
├── main.go
├── go.mod
├── go.sum
└── README.md
```

## Data Storage

The SQLite database is stored at:

``` text
~/.studyctl/studyctl.db
```

The application creates the directory automatically.

``` text
~/.studyctl/
└── studyctl.db
```

### Sessions Table

The sessions table contains:

  Column               Description
  -------------------- --------------------
  `id`                 Unique session ID
  `task`               Activity name
  `type`               `study` or `work`
  `start_time`         Session start
  `end_time`           Session end
  `duration_seconds`   Completed duration
  `notes`              Optional notes

Example:

``` text
id: 42
task: Machine Learning
type: study
start_time: 2026-10-03 18:30:00
end_time: 2026-10-03 19:42:00
duration_seconds: 4320
notes: Finished Transformer architecture
```

### Goals Table

The goals table contains:

  Column             Description
  ------------------ ----------------
  `id`               Unique goal ID
  `date`             Goal date
  `target_seconds`   Daily target

## Command Reference

The complete command structure is:

``` text
studyctl
├── start
├── stop
├── status
├── log
├── today
├── week
├── month
├── history
├── delete
├── goal
│   ├── set
│   └── show
├── version
└── completion
    ├── bash
    └── zsh
```

------------------------------------------------------------------------

# `studyctl`

Display the main help page.

``` bash
studyctl --help
```

## Usage

``` text
studyctl [command]
```

## Examples

``` bash
studyctl start "Machine Learning"
studyctl stop
studyctl today
studyctl week
studyctl month
studyctl goal set 8h
```

------------------------------------------------------------------------

# `studyctl start`

Starts a new study or work session.

## Usage

``` bash
studyctl start <task> [flags]
```

## Arguments

  Argument     Required Description
  ---------- ---------- ----------------------
  `<task>`          Yes Name of the activity

## Flags

  Flag           Default   Description
  -------------- --------- ---------------------------------
  `-t, --type`   `study`   Session type: `study` or `work`

## Examples

Start studying:

``` bash
studyctl start "Machine Learning"
```

Start AWS preparation:

``` bash
studyctl start "AWS MLA"
```

Start work:

``` bash
studyctl start "Client Project" --type work
```

Short form:

``` bash
studyctl start "FastAPI" -t study
```

## Behavior

Only one active session is allowed at a time.

If another session is running:

``` text
Error: a session is already running
```

Check the current session:

``` bash
studyctl status
```

Stop it:

``` bash
studyctl stop
```

------------------------------------------------------------------------

# `studyctl stop`

Stops the currently active session.

## Usage

``` bash
studyctl stop [flags]
```

## Flags

  Flag              Description
  ----------------- -------------------------------------
  `--note <text>`   Add a note to the completed session

## Examples

``` bash
studyctl stop
```

With a note:

``` bash
studyctl stop --note "Finished attention mechanism"
```

## Example Output

``` text
Session completed
----------------------------
Task     : Machine Learning
Type     : study
Duration : 1h 12m
Note     : Finished attention mechanism
----------------------------
```

If there is no active session:

``` text
Error: no active session
```

------------------------------------------------------------------------

# `studyctl status`

Displays the currently active session.

## Usage

``` bash
studyctl status
```

## Example

``` text
Current session
----------------------------
Task      : Machine Learning
Started   : 18:42:11
Elapsed   : 47m 31s
```

If nothing is running:

``` text
No active session.
```

## Important Behavior

The active session is not treated as completed historical time until
`stop` is executed.

This prevents an accidentally abandoned session from permanently
inflating reports.

------------------------------------------------------------------------

# `studyctl log`

Manually records a completed session.

This is useful when time was spent on an activity without running
`studyctl start`.

## Usage

``` bash
studyctl log <task> --duration <duration> [flags]
```

## Arguments

  Argument     Required Description
  ---------- ---------- ---------------
  `<task>`          Yes Activity name

## Required Flags

  Flag           Description
  -------------- --------------------------------
  `--duration`   Duration such as `90m` or `2h`

## Optional Flags

  Flag       Default   Description
  ---------- --------- -------------------
  `--type`   `study`   `study` or `work`
  `--note`   empty     Session note

## Supported Duration Formats

``` text
30m
45m
90m
2h
4h
7.5h
```

## Examples

``` bash
studyctl log "Transformer Architecture" --duration 90m
```

``` bash
studyctl log "AWS MLA" \
    --duration 2h \
    --type study
```

``` bash
studyctl log "Client Project" \
    --duration 3h \
    --type work \
    --note "API development"
```

The manually logged session is inserted as an already completed session.

------------------------------------------------------------------------

# `studyctl today`

Displays today's productivity.

## Usage

``` bash
studyctl today
```

## Information Displayed

-   Study time
-   Work time
-   Total time
-   Hourly activity
-   Task breakdown

## Example

``` text
Today — 03 Oct 2026
============================================

Study : 5h 42m
Work  : 2h 13m
Total : 7h 55m

Hourly timeline
--------------------------------------------
09:00  Study 58m     Work 0m
10:00  Study 1h 01m  Work 0m
11:00  Study 47m     Work 0m
13:00  Study 0m      Work 54m
14:00  Study 0m      Work 59m
15:00  Study 43m     Work 20m

Task breakdown
--------------------------------------------
AWS MLA                        study  2h 06m
Machine Learning               study  2h 18m
FastAPI                        study  1h 18m
Client Project                 work   1h 53m
```

------------------------------------------------------------------------

# `studyctl week`

Displays productivity for the current week.

## Usage

``` bash
studyctl week
```

## Information Displayed

-   Activity for each day
-   Study total
-   Work total
-   Overall total

## Example

``` text
Week — 28 Sep to 04 Oct
============================================
Mon 28  7h 21m
Tue 29  8h 04m
Wed 30  6h 47m
Thu 01  9h 12m
Fri 02  7h 56m
Sat 03  7h 55m
Sun 04  0m
--------------------------------------------
Study : 32h 14m
Work  : 15h 01m
Total : 47h 15m
```

The week begins on Monday.

------------------------------------------------------------------------

# `studyctl month`

Displays productivity for the current month.

## Usage

``` bash
studyctl month
```

## Information Displayed

-   Every day of the current month
-   Active days
-   Study total
-   Work total
-   Overall total
-   Average time per active day

## Example

``` text
October 2026
============================================================
01 Thu  7h 42m   ███████████████░░░░░
02 Fri  8h 11m   █████████████████░░░
03 Sat  7h 55m   ███████████████░░░░░
04 Sun  0m       ░░░░░░░░░░░░░░░░░░░░
05 Mon  9h 02m   ██████████████████░░
...
------------------------------------------------------------
Active days : 18 / 31
Study       : 96h 24m
Work        : 48h 17m
Total       : 144h 41m
Daily avg   : 8h 02m
```

The monthly report uses completed sessions only.

------------------------------------------------------------------------

# `studyctl history`

Displays today's completed sessions.

## Usage

``` bash
studyctl history
```

## Example

``` text
Today's sessions
====================================================
#41  AWS MLA                    study  09:01 → 09:58  57m
#42  Machine Learning           study  10:03 → 11:04  1h 01m
#43  Client Project             work   13:10 → 14:04  54m
#44  FastAPI                    study  15:02 → 15:45  43m
```

## Session IDs

The ID can be used with `delete`.

For example:

``` bash
studyctl delete 43
```

------------------------------------------------------------------------

# `studyctl delete`

Permanently deletes a recorded session.

## Usage

``` bash
studyctl delete <id>
```

## Arguments

  Argument     Required Description
  ---------- ---------- -------------
  `<id>`            Yes Session ID

## Example

First inspect history:

``` bash
studyctl history
```

Then:

``` bash
studyctl delete 43
```

Output:

``` text
Session #43 deleted.
```

## Warning

Deletion is permanent from the SQLite database unless a database backup
exists.

------------------------------------------------------------------------

# `studyctl goal`

Manages daily productivity targets.

## Usage

``` bash
studyctl goal <command>
```

## Subcommands

``` text
set
show
```

------------------------------------------------------------------------

# `studyctl goal set`

Sets today's productivity target.

## Usage

``` bash
studyctl goal set <duration>
```

## Arguments

  Argument       Description
  -------------- -----------------
  `<duration>`   Target duration

## Supported Formats

``` text
8h
10h
6h
90m
7.5h
```

## Examples

``` bash
studyctl goal set 8h
```

``` bash
studyctl goal set 10h
```

``` bash
studyctl goal set 90m
```

Setting the goal again replaces today's previous target.

------------------------------------------------------------------------

# `studyctl goal show`

Displays today's target and progress.

## Usage

``` bash
studyctl goal show
```

## Example

``` text
Today's goal
--------------------------------------------
Target   : 8h 00m
Current  : 5h 42m
Progress : 71%
[█████████████████████░░░░░░░░░]
Remaining: 2h 18m
```

If the goal is reached:

``` text
Goal completed.
```

If no goal exists:

``` text
No daily goal configured.
```

------------------------------------------------------------------------

# `studyctl version`

Displays the application version.

## Usage

``` bash
studyctl version
```

## Example

``` text
studyctl v0.1.0
```

------------------------------------------------------------------------

# `studyctl completion`

Generates shell completion scripts.

## Usage

``` bash
studyctl completion <shell>
```

## Supported Shells

``` text
bash
zsh
```

## Bash

Generate:

``` bash
studyctl completion bash
```

Install:

``` bash
studyctl completion bash \
    > ~/.local/share/bash-completion/completions/studyctl
```

Restart the shell or reload completion configuration.

## Zsh

Generate:

``` bash
studyctl completion zsh
```

The generated completion script can be sourced or installed according to
the Zsh completion configuration.

------------------------------------------------------------------------

# Global Help

Every command supports:

``` bash
--help
```

or:

``` bash
-h
```

Examples:

``` bash
studyctl --help
studyctl start --help
studyctl stop --help
studyctl log --help
studyctl goal --help
studyctl goal set --help
```

------------------------------------------------------------------------

# Typical Daily Workflow

## Start the day

Set a target:

``` bash
studyctl goal set 8h
```

Start studying:

``` bash
studyctl start "AWS MLA"
```

Check progress:

``` bash
studyctl status
```

Stop:

``` bash
studyctl stop --note "Completed SageMaker section"
```

Start another task:

``` bash
studyctl start "Machine Learning"
```

Stop:

``` bash
studyctl stop
```

Check the day:

``` bash
studyctl today
```

Check the goal:

``` bash
studyctl goal show
```

Review the detailed history:

``` bash
studyctl history
```

------------------------------------------------------------------------

# Work + Study Workflow

A work session:

``` bash
studyctl start "Client Project" --type work
```

Stop:

``` bash
studyctl stop --note "Completed API integration"
```

Study:

``` bash
studyctl start "AWS MLA" --type study
```

Stop:

``` bash
studyctl stop
```

The reports keep the categories separate:

``` text
Study : 6h 20m
Work  : 2h 10m
Total : 8h 30m
```

------------------------------------------------------------------------

# Data Model

The core entity is a session.

``` text
Session
│
├── ID
├── Task
├── Type
├── Start Time
├── End Time
├── Duration
└── Notes
```

A session starts as:

``` text
start_time = current time
end_time   = NULL
```

When stopped:

``` text
end_time = current time
duration = end_time - start_time
```

This allows the application to identify an active session by:

``` sql
WHERE end_time IS NULL
```

Only one active session is permitted.

------------------------------------------------------------------------

# Time Accounting Rules

## Completed Sessions

Completed sessions contribute to:

``` text
today
week
month
history
goal progress
```

## Active Session

The active session is shown by:

``` bash
studyctl status
```

but is not included in completed historical totals until stopped.

## Manual Sessions

`studyctl log` creates a completed session immediately.

## Time Zones

Times are interpreted using the local system time zone.

The SQLite queries use local time when grouping sessions by date and
hour.

------------------------------------------------------------------------

# Error Handling

Expected errors include:

``` text
a session is already running
no active session
task cannot be empty
type must be 'study' or 'work'
duration must look like 8h or 90m
session not found
```

The application should fail with a useful message rather than silently
modifying data.

------------------------------------------------------------------------

# Development

Format source code:

``` bash
gofmt -w .
```

Update dependencies:

``` bash
go mod tidy
```

Run locally:

``` bash
go run . --help
```

Run a command:

``` bash
go run . today
```

Build:

``` bash
go build -o studyctl .
```

------------------------------------------------------------------------

# Testing

Run all Go tests:

``` bash
go test ./...
```

Run with the race detector:

``` bash
go test -race ./...
```

A future test suite should cover:

-   Starting a session
-   Preventing multiple active sessions
-   Stopping a session
-   Duration calculation
-   Manual logging
-   Study/work separation
-   Daily aggregation
-   Weekly aggregation
-   Monthly aggregation
-   Goal calculation
-   Session deletion
-   Invalid duration input
-   Invalid session type
-   Empty task names

------------------------------------------------------------------------

# Release Build

Build a Linux binary:

``` bash
CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o studyctl .
```

The result is:

``` text
studyctl
```

Install:

``` bash
sudo install -m 755 studyctl /usr/local/bin/studyctl
```

Verify:

``` bash
which studyctl
studyctl version
```

------------------------------------------------------------------------

# Database Backup

The application database is:

``` text
~/.studyctl/studyctl.db
```

Create a backup:

``` bash
cp ~/.studyctl/studyctl.db \
   ~/.studyctl/studyctl-backup.db
```

Restore:

``` bash
cp ~/.studyctl/studyctl-backup.db \
   ~/.studyctl/studyctl.db
```

Before restoring, stop any running `studyctl` command.

------------------------------------------------------------------------

# Recommended Future Features

The current command set forms the v1 foundation. Future releases can
add:

## Dashboard

``` bash
studyctl dashboard
```

A terminal dashboard showing:

-   Current session
-   Today's progress
-   Daily goal
-   Study/work split
-   Weekly graph
-   Recent sessions

## Yearly Report

``` bash
studyctl year
```

Monthly totals for an entire year.

## Custom Date Ranges

``` bash
studyctl report \
    --from 2026-10-01 \
    --to 2026-10-31
```

## Project Categories

Instead of treating every task independently:

``` bash
studyctl start "Transformer Architecture" \
    --project "ML Engineering"
```

Possible projects:

``` text
ML Engineering
AWS MLA
Interview Preparation
Client Project
Personal Project
```

## Export

CSV:

``` bash
studyctl export --format csv
```

JSON:

``` bash
studyctl export --format json
```

## Idle Detection

A Linux background daemon could detect periods where the machine is
idle.

Potential command:

``` bash
studyctl daemon
```

The application could then ask whether idle time should be excluded from
a running session.

## Desktop Notifications

For example:

``` text
Study session running for 2h 00m.
```

or:

``` text
Daily goal reached: 8h.
```

## Automatic Application Tracking

A future Linux-specific feature could record active applications or
windows while a session is running.

This should be opt-in because application/window tracking creates
additional privacy considerations.

------------------------------------------------------------------------

# Design Principles

## Local First

The application should work without an internet connection.

## Simple CLI

Common actions should require only one command.

``` bash
studyctl start "ML"
studyctl stop
studyctl today
```

## Reliable Data

Time calculations should be based on timestamps stored in SQLite rather
than terminal output.

## Explicit Tracking

The user controls when a session starts and stops.

## Extensible Architecture

Business logic should remain outside Cobra command implementations so
future interfaces can reuse it.

------------------------------------------------------------------------

# Command Cheat Sheet

``` text
START
studyctl start "Machine Learning"
studyctl start "Client Project" --type work

STOP
studyctl stop
studyctl stop --note "Finished chapter"

CURRENT
studyctl status

MANUAL LOG
studyctl log "AWS MLA" --duration 90m
studyctl log "Client Project" --duration 2h --type work

REPORTS
studyctl today
studyctl week
studyctl month
studyctl history

GOALS
studyctl goal set 8h
studyctl goal show

DATA
studyctl delete 42

SYSTEM
studyctl version
studyctl completion bash
studyctl completion zsh

HELP
studyctl --help
studyctl <command> --help
```

# Version

Current documentation target:

``` text
studyctl v0.1.0
```

The v0.1 release focuses on reliable session tracking, local
persistence, reporting, and daily goals.