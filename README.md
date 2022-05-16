# go-MyGit [![Paypal donate](https://www.paypalobjects.com/en_US/i/btn/btn_donate_LG.gif)](https://www.paypal.com/donate/?business=HZF49NM9D35SJ&no_recurring=0&currency_code=CAD)

Command line tool for easy mass configuration of git remotes, and working with Github/Gites repositories API.

> This replaces [mygit](https://github.com/J-Siu/mygit), implemented in Bash.

### Table Of Content
<!-- TOC -->

- [Table Of Content](#table-of-content)
- [Highlight](#highlight)
  - [Push multiple repository](#push-multiple-repository)
  - [Set Private](#set-private)
  - [Update Description](#update-description)
  - [Update Topics](#update-topics)
- [What It Does](#what-it-does)
- [What It Does Not](#what-it-does-not)
- [Features](#features)
- [Limitation](#limitation)
- [Usage](#usage)
  - [Debug](#debug)
  - [Selector](#selector)
    - [-g/--group](#-g--group)
    - [-r/--remote](#-r--remote)
  - [Git Base Commands](#git-base-commands)
    - [init](#init)
    - [push](#push)
      - [--tags](#--tags)
      - [--all](#--all)
- [Configuration](#configuration)
- [Packages Used](#packages-used)
- [Binary](#binary)
- [Repository](#repository)
- [Contributors](#contributors)
- [Change Log](#change-log)
- [License](#license)

<!-- /TOC -->

### Highlight

Following are highlight of some `go-mygit` functions.

#### Push multiple repository

```sh
# Push to all remotes
go-mygit push
# Support path operation
go-mygit push docker_*
```

#### Set Private

```sh
# Set current repository to private on server
go-mygit repo set private true
# Support path operation
go-mygit repo set private false docker_*
```

#### Update Description

```sh
# Update description on server of current repository/directory
go-mygit repo set description "This is a new description"
```

#### Update Topics

```sh
# Update topics on server of current repository/directory
go-mygit repo set topic golang go project
```

### What It Does

> These are the reasons "mygit" got created.

- Set up same set of git remote repeatedly
- Pushing same repo to multiple git servers which are not mirrored
- Update some repository info on git server

### What It Does Not

- Replacing `git` command. (`git` command is required for git function to work.)
- Replacing Github cli `gh` (`go-mygit` only cover very few api in comparison.)

### Features

- Configuration File
  - [x] remotes
  - [x] groups
  - [x] secrets
- Selector for git servers
  - [x] -g/--group
  - [x] -r/--remote
- Base(git) Commands
  - [x] init
  - [x] push
  - [x] remote
    - [x] add
    - [x] list
    - [x] remove
- Repository(api)
  - [x] list all repo on server
  - [x] create repo on server
  - [x] get / set
    - [x] description
    - [x] private
    - [x] public key(get only)
    - [x] secret
    - [x] topic
    - [x] visibility
  - [x] delete
    - [x] repository
    - [x] secret

### Limitation

- Current supported git servers
  - github
  - gitea
  - gogs(not tested)

### Usage

```sh
Git automation script support group action.

Usage:
  go-mygit [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config       Print configurations
  help        Help about any command
  init        Git init and set remotes
  push        Push to all remote repositories
  remote      remote(git) commands
  repository  Repository commands

Flags:
      --config string        config file (default is $HOME/.go-mygit.json)
  -d, --debug                Enable debug
  -g, --group stringArray    Specify group
  -h, --help                 help for go-mygit
  -r, --remote stringArray   Specify remotes

Use "go-mygit [command] --help" for more information about a command.
```

#### Debug

Use `-d` for debug output.

Example:
```sh
go-mygit -d repo des get
```

#### Selector

`go-mygit` allow command applied to groups or remotes through the use of `-g/--group` and `-r/--remote`. This applies to all commands except `remote` and `group` mentioned above.

`-g/--group` and `-r/--remote` must be placed right after `mygit` and before any command.

##### -g/--group

```sh
go-mygit -g external <command>
go-mygit -g external -g internal <command>
```

##### -r/--remote

```sh
go-mygit -r gh <command>
go-mygit -r gh -r server3 <command>
```

`-g/--group` and `-r/--remote` can be used at the same time.

```sh
go-mygit -g external -r server2 <command>
```

#### Git Base Commands

`init`, `push` are git base commands.

##### init

`go-mygit init` will clear all existing remote and add remote base on `-g`/`-r` selector. If no group nor remote are specified, all configured remotes will be added.

`go-mygit init` by default use current directory name as repository name. Repository name can be specified in the format `go-mygit init <repository>`. File `.go-mygit.json` containing the repository name will be created, which is used by API based commands.

Before `go-mygit init`:

```sh
$ git remote -v
origin  https://github.com/J-Siu/mygit.git (fetch)
origin  https://github.com/J-Siu/mygit.git (push)
```

`go-mygit init` without selector:

```sh
$ go-mygit init
Reinitialized existing Git repository in /tmp/mygit/.git/

$ git remote -v
gh      git@github.com:/username1/mygit.git (fetch)
gh      git@github.com:/username1/mygit.git (push)
server2 git@server2:/username2/mygit.git (fetch)
server2 git@server2:/username2/mygit.git (push)
server3 git@server3:/username3/mygit.git (fetch)
server3 git@server3:/username3/mygit.git (push)
```

`go-mygit init` with group internal:

```sh
$ go-mygit --group internal init
Reinitialized existing Git repository in /tmp/mygit/.git/

$ git remote -v
server2 git@server2:/username2/mygit.git (fetch)
server2 git@server2:/username2/mygit.git (push)
server3 git@server3:/username3/mygit.git (fetch)
server3 git@server3:/username3/mygit.git (push)
```

`go-mygit init` with repository name:

```sh
$ go-mygit init mygit2
Reinitialized existing Git repository in /tmp/mygit/.git/

$ git remote -v
gh      git@github.com:/username1/mygit2.git (fetch)
gh      git@github.com:/username1/mygit2.git (push)
server2 git@server2:/username2/mygit2.git (fetch)
server2 git@server2:/username2/mygit2.git (push)
server3 git@server3:/username3/mygit2.git (fetch)
server3 git@server3:/username3/mygit2.git (push)
```

##### push

`go-mygit push` base on `-g`/`-r` selector. If no group nor remote are specified, all configured remotes will be pushed in sequence.

```sh
go-mygit push
```

```sh
go-mygit -r gh push
```

`go-mygit push` support options `--tags` and `--all`

###### --tags

If `--tags` is used, `go-mygit push` will push all tags.

```sh
go-mygit push --tags
```

###### --all

If `--all` is used, `mygit push` will push all branches(`--all`).

```sh
go-mygit -r gh push --all
```

### Configuration

See [go-mygit.json.md](go-mygit.json.md)
### Packages Used

- [go-gitapi](https://github.com/J-Siu/go-gitapi)
- [go-helper](https://github.com/J-Siu/go-helper)
- [cobra](//github.com/spf13/cobra)
- [viper](//github.com/spf13/viper)

### Binary

https://github.com/J-Siu/go-mygit/releases

### Repository

- [go-mygit](https://github.com/J-Siu/go-mygit)

### Contributors

- [John, Sing Dao, Siu](https://github.com/J-Siu)

### Change Log

- v1.0.0
  - Feature complete
- v1.0.1
  - Fix repo new endpoint
- v2.0.0
  - Command line restructure
  - Clean up func name
  - Clean up file name
  - Fix typos
- v2.0.1
  - upgrade go-helper and go-gitapi for bugfix
- v2.1.0
  - Add repo name support for all repoGet* commands
- v2.2.0
  - Support deletion of github repository action secret
- v2.3.0
  - Add repo(dir) name support for init, push, repo new, and all remote* commands
- v2.4.0
  - Add global --noskip flag
  - Improve commands Use, Short, Long
  - Improve repo/dir handling from command line
  - lib.GitApiFromRemote() -> Remote.GetGitApi()
- v2.4.1
  - Fix `goreleaser`
- v2.4.2
  - Proper go mod path for v2

### License

The MIT License (MIT)

Copyright © 2022 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
