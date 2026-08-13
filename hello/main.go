package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("名前を入力してください:")
	var name string
	fmt.Scanln(&name)
	fmt.Printf("こんにちは、%sさん！\n", name)
	os.Exit(0)
}
