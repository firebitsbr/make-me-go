package main

//go:generate go run embed/embed-assets.go

import (
	"fmt"

	makeMeGo "github.com/kcartlidge/make-me-go"
	generated "github.com/kcartlidge/make-me-go/example/generated"
)

func main() {
	// Go generate should have been run before this source is built.
	// The result (from embed/embed-assets.go) will be the contents
	// of the assets folder tree being gathered into a slice in the
	// generated/assets.go source file.
	fmt.Printf("You have %d asset files.\n", len(generated.Assets))

	// This recreates the full initial assets folder structure inside
	// the recreated-assets folder.
	err := makeMeGo.WriteAssets("recreated-assets", generated.Assets)
	if err != nil {
		fmt.Println(err)
	}
}
