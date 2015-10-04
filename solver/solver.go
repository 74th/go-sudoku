// Package solver 数独ソルバー
package solver

import (
	"errors"
	"sudoku"
)

// 正解
var _answer *solvingSudoku
var _finalPoint int

// Solve 解く
func Solve(s *sudoku.Sudoku) (sudoku.Sudoku, int) {
	_finalPoint = 0
	sol := convertSolvingSudoku(*s)
	// ルールの逆をつくる
	reverseRule(sol)
	// 最初の候補絞り
	reduceTable(sol)

	// 解く
	result := testCandidate(sol, 1, 0)
	if !result {
		panic("解が出ませんでした")
	}

	// 解を埋める
	ans := sudoku.Sudoku{
		Rules: s.Rules,
		Table: s.Table}
	setSolvingSudokuAns(&ans, _answer)
	return ans, _finalPoint
}

// テストする
// この関数に来るときには、試しに埋めている
func testCandidate(sudoku *solvingSudoku, depth int, point int) bool {
	// 10回簡単に解が決まらないか、試行する
	for i := 0; i < 81; i++ {
		solves := solveOneCadidate(sudoku)
		solves = append(solves, solveOneAppeare(sudoku)...)
		// 解が増えない場合、break
		if len(solves) == 0 {
			break
		}
		// 解が出ている場合、その解で候補を狭める
		for _, solve := range solves {
			//fmt.Println(depth, ":(", solve.x, ",", solve.y, ")->", solve.num)
			// 深さを10点とする
			point = point + depth*(i+1)
			reduce(sudoku, solve.X, solve.Y)
		}
		//printTable(sudoku)
	}

	// 最小の候補を導出する
	leastCandidateMass, num, err := searchLeastCandidateMass(sudoku)
	point = point + num*100
	if err != nil {
		// エラーの場合、候補が無くなった、誤りのマスがあった
		// この試行は失敗とする
		//fmt.Println(depth, ":REVERSE")
		return false
	}
	if leastCandidateMass == nil {
		// 全て埋まった場合、この試行は成功とする
		// 解答チェック（同時解答で駄目なパターンあり）
		result := collectAnswer(sudoku)
		if result {
			_answer = sudoku
			_finalPoint = point
			return true
		}
		return false
	}
	// まだ候補がある場合、すべての候補でテストする
	for n := 1; n < 10; n++ {
		if leastCandidateMass.candidate[n] {
			// コピーを作成
			var newsolvingSudoku *solvingSudoku
			newsolvingSudoku = &solvingSudoku{
				Rules:   sudoku.Rules,
				stricts: sudoku.stricts,
				Table:   sudoku.Table}
			// この候補以外falseにする
			mass := &newsolvingSudoku.Table[leastCandidateMass.X][leastCandidateMass.Y]
			for n2 := 1; n2 < 10; n2++ {
				if n != n2 {
					mass.candidate[n2] = false
				}
			}
			// reduceして次の解に進む
			reduce(newsolvingSudoku, leastCandidateMass.X, leastCandidateMass.Y)
			//fmt.Println(depth, ":Test(", mass.x, ",", mass.y, ")->", n)
			result := collectAnswer(newsolvingSudoku)
			if !result {
				// 同時回答で謝るケースあり
				return false
			}
			result = testCandidate(newsolvingSudoku, depth+1, point)
			if result {
				// 正解と帰ってきた場合、trueを戻す
				return true
			}
		}
	}
	// この解では得られなかった場合
	//fmt.Println(depth, ":ENDREVERSE")
	return false
}

// 各グループで、その数字が1つのマスでしか候補でなければ、解とする
func solveOneAppeare(sudoku *solvingSudoku) []solvingMass {
	solves := make([]solvingMass, 0, 0)
	table := &sudoku.Table
	for _, rule := range sudoku.Rules {
		for n := 1; n < 10; n++ {
			isAns := false
			var ansMass *solvingMass
			for _, mass := range rule.List {
				if table[mass.X][mass.Y].candidate[n] {
					if isAns {
						isAns = false
						break
					} else {
						isAns = true
						ansMass = &table[mass.X][mass.Y]
					}
				}
			}
			if isAns {
				ansMass.Num = n
				ansMass.isSolve = true
				solves = append(solves, *ansMass)
			}
		}
	}
	return solves
}

// そのマスで、候補が1つの数字しかなければ、解とする
func solveOneCadidate(sudoku *solvingSudoku) []solvingMass {
	solves := make([]solvingMass, 0, 0)
	table := &sudoku.Table
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
					mass.Num = num
					mass.isSolve = true
					solves = append(solves, *mass)
				}
			}
		}
	}
	return solves
}

// 候補を削除する
func reduceTable(sudoku *solvingSudoku) {
	table := &sudoku.Table
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			if table[x][y].isSolve {
				reduce(sudoku, x, y)
			}
		}
	}
}

// 候補を探す
func reduce(sudoku *solvingSudoku, x int, y int) {
	table := &sudoku.Table
	num := table[x][y].Num
	rules := sudoku.stricts[x][y]
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		for j := 0; j < 9; j++ {
			mass := rule.List[j]
			table[mass.X][mass.Y].candidate[num] = false
		}
	}
}

// solvingSudoku ナンプレ全体
type solvingSudoku struct {
	Table   [9][9]solvingMass
	Rules   []sudoku.Rule
	stricts [9][9][]sudoku.Rule
}

// solvingMass マス
type solvingMass struct {
	X, Y      int
	Num       int
	candidate [10]bool
	isSolve   bool
}

// 候補の少ないマスを探す
func searchLeastCandidateMass(solvingSudoku *solvingSudoku) (*solvingMass, int, error) {

	table := &solvingSudoku.Table
	leastNum := 9
	leastMass := make([]solvingMass, 0, 81)
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
					return nil, 0, errors.New("候補がゼロのマスが発見されました")
				}
			}
		}
	}

	// 最も候補の少ないマスがない
	// →解けた
	if len(leastMass) == 0 {
		return nil, 0, nil
	}
	// 1個選定する
	// TODO
	return &leastMass[0], leastNum, nil
}

// ルールを裏返す
func reverseRule(solvingSudoku *solvingSudoku) {
	ruleLen := len(solvingSudoku.Rules)
	for i := 0; i < ruleLen; i++ {
		rule := &solvingSudoku.Rules[i]
		for i := 0; i < len(rule.List); i++ {
			mass := rule.List[i]
			solvingSudoku.stricts[mass.X][mass.Y] = append(solvingSudoku.stricts[mass.X][mass.Y], *rule)
		}
	}
}

// 回答の整合性チェック
func collectAnswer(solvingSudoku *solvingSudoku) bool {
	table := &solvingSudoku.Table
	for _, rule := range solvingSudoku.Rules {
		nums := [9]bool{true, true, true, true, true, true, true, true, true}
		for _, mass := range rule.List {
			num := table[mass.X][mass.Y].Num
			if num == 0 {
				continue
			}
			if nums[num-1] {
				nums[num-1] = false
			} else {
				//fmt.Println("(", mass.x, ",", mass.y, ")にて重複した解答発見")
				return false
			}
		}
	}
	return true
}

// 数独から解く形式への変換
func convertSolvingSudoku(s sudoku.Sudoku) *solvingSudoku {
	result := solvingSudoku{Rules: s.Rules}
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			result.Table[x][y] = solvingMass{
				X:       x,
				Y:       y,
				Num:     s.Table[x][y].Num,
				isSolve: s.Table[x][y].Num != 0}
			for i := 1; i < 10; i++ {
				result.Table[x][y].candidate[i] = true
			}
		}
	}
	return &result
}

// 数独から解く形式への変換
func setSolvingSudokuAns(s *sudoku.Sudoku, sol *solvingSudoku) {
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			s.Table[x][y].Num = sol.Table[x][y].Num
		}
	}
}
