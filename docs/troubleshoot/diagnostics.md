---
title: Diagnostics and logging
description: Enable the gh-app-auth diagnostic log safely, read its events, and share it without leaking secrets.
sidebar:
  order: 2
---

gh-app-auth can write a diagnostic log of every credential flow. It is off by default and never
writes to the terminal, so it does not disturb Git's use of standard output.

## Turn it on

Set `GH_APP_AUTH_DEBUG_LOG` to an **explicit absolute file path**, then run the failing command:

```bash
export GH_APP_AUTH_DEBUG_LOG="$HOME/gh-app-auth-debug.log"
git clone https://github.com/example-org/private-repo.git
```

Any non-empty value enables logging, and the value is the file name. A value of `1` creates a file
named `1` in the current directory, not a default log. Unset the variable when you are done.

The log file is created with mode `0600` and appended to. Missing parent directories are created with
mode `0700`. If the directory cannot be created, the log goes to `gh-app-auth-debug.log` in the system
temporary directory, and if the file cannot be opened, logging is silently disabled. `--debug` is not
a substitute in v0.2.0.

## Read the log

Each line has a timestamp, an event name and an operation ID, followed by `key=value` fields:

```text
[2026-01-01T10:00:00.000Z] SESSION_START [session_1_42_op1] pid=42 args=[...]
[2026-01-01T10:00:00.010Z] FLOW_START [session_1_42_op2] operation=git_credential flow=START
[2026-01-01T10:00:00.020Z] FLOW_STEP [session_1_42_op3] step=match_app flow=STEP url=...
[2026-01-01T10:00:00.030Z] FLOW_SUCCESS [session_1_42_op9] operation=git_credential flow=SUCCESS
```

| Event | Meaning |
| --- | --- |
| `SESSION_START` / `SESSION_END` | One process. Records the process ID and its arguments |
| `FLOW_START` / `FLOW_SUCCESS` / `FLOW_ERROR` | The begin and end of one credential operation |
| `FLOW_STEP` | One step: reading input, matching, generating credentials |

Useful steps to look for: `no_pattern_match`, `app_matched` and `app_matched_by_pattern` (which App was
chosen and by which mechanism), `no_match_exit` (nothing matched, so the helper stayed silent) and
`host_only_query` (Git provided a host without a path, so the helper returned empty output).

To follow one failure, search for the process ID from `SESSION_START` and read that session
top to bottom.

## Share a log safely

- Tokens are hashed rather than printed, and obvious secret-looking values are redacted. Do not
  rely on that alone.
- The log still contains **command-line arguments, repository URLs and file paths**. Read it before
  attaching it anywhere.
- Do not enable `GIT_TRACE`, `GIT_CURL_VERBOSE` or similar Git tracing together with real credentials
  unless you must. They can print authorization headers.
- Delete the log afterwards: `rm -f "$GH_APP_AUTH_DEBUG_LOG"`.

## Related

[Environment variables](../reference/environment-variables.md) ·
[troubleshooting](index.md) · [security model](../understand/security-model.md)
