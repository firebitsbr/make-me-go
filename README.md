# make-me-go
Use ```go generate``` to automatically create Go source to recreate a folder
tree of assets.

## What does make-me-go do?

Imagine you have a Go executable and it depends on associated asset files
like web templates, sample documents or default config files.

Your usual options are:

1. Include the assets files in a deployment zip/package.
1. Write Go code that creates the asset files programmatically.

**Now you have a better option.**

Go includes the ```go generate``` *command*. You can add a ```generate``` *comment*
to the top of your ```main.go``` file and when you issue a ```go generate```
*command* it will find that ```generate``` *comment* and run it. This is independent
of the main build process and can function similar to grunt, gulp or equivalent
but *using native Go code*.

What **make-me-go** does is provide a package you can call during this generate
phase which will take a given assets folder and recursively gather all it's
contents. It will then convert those assets into hard-coded Go source which will
be built as part of your normal codebase. That Go source can be used to recreate
the original assets without you having to include them with your executable.

## It's easier than it sounds.

1. Drop all your assets into a folder structure.
1. Add a ```generate``` comment to your ```main.go``` file.
1. Add a package with a single file which is referred to by the step above and
which pulls in the **make-me-go** package and calls it.
1. Run ```go generate``` to create your asset's Go source code.
1. Your own Go code can now fetch the asset as a byte array from the newly-created
Go source code's slice (keyed by path) or you can ask **make-me-go** to recreate
the whole lot, mirroring the original folder path, in one call.

## There's an example that makes it clear.

In your terminal/command prompt navigate into the ```example``` folder and run
```go generate```. The contents of the ```generated``` subfolder will be updated
to reflect the contents of the ```assets``` folder.
You can then build and run your executable and it will update the folder containing
the ```recreated-assets```.

``` sh
cd example
go generate
go build -o example
./example
```

*If you are using Windows change ```example``` above to ```example.exe```.*

---

### About those unit tests ...

*There are no unit tests.*

I needed this in a rush for a one-off bit of processing and didn't expect
to be keeping it. However it proved so useful I'm sharing it anyway.
In my real-world usage it works exactly as expected and I probably won't be
revisiting it to add test coverage. Works for me. YMMV. Help yourself.
