// ICM de-noising
package icm

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Node struct {
	yVal, xVal int
	// x, y       int
}

func (n *Node) X() int {
	return n.xVal
}

func (n *Node) Y() int {
	return n.yVal
}

func (n *Node) Set(val int) {
	n.xVal = val
}

func (u *UGM) Node(x, y int) *Node {
	return u.image[x][y]
}

type UGM struct {
	image                    [][]*Node
	h, beta, eta             float64
	h_sum, beta_sum, eta_sum int
}

func NewUGM(img image.Image, h, beta, eta float64) *UGM {

	// img.Rect.Dx()

	Xmin, Xmax := img.Bounds().Min.X, img.Bounds().Max.X
	Ymin, Ymax := img.Bounds().Min.Y, img.Bounds().Max.Y

	var black, white color.Color
	mat := make([][]*Node, Xmax-Xmin)
	for x := Xmin; x < Xmax; x++ {
		mat[x] = make([]*Node, Ymax)
		for y := Ymin; y < Ymax; y++ {
			if x == 0 && y == 0 {
				black = img.At(x, y)
			}
			col := img.At(x, y)

			var yVal int
			if col == black {
				yVal = -1
			} else {
				if white == nil {
					white = col
				}
				yVal = 1
			}
			mat[x][y] = &Node{
				yVal: yVal,
				xVal: yVal,
			}
		}
	}

	return &UGM{
		image: mat,
		h:     h,
		beta:  beta,
		eta:   eta,
	}
}

func (u *UGM) getNode(x, y int) *Node {
	if x < 0 || x >= len(u.image) {
		return nil
	}

	if y < 0 || y >= len(u.image[0]) {
		return nil
	}

	return u.image[x][y]

}

func (u *UGM) Left(x, y int) *Node {
	if y <= 0 {
		return nil
	}
	return u.image[x][y-1]
}

func (u *UGM) Right(x, y int) *Node {
	if y >= len(u.image[0])-1 {
		return nil
	}
	return u.image[x][y+1]
}

func (u *UGM) Up(x, y int) *Node {
	if x <= 0 {
		return nil
	}
	return u.image[x-1][y]
}

func (u *UGM) Down(x, y int) *Node {
	if x >= len(u.image)-1 {
		return nil
	}
	return u.image[x+1][y]
}

func (u *UGM) neighbours(x, y int) []*Node {

	neighbours := make([]*Node, 0, 4)

	xj := u.Up(x, y)
	if xj != nil {
		neighbours = append(neighbours, xj)
	}
	xj = u.Left(x, y)
	if xj != nil {
		neighbours = append(neighbours, xj)
	}
	xj = u.Right(x, y)
	if xj != nil {
		neighbours = append(neighbours, xj)
	}
	xj = u.Down(x, y)
	if xj != nil {
		neighbours = append(neighbours, xj)
	}
	return neighbours
}

func (u *UGM) X() int {
	return len(u.image)
}

func (u *UGM) Y() int {
	if len(u.image) == 0 {
		return 0
	}
	return len(u.image[0])
}

// We run through the materix and calculate the H and Eta terms
// as in the ICM, but we only sum up the nodes to the right and below
// each node, so as to prevent counting edges twice in the summation
func (u *UGM) E() float64 {
	E := 0.0
	for x, row := range u.image {
		for y, node := range row {
			beta_sum := 0
			right := u.Right(x, y)
			if right != nil {
				beta_sum += node.X() * right.X()
			}
			down := u.Down(x, y)
			if down != nil {
				beta_sum += node.X() * down.X()
			}
			eta_sum := node.X() * node.Y()

			E += u.h*float64(node.X()) - u.beta*float64(beta_sum) - u.eta*float64(eta_sum)
		}
	}

	return E
}

// The values grow from the top right of the picture.
// for this purpose, X is the row index,
// Y is the Column index (counter intuitive i know)
func (u *UGM) ICM(iter int) {
	for i := 0; i < iter; i++ {

		for x, row := range u.image {
			for y, node := range row {
				beta_minus := 0
				beta_plus := 0
				for _, n := range u.neighbours(x, y) {
					beta_minus += -1 * n.X()
					beta_plus += 1 * n.X()
				}

				eta_minus, eta_plus := -1*node.Y(), 1*node.Y()

				E_minus := u.h*-1 - u.beta*float64(beta_minus) - u.eta*float64(eta_minus)
				E_plus := u.h*1 - u.beta*float64(beta_plus) - u.eta*float64(eta_plus)

				// fmt.Printf("(%3d,%3d) plus: %1f vs. %1f minus\n", x, y, E_plus, E_minus)

				if E_minus < E_plus {
					node.Set(-1)
				} else {
					node.Set(1)
				}
			}
		}
	}
}

func (u *UGM) WriteToFile(path string, img image.Image) {
	// make new file
	newImg := image.NewGray16(img.Bounds())
	for x := 0; x < u.X(); x++ {
		for y := 0; y < u.Y(); y++ {
			if u.Node(x, y).X() == -1 {
				newImg.SetGray16(x, y, color.Black)
			} else {
				newImg.SetGray16(x, y, color.White)
			}
		}
	}
	// create/overwrite file
	imgFile, err := os.Create(path)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	defer imgFile.Close() // close when exiting..

	// write png to file
	if err := png.Encode(imgFile, newImg); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
