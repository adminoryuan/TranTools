package main

import (
	"sync"
)

func main() {
	tools := Trantools{}

	tools.Start()

	a := sync.WaitGroup{}
	a.Add(1)
	a.Wait()

}
