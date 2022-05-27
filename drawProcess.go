package main

import (
	"fmt"
)

type DrawProcess struct {
	maxProcess int64
	ProSignle  chan int64
	currVal    int64
}

func Init(MaxProcess int64) DrawProcess {
	return DrawProcess{
		maxProcess: MaxProcess,
		ProSignle:  make(chan int64, 10),
		currVal:    0,
	}
}

func (d *DrawProcess) Draw() {
	fmt.Println(d.maxProcess)
	SendLens := (d.maxProcess / 1024) / 100
	alerSenlen := 0
	for {
		select {
		case c := <-d.ProSignle:
			if c-int64(alerSenlen) < SendLens {

				alerSenlen += 1

				d.currVal += c
				continue
			}
			alerSenlen += 1

			d.currVal += c

			process := 0

			for i := int64(0); i < int64(alerSenlen); i += SendLens {
				fmt.Printf("*")
				process += 1
			}
			if process > 100 {
				process = 100
			}
			fmt.Printf("%d%%\n", process)

			if process == 100 {
				return
			}

		default:
			continue
		}

	}

}
