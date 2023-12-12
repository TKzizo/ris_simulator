package reducedcomplex

import "fmt"

type Cmatrix struct {
	Row, Col int
	Data     [][]complex128
}

func (c *Cmatrix) Init(row, col int) {
	c.Row = row
	c.Col = col
	c.Data = make([][]complex128, row, row)
	for i := 0; i < row; i++ {
		c.Data[i] = make([]complex128, col, col)
	}

}

func Mul(a, b Cmatrix) Cmatrix {
	if a.Col != b.Row {
		panic("mismatch of matrices sizes")
	}

	c := Cmatrix{}
	c.Init(a.Row, b.Col)

	for i, row := range a.Data {
		for x := 0; x < b.Col; x++ {
			for y := 0; y < b.Row; y++ {
				c.Data[i][x] += row[y] * b.Data[y][x]

			}
		}
	}
	return c
}

func Add(a, b Cmatrix) Cmatrix {
	if a.Col != b.Col || a.Row != b.Row {
		panic("mismatch of matrices sizes")
	}

	c := Cmatrix{}
	c.Init(a.Row, a.Col)

	for x := 0; x < a.Row; x++ {
		for y := 0; y < a.Col; y++ {
			c.Data[x][y] = a.Data[x][y] + b.Data[x][y]
		}
	}

	return c
}

func Sub(a, b Cmatrix) Cmatrix {
	if a.Col != b.Col || a.Row != b.Row {
		panic("mismatch of matrices sizes")
	}

	c := Cmatrix{}
	c.Init(a.Row, b.Col)

	for x := 0; x < a.Row; x++ {
		for y := 0; y < a.Col; y++ {
			c.Data[x][y] = a.Data[x][y] - b.Data[x][y]
		}
	}

	return c
}

func Transpose(a Cmatrix) Cmatrix {
	c := Cmatrix{}
	c.Init(a.Col, a.Row)

	for x := 0; x < a.Row; x++ {
		for y := 0; y < a.Col; y++ {
			c.Data[y][x] = a.Data[x][y]
		}
	}
	return c
}

func Scale(a Cmatrix, s complex128) Cmatrix {

	c := Cmatrix{}
	c.Init(a.Row, a.Col)

	for x := 0; x < a.Row; x++ {
		for y := 0; y < a.Col; y++ {
			c.Data[x][y] = a.Data[x][y] * s
		}
	}

	return c
}

func (c Cmatrix) String() string {
	fmt.Printf("Rows: %d,Cols: %d\n", c.Row, c.Col)
	ret := ""
	for _, line := range c.Data {
		ret = ret + fmt.Sprintln(line)
	}
	return ret
}
