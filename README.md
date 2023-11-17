# Task Gopher
> CLI task manager with a web server

![Build](https://github.com/czonios/task-gopher/actions/workflows/go.yml/badge.svg)
![Lint](https://github.com/czonios/task-gopher/actions/workflows/lint.yml/badge.svg)

This project is a CLI task manager/todo application, with a web server built using [Echo][echo], and an SQLite database to hold the task data.
It is built as a reimplementaion and continuation of [TaskCLI](https://github.com/charmbracelet/taskcli/tree/main) which itself is inspired by [Task Warrior](https://taskwarrior.org/). It is not a fork, but rather a reimagination of the implementation from scratch. Some components have been copied (e.g. the CLI Kanban command and helper functions). After the initial implementation, some conventions from the TaskCLI were used (e.g. using a `status` enum, instead of a bool (todo/done) used originally) when they were more convenient than the original implementation.

The differences with TaskCLI, in a glance:
- implemented a server for remotely accessing the tasks
- added an optional description to the tasks
- implemented a `deldb` command to clear the database of tasks, that cannot be run remotely

## Setup

### Requirements
- Go: the Go language, use `go version` command to check if it is installed. This has been tested on `go1.21.4`
- If you want to have a task-gopher server that you can access from other devices, then I suggest using [ZeroTier][zerotier], which allows you to add devices to a virtual network so you can view them as if they are on your local network, with static IP addresses, as long as they are connected to the internet. I prefer it because it is simple, open source, and free for personal use.

### Build and run
##### Clone repository
Clone this repo in the correct directory - **IMPORTANT**: we use the Go convention of holding packages in `$HOME/go/src/github.com/<username>/<package>`:
```sh
git clone https://github.com/czonios/task-gopher.git $HOME/go/src/github.com/czonios
```

##### Set environment variables
Create a `.env` file in `$HOME/go/src/github.com/czonios/task-gopher`
Add the following environment variables:
- `ADDRESS` the address of the server
- `PORT` the port the server runs on

##### Start the server
The following commands start the task-gopher server (on the device that will hold the database). Don't forget to set the `ADDRESS` and `PORT` of the server as environment variables in `$HOME/go/src/github.com/czonios/task-gopher/.env` for this to work! Since this is the server instance, you can use `http://localhost` for the `ADDRESS`.
```sh
cd $HOME/go/src/github.com/czonios/task-gopher
go install ./...
task-gopher serve
```

##### Start a client
The client can be either in the same machine as the server, or in any other machine that can ping the server. Don't forget to set the `ADDRESS` and `PORT` of the server as environment variables in `$HOME/go/src/github.com/czonios/task-gopher/.env` for this to work!
```sh
cd $HOME/go/src/github.com/czonios/task-gopher
go install ./...
task-gopher --help # will list available commands
```

## Checklist

### Tasklist server
- [x] add [Echo] server
  - [x] add routes for all CLI commands
    - [x] addTask
    - [x] updateTask
    - [x] deleteTask
    - [x] listTasks
    - [x] kanban
  - [x] add .env file with server addr, port, credentials (if needed)

### Mobile app
- [ ] create basic mobile app using [Go app][gomobile] or [Fyne][fyne] or [Wails][wails]

## Project layout

`main.go` - defines task struct, handles initial setup including opening a db and setting data path for our app
`cmds.go` - does all Cobra commands and setup for CLI

[lipgloss]: https://github.com/charmbracelet/lipgloss
[charm]: https://github.com/charmbracelet/charm
[cobra]: https://github.com/spf13/cobra
[kanban-video]: https://www.youtube.com/watch?v=ZA93qgdLUzM&list=PLLLtqOZfy0pcFoSIeGXO-SOaP9qLqd_H6
[fyne]: https://fyne.io/
[wails]: https://wails.io/
[gomobile]: https://pkg.go.dev/golang.org/x/mobile/app
[zerotier]: https://www.zerotier.com/
[echo]: https://echo.labstack.com/