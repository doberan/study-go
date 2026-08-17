package main

import (
	"database/sql"
	"errors"
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
	repository := &TodoRepository{db: db}
	p := tea.NewProgram(initialModel(repository))
	if _, err := p.Run(); err != nil {
		fmt.Println("エラー:", err)
		return
	}
}

func commandList(repository *TodoRepository) {
	fmt.Println("Todo一覧")

	todos, err := renderTodoList(repository)
	if err != nil {
		fmt.Println("Todoの取得に失敗しました:", err)
		return
	}
	if len(todos) == 0 {
		fmt.Println("Todoは存在しません")
		return
	}

	for _, todo := range todos {
		var status string
		if todo.Done {
			status = "☒"
		} else {
			status = "☐"
		}
		fmt.Printf("| %d | %s | %s |\n", todo.ID, todo.Description, status)
	}
}

func commandAdd(repository *TodoRepository) {
	var description string
	fmt.Print("Todoの説明を入力してください:")
	fmt.Scanln(&description)

	err := addTodo(repository, description)
	if err != nil {
		fmt.Println("Todoの追加に失敗しました:", err)
		return
	}
	fmt.Println("Todoを追加しました。")
}

func commandDelete(repository *TodoRepository) {
	var id int
	fmt.Print("削除するTodoのID:")
	fmt.Scanln(&id)

	err := deleteTodo(repository, id)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("指定したTodoは存在しません")
		return
	}
	if err != nil {
		fmt.Println("Todoの削除に失敗しました:", err)
		return
	}
	fmt.Println("Todoを削除しました。")
}

func commandHelp() {
	fmt.Println("使用可能なコマンド:")
	fmt.Println("/list - Todo一覧を表示")
	fmt.Println("/add - Todoを追加")
	fmt.Println("/delete - Todoを削除")
	fmt.Println("/toggle - Todoのステータスを変更")
	fmt.Println("/exit - プログラムを終了")
	fmt.Println("/help - コマンドの説明を表示")
}

func commandToggle(repository *TodoRepository) {
	var id int
	fmt.Print("ステータスを変更するTodoのID:")
	fmt.Scanln(&id)

	todo, err := toggleTodoStatus(repository, id)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("指定したTodoは存在しません")
		return
	}
	if err != nil {
		fmt.Println("Todoのステータスの更新に失敗しました:", err)
		return
	}

	var status string
	if todo.Done {
		status = "完了"
	} else {
		status = "未完了"
	}
	fmt.Printf("Todoのステータスを%sに更新しました。\n", status)
}
