package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
		return m, updateTodo(m.repository, todo)
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
		return m, updateTodo(m.repository, todo)
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

		return m, createTodo(m.repository, description)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) updateDeleteScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		todo := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo
		return m, deleteTodo(m.repository, todo.ID)
	case "n":
		m.screen = screenList
		return m, nil
	}
	return m, nil
}
