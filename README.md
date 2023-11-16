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
- [ ] add CLI
  - [ ] add task
  - [ ] delete task
  - [ ] edit task
  - [ ] get tasks

### Add a little... *je ne sais quoi*
- [ ] print to table layout with [Lip Gloss][lipgloss]
- [ ] print to Kanban layout with [Lip Gloss][lipgloss]

## Project layout

`db.go` - here is our custom `task` struct and data layer
`main.go` - handles initial setup including opening a db and setting data path for our app
`cmds.go` - does all Cobra commands and setup for CLI

[lipgloss]: https://github.com/charmbracelet/lipgloss
[charm]: https://github.com/charmbracelet/charm
[cobra]: https://github.com/spf13/cobra
[kanban-video]: https://www.youtube.com/watch?v=ZA93qgdLUzM&list=PLLLtqOZfy0pcFoSIeGXO-SOaP9qLqd_H6