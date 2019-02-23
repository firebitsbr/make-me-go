package main

import makemego "github.com/kcartlidge/make-me-go"

// This is the package that the parent main.go file calls into from it's
// go generate comment line. It MUST be a main package in it's own right.

func main() {
	// This does all assets (recursively) in one file.
	// You could split this into separate calls for scripts and styles
	// for example, thus generating multiple Go source files each with
	// their own slice of asset contents.
	makemego.Generate("./assets", "generated", "generated", "Assets", "Auto-generated assets code.")
}
