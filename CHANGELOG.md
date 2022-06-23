## Changelog

- Fix duplicated URL scheme when using `cli config set -a http://xy.z`
- In previous versions, the URL would have been set to `http://http://xy.z`
- Now, the subcommand auto detects if the `http://` prefix is set
- If the prefix is abundant, the URL is updated to use the default scheme `http`
