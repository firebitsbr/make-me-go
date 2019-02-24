package makemego

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const bytesPerLine = 40
const charsPerLine = 130

// writeByteArrayAsGoCode writes the bytes as text numerics in chunks of 40.
func writeByteArrayAsGoCode(out *os.File, content []byte) {
	out.WriteString("\t")
	for i, c := range content {
		if i%bytesPerLine == 0 {
			out.WriteString("\n\t\t")
		}
		out.WriteString(strconv.Itoa(int(c)) + ", ")
	}
}

// writeByteArrayAsGoCodeUtf8 writes the bytes as UTF8 text
func writeByteArrayAsGoCodeUtf8(out *os.File, content []byte) {
	out.WriteString("\n")

	s := string(content[:])
	utf8lines := chop(s, charsPerLine)
	mx := len(utf8lines) - 1
	for i, line := range utf8lines {
		line = strings.Replace(line, "\\", "\\\\", -1)
		line = strings.Replace(line, "\n", "\\n", -1)
		line = strings.Replace(line, "\"", "\\\"", -1)

		out.WriteString("\t\t\"")
		out.WriteString(line)
		if i < mx {
			out.WriteString("\" + \n\t")
		} else {
			out.WriteString("\"")
		}
	}
}

func chop(original string, width int) []string {
	var s []string
	for len(original) > width {
		i := width
		for i >= width-utf8.UTFMax && !utf8.RuneStart(original[i]) {
			i--
		}
		s = append(s, original[:i])
		original = original[i:]
	}
	if len(original) > 0 {
		s = append(s, original)
	}
	return s
}
