package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := initDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	repository := &TodoRepository{db: db}
	p := tea.NewProgram(initialModel(repository))
	if _, err := p.Run(); err != nil {
		fmt.Println("エラー:", err)
		return
	}
}
