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
	service TodoService

	todoList list.Model // 表示するTodo一覧

	err error

	screen    screen
	textInput textinput.Model
}

func (m model) Init() tea.Cmd {
	return LoadTodos(m.service)
}

func NewModel(service TodoService) model {
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

type TodosLoadedMsg struct {
	todos []Todo
	err   error
}

func LoadTodos(service TodoService) tea.Cmd {
	return func() tea.Msg {
		todos, err := service.List()

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

func updateTodo(service TodoService, todo Todo) tea.Cmd {
	return func() tea.Msg {
		todo, err := service.Update(todo)
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

func createTodo(service TodoService, description string) tea.Cmd {
	return func() tea.Msg {
		todo, err := service.Create(description)

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

func deleteTodo(service TodoService, id int) tea.Cmd {
	return func() tea.Msg {
		err := service.Delete(id)
		return TodoDeleteMsg{
			id:  id,
			err: err,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TodosLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.todoList.SetItems(newTodoItems(msg.todos))
		return m, nil
	case TodoUpdateMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.screen = screenList
		items := m.todoList.Items()
		for i, item := range items {
			current := item.(todoItem)
			if current.todo.ID == msg.todo.ID {
				items[i] = todoItem{
					todo: msg.todo,
				}
				break
			}
		}
		m.todoList.SetItems(items)
		m.textInput.Reset()
		m.textInput.Blur()
		return m, nil
	case TodoCreateMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		items := m.todoList.Items()

		items = append(items, todoItem{todo: msg.todo})

		m.todoList.SetItems(items)

		m.screen = screenList
		m.textInput.Reset()
		m.textInput.Blur()
		return m, nil
	case TodoDeleteMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		items := m.todoList.Items()
		for i, item := range items {
			current := item.(todoItem)
			if current.todo.ID == msg.id {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		m.todoList.SetItems(items)

		if len(items) > 0 && m.todoList.Cursor() >= len(items) {
			m.todoList.Select(len(items) - 1)
		}

		m.screen = screenList
		return m, nil
	case tea.KeyMsg:
		switch m.screen {
		case screenList:
			return m.updateListScreen(msg, cmd)
		case screenAdd:
			return m.updateAddScreen(msg)
		case screenEdit:
			return m.updateEditScreen(msg)
		case screenDelete:
			return m.updateDeleteScreen(msg)
		}
	}
	return m, cmd
}

func (m model) updateListScreen(msg tea.KeyMsg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	m.todoList, cmd = m.todoList.Update(msg)

	switch msg.String() {
	case "q":
		return m, tea.Quit
	case " ":
		if len(m.todoList.Items()) == 0 {
			return m, nil
		}
		todo := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo
		todo.Done = !todo.Done
		return m, updateTodo(m.service, todo)
	case "a":
		m.screen = screenAdd
		m.textInput.Focus()

		return m, textinput.Blink
	case "e":
		if len(m.todoList.Items()) == 0 {
			return m, nil
		}
		todo := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo

		m.textInput.SetValue(todo.Description)
		m.textInput.Focus()
		m.screen = screenEdit

		return m, textinput.Blink
	case "d":
		if len(m.todoList.Items()) == 0 {
			return m, nil
		}
		m.screen = screenDelete
		return m, nil
	}
	return m, nil
}

func (m model) updateEditScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenList
		m.textInput.Reset()
		m.textInput.Blur()
		return m, nil
	case "enter":
		description := m.textInput.Value()

		if description == "" {
			return m, nil
		}

		todo := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo
		todo.Description = description
		return m, updateTodo(m.service, todo)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) updateAddScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenList
		m.textInput.Reset()
		m.textInput.Blur()

		return m, nil
	case "enter":
		description := m.textInput.Value()

		if description == "" {
			return m, nil
		}

		return m, createTodo(m.service, description)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) updateDeleteScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		todo := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo
		return m, deleteTodo(m.service, todo.ID)
	case "n":
		m.screen = screenList
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case screenAdd:
		return m.viewAdd()
	case screenDelete:
		return m.viewDelete()
	case screenEdit:
		return m.viewEdit()
	default:
		return m.viewList()
	}
}

func (m model) viewEdit() string {
	s := "Edit Todo\n\n"
	s += m.textInput.View()
	s += "\n\nEnter: 保存	Esc: キャンセル"

	return s
}

func (m model) viewAdd() string {
	s := "Add Todo\n\n"
	s += m.textInput.View()
	s += "\n\nEnter: 追加	Esc: キャンセル"

	return s
}

func (m model) viewList() string {
	s := m.todoList.View()

	s += "\n"
	s += "↑↓/jk: 移動  Space: 完了/未完了  a: 追加  e: 編集  d: 削除  q: 終了\n"

	return s
}

func (m model) viewDelete() string {
	description := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo.Description
	s := "Delete Todo\n\n"
	s += description + " を削除しますか？\n\n"
	s += "y: 削除 n: キャンセル\n"

	return s
}
