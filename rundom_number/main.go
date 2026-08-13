package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var result int
	var input int
	var count int
	result = rand.Intn(100) + 1 // 1から100までのランダムな整数を生成
	count = 0
	for {
		fmt.Print("数値を入力してください:")
		fmt.Scanln(&input)
		count++
		fmt.Printf("%d回目の入力です。\n", count)
		fmt.Printf("入力された数値: %d\n", input)

		if input == result {
			fmt.Println("おめでとうございます！正解です！")
			break
		} else if input < result {
			fmt.Println("残念！入力された数値は答えの数値より小さいです。")
		} else {
			fmt.Println("残念！入力された数値は答えの数値より大きいです。")
		}
	}
}