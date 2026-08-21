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
	repository := &MySQLTodoRepository{db: db}
	service := NewTodoService(repository)
	p := tea.NewProgram(NewModel(service))
	if _, err := p.Run(); err != nil {
		fmt.Println("エラー:", err)
		return
	}
}
