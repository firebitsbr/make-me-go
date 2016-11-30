package main

//go:generate go run embed/embed-assets.go

import (
	"fmt"

	makeMeGo "github.com/kcartlidge/make-me-go"
	generated "github.com/kcartlidge/make-me-go/example/generated"
)

func main() {
	fmt.Printf("You have %d asset files.\n", len(generated.Assets))
	err := makeMeGo.WriteAssets("recreated-assets", generated.Assets)
	if err != nil {
		fmt.Println(err)
	}
}
