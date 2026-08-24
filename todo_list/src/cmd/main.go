package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"

	"todo_list/internal/domain"
	"todo_list/internal/infrastructure"
	"todo_list/internal/ui"
	"todo_list/internal/usecase"
)

func main() {
	db := infrastructure.InitDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	repository := domain.NewMySQLTodoRepository(db)
	service := usecase.NewTodoService(repository)
	p := tea.NewProgram(ui.NewModel(service))
	if _, err := p.Run(); err != nil {
		fmt.Println("エラー:", err)
		return
	}
}
