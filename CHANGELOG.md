## Changelog

### Workspace Subcommand
#### Bugfix
- Fixed one bug which occurs when deleting a newly-created HMS project right after it's creation (*without initial push*)
- This caused the new project to be nameless and to have no icon which looked ugly in the web-UI

#### Pre-Push Linting Hook
- Added the `LintOnPush` configuration parameter
- Can be configured using `shome config set ... -l` to enable it
- And `shome config set ...` (*no `-l`*) to disable it
- Automatically lints a Homescript using its local state before pushing
- If errors are detected, warnings are emitted, but the push is continued


