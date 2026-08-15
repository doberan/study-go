package main

import (
	"database/sql"
	"fmt"
)

type Todo struct {
	ID          int
	Description string
	Done        bool
}

// コマンドの説明を表示
func renderHelp() {
	fmt.Println("使用可能なコマンド:")
	fmt.Println("/list - Todo一覧を表示")
	fmt.Println("/add - Todoを追加")
	fmt.Println("/delete - Todoを削除")
	fmt.Println("/toggle - Todoのステータスを変更")
	fmt.Println("/exit - プログラムを終了")
	fmt.Println("/help - コマンドの説明を表示")
}

func addTodo(db *sql.DB) {
	var description string
	fmt.Print("Todoの説明を入力してください:")
	fmt.Scanln(&description)

	newTodo := Todo{Description: description, Done: false}
	_, err := db.Exec(`INSERT INTO todos (description, done) VALUES (?, ?)`,
		newTodo.Description, newTodo.Done)
	if err != nil {
		fmt.Println("Todoの追加に失敗しました:", err)
		return
	}
	fmt.Println("Todoを追加しました。")
}

func deleteTodo(db *sql.DB) {
	var id int
	fmt.Print("削除するTodoのID:")
	fmt.Scanln(&id)
	result, err := db.Exec(`DELETE FROM todos WHERE id = ?`, id)

	if err != nil {
		fmt.Println("Todoの削除に失敗しました:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		fmt.Println("Todoの削除に失敗しました:", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Println("指定されたIDのTodoは存在しません。")
		return
	}
	fmt.Println("Todoを削除しました。")
}

// Todo一覧を表示
func renderTodoList(repository *TodoRepository) {
	fmt.Println("Todo一覧")

	todos, err := repository.FindAll()

	if err != nil {
		fmt.Println("Todoの取得に失敗しました:", err)
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

func toggleTodoStatus(db *sql.DB) {
	var id int
	fmt.Print("ステータスを変更するTodoのID:")
	fmt.Scanln(&id)

	result, err := db.Exec(`UPDATE todos SET done = NOT done WHERE id = ?`, id)

	if err != nil {
		fmt.Println("Todoのステータス変更に失敗しました:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		fmt.Println("Todoのステータス変更に失敗しました:", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Println("指定されたIDのTodoは存在しません。")
		return
	}

	var done bool
	err = db.QueryRow(`SELECT done FROM todos WHERE id = ?`, id).Scan(&done)

	if err != nil {
		fmt.Println("Todoのステータス取得に失敗しました:", err)
		return
	}

	if done {
		fmt.Println("Todoのステータスを完了に変更しました。")
	} else {
		fmt.Println("Todoのステータスを未完了に変更しました。")
	}
}
