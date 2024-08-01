package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	message := "Hello, OTUS!"
	reverseString := reverse.String(message)
	fmt.Println(reverseString)
}
