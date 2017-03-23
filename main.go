package main

import (
	"encoding/base32"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

var (
	encoding = flag.String("encoding", "", "Which base32 encoding to use - std, hex, or a 32-byte custom value")
	decode   = flag.Bool("d", false, "Decode instead of encode")
)

func main() {
	flag.Parse()

	var enc *base32.Encoding
	switch *encoding {
	case "", "std":
		enc = base32.StdEncoding
	case "hex":
		enc = base32.HexEncoding
	default:
		if len(*encoding) == 32 {
			enc = base32.NewEncoding(*encoding)
		} else {
			panic(fmt.Errorf("no encoding known for %q", *encoding))
		}
	}

	if isatty.IsTerminal(os.Stdin.Fd()) {
		panic(fmt.Errorf("trying to read from stdin, but stdin is a terminal"))
	}

	if *decode {
		if _, err := io.Copy(os.Stdout, base32.NewDecoder(enc, os.Stdin)); err != nil {
			panic(err)
		}
	} else {
		wr := base32.NewEncoder(enc, os.Stdout)
		if _, err := io.Copy(wr, os.Stdin); err != nil {
			panic(err)
		}
		if err := wr.Close(); err != nil {
			panic(err)
		}
	}

	fmt.Printf("\n")
}
