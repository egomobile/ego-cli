# ego-cli

Command Line Interface, which is designed to handle things, like Dev(Op) and other common tasks, much faster, written in [Go](https://golang.org/).

> ⚠️⚠️⚠️ **This project will be a refactored version of the legacy npm application [ego-cli](https://github.com/egodigital/ego-cli), written in Go as native binaries.** ⚠️⚠️⚠️

## Execute

Run this from project directory:

```bash
go run . ip -4 -6
```

to detect current IP addresses, e.g.

## Available commands

A (non complete) list of some interesting commands:

```
chuck         # Tries to output a random Chuck Norris joke by using icndb.com
git-pull      # Does a "git pull", in the current branch for all remotes in one command
git-push      # Does a "git push", in the current branch for all remotes in one command
git-sync      # Does a git "pull" and "push" in one command, in the current branch for all remotes
local-ip      # Lists all IP addresses of all known network interfaces service
node-install  # Deletes "node_modules" and executes "npm install" with optional "npm update" and "npm audit fix"
public-ip     # Tries to detect public IPv4 and IPv6 address(es) by using ipify.org service
serve         # Starts a HTTP server, serving the files in the current directory
```

To list all available commands, simply use `--help` option.

## Credits

* [Bootstrap 5](https://getbootstrap.com/docs/5.0/getting-started/introduction/)
* [Go Release Binaries action](https://github.com/marketplace/actions/go-release-binaries)
* [Modules, used by that stoftware](https://github.com/egomobile/ego-cli/blob/master/go.mod)
