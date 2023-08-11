package main

import (
	"fmt"

	"github.com/ggc1412/hello-go/myDict"
)

func main() {
	dictionary := myDict.Dictionary{}
	baseWord := "hello"
	dictionary.Add(baseWord, "First")
	dictionary.Search(baseWord)
	err := dictionary.Delete(baseWord)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(baseWord)
	}
	err1 := dictionary.Delete(baseWord)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(baseWord)
	}
}
