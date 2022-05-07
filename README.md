## go-MyGit - Git action with group, implement in Go

`go-mygit` is a command line tools for easy configuration of git remote, github/gites repositories.

> This is a major upgrade of [mygit](https://github.com/J-Siu/mygit) which was written in bash.

<!-- TOC -->

- [Who & Why](#who--why)
- [Features](#features)
- [Limitation](#limitation)
- [Usage](#usage)
  - [go-mygit](#go-mygit)
  - [go-mygit config](#go-mygit-config)
  - [go-mygit remote](#go-mygit-remote)
  - [go-mygit repository](#go-mygit-repository)
  - [go-mygit repository get](#go-mygit-repository-get)
  - [go-mygit repository set](#go-mygit-repository-set)
  - [Debug](#debug)
  - [Selector](#selector)
    - [-g/--group](#-g--group)
    - [-r/--remote](#-r--remote)
  - [Git Base Commands](#git-base-commands)
    - [init](#init)
    - [push](#push)
      - [--tags](#--tags)
      - [--all](#--all)
- [Configuration File](#configuration-file)
- [Repository](#repository)
- [Contributors](#contributors)
- [Change Log](#change-log)
- [License](#license)

<!-- /TOC -->

### Who & Why

- Creating repositories for same set of remote servers repeatedly
- Setting up repositories on multiple machines repeatedly
- Working with repositories that push to same set of git servers

### Features

- Configuration File
  - [x] remotes
  - [x] groups
  - [x] secrets
- Selector
  - [x] -g/--group
  - [x] -r/--remote
- Base(git) Commands
  - [x] init
  - [x] push
  - [x] remote
    - [x] add
    - [x] list
    - [x] remove
- repository(api)
  - [x] delete
  - [x] get
    - [x] all
    - [x] description
    - [x] private
    - [x] publickey
    - [x] topics
    - [x] visibility
  - [x] list
  - [x] new
  - [x] set
    - [x] description
    - [x] private
    - [x] secrets
    - [x] topics
    - [x] visibility

### Limitation

- Current supported git servers
  - gitea
  - github
  - gogs
- API commands must be executed at root of repository

### Usage

#### go-mygit
```sh
Git automation script support group action.

Usage:
  go-mygit [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Print configurations
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
#### go-mygit config
```sh
Print configurations

Usage:
  go-mygit config [command]

Aliases:
  config, c, conf

Available Commands:
  all         Print all configurations
  group       Print groups configuration
  remote      Print remotes configuration
  secret      Print secret configuration

Flags:
  -h, --help   help for config

Global Flags:
      --config string        config file (default is $HOME/.go-mygit.json)
  -d, --debug                Enable debug
  -g, --group stringArray    Specify group
  -r, --remote stringArray   Specify remotes

Use "go-mygit config [command] --help" for more information about a command.
```
#### go-mygit remote
```sh
remote(git) commands

Usage:
  go-mygit remote [command]

Aliases:
  remote, rmt

Available Commands:
  add         Add git remotes base on configuration and flags
  list        List git remotes in current repository
  remove      Delete remotes in current repository

Flags:
  -h, --help   help for remote

Global Flags:
      --config string        config file (default is $HOME/.go-mygit.json)
  -d, --debug                Enable debug
  -g, --group stringArray    Specify group
  -r, --remote stringArray   Specify remotes

Use "go-mygit remote [command] --help" for more information about a command.
```
#### go-mygit repository
```sh
Repository commands

Usage:
  go-mygit repository [command]

Aliases:
  repository, repo

Available Commands:
  delete      Delete remote repository
  get         get info
  list        List all remote repositories
  new         Create remote repository
  set         set info

Flags:
  -h, --help   help for repository

Global Flags:
      --config string        config file (default is $HOME/.go-mygit.json)
  -d, --debug                Enable debug
  -g, --group stringArray    Specify group
  -r, --remote stringArray   Specify remotes

Use "go-mygit repository [command] --help" for more information about a command.
```
#### go-mygit repository get
```sh
get info

Usage:
  go-mygit repository get [command]

Aliases:
  get, g

Available Commands:
  all         get all info(json)
  description get description
  private     get private status
  publickey   get public key
  topics      get topics
  visibility  get visibility

Flags:
  -h, --help   help for get

Global Flags:
      --config string        config file (default is $HOME/.go-mygit.json)
  -d, --debug                Enable debug
  -g, --group stringArray    Specify group
  -r, --remote stringArray   Specify remotes

Use "go-mygit repository get [command] --help" for more information about a command.
```
#### go-mygit repository set
```sh
set info

Usage:
  go-mygit repository set [command]

Aliases:
  set, s

Available Commands:
  description set description
  private     set private status
  secrets     set action secrets
  topics      set topics
  visibility  set visibility status

Flags:
  -h, --help   help for set

Global Flags:
      --config string        config file (default is $HOME/.go-mygit.json)
  -d, --debug                Enable debug
  -g, --group stringArray    Specify group
  -r, --remote stringArray   Specify remotes

Use "go-mygit repository set [command] --help" for more information about a command.
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

### Configuration File

Following configuration will be used in all examples:

```json
{
	"secrets": [
		{
			"name": "",
			"value": ""
		},
		{
			"name": "",
			"value": ""
		}
	],
	"remotes": [
		{
			"group": "",
			"name": "GitHub",
			"private": false,
			"ssh": "",
			"token": "",
			"entrypoint": "https://api.github.com",
			"vendor": "github"
		},
		{
			"group": "",
			"name": "MyServer",
			"private": true,
			"ssh": "",
			"token": "",
			"entrypoint": "https://gitea.myserver.com/api/v1",
			"user": "",
			"vendor": "gitea"
		}
	]
}
```

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

### License

The MIT License (MIT)

Copyright © 2022 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
