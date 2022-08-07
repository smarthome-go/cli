## Changelog

### Configuration File
- The configuration file has been moved from `[config]/homescript.yml` to `[config]/smarthome-cli.yml`
- *Note*: This means you have to copy (`cp`) the old file to the new location
- Otherwise, you can run `smarthome-cli config set` again

### REPL
- The leading prompt is now changed to be `username@hostname` instead of `username@homescript`
- This is because the CLI has long replaced `homescript`, an old variant of this program

#### Sample REPL Session on 2.13.0
```
Welcome to Homescript interactive v2.13.0. CLI commands and comments start with #
Server: v0.0.51:go1.19 on http://localhost:8082
admin@localhost> switch('
Error: Program terminated abnormally with exit-code 1
SyntaxError at repl:1:8

  1  | switch('
              ^

String literal never closed
SyntaxError at repl:1:8

  1  | switch('
              ^

Expected expression, found 'Unknown'
admin@localhost [1][0.00s]>
```

