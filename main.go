package main

import (
	"fmt"

	"github.com/coding-for-fun-org/go-playground/pkg/dictionary"
)

func main() {
	d := dictionary.Dictionary{}
	err := d.Add("hello", "greeting")
	if err != nil {
		fmt.Println(err)
	}

	d.Add("what", "asking for information")
	d.Add("world", "planet earth")

	helloDefinition, errHello := d.Search("hello")
	if err != nil {
		fmt.Println(errHello)
	} else {
		fmt.Println(helloDefinition)
	}

	errHelloUpdate := d.Update("hello", "greeting2")
	if errHelloUpdate != nil {
		fmt.Println(errHelloUpdate)
	}

	newHelloDefinition, errNewHello := d.Search("hello")
	if errNewHello != nil {
		fmt.Println(errNewHello)
	} else {
		fmt.Println(newHelloDefinition)
	}

	whatDefinition, errWhat := d.Search("what")
	if errWhat != nil {
		fmt.Println(errWhat)
	} else {
		fmt.Println(whatDefinition)
	}

	d.Delete("what")
	if _, errWhat := d.Search("what"); errWhat != nil {
		fmt.Println("successfully deleted")
	}

	world2Definition, errWorld2 := d.Search("world2")
	if errWorld2 != nil {
		fmt.Println(errWorld2)
	} else {
		fmt.Println(world2Definition)
	}
}
