package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/scanner"
)

const (
	DIM_1 = 3
	DIM_2 = DIM_1 * DIM_1
	DIM_3 = DIM_1 * DIM_2
	DIM_4 = DIM_1 * DIM_3
)

type set [DIM_2]bool
type grid [DIM_4]set

/* Helpers */

// row returns the indices, excluding idx, in the same row as idx
func row(idx int) []int {
	out := make([]int, DIM_2-1)
	start := (idx / DIM_2) * DIM_2
	p := 0
	for n := 0; n < DIM_2; n++ {
		i := start + n
		if i == idx {
			continue
		}
		out[p] = i
		p++
	}
	return out
}

// column returns the indices, excluding idx, in the same column as idx
func column(idx int) []int {
	out := make([]int, DIM_2-1)
	start := idx % DIM_2
	p := 0
	for n := 0; n < DIM_2; n++ {
		i := start + DIM_2*n
		if i == idx {
			continue
		}
		out[p] = i
		p++
	}
	return out
}

// square returns the indices, excluding idx, in the same square as idx
func square(idx int) []int {
	out := make([]int, DIM_2-1)
	start := (idx/DIM_3)*DIM_3 + ((idx%DIM_2)/DIM_1)*DIM_1
	p := 0
	for n := 0; n < DIM_1; n++ {
		for m := 0; m < DIM_1; m++ {
			i := start + DIM_2*n + m
			if i == idx {
				continue
			}
			out[p] = i
			p++
		}
	}
	return out
}

func (s set) size() int {
	var n int
	for _, b := range s {
		if b {
			n++
		}
	}
	return n
}

func (s set) candidates() []int {
	out := make([]int, s.size())
	i := 0
	for n, b := range s {
		if b {
			out[i] = n
			i++
		}
	}
	return out
}

func (s set) value() int {
	if s.size() > 1 {
		return -1
	}
	for n, b := range s {
		if b {
			return n
		}
	}
	panic("unreachable")
}

/* Co-routines implementing the actual algorithm */

func spread(g grid, idx int, val int) (grid, bool) {
	for _, f := range []func(int) []int{row, column, square} {
		for _, i := range f(idx) {
			if g[i][val] {
				g[i][val] = false
				switch g[i].size() {
				case 0:
					return g, false
				case 1:
					var b bool
					g, b = assign(g, i, g[i].value())
					if !b {
						return g, false
					}
				}
			}
		}
	}
	return g, true
}

func assign(g grid, idx int, val int) (grid, bool) {
	for i := range g[idx] {
		g[idx][i] = false
	}
	g[idx][val] = true

	var b bool
	g, b = spread(g, idx, val)
	if !b {
		return g, false
	}

	for _, f := range []func(int) []int{row, column, square} {
		for _, i := range f(idx) {
			if g[i].value() == g[idx].value() {
				return g, false
			}
		}
	}

	return g, true
}

func solve(g grid) (grid, bool) {
	for i := 0; i < len(g); i++ {
		if g[i].size() >= 2 {
			for _, n := range g[i].candidates() {
				new_g, b := assign(g, i, n)
				if !b {
					g[i][n] = false
					continue
				}
				new_g, b = solve(new_g)
				if !b {
					g[i][n] = false
					continue
				}
				return solve(new_g)
			}
			return g, false
		}
	}
	return g, true
}

/* Output */

func (s set) String() string {
	if s.size() > 1 {
		return "."
	}
	for n, b := range s {
		if b {
			return fmt.Sprintf("%d", n+1)
		}
	}
	return fmt.Sprintf("?")
}

func (g grid) String() string {
	str := ""
	for i, s := range g {
		str = fmt.Sprintf("%s%s", str, s)
		if i%DIM_2 == DIM_2-1 {
			str = fmt.Sprintf("%s\n", str)
		} else {
			str = fmt.Sprintf("%s ", str)
		}
	}
	return str
}

/* Input */

func makeGrid(in [DIM_4]int) grid {
	var g grid
	for i := range g {
		var s set
		g[i] = s
		if val := in[i]; val == -1 {
			for j := range s {
				g[i][j] = true
			}
		} else {
			g[i][val] = true
		}
	}
	return g
}

func parseSudoku(str string) ([DIM_4]int, error) {
	var out [DIM_4]int
	sr := &scanner.Scanner{}
	sr = sr.Init(strings.NewReader(str))
	sr.Mode = scanner.ScanInts
	i := 0
	for r := sr.Scan(); r != scanner.EOF; r = sr.Scan() {
		n := 0
		if r == scanner.Int {
			n64, err := strconv.ParseInt(sr.TokenText(), 0, 0)
			n = int(n64)
			if err != nil {
				return out, fmt.Errorf("error parsing sudoku: %v", err)
			}
		}

		// Sets are implemented by indexing an array of booleans: subtract one
		if i == len(out) {
			return out, fmt.Errorf("error parsing sudoku: too many numbers")
		}
		out[i] = n - 1

		i++
	}
	return out, nil
}

func main() {
	str := `9 0 0 0 0 0 0 0 1
 0 0 0 0 0 3 0 8 5
 0 0 1 0 2 0 0 0 0
 0 0 0 5 0 7 0 0 0
 0 0 4 0 0 0 1 0 0
 0 9 0 0 0 0 0 0 0
 5 0 0 0 0 0 0 7 3
 0 0 2 0 1 0 0 0 0
 0 0 0 0 4 0 0 0 9`

	in, err := parseSudoku(str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
		os.Exit(1)
	}

	g := makeGrid(in)
	fmt.Println(g)
	g, b := solve(g)

	if !b {
		fmt.Println("impossible sudoku")
	} else {
		fmt.Println(g)
	}
}
