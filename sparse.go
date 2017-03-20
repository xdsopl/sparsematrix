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

type Vector struct {
	dim int
	ones []int
}

func NewVector(dim int) Vector {
	return Vector{dim, make([]int, 0)}
}

func (p Vector) Clone() Vector {
	ones := make([]int, len(p.ones))
	copy(ones, p.ones)
	return Vector{p.dim, ones}
}

func (p *Vector) Add(a Vector) Vector {
	if p.dim != a.dim {
		panic("p.dim != a.dim");
	}
	ones := append(p.ones, a.ones...)
	sort.Ints(ones)
	sum := false
	l := 0
	for k := 0; k < len(ones); k++ {
		if ones[k] == ones[l] {
			sum = !sum
		} else {
			if sum { l++ }
			ones[l] = ones[k]
			sum = true
		}
	}
	if sum { l++ }
	return Vector{p.dim, ones[:l]}
}

type RowVecMat struct {
	rows, cols int
	vecs []Vector
}

func NewRowVecMat(rows, cols int) RowVecMat {
	vecs := make([]Vector, rows)
	for i := range vecs { vecs[i] = NewVector(cols) }
	return RowVecMat{rows, cols, vecs}
}

func IdentityRowVecMat(dim int) RowVecMat {
	vecs := make([]Vector, dim)
	for i := 0; i < dim; i++ {
		vecs[i] = Vector{dim, make([]int, 1)}
		vecs[i].ones[0] = i
	}
	return RowVecMat{dim, dim, vecs}
}

func (p RowVecMat) Clone() RowVecMat {
	vecs := make([]Vector, p.rows)
	for i, vec := range p.vecs { vecs[i] = vec.Clone() }
	return RowVecMat{p.rows, p.cols, vecs}
}

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

func (p Matrix) Clone() Matrix {
	ones := make([]Coordinate, len(p.ones))
	copy(ones, p.ones)
	return Matrix{p.rows, p.cols, ones}
}

func (p RowVecMat) ConvertMatrix() Matrix {
	count := 0
	for _, vec := range p.vecs { count += len(vec.ones) }
	ones := make([]Coordinate, count)
	count = 0
	for col, vec := range p.vecs {
		for _, row := range vec.ones {
			ones[count] = Coordinate{row, col}
			count++
		}
	}
	return Matrix{p.rows, p.cols, ones}
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

func (left Matrix) Concatenate(right Matrix) Matrix {
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

func (p Matrix) Transpose() Matrix {
	ones := make([]Coordinate, len(p.ones))
	for index, one := range p.ones {
		ones[index].row = one.col
		ones[index].col = one.row
	}
	return Matrix{p.cols, p.rows, ones}
}

func (left Matrix) Multiply(right Matrix) Matrix {
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

func (p *RowVecMat) Swap(i, j int) {
	if i < 0 || i >= p.cols || j < 0 || j >= p.cols {
		panic("i < 0 || i >= p.cols || j < 0 || j >= p.cols")
	}
	p.vecs[i], p.vecs[j] = p.vecs[j], p.vecs[i]
}

func (p *RowVecMat) Add(i, j int) {
	if i < 0 || i >= p.cols || j < 0 || j >= p.cols {
		panic("i < 0 || i >= p.cols || j < 0 || j >= p.cols")
	}
	p.vecs[i] = p.vecs[i].Add(p.vecs[j])
}

func (p Matrix) IsIdentity() bool {
	if p.cols != p.rows || len(p.ones) != p.cols {
		return false
	}
	sort.Sort(ByRowCol(p.ones))
	for idx, one := range p.ones {
		if one.col != one.row || idx != one.col {
			return false
		}
	}
	return true
}

func (p Matrix) WriteImage(name string) {
	img := image.NewGray(image.Rect(0, 0, p.cols, p.rows))
	for _, one := range p.ones {
		img.Set(one.col, one.row, color.White)
	}
	file, err := os.Create(name)
	if err != nil { panic(err) }
	if err := png.Encode(file, img); err != nil { panic(err) }
	fmt.Println("Wrote " + name)
}

func (p RowVecMat) WriteImage(name string) {
	img := image.NewGray(image.Rect(0, 0, p.cols, p.rows))
	for row, vec := range p.vecs {
		for _, col := range vec.ones {
			img.Set(col, row, color.White)
		}
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
	GT := IdentityMatrix(N).Concatenate(P).Transpose()
	H := P.Transpose().Concatenate(IdentityMatrix(N))
	if (N < 1000) {
		GT.WriteImage("GT.png")
		H.WriteImage("H.png")
	}
	HGT := H.Multiply(GT)
	fmt.Println("HammingWeight of H*GT =", HGT.HammingWeight())

	/*
	Finding the inverse for a huge (random but regular) GF(2) sparse matrix
	is expensive. So let's create one and it's inverse at the same time using
	the following idea: $(\prod^{N}_{i}{E_i})^{-1}=(\prod^{N}_{i}{E_i^T})^T$
	*/
	A := IdentityRowVecMat(N)
	BT := IdentityRowVecMat(N)
	for n := 0; n < 2*N; n++ {
		var i, j int
		for i == j { i, j = rnd.Intn(N), rnd.Intn(N) }
		A.Swap(i, j)
		BT.Swap(j, i)
	}
	for n := 0; n < N/2; n++ {
		var i, j int
		for i == j { i, j = rnd.Intn(N), rnd.Intn(N) }
		A.Add(i, j)
		BT.Add(j, i)
	}
	AB := A.ConvertMatrix().Multiply(BT.ConvertMatrix().Transpose())
	if (N < 1000) {
		A.WriteImage("A.png")
		BT.WriteImage("BT.png")
		AB.WriteImage("AB.png")
	}
	fmt.Println("AB IsIdentity =", AB.IsIdentity())
}

