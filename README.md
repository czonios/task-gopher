# task-gopher

This project is inspired by the incredible work on Task Warrior, an open source
CLI task manager. I use this project quite a bit for managing my projects
without leaving the safety and comfort of my terminal. (⌐■_■)

We built a kanban board TUI in a previous [tutorial][kanban-video], so the
idea here is that we're going to build a task management CLI with [Cobra][cobra] that has Lip Gloss
styles *and* can be viewed using our kanban board.

## Setup

##### Clone repository
Clone this repo in the correct directory - **IMPORTANT**:
```sh
git clone https://github.com/czonios/task-gopher.git $HOME/go/src/github.com/czonios
```

##### Set environment variables
Create a `.env` file in `$HOME/go/src/github.com/czonios/task-gopher`
Add the following environment variables:
- `ADDRESS` the address of the server
- `PORT` the port the server runs on

##### Start the server
The following commands start the task-gopher server (on the device that will hold the database).
Don't forget to set the `ADDRESS` and `PORT` of the server as environment variables in `$HOME/go/src/github.com/czonios/task-gopher/.env` for this to work! Since this is the server instance, you can use `http://localhost` for the `ADDRESS`.
```sh
cd $HOME/go/src/github.com/czonios/task-gopher
go install ./...
task-gopher serve
```

##### Start a client
Don't forget to set the `ADDRESS` and `PORT` of the server as environment variables in `$HOME/go/src/github.com/czonios/task-gopher/.env` for this to work!
```sh
cd $HOME/go/src/github.com/czonios/task-gopher
go install ./...
task-gopher --help # will list available commands
```

## Checklist

### Data storage
- [x] set up a (SQLite?) database
  - [x] open DB
  - [x] add task
  - [x] delete task
  - [x] edit task
  - [x] get tasks

### Make a CLI with Cobra
- [x] add CLI
  - [x] add task
  - [x] delete task
  - [x] edit task
  - [x] get tasks

### Add a little... *je ne sais quoi*
- [x] print to table layout with [Lip Gloss][lipgloss]
- [x] print to Kanban layout with [Lip Gloss][lipgloss]

### Tasklist server
- [x] add server
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