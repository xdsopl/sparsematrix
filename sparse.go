/*
sparse - GF(2) sparse matrix fun
Written in 2017 by <Ahmet Inan> <xdsopl@gmail.com>
To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide. This software is distributed without any warranty.
You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.
*/

package main
import (
	"os"
	"fmt"
	"sort"
	"math/rand"
	"image"
	"image/png"
	"image/color"
	"time"
)

type Coordinate struct {
	row, col int
}

type ByColRow []Coordinate
func (a ByColRow) Len() int { return len(a) }
func (a ByColRow) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByColRow) Less(i, j int) bool {
	return a[i].col < a[j].col ||
		a[i].col == a[j].col &&
		a[i].row < a[j].row
}

type ByRowCol []Coordinate
func (a ByRowCol) Len() int { return len(a) }
func (a ByRowCol) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRowCol) Less(i, j int) bool {
	return a[i].row < a[j].row ||
		a[i].row == a[j].row &&
		a[i].col < a[j].col
}

type Matrix struct {
	rows, cols int
	ones []Coordinate
}

func NewMatrix(rows, cols int) Matrix {
	return Matrix{rows, cols, make([]Coordinate, 0)}
}

func IdentityMatrix(dimension int) Matrix {
	diagonal := make([]Coordinate, dimension)
	for i := 0; i < dimension; i++ {
		diagonal[i] = Coordinate{i, i}
	}
	return Matrix{dimension, dimension, diagonal}
}

func CloneMatrix(source Matrix) Matrix {
	ones := make([]Coordinate, len(source.ones))
	copy(ones, source.ones)
	return Matrix{source.rows, source.cols, ones}
}

func (p *Matrix) AddUnchecked(row, col int) {
	if row < 0 || row >= p.rows || col < 0 || col >= p.cols {
		panic("row < 0 || row >= p.rows || col < 0 || col >= p.cols")
	}
	p.ones = append(p.ones, Coordinate{row, col})
}

func (p *Matrix) RemoveDuplicates() {
	sort.Sort(ByRowCol(p.ones))
	j := 0
	for i := 0; i < len(p.ones); i++ {
		if p.ones[i] != p.ones[j] {
			j++
			p.ones[j] = p.ones[i]
		}
	}
	p.ones = p.ones[:j + 1]
}

func Concatenate(left, right Matrix) Matrix {
	if left.rows != right.rows {
		panic("left.rows != right.rows")
	}
	rows := left.rows
	cols := left.cols + right.cols
	ones := make([]Coordinate, len(left.ones) + len(right.ones))
	copy(ones, left.ones)
	for i := 0; i < len(right.ones); i++ {
		ones[len(left.ones) + i].col = right.ones[i].col + left.cols
		ones[len(left.ones) + i].row = right.ones[i].row
	}
	return Matrix{rows, cols, ones}
}

func (p Matrix) HammingWeight() int {
	return len(p.ones)
}

func (p *Matrix) HammingWeightsOfRows() []int {
	weights := make([]int, p.rows)
	for _, one := range p.ones {
		weights[one.row]++
	}
	return weights
}

func (p *Matrix) HammingWeightsOfCols() []int {
	weights := make([]int, p.cols)
	for _, one := range p.ones {
		weights[one.col]++
	}
	return weights
}

func MinMax(a []int) (int, int) {
	min, max := a[0], a[0]
	for _, e := range a {
		if e > max { max = e }
		if e < min { min = e }
	}
	return min, max;
}

func Transpose(source Matrix) Matrix {
	ones := make([]Coordinate, len(source.ones))
	for index, one := range source.ones {
		ones[index].row = one.col
		ones[index].col = one.row
	}
	return Matrix{source.cols, source.rows, ones}
}

func Multiply(left, right Matrix) Matrix {
	if left.cols != right.rows {
		panic("left.cols != right.rows")
	}
	rows := left.rows
	cols := right.cols
	ones := make([]Coordinate, 0)
	sort.Sort(ByRowCol(left.ones))
	sort.Sort(ByColRow(right.ones))
	for lBegin, lEnd := 0, 0; lBegin < len(left.ones); lBegin = lEnd {
		row := left.ones[lBegin].row
		for r := 0; r < len(right.ones); {
			col := right.ones[r].col
			sum := false
			for l := lBegin; l < len(left.ones) && r < len(right.ones); {
				if left.ones[l].row != row || right.ones[r].col != col {
					if lEnd < l { lEnd = l }
					break
				}
				if left.ones[l].col > right.ones[r].row {
					r++
					continue
				}
				if left.ones[l].col < right.ones[r].row {
					l++
					continue
				}
				sum = !sum
				l++
				r++
			}
			if sum { ones = append(ones, Coordinate{row, col}) }
			for ; r < len(right.ones) && col == right.ones[r].col; r++ {}
		}
		for ; lEnd < len(left.ones) && row == left.ones[lEnd].row; lEnd++ {}
	}
	return Matrix{rows, cols, ones}
}

func (p *Matrix) WriteImage(name string) {
	img := image.NewGray(image.Rect(0, 0, p.cols, p.rows))
	for _, one := range p.ones {
		img.Set(one.col, one.row, color.White)
	}
	file, err := os.Create(name)
	if err != nil { panic(err) }
	if err := png.Encode(file, img); err != nil { panic(err) }
	fmt.Println("Wrote " + name)
}

func main() {
	N := 500
	P := NewMatrix(N, N)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < N; i++ {
		P.AddUnchecked(rnd.Intn(P.rows), rnd.Intn(P.cols))
	}
	P.RemoveDuplicates()
	fmt.Println("HammingWeight of P =", P.HammingWeight())
	MinRowWeight, MaxRowWeight := MinMax(P.HammingWeightsOfRows())
	MinColWeight, MaxColWeight := MinMax(P.HammingWeightsOfCols())
	fmt.Println("(Min, Max) of HammingWeightsOfRows of P =", MinRowWeight, MaxRowWeight)
	fmt.Println("(Min, Max) of HammingWeightsOfCols of P =", MinColWeight, MaxColWeight)
	GT := Transpose(Concatenate(IdentityMatrix(N), P))
	GT.WriteImage("GT.png")
	H := Concatenate(Transpose(P), IdentityMatrix(N))
	H.WriteImage("H.png")
	HGT := Multiply(H, GT)
	fmt.Println("HammingWeight of H*GT =", HGT.HammingWeight())
}

