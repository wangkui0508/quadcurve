package main

import (
	"os"
	"github.com/wangkui0508/quadcurve/qc"
)

func main() {
	cmd := qc.DrawQuadCurveCmd()
	cmd.SetArgs(os.Args[2:])
	err := cmd.Execute()
	if err != nil {
		println(err.Error())
	}
}
