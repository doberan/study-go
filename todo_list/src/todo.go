package main

import (
	"database/sql"
	"errors"
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

func addTodo(repository *TodoRepository) {
	var description string
	fmt.Print("Todoの説明を入力してください:")
	fmt.Scanln(&description)

	newTodo := Todo{Description: description, Done: false}

	err := repository.Create(newTodo)
	if err != nil {
		fmt.Println("Todoの追加に失敗しました:", err)
		return
	}
	fmt.Println("Todoを追加しました。")
}

func deleteTodo(repository *TodoRepository) {
	var id int
	fmt.Print("削除するTodoのID:")
	fmt.Scanln(&id)

	err := repository.Delete(id)

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

// Todo一覧を表示
func renderTodoList(repository *TodoRepository) {
	fmt.Println("Todo一覧")

	todos, err := repository.FindAll()

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

func toggleTodoStatus(repository *TodoRepository) {
	var id int
	fmt.Print("ステータスを変更するTodoのID:")
	fmt.Scanln(&id)

	todo, findErr := repository.FindByID(id)
	if errors.Is(findErr, sql.ErrNoRows) {
		fmt.Println("指定されたTodoは存在しません")
		return
	}
	if findErr != nil {
		fmt.Println("Todoの取得に失敗しました:", findErr)
		return
	}

	todo.Done = !todo.Done

	updateErr := repository.Update(todo)
	if errors.Is(updateErr, sql.ErrNoRows) {
		fmt.Println("指定されたTodoは存在しません")
		return
	}
	if updateErr != nil {
		fmt.Println("Todoの更新に失敗しました:", updateErr)
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
