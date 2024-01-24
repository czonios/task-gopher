# Task Gopher

> A CLI task management tool with remote access capabilities.

![Build](https://github.com/czonios/task-gopher/actions/workflows/go.yml/badge.svg)
![Lint](https://github.com/czonios/task-gopher/actions/workflows/lint.yml/badge.svg)

<p align="center">
  <img width=300 src="./assets/logo.png">
</p>

This project is a CLI task manager/todo application, with a web server built using [Echo][echo], and an SQLite database to hold the task data.
It is built as a reimplementaion and continuation of [TaskCLI](https://github.com/charmbracelet/taskcli/tree/main) which itself is inspired by [Task Warrior](https://taskwarrior.org/). It is not a fork, but rather a reimagination from scratch. Some components have been copied (e.g. the CLI Kanban command and helper functions). After the initial implementation, some conventions from the TaskCLI were used (e.g. using a `status` enum, instead of a bool (todo/done) used originally) when they were more convenient than the original implementation.

The differences with TaskCLI, at a glance:

- implemented an [Echo][echo] server for remotely accessing the tasks (multiple clients)
- implemented updating all clients through [Gorilla WebSocket](https://github.com/gorilla/websocket)
- added [Docker](https://docs.docker.com/get-docker/) containers with build and run scripts for the app
- added optional extra functionality such as a `task_type` enum and `description` fields to the tasks
- implemented an extra CLI command, `deldb`, to clear the database of tasks

## Setup

### Requirements

- [Go](https://go.dev/learn/): the Go language, use `go version` command to check if it is installed. This has been tested on `go1.21.4`. If you use Docker, you don't need to install it.
- (optional) [Docker](https://docs.docker.com/get-docker/): if you prefer to run the app in a container, scripts are included for building and running the app in a Docker container
- (optional) Server node: If you want to have a task-gopher server that you can access from other devices, then I suggest using [ZeroTier][zerotier], which allows you to add devices to a virtual network so you can view them as if they are on your local network, with static IP addresses, as long as they are connected to the internet. I prefer it because it is simple, open source, and free for personal use.

### Build and run

#### Clone repository

Clone this repo - we use the Go convention of holding packages from GitHub in `$HOME/go/src/github.com/<username>/<package>`:

```sh
git clone https://github.com/czonios/task-gopher.git $HOME/go/src/github.com/czonios/task-gopher
```

#### Set environment variables

The application tries to read a `.env` file in the root directory of the project and load the environment variables it contains. The `.env` file is optional, but the following environment variables must be set:

- `ADDRESS` the address of the server
- `PORT` the port the server runs on

#### Start the server

The following commands start the task-gopher server (on the device that will hold the database). Don't forget to set the `ADDRESS` and `PORT` of the server as environment variables in `.env` for this to work! Since this is the server instance, you can use `http://localhost` for the `ADDRESS`.

##### Option 1: using Go

```sh
# cd to root directory of project
cd $HOME/go/src/github.com/czonios/task-gopher
go install ./...
task-gopher serve
# NOTE: $HOME/go/bin should be in your PATH
# otherwise you can run go run ./... serve
```

##### Option 2: using Docker

```sh
# cd to root directory of project
cd $HOME/go/src/github.com/czonios/task-gopher
cd build
docker compose up
```

#### Start a client

The client can be either in the same machine as the server, or in any other machine that can ping the server.
Don't forget to set the `ADDRESS` and `PORT` of the server as environment variables in `$HOME/go/src/github.com/czonios/task-gopher/.env` for this to work!

##### Option 1: using Go

```sh
# cd to root directory of project
cd $HOME/go/src/github.com/czonios/task-gopher
go install ./...
task-gopher --help # will list available commands
# NOTE: $HOME/go/bin should be in your PATH
# otherwise you can run go run ./... --help
```

##### Option 2: using Docker

```sh
# coming soon!
```

## Meta

Christos A. Zonios – [czonios.github.io](https://czonios.github.io) – c.zonios (at) uoi (dot) gr

Distributed under the MIT license. See ``LICENSE`` for more information.

[https://github.com/czonios/task-gopher](https://github.com/czonios/task-gopher)

## Contributing

1. Fork it (<https://github.com/czonios/task-gopher/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Run the [tests](#tests)
4. Commit your changes (`git commit -am 'Add some fooBar'`)
5. Push to the branch (`git push origin feature/fooBar`)
6. Create a new Pull Request

### Project layout

```sh
task-gopher
├── LICENSE
├── README.md
├── assets
│   └── logo.png
├── build
│   ├── build_docker_img.sh     # script for building the docker image
│   └── dockerfile
├├── cmd
│   └── task-gopher
│       ├── cli.go              # Cobra commands and setup for CLI
│       ├── server.go           # server and routes to interract with the task manager
│       └── task-gopher.go      # main function, Task struct, handles initial setup
├── data
│   └── tasks.db                # created by the server
├── go.mod
└── scripts
    ├── run_docker_daemon.sh    # script for running the docker container as a daemon on startup
    └── run_docker_img.sh       # script for normal docker run
```

### Testing

```sh
go test ./...
```

## Next steps

### Docker containers

- [x] server container
- [ ] app container
- [ ] documentation for containers

### Tests

- [x] `task-gopher.go`
- [ ] `server.go`

### Mobile app

- [ ] create basic mobile app using [Go app][gomobile] or [Fyne][fyne] or [Wails][wails]

[fyne]: https://fyne.io/
[wails]: https://wails.io/
[gomobile]: https://pkg.go.dev/golang.org/x/mobile/app
[zerotier]: https://www.zerotier.com/
[echo]: https://echo.labstack.com/
