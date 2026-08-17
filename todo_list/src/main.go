package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := initDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	repository := &TodoRepository{db: db}

	fmt.Println("Hello World!")

	for {
		fmt.Print("コマンド入力:")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "/list":
			renderTodoList(repository)
		case "/add":
			addTodo(repository)
		case "/delete":
			deleteTodo(repository)
		case "/toggle":
			toggleTodoStatus(repository)
		case "/exit":
			fmt.Println("終了します。")
			return
		case "/help":
			renderHelp()
		default:
			fmt.Println("無効なコマンドです。")
			renderHelp()
		}
	}
}
