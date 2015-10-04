// メイン関数
package main

import (
	"flag"
	"fmt"
	"sudoku"
	"sudoku/solver"
	"time"
)

// メイン関数
func main() {
	flag.Parse()

	// 開始時刻の記録
	startTime := time.Now()

	s := sudoku.ReadTable()
	sudoku.ApplyNormalRule(&s)

	// 問題の表示
	sudoku.PrintTable(&s)

	// 解く
	ans, point := solver.Solve(&s)

	dulation := time.Since(startTime)
	fmt.Println("解が出ました！ ", dulation)
	fmt.Println("難易度: ", point)

	// 回答の表示
	sudoku.PrintTable(&ans)

}
