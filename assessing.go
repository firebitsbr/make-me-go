package makeMeGo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// UTF8Extentions ... All matching extentions (eg ".txt", ".html") will be written as UTF8 and not bytes.
var UTF8Extentions = []string{
	".txt",
	".md",
	".markdown",
	".html",
	".css",
	".js",
	".go",
	".rtf",
}

// Generate scans the given 'assetsFolder' (recursively) and writes Go code into
// the 'codeFolder' as a slice of Go byte arrays with the file and collection both
// named by the 'collectionName' parameter.
// The 'codeFolder' path will be created if it does not already exist.
// The collectionName and resultant Go filename will have spaces replaced with
// underscores.
func Generate(assetsFolder, codeFolder, collectionName string) {
	fmt.Println("Generating", collectionName, "Go code.")

	// Normalise the namings.
	collectionName = strings.TrimSpace(strings.Replace(collectionName, " ", "_", -1))
	codeFilename := strings.ToLower(collectionName)

	// Gather details of all asset files.
	// The prefix is a small kludge related to path separators.
	prefix := path.Dir(path.Join(assetsFolder, "dummy"))
	matchingFiles := make(map[string][]byte)
	walkFn := func(currentPath string, info os.FileInfo, err error) error {
		// Process a single file.
		if info.IsDir() {
			return err
		}
		content, err := ioutil.ReadFile(currentPath)
		shortPath := path.Clean(currentPath[len(prefix)+1:])
		matchingFiles[shortPath] = content
		return err
	}
	filepath.Walk(assetsFolder, walkFn)

	// Create an index of the files to guarantee the order of the collection.
	// Without this, each call to go generate may result in the same content
	// appearing in a different sequence, causing spurious code commits.
	// In essence, Go specifically does NOT guarantee slice order.
	var idx []string
	for k := range matchingFiles {
		idx = append(idx, k)
	}
	sort.Strings(idx)

	// Create the dstination path (if not existing) and the output Go file.
	os.MkdirAll(codeFolder, os.ModePerm)
	out, _ := os.Create(path.Join(codeFolder, codeFilename+".go"))
	defer out.Close()

	// Write a summary at the top of the Go file.
	out.WriteString("package generated \n\n")
	out.WriteString("// " + collectionName + " ... Created via 'go generate' - manual edits will be lost.\n\n")
	for _, idx := range idx {
		out.WriteString("// Includes: " + idx + "\n")
	}
	out.WriteString("\n")

	// Write out all the assets in this collection as byte arrays.
	// Go always uses *nix path separators internally.
	validExts := make(map[string]bool)
	for _, v := range UTF8Extentions {
		validExts[v] = true
	}
	out.WriteString("var " + collectionName + " = map[string][]byte {\n")
	for _, filename := range idx {
		content := matchingFiles[filename]
		e := strings.ToLower(path.Ext(filename))
		_, ok := validExts[e]
		if ok {
			out.WriteString("  \"" + strings.Replace(filename, "\\", "/", -1) + "\": []byte(")
			writeByteArrayAsGoCodeUtf8(out, content)
			out.WriteString("),\n")
		} else {
			out.WriteString("  \"" + strings.Replace(filename, "\\", "/", -1) + "\": []byte{")
			writeByteArrayAsGoCode(out, content)
			out.WriteString("\n  },\n")
		}
	}
	out.WriteString("}\n")
}
