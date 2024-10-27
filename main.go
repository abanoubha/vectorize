package main

import (
	"fmt"
	"os"
)

func main(){
	// imgbuff, err := os.ReadFile("test_cases/nursery-cover.png")
	_, err := os.ReadFile("test_cases/nursery-cover.png")
	if err != nil {
		fmt.Println("can not read file, or file not found")
	}

	fmt.Println("image loaded successfully")
}
