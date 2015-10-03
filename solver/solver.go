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
	fmt.Println(table.stricts[0][0])
	solve(&table)
}

// Solve
func solve(numTable *NumTable) {
	printTable(numTable)
	reduceTable(numTable)
	fmt.Println(numTable.table[5][8])
	for i := 0; i < 10; i++ {
		solve := oneCandidate(numTable)
		fmt.Println(numTable.table[5][8])
		if !solve {
			break
		}
		printTable(numTable)
	}
	fmt.Println(numTable.table[5][8])
}

// 1つ解く
func oneCandidate(numTable *NumTable) bool {
	solveLatestOne := false
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
					fmt.Println("(", x, ",", y, ")->", num)
					reduce(numTable, x, y)
					fmt.Println(numTable.table[5][8])
					solveLatestOne = true
				}
			}
		}
	}
	return solveLatestOne
}

// テーブルプリント
func printTable(numTable *NumTable) {
	table := &numTable.table
	fmt.Println("---")
	for x := 0; x < 9; x++ {
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
	mass := &table[x][y]
	// タテ列排除
	num := mass.num
	for k := 0; k < 9; k++ {
		table[x][k].candidate[num] = false
	}
	// ヨコ列排除
	for k := 0; k < 9; k++ {
		table[k][y].candidate[num] = false
	}
	// 9マス調査
	tate := x - (x % 3)
	yoko := y - (y % 3)
	fmt.Println("x", x, ",y", y, ",Tate:", tate, ",Yoko", yoko)
	for x2 := 0; x2 < 3; x2++ {
		for y2 := 0; y2 < 3; y2++ {
			table[tate+x2][yoko+y2].candidate[num] = false
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
	x, y      int
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
			for k := 0; k < 9; k++ {
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
