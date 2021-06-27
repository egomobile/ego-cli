# ego-cli

Command Line Interface, which is designed to handle things, like Dev(Op) and other common tasks, much faster.

> ⚠️⚠️⚠️ **This project will be a refactored version of the legacy npm application [ego-cli](https://github.com/egodigital/ego-cli), written in Go as native binaries.** ⚠️⚠️⚠️

## Execute

Run this from project directory:

```bash
go run .
```

## Available commands

A (non complete) list of some interesting commands:

```
chuck  # Tries to output a random Chuck Norris joke by using icndb.com service
ip     # Tries to detect public IPv4 and IPv6 address(es) by using ipify.org service
```

To list all available commands, simply use `--help` option.

## Credits

* [Go Release Binaries action](https://github.com/marketplace/actions/go-release-binaries)
* [Modules, used by that stoftware](./go.mod)
