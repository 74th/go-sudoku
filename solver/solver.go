package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	flag.Parse()
	table := readTable()
	applyNormalRule(&table)
	reverseRule(&table)
	solve(&table)
}

// Solve
func solve(numTable *NumTable) {
	printTable(numTable)
	// 最初の候補絞り
	reduceTable(numTable)
	// 単純に決める
	for i := 0; i < 10; i++ {
		solves := solveOneCadidate(numTable)
		solves = append(solves, solveOneAppeare(numTable)...)
		if len(solves) == 0 {
			break
		}
		for _, solve := range solves {
			fmt.Println("(", solve.x, ",", solve.y, ")->", solve.num)
			reduce(numTable, solve.x, solve.y)
		}
		printTable(numTable)
	}
}

// 各グループで、その数字が1つのマスでしか候補でなければ、解とする
func solveOneAppeare(numTable *NumTable) []AbleNum {
	solves := make([]AbleNum, 0, 0)
	table := &numTable.table
	for _, group := range numTable.groups {
		for n := 1; n < 10; n++ {
			isAns := false
			var ansMass *AbleNum
			for _, mass := range group.list {
				if table[mass.x][mass.y].candidate[n] {
					if isAns {
						isAns = false
						break
					} else {
						isAns = true
						ansMass = &table[mass.x][mass.y]
					}
				}
			}
			if isAns {
				ansMass.num = n
				ansMass.isSolve = true
				solves = append(solves, *ansMass)
			}
		}
	}
	return solves
}

// そのマスで、候補が1つの数字しかなければ、解とする
func solveOneCadidate(numTable *NumTable) []AbleNum {
	solves := make([]AbleNum, 0, 0)
	table := &numTable.table
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			if !table[x][y].isSolve {
				mass := &table[x][y]

				// trueを1つ探す
				num := 0
				for n := 1; n < 10; n++ {
					if mass.candidate[n] {
						if num == 0 {
							num = n
						} else {
							num = 0
							break
						}
					}
				}
				if num != 0 {
					mass.num = num
					mass.isSolve = true
					solves = append(solves, *mass)
				}
			}
		}
	}
	return solves
}

// テーブルプリント
func printTable(numTable *NumTable) {
	table := &numTable.table
	fmt.Println("---")
	fmt.Println("  012345678")
	fmt.Println()
	for x := 0; x < 9; x++ {
		fmt.Print(x, " ")
		for y := 0; y < 9; y++ {
			fmt.Print(table[x][y].num)
		}
		fmt.Println()
	}
	fmt.Println("---")
}

// 候補を削除する
func reduceTable(numTable *NumTable) {
	table := &numTable.table
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			if table[x][y].isSolve {
				reduce(numTable, x, y)
			}
		}
	}
}

// 候補を探す
func reduce(numTable *NumTable, x int, y int) {
	table := &numTable.table
	num := table[x][y].num
	groups := numTable.stricts[x][y]
	for i := 0; i < len(groups); i++ {
		group := groups[i]
		for j := 0; j < 9; j++ {
			mass := group.list[j]
			table[mass.x][mass.y].candidate[num] = false
		}
	}
}

// NumTable ナンプレ全体
type NumTable struct {
	table   [9][9]AbleNum
	groups  []NumPlaGroup
	stricts [9][9][]NumPlaGroup
}

// NumPlaGroup ナンプレグループ
type NumPlaGroup struct {
	list [9]NumPlaMass
}

// NumPlaMass ナンプレマス
type NumPlaMass struct {
	x, y int
}

// AbleNum マス
type AbleNum struct {
	NumPlaMass
	candidate [10]bool
	isSolve   bool
	num       int
}

// テーブル解読
func readTable() NumTable {
	numTable := NumTable{}
	table := &numTable.table

	fs, _ := os.Open(flag.Arg(0))
	reader := bufio.NewReaderSize(fs, 4096)
	for x := 0; x < 9; x++ {
		lineByte, _, _ := reader.ReadLine()
		lineStr := string(lineByte)
		for y := 0; y < 9; y++ {
			mass := &table[x][y]
			mass.x = x
			mass.y = y
			mass.num, _ = strconv.Atoi(lineStr[y : y+1])
			for k := 1; k < 10; k++ {
				mass.candidate[k] = true
			}
			if mass.num == 0 {
				mass.isSolve = false
			} else {
				mass.isSolve = true
			}

		}
	}
	return numTable
}

// 標準ルール
func applyNormalRule(numTable *NumTable) {
	groups := make([]NumPlaGroup, 0, 27)
	groups = append(groups, makeTateRule()...)
	groups = append(groups, makeYokoRule()...)
	groups = append(groups, makeBoxRule()...)
	numTable.groups = groups
}

// 縦ルール
func makeTateRule() []NumPlaGroup {
	groups := make([]NumPlaGroup, 0, 9)
	for x := 0; x < 9; x++ {
		group := NumPlaGroup{}
		for y := 0; y < 9; y++ {
			group.list[y].x = x
			group.list[y].y = y
		}
		groups = append(groups, group)
	}
	return groups
}

// 横ルール
func makeYokoRule() []NumPlaGroup {
	groups := make([]NumPlaGroup, 0, 9)
	for y := 0; y < 9; y++ {
		group := NumPlaGroup{}
		for x := 0; x < 9; x++ {
			group.list[x].x = x
			group.list[x].y = y
		}
		groups = append(groups, group)
	}
	return groups
}

// 3x3 マスずつのルール
func makeBoxRule() []NumPlaGroup {
	groups := make([]NumPlaGroup, 0, 9)
	for x1 := 0; x1 < 3; x1++ {
		for y1 := 0; y1 < 3; y1++ {
			group := NumPlaGroup{}
			i := 0
			for x2 := 0; x2 < 3; x2++ {
				for y2 := 0; y2 < 3; y2++ {
					group.list[i].x = x1*3 + x2
					group.list[i].y = y1*3 + y2
					i = i + 1
				}
			}
			groups = append(groups, group)
		}
	}
	return groups
}

// ルールを裏返す
func reverseRule(numTable *NumTable) {
	groupLen := len(numTable.groups)
	for i := 0; i < groupLen; i++ {
		group := &numTable.groups[i]
		for i := 0; i < len(group.list); i++ {
			mass := group.list[i]
			numTable.stricts[mass.x][mass.y] = append(numTable.stricts[mass.x][mass.y], *group)
		}
	}
}
