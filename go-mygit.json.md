### Location

File location: `$HOME/.go-mygit.json`.

### Sample Config

Config file is in json format.

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
      "user": ""
      "vendor": "github"
    },
    {
      "group": "",
      "name": "MyServer",
      "private": true,
      "ssh": "",
      "token": "",
      "entrypoint": "https://<domain>/api/v1",
      "user": "",
      "vendor": "gitea",
      "skipverify": true
    }
  ]
}
```

### Section - secrets

```json
"secrets": [
  {
    "name": "",
    "value": ""
  },
  ...
],
```

`secrets` section hold name, value pair  for Github action secret.

`go-mygit repo set secrets` will push all secret pairs to server.

One use case is to push docker username and docker api token for docker projects with Github workflow. For example if using [publish-docker.yml](https://github.com/J-Siu/github-workflows/blob/main/publish-docker.yml), you can setup following:

```json
"secrets": [
  {
    "name": "DOCKER_HUB_USERNAME",
    "value": "<docker username>"
  },
  {
    "name": "DOCKER_HUB_ACCESS_TOKEN",
    "value": "<docker token>"
  }
],
```

### Section - remotes

```json
  "remotes": [
  {
    "group": "external",
    "name": "gh",
    "private": false,
    "ssh": "git@github.com",
    "token": "",
    "entrypoint": "https://api.github.com",
    "user": "J-Siu",
    "vendor": "github",
    "skipverify": true
  },
  ...
]
```

`remotes` section hold both git remote info and api info.

- "group":

  Group name of remote. See [Selector](#selector) below.

- "name":

  Name of remote. It is the name of the git remote.

- "private":

  true/false, default private status when creating repository on server(`go-mygit repo new`).

- "ssh":

  SSH portion of git remote url. The part before ':'. For Github, it is always "git@github.com". If it is set up in `./ssh/config`, the `Host` can be used. Example:

  `./ssh/config`
  ```sh
  Host gh
    HostName github.com
    User git
  ```

  "gh" can be used.

- "token"

  Github/Gitea API token

- "entrypoint"

  API endpoint.

  Github should always use "https://api.github.com".

  Gitea server usually "https://<your_domain>/api/v1".

- "user"

  Server username. Used in both API and git remote url construction.

- "vendor"

  Either "github" or "gitea"

- "skipverify"

  Skip SSL/TLS verify for API call.

### Selector

Except `config`, all `go-mygit` commands support both "-g/--group \<group>" and "-r/--remote \<name>" options.

For example, when using "go-mygit init", you can supply "-r":

```sh
go-mygit init -r gh
```

Then only the "gh" remote will be setup. Same idea apply for "-g".
