# Welcome to Katapult CLI!

So you love Katapult, but you want to be able to use it from the command line? This is the tool for you!

## CLI Support

Katapult CLI aims to allow you to do all the actions you can do from the UI but inside of your terminal. As of right now, the following options are supported:

- [Organisation actions](organisation-actions.md)
- [Network actions](network-actions.md)
- [Data centre actions](data-centre-actions.md)
- [Virtual machine actions](virtual-machine-actions.md)

## Output Types
All commands in the CLI support outputting YAML, JSON, and text (with custom templating support). To set the output type, you can use `-o <yaml/json/text>`.

For advanced use, you can also use `-t` to provide a custom Go template. This will contain the API response object for what you are trying to access in the form that it is parsed by go-katapult.

## Setup
To setup the Katapult CLI, you will want to install the package for your respective package manager:

### Brew (macOS/Linux)
```
$ brew tap krystal/tap
$ brew install krystal/tap/katapult-cli
```

### Manual Installation (Windows/other Linux distros)
For manual installation, you can download the binary from the GitHub release.
