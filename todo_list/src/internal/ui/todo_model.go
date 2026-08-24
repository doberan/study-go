package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"todo_list/internal/usecase"
)

type screen int

const (
	screenList screen = iota
	screenAdd
	screenDelete
	screenEdit
)

type model struct {
	service usecase.TodoService

	todoList list.Model // 表示するTodo一覧

	err error

	screen    screen
	textInput textinput.Model
}

func (m model) Init() tea.Cmd {
	return LoadTodos(m.service)
}

func NewModel(service usecase.TodoService) model {
	ti := textinput.New()
	ti.Placeholder = "Todoの説明を入力してください"

	items := []list.Item{}

	todoList := list.New(items, todoDelegate{}, 50, 10)
	todoList.Title = "Todo List"
	todoList.SetShowHelp(false)

	return model{
		service:   service,
		todoList:  todoList,
		err:       nil,
		screen:    screenList,
		textInput: ti,
	}
}
