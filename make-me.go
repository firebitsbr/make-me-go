package makeMeGo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Generate scans the given 'assetsFolder' (recursively) and writes Go code into
// the 'codeFolder' as a slice of Go byte arrays with the file and collection both
// named by the 'collectionName' parameter.
// The 'codeFolder' path will be created if it does not already exist.
func Generate(assetsFolder, codeFolder, collectionName string) {
	fmt.Println("Generating", collectionName, "Go code.")
	codeFilename := strings.ToLower(collectionName)
	prefix := path.Dir(path.Join(assetsFolder, "dummy"))
	matchingFiles := make(map[string][]byte)
	walkFn := func(currentPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		content, err := ioutil.ReadFile(currentPath)
		shortPath := path.Clean(currentPath[len(prefix)+1:])
		matchingFiles[shortPath] = content
		return err
	}
	filepath.Walk(assetsFolder, walkFn)

	var idx []string
	for k := range matchingFiles {
		idx = append(idx, k)
	}
	sort.Strings(idx)

	os.MkdirAll(codeFolder, os.ModePerm)
	out, _ := os.Create(path.Join(codeFolder, codeFilename+".go"))
	defer out.Close()

	out.WriteString("package generated \n\n")
	out.WriteString("// " + collectionName + " ... Created via 'go generate' - manual edits will be lost.\n\n")

	for _, idx := range idx {
		out.WriteString("// Includes: " + idx + "\n")
	}
	out.WriteString("\n")

	out.WriteString("var " + collectionName + " = map[string][]byte {\n")
	for _, filename := range idx {
		content := matchingFiles[filename]
		out.WriteString("  \"" + strings.Replace(filename, "\\", "/", -1) + "\": []byte{")
		writeByteArrayAsGoCode(out, content)
		out.WriteString("\n  },\n")
	}
	out.WriteString("}\n")
}

// WriteAssets recreates the assets in the collection (from the generated
// Go code) into the outputFolder, respecting the original structure.
func WriteAssets(outputFolder string, collection map[string][]byte) error {
	for filename, contents := range collection {
		pathname := path.Join(outputFolder, filename)
		err := os.MkdirAll(path.Dir(pathname), os.ModePerm)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(pathname, contents, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeByteArrayAsGoCode ... Sends the bytes as text in lines of 40.
func writeByteArrayAsGoCode(out *os.File, content []byte) {
	for i, c := range content {
		if i%40 == 0 {
			out.WriteString("\n    ")
		}
		out.WriteString(strconv.Itoa(int(c)) + ",")
	}
}
