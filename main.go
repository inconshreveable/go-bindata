// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"unicode"
        "path"
        "path/filepath"
)

const (
	AppName    = "bindata"
	AppVersion = "2.0.0"
)

var (
	in           []string // flag.Args()
        outDir       = flag.String("o", ".", "Optional path to the output directory.")
	pkgname      = flag.String("p", "assets", "Optional name of the package to generate.")
	version      = flag.Bool("v", false, "Display version information.")
)

func main() {
	parseArgs()

        inputs := make([]Input, 0)
        for _, dirpath := range in {
            err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
                if info.IsDir() {
                    return nil
                }

                if err != nil {
                    return nil
                }

                fs, err := os.Open(path)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "[e] %s\n", err)
                        os.Exit(1)
                }

                inputs = append(inputs, Input{ path, fs })

                return nil
            })

            if err != nil {
                fmt.Fprintf(os.Stderr, "[e] %s\n", err)
                os.Exit(1)
            }
        }

        rfd, err := os.Create(path.Join(*outDir, "assets_release.go"))
        if err != nil {
                fmt.Fprintf(os.Stderr, "[e] %s\n", err)
                return
        }

        defer rfd.Close()

        dfd, err := os.Create(path.Join(*outDir, "assets_debug.go"))
        if err != nil {
                fmt.Fprintf(os.Stderr, "[e] %s\n", err)
        }

        defer dfd.Close()

        translate(inputs, rfd, dfd, *pkgname)

        fmt.Fprintln(os.Stdout, "[i] Done.\n")
}

// parseArgs processes and verifies commandline arguments.
func parseArgs() {
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "%s v%s (Go runtime %s)\n",
			AppName, AppVersion, runtime.Version())
		os.Exit(0)
	}

	if len(*pkgname) == 0 {
		fmt.Fprintln(os.Stderr, "[w] No package name specified. Using 'assets'.\n")
		*pkgname = "assets"
	} else {
		if unicode.IsDigit(rune((*pkgname)[0])) {
			// Identifier can't start with a digit.
			*pkgname = "_" + *pkgname
		}
	}

        in = flag.Args()
}
