package ui

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
