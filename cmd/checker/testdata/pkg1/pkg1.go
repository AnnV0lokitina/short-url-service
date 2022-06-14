package pkg1

import (
	"fmt"
	"os"
)

func main1() {
	fmt.Println("test")
	x := 0
	os.Exit(x)
}

func main() {
	fmt.Println("test")
	x := 0
	os.Exit(x)
	main1()
}
