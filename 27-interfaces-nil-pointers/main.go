package main

import (
	"bytes"
	"io"
)

func f(out io.Writer) {
	// even though out's value is nil it's type is *bytes.Buffer in the itable. So out is not nil.
	if out != nil {
		out.Write([]byte("done\n")) // This program is suppose to crash here to show the problem with nil pointers.
	}
}

func main() {
	var buf *bytes.Buffer // to fix this issue change this to
	// var buf io.Writer
	// buf = new(bytes.Buffer)
	f(buf)
}
