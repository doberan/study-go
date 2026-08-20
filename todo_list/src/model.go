package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenList screen = iota
	screenAdd
	screenDelete
	screenEdit
)

type model struct {
	repository *TodoRepository

	todoList list.Model // 表示するTodo一覧

	err error

	screen    screen
	textInput textinput.Model
}

type TodosLoadedMsg struct {
	todos []Todo
	err   error
}

func LoadTodos(repository *TodoRepository) tea.Cmd {
	return func() tea.Msg {
		todos, err := repository.FindAll()

		return TodosLoadedMsg{
			todos: todos,
			err:   err,
		}
	}
}

type TodoUpdateMsg struct {
	todo Todo
	err  error
}

func updateTodo(repository *TodoRepository, todo Todo) tea.Cmd {
	return func() tea.Msg {
		err := repository.Update(todo)
		return TodoUpdateMsg{
			todo: todo,
			err:  err,
		}
	}
}

type TodoCreateMsg struct {
	todo Todo
	err  error
}

func createTodo(repository *TodoRepository, description string) tea.Cmd {
	return func() tea.Msg {
		todo := Todo{
			Description: description,
			Done:        false,
		}
		err := repository.Create(&todo)

		return TodoCreateMsg{
			todo: todo,
			err:  err,
		}
	}
}

type TodoDeleteMsg struct {
	id  int
	err error
}

func deleteTodo(repository *TodoRepository, id int) tea.Cmd {
	return func() tea.Msg {
		err := repository.Delete(id)
		return TodoDeleteMsg{
			id:  id,
			err: err,
		}
	}
}

func initialModel(repository *TodoRepository) model {
	ti := textinput.New()
	ti.Placeholder = "Todoの説明を入力してください"

	items := []list.Item{}

	todoList := list.New(items, todoDelegate{}, 50, 10)
	todoList.Title = "Todo List"
	todoList.SetShowHelp(false)

	return model{
		repository: repository,
		todoList:   todoList,
		err:        nil,
		screen:     screenList,
		textInput:  ti,
	}
}

func (m model) Init() tea.Cmd {
	return LoadTodos(m.repository)
}
