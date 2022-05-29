## Changelog

### Homescript Arguments
- Since Homescript `v0.9.0-beta`, scripts can use arguments. However, those arguments must be included before runtime. 
- Added support for local Homescript argument parsing
- When running / executing a Homescript via the CLI, Homescript arguments can be specified
- Example for local arguments: `shome ws run key:value key2:value2`
- When specifying the key-value pairs for the Homescript arguments, they can often be included by using the syntax from above
- The argument's key and value must be separated using a `:` colon and mustn't contain colons themselves 

