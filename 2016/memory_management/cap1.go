package main

import "fmt"

func main() {
	// START OMIT
	a := make([]int, 10)
	a = append(a, 1)
	fmt.Println("cap a:", cap(a))
	// END OMIT
}
