package qc

import (
	"fmt"
	"errors"
	"strconv"
	"bufio"
	"os"

	"github.com/Arafatk/glot"
	"github.com/spf13/cobra"

	"github.com/wangkui0508/float128"
)

const (
	MaxArea = int64(900)*int64(10000_0000)*int64(10000_0000)
)

var (
	oneThird = float128.F128FromF64(float64(1)/3.0)
	half = float128.F128FromF64(0.5)
	maxAreaF = float128.F128FromI64(MaxArea)
)

type QuadCurve struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
}

func NewQuadCurveFromTriplePoints(x1, y1, x2, y2, x3, y3 float64) (QuadCurve, error) {
	var curve QuadCurve
	a1 := (x1-x2)*(x1-x3)
	if a1 == 0 {
		return curve, errors.New("a1 is too large")
	}
	a1 = y1/a1
	b1 := -a1 * (x2 + x3)
	c1 := a1 * x2 * x3

	a2 := (x1-x2)*(x2-x3)
	if a2 == 0 {
		return curve, errors.New("a2 is too large")
	}
	a2 = -y2/a2
	b2 := -a2 * (x1 + x3)
	c2 := a2 * x1 * x3

	a3 := (x1-x3)*(x3-x2)
	if a3 == 0 {
		return curve, errors.New("a3 is too large")
	}
	a3 = -y3/a3
	b3 := -a3 * (x1 + x2)
	c3 := a3 * x1 * x2

	curve.A = a1+a2+a3
	curve.B = b1+b2+b3
	curve.C = c1+c2+c3
	return curve, nil
}

func (curve QuadCurve) CalcY(x float64) float64 {
	return curve.A*x*x + curve.B*x + curve.C
}

func (curve QuadCurve) CalcArea(x int64, scale int64) int64 {
	if scale <= 0 {
		return -1
	}
	A := float128.F128FromF64(curve.A)
	B := float128.F128FromF64(curve.B)
	C := float128.F128FromF64(curve.C)
	X := float128.F128FromI64(x)
	S := float128.F128FromI64(scale)
	res := C.Mul(X)
	X2 := X.Mul(X)
	res = res.Add( B.Mul(X2).Mul(half) )
	X3 := X2.Mul(X)
	res = res.Add( A.Mul(X3).Mul(oneThird) )
	res = res.Div( S )
	if res.GTE(maxAreaF) {
		return -1
	}
	return res.ToI64()
}

func (curve QuadCurve) Draw(start, stop, step float64, filename string) {
	dimensions := 2
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)
	fct := func(x float64) float64 { return (curve.CalcY(x)) }
	groupName := "Quadratic Curve"
	style := "lines"
	pointsX := make([]float64,0,100)
	for f:=start; f<stop; f=f+step {
		pointsX=append(pointsX, f)
	}
	plot.AddFunc2d(groupName, style, pointsX, fct)
	plot.SavePlot(filename)
	fmt.Printf("Plot was saved to %s\n", filename)
}

func DrawQuadCurveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draw-quad-curve [x1] [y1] [x2] [y2] [x3] [y3] [filename]",
		Short: "Draw a Quadratic Curve through three points.",
		Long: `Draw a Quadratic Curve through (x1,y1), (x2,y2) and (x3,y3), and then save it to filename.",.
Example: ./quadcurve draw-quad-curve 1.0 5.0 2.0 3.0 3.0 3.0 a.png
`,
		Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var v [6]float64
			for i:=0; i<6; i++ {
				v[i], err = strconv.ParseFloat(args[i], 64)
				if err != nil {
					return err
				}
			}
			curve, err := NewQuadCurveFromTriplePoints(v[0],v[1],v[2],v[3],v[4],v[5])
			if err != nil {
				return err
			}
			fmt.Printf("The curve is %0.4f*x**2 + %0.4f*x + %0.4f\n", curve.A, curve.B, curve.C)
			start, stop := v[0], v[4]
			step := (stop-start)/200.0
			curve.Draw(start, stop, step, args[6])
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return nil
		},
	}

	return cmd
}

