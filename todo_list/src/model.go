package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	repository *TodoRepository
	todos      []Todo // 表示するTodo一覧
	cursor     int    // 現在選択しているTodoの位置
	err        error
}

type TodosLoadedMsg struct {
	todos []Todo
	err   error
}

func initialModel(repository *TodoRepository) model {

	return model{
		repository: repository,
		todos: []Todo{
			{ID: 1, Description: "Goを勉強する", Done: false},
			{ID: 2, Description: "Bubble Teaを学ぶ", Done: false},
			{ID: 3, Description: "Bubble Teaを作る", Done: false},
			{ID: 3, Description: "Bubble Teaを作る", Done: false},
			{ID: 3, Description: "寝る", Done: true},
		},
		cursor: 0,
		err:    nil,
	}
}

func (m model) Init() tea.Cmd {
	return LoadTodos(m.repository)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TodosLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.todos = msg.todos
		return m, nil
	case TodoUpdateMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.todos[m.cursor] = msg.todo
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case "enter":
			todo := m.todos[m.cursor]
			todo.Done = !todo.Done
			return m, updateTodo(m.repository, todo)
		}

	}
	return m, nil
}

func (m model) View() string {
	s := "Todo List \n\n"

	for i, todo := range m.todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		status := "□"
		if todo.Done {
			status = "☒"
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, status, todo.Description)
	}
	s += "\n↑↓: 移動	Enter: 完了/未完了	q: 終了\n"
	return s
}
