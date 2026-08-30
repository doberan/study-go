package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.screen {
	case screenDelete:
		return m.viewDelete()
	default:
		return m.viewList()
	}
}

func (m model) viewList() string {
	var searchBox string

	boxWidth := m.width - 2
	if boxWidth < 20 {
		boxWidth = 20
	}

	headerHeight := 0
	if m.screen == screenAdd || m.screen == screenEdit {
		headerHeight = 3
	}

	listHeight := m.height - headerHeight - 3
	if listHeight < 3 {
		listHeight = 3
	}
	m.todoList.SetHeight(listHeight)

	content := m.todoList.View()
	help := "[↑↓/jk: 移動][Space: 完了/未完了][a: 追加][e: 編集][d: 削除][q: 終了]"

	var listStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(boxWidth)

	listBox := listStyle.Render(content)

	if m.screen == screenAdd || m.screen == screenEdit {
		searchContent := m.textInput.View()
		help = "[Enter: 保存][Esc: キャンセル]"

		var searchStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(boxWidth)

		searchBox = searchStyle.Render(searchContent)
	}

	return lipgloss.JoinVertical(lipgloss.Left, searchBox, listBox, help)
}

func (m model) viewDelete() string {
	description := m.todoList.Items()[m.todoList.Cursor()].(todoItem).todo.Description
	s := "Delete Todo\n\n"
	s += description + " を削除しますか？\n\n"
	s += "y: 削除 n: キャンセル\n"

	return s
}
