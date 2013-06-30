// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"compress/gzip"
	"fmt"
	"io"
)

type Input struct {
    name string
    rd io.ReadCloser
}

// translate translates the input file to go source code.
func translate(inputs []Input, releaseOutput, debugOutput io.Writer, pkgname string) {
	translate_nomemcpy_comp(inputs, releaseOutput, debugOutput, pkgname)
}

// input -> gzip -> gowriter -> output.
func translate_nomemcpy_comp(inputs []Input, releaseOutput, debugOutput io.Writer, pkgname string) {
	fmt.Fprintf(debugOutput, `// +build !release

package %s

import (
	"io/ioutil"
)

func ReadAsset(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}
`, pkgname)
	fmt.Fprintf(releaseOutput, `// +build release

package %s

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

var assets = map[string] string { 
`, pkgname)

	for _, in := range inputs {
		fmt.Fprintf(releaseOutput, "\n\t\"%s\": \"", in.name)
		gz := gzip.NewWriter(&StringWriter{Writer: releaseOutput})
		io.Copy(gz, in.rd)
		gz.Close()
		in.rd.Close()
		fmt.Fprintf(releaseOutput, `",`)

	}

	fmt.Fprintf(releaseOutput, `
}

func ReadAsset(name string) ([]byte, error) {
	contents, ok := assets[name]
	if !ok {
	    return nil, fmt.Errorf("Asset %%s not compiled into binary", name)
	}
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&contents))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(contents)
	bx.Cap = bx.Len

	gz, err := gzip.NewReader(bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	io.Copy(&buf, gz)
	gz.Close()

	return buf.Bytes(), nil
}`)

}
