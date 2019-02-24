package makemego

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
	".csv",
	".markdown",
	".html",
	".css",
	".less",
	".sass",
	".stylus",
	".js",
	".go",
	".sh",
	".bat",
	".rtf",
	".xml",
	".xhtml",
}

// Generate scans the given 'assetsFolder' (recursively) and writes Go code into
// the 'codeFolder' as a slice of Go byte arrays with the file and collection both
// named by the 'collectionName' parameter in the 'packageName' package (eg 'main').
// The 'codeFolder' path will be created if it does not already exist.
// The names will have spaces replaced with underscores.
func Generate(
	assetsFolder,
	codeFolder,
	packageName,
	collectionName,
	hint string,
) {
	fmt.Println("Generating", packageName+"."+collectionName, "Go code.")

	// Normalise the namings.
	packageName = strings.TrimSpace(strings.Replace(packageName, " ", "_", -1))
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
		if info.Name() == ".DS_Store" {
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

	// Create the destination path (if not existing) and the output Go file.
	os.MkdirAll(codeFolder, os.ModePerm)
	out, _ := os.Create(path.Join(codeFolder, codeFilename+".go"))
	defer out.Close()

	// Write a summary at the top of the Go file.
	out.WriteString("package " + packageName + " \n\n")
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
	out.WriteString("// " + collectionName + " ... " + hint + "\n")
	out.WriteString("var " + collectionName + " = map[string][]byte{\n")
	for i, filename := range idx {
		if i > 0 {
			out.WriteString("\n")
		}
		content := matchingFiles[filename]
		e := strings.ToLower(path.Ext(filename))
		_, ok := validExts[e]
		if ok {
			out.WriteString("\t\"" + strings.Replace(filename, "\\", "/", -1) + "\": []byte(")
			writeByteArrayAsGoCodeUtf8(out, content)
			out.WriteString("),\n")
		} else {
			out.WriteString("\t\"" + strings.Replace(filename, "\\", "/", -1) + "\": []byte{")
			writeByteArrayAsGoCode(out, content)
			out.WriteString("\n\t},\n")
		}
	}
	out.WriteString("}\n")
}
