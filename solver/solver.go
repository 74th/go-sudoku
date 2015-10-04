package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	table := readTable()

	startTime := time.Now()

	applyNormalRule(&table)
	reverseRule(&table)
	solve(&table)

	dulation := time.Since(startTime)
	fmt.Println("解が出ました！ ", dulation)
}

// Solve
func solve(numTable *NumTable) {
	printTable(numTable)
	// 最初の候補絞り
	reduceTable(numTable)
	// 解く
	result := testCandidate(numTable, 0)
	if !result {
		panic("解が出ませんでした")
	}
}

// テストする
// この関数に来るときには、試しに埋めている
func testCandidate(numTable *NumTable, depth int) bool {
	// 10回簡単に解が決まらないか、試行する
	for i := 0; i < 81; i++ {
		solves := solveOneCadidate(numTable)
		solves = append(solves, solveOneAppeare(numTable)...)
		// 解が増えない場合、break
		if len(solves) == 0 {
			break
		}
		// 解が出ている場合、その解で候補を狭める
		for _, solve := range solves {
			fmt.Println(depth, ":(", solve.x, ",", solve.y, ")->", solve.num)
			reduce(numTable, solve.x, solve.y)
		}
		printTable(numTable)
	}

	// 最小の候補を導出する
	leastCandidateMass, err := searchLeastCandidateMass(numTable)
	if err != nil {
		// エラーの場合、候補が無くなった、誤りのマスがあった
		// この試行は失敗とする
		return false
	}
	if leastCandidateMass == nil {
		// 全て埋まった場合、この試行は成功とする
		// 解答チェック（同時解答で駄目なパターンあり）
		return collectAnswer(numTable)
	}
	// まだ候補がある場合、すべての候補でテストする
	for n := 1; n < 10; n++ {
		if leastCandidateMass.candidate[n] {
			// コピーを作成
			var newNumTable *NumTable
			newNumTable = &NumTable{
				groups:  numTable.groups,
				stricts: numTable.stricts,
				table:   numTable.table}
			// この候補以外falseにする
			mass := &newNumTable.table[leastCandidateMass.x][leastCandidateMass.y]
			for n2 := 1; n2 < 10; n2++ {
				if n != n2 {
					mass.candidate[n2] = false
				}
			}
			// reduceして次の解に進む
			reduce(newNumTable, leastCandidateMass.x, leastCandidateMass.y)
			fmt.Println(depth, ":Test(", mass.x, ",", mass.y, ")->", n)
			result := collectAnswer(newNumTable)
			if !result {
				// 同時回答で謝るケースあり
				return false
			}
			result = testCandidate(newNumTable, depth+1)
			if result {
				// 正解と帰ってきた場合、trueを戻す
				return true
			}
		}
	}
	// この解では得られなかった場合
	return false
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

// TryNum ナンプレ
type TryNum struct {
	NumPlaMass
	num int
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

// 候補の少ないマスを探す
func searchLeastCandidateMass(numTable *NumTable) (*AbleNum, error) {

	table := &numTable.table
	leastNum := 9
	leastMass := make([]AbleNum, 0, 81)
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			if !table[x][y].isSolve {
				candidates := &table[x][y].candidate
				num := 0
				for n := 1; n < 9; n++ {
					if candidates[n] {
						num++
					}
				}
				if num == leastNum {
					leastMass = append(leastMass, table[x][y])
				} else if num < leastNum {
					leastMass = leastMass[:0]
					leastMass = append(leastMass, table[x][y])
					leastNum = num
				} else if leastNum == 0 {
					return nil, errors.New("候補がゼロのマスが発見されました")
				}
			}
		}
	}

	// 最も候補の少ないマスがない
	// →解けた
	if len(leastMass) == 0 {
		return nil, nil
	}
	// 1個選定する
	// TODO
	return &leastMass[0], nil
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

// 回答の整合性チェック
func collectAnswer(numTable *NumTable) bool {
	table := &numTable.table
	for _, group := range numTable.groups {
		nums := [9]bool{true, true, true, true, true, true, true, true, true}
		for _, mass := range group.list {
			num := table[mass.x][mass.y].num
			if num == 0 {
				continue
			}
			if nums[num-1] {
				nums[num-1] = false
			} else {
				printTable(numTable)
				fmt.Println("(", mass.x, ",", mass.y, ")にて重複した解答発見")
				return false
			}
		}
	}
	return true
}
