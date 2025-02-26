package main

import (
	"bytes"
	"fmt"
	"io"
)

/*
common built in interfaces:
io.Writer - just has a Write() method. Applies to bytes.Buffer, os.Stdout, and files
io.Reader - just has Read() method.
io.Closer - just has Close()
io.ReadWriter - has Read() and Write()
io.ReadWriteCloser - hash Read(), Write() and Close()
fmt.Stringer - just has String() method.
fmt.Value - has String() and Set() methods.
sort.Interface - has Len(), Less() and Swap() methods.
builtin.error - just has Error() method.
*/

type dummyWriter struct {
	origWriter io.Writer
	count      int64
}

func (d *dummyWriter) Write(p []byte) (n int, e error) {
	d.count += int64(len(p))
	return d.origWriter.Write(p)
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	dw := dummyWriter{origWriter: w}
	return &dw, &dw.count
}

func main() {
	b := new(bytes.Buffer)
	// shows power of interfaces by wrapping a writer object into another interface so can do custom logic
	cb1, lenWritten1 := CountingWriter(b)
	cb2, lenWritten2 := CountingWriter(b)
	fmt.Fprintf(cb1, "hello %s", "world")
	fmt.Fprintf(cb2, "hello %s", "abc")
	fmt.Printf("lenWritter1: %v\n", *lenWritten1)
	fmt.Printf("lenWritter2: %v\n", *lenWritten2)
	fmt.Printf("buffer length: %v\n", b.Len())
}
