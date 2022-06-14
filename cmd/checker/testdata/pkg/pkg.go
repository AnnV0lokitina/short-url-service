package main

import (
	"fmt"
	"os"
)

func main1() {
	fmt.Println("test")
	x := 0
	fmt.Println(x)
	os.Exit(x)
}

func main() {
	fmt.Println("test")
	x := 0
	fmt.Println(x)
	// os.Exit(x)
	main1()
}
