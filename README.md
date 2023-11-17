# task-gopher

This project is inspired by the incredible work on Task Warrior, an open source
CLI task manager. I use this project quite a bit for managing my projects
without leaving the safety and comfort of my terminal. (⌐■_■)

We built a kanban board TUI in a previous [tutorial][kanban-video], so the
idea here is that we're going to build a task management CLI with [Cobra][cobra] that has Lip Gloss
styles *and* can be viewed using our kanban board.

*Note: We walk through the code explaining each and every piece of it in the
[corresponding video](https://youtu.be/yiFhQGJeRJk) for this tutorial. Enjoy!!*

Here's the plan:

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
- [ ] add server
  - [ ] add routes for all CLI commands
  - [ ] option to create/use local DB or server DB
  - [ ] add .env file with server addr, port, credentials (if needed)

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