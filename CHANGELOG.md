## Changelog

### Support For Authentication Tokens
As of Smarthome version `0.0.55`, you can fully leverage the potential of `authentication tokens`, thus allowing clients like this to authenticate themselves securely.  
- Every user can generate infinite authentication tokens
- Each token is unique and can authenticate a user
- Clients can use the token to authenticate without knowing a username or a password
- The CLI can use such a token... (*by modifying the config file*)
  - If the `token_authentication` variable (under connection) is set to `true`
  - By placing the token inside the `token` variable (under credentials)

### Configuration File
- The configuration file has been moved from `[config]/smarthome-cli.yml` to `[config]/smarthome-cli/config.toml`
- *Note*: This means that you have to manually translate your old `YAML` configuration into `TOML`

#### Newly Created Config File
This config file is generated the first time you start the CLI after the upgrade.  
If you want to edit this file manually but forgot its location, simply execute `shome config get`, the file's location will be printed in the console.

```toml
[connection]
  smarthome_url = "http://localhost"
  token_authentication = false

[credentials]
  password = ""
  token = ""
  username = ""

[homescript]
  lint_on_push = true
```
