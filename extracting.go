package makemego

import (
	"io/ioutil"
	"os"
	"path"
)

// WriteAssets recreates the assets in the collection (from the generated
// Go code) into the outputFolder, respecting the original structure.
func WriteAssets(outputFolder string, collection map[string][]byte) error {
	for filename, contents := range collection {
		pathname := path.Join(outputFolder, filename)

		// Ensure the full folder path exists then create the file inside it.
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
