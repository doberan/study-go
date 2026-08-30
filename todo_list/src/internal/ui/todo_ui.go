package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"todo_list/internal/domain"
)

type todoItem struct {
	todo domain.Todo
}

func (i todoItem) FilterValue() string {
	return i.todo.Description
}

func newTodoItems(todos []domain.Todo) []list.Item {
	items := make([]list.Item, 0, len(todos))

	for _, todo := range todos {
		items = append(items, todoItem{
			todo: todo,
		})
	}

	return items
}

type todoDelegate struct{}

func (d todoDelegate) Height() int {
	return 1
}

func (d todoDelegate) Spacing() int {
	return 0
}

func (d todoDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d todoDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	todo := item.(todoItem)

	cursor := " "
	if index == m.Index() {
		cursor = ">"
	}

	status := "□"
	if todo.todo.Done {
		status = "☒"
	}

	fmt.Fprintf(w, "%s %s | %s", cursor, status, todo.todo.Description)
}

var _ list.Item = todoItem{}
var _ list.ItemDelegate = todoDelegate{}
