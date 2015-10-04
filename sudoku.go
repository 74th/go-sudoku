package sudoku

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Sudoku 数独
type Sudoku struct {
	Table [9][9]Mass
	Rules []Rule
}

// Rule 数独ルール
type Rule struct {
	List [9]Mass
}

// Mass マス
type Mass struct {
	X, Y int
	Num  int
}

// ReadTable テーブル解読
func ReadTable() Sudoku {
	sudoku := Sudoku{}
	table := &sudoku.Table

	fs, _ := os.Open(flag.Arg(0))
	reader := bufio.NewReaderSize(fs, 4096)
	for x := 0; x < 9; x++ {
		lineByte, _, _ := reader.ReadLine()
		lineStr := string(lineByte)
		for y := 0; y < 9; y++ {
			mass := &table[x][y]
			mass.X = x
			mass.Y = y
			mass.Num, _ = strconv.Atoi(lineStr[y : y+1])
		}
	}
	return sudoku
}

// ApplyNormalRule 標準ルール
func ApplyNormalRule(s *Sudoku) {
	rules := make([]Rule, 0, 27)
	rules = append(rules, makeTateRule()...)
	rules = append(rules, makeYokoRule()...)
	rules = append(rules, makeBoxRule()...)
	s.Rules = rules
}

// 縦ルール
func makeTateRule() []Rule {
	rules := make([]Rule, 0, 9)
	for x := 0; x < 9; x++ {
		rule := Rule{}
		for y := 0; y < 9; y++ {
			rule.List[y].X = x
			rule.List[y].Y = y
		}
		rules = append(rules, rule)
	}
	return rules
}

// 横ルール
func makeYokoRule() []Rule {
	rules := make([]Rule, 0, 9)
	for y := 0; y < 9; y++ {
		rule := Rule{}
		for x := 0; x < 9; x++ {
			rule.List[x].X = x
			rule.List[x].Y = y
		}
		rules = append(rules, rule)
	}
	return rules
}

// 3x3 マスずつのルール
func makeBoxRule() []Rule {
	rules := make([]Rule, 0, 9)
	for x1 := 0; x1 < 3; x1++ {
		for y1 := 0; y1 < 3; y1++ {
			rule := Rule{}
			i := 0
			for x2 := 0; x2 < 3; x2++ {
				for y2 := 0; y2 < 3; y2++ {
					rule.List[i].X = x1*3 + x2
					rule.List[i].Y = y1*3 + y2
					i = i + 1
				}
			}
			rules = append(rules, rule)
		}
	}
	return rules
}

// PrintTable テーブルプリント
func PrintTable(s *Sudoku) {
	table := &s.Table
	fmt.Println("---")
	fmt.Println("  012345678")
	fmt.Println()
	for x := 0; x < 9; x++ {
		fmt.Print(x, " ")
		for y := 0; y < 9; y++ {
			fmt.Print(table[x][y].Num)
		}
		fmt.Println()
	}
	fmt.Println("---")
}
