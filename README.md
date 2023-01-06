# Storiesque

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/noetarbouriech/storiesque?filename=backend%2Fgo.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/noetarbouriech/storiesque/backend)](https://goreportcard.com/report/github.com/noetarbouriech/storiesque/backend)
[![TypeScript](https://badgen.net/badge/icon/typescript?icon=typescript&label)](https://typescriptlang.org)

A fullstack application for creating and reading [Gamebooks](https://en.wikipedia.org/wiki/Gamebook) online. This application is composed of a REST API made in Golang and of a FrontEnd made in Sveltekit.It was made as a project during my third year at Polytech Montpellier.

![Screenshot Storiesque](https://user-images.githubusercontent.com/78071629/211106925-040ce940-11e9-4edb-91fc-287cc4320758.png)

## Installation

### With docker (recommended)

1. Clone the repo
```
git clone git@github.com:noetarbouriech/storiesque.git
```

2. Edit the `.env` file

2. Run the following command
```
sudo docker compose -f docker-compose.prod.yml up -d
```

Note: *Some files and folders may have a permission issue. Use of chmod may be required.*

### Without docker

TODO
