package makeMeGo

import (
	"os"
	"strconv"
	"strings"
)

const bytesPerLine = 40

// writeByteArrayAsGoCode writes the bytes as text numerics in chunks of 40.
func writeByteArrayAsGoCode(out *os.File, content []byte) {
	for i, c := range content {
		if i%bytesPerLine == 0 {
			out.WriteString("\n    ")
		}
		out.WriteString(strconv.Itoa(int(c)) + ",")
	}
}

// writeByteArrayAsGoCodeUtf8 writes the bytes as UTF8 text
func writeByteArrayAsGoCodeUtf8(out *os.File, content []byte) {
	out.WriteString("\n\"")
	s := string(content[:])
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	out.WriteString(s)
	out.WriteString("\"")
}
