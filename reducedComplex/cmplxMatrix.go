package reducedcomplex

type Cmatrix struct {
	row, col int
	Data     [][]complex128
}

func (c *Cmatrix) Init(row, col int) {
	c.row = row
	c.col = col
	c.Data = make([][]complex128, row, row)
	for i := 0; i < row; i++ {
		c.Data[i] = make([]complex128, col, col)
	}

}

func Mul(a, b Cmatrix) Cmatrix {
	if a.col != b.row {
		panic("mismatch of matrices sizes")
	}

	c := Cmatrix{}
	c.Init(a.row, b.col)

	for i, row := range a.Data {
		for x := 0; x < b.col; x++ {
			for y := 0; y < b.row; y++ {
				c.Data[i][x] += row[y] * b.Data[y][x]

			}
		}
	}
	return c
}

func Add(a, b Cmatrix) Cmatrix {
	if a.col != b.col || a.row != b.row {
		panic("mismatch of matrices sizes")
	}

	c := Cmatrix{}
	c.Init(a.row, a.col)

	for x := 0; x < a.row; x++ {
		for y := 0; y < a.col; y++ {
			c.Data[x][y] = a.Data[x][y] + b.Data[x][y]
		}
	}

	return c
}

func Sub(a, b Cmatrix) Cmatrix {
	if a.col != b.col || a.row != b.row {
		panic("mismatch of matrices sizes")
	}

	c := Cmatrix{}
	c.Init(a.row, b.col)

	for x := 0; x < a.row; x++ {
		for y := 0; y < a.col; y++ {
			c.Data[x][y] = a.Data[x][y] - b.Data[x][y]
		}
	}

	return c
}

func Transpose(a Cmatrix) Cmatrix {
	c := Cmatrix{}
	c.Init(a.col, a.row)

	for x := 0; x < a.row; x++ {
		for y := 0; y < a.col; y++ {
			c.Data[y][x] = a.Data[x][y]
		}
	}
	return c
}

func Scale(a Cmatrix, s complex128) Cmatrix {

	c := Cmatrix{}
	c.Init(a.row, a.col)

	for x := 0; x < a.row; x++ {
		for y := 0; y < a.col; y++ {
			c.Data[x][y] = a.Data[x][y] * s
		}
	}

	return c
}
