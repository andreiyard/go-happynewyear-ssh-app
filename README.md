# go-happynewyear-ssh-app

This is a simple SSH application that shows animated screen when you connect to the app via SSH.

## Dependencies

This project is implemented using two main dependencies:
- [Bubble Tea (Go TUI framework)](https://github.com/charmbracelet/bubbletea)
- [Wish (Go SSH app framework)](https://github.com/charmbracelet/wish)

## How to run:

There are two ways to run this app:
1. Standalone local terminal app
```bash
go run ./cmd/app-local
```
2. SSH server
```bash
go run ./cmd/app-ssh
```

## Plan:

~~1. Develop TUI animation using bubbletea~~  
~~2. Wrap resulting application in wish to be able to serve via SSH~~  
3. Improve and polish visuals
