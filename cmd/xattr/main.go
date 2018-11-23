package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/falzm/go-xattr"
)

var (
	flagHelp bool
)

func init() {
	flag.BoolVar(&flagHelp, "h", false, "display this help and exit")
	flag.Parse()

	if flagHelp {
		printUsage(os.Stdout)
	}

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "missing arguments")
		printUsage(os.Stdout)
	}
}

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "list", "ls":
		list(args)

	case "get":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: missing arguments", cmd)
			printUsage(os.Stdout)
		}
		get(args[0], args[1:])

	case "set":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "%s: missing arguments", cmd)
			printUsage(os.Stdout)
		}
		set(args[0], args[1], args[2:])

	case "remove", "rm":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s: missing arguments", cmd)
			printUsage(os.Stdout)
		}
		remove(args[0], args[1:])

	case "clear":
		clear(args)

	default:
		fmt.Fprintf(os.Stderr, "unsupported command %q", cmd)
		printUsage(os.Stdout)
	}
}

func dieOnError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("error: %s\n", format), a...)
	os.Exit(1)
}

func list(files []string) {
	if len(files) == 1 {
		attrs, err := xattr.List(files[0])
		if err != nil {
			dieOnError("unable to list extended attributes: %s", err)
		}

		for _, attr := range attrs {
			fmt.Println(attr)
		}

		return
	}

	for _, file := range files {
		attrs, err := xattr.List(file)
		if err != nil {
			dieOnError("unable to list extended attributes: %s", err)
		}

		fmt.Printf("%s:\n", file)
		for _, attr := range attrs {
			fmt.Println("  ", attr)
		}
	}
}

func get(attr string, files []string) {
	if len(files) == 1 {
		a, err := xattr.Get(files[0], attr)
		if err != nil {
			if xattr.IsNotExist(err) {
				return
			}

			dieOnError("unable to get extended attribute %q: %s", attr, err)
		}

		fmt.Println(string(a))

		return
	}

	for _, file := range files {
		a, err := xattr.Get(file, attr)
		if err != nil {
			if !xattr.IsNotExist(err) {
				dieOnError("unable to get extended attribute %q: %s", attr, err)
			}

			a = []byte{}
		}

		fmt.Printf("%s: %s\n", file, a)
	}
}

func set(attr, value string, files []string) {
	for _, file := range files {
		if err := xattr.Set(file, attr, []byte(value)); err != nil {
			dieOnError("unable to set extended attribute %q: %s", attr, err)
		}
	}

	fmt.Println("OK")
}

func remove(attr string, files []string) {
	for _, file := range files {
		if err := xattr.Remove(file, attr); err != nil {
			if xattr.IsNotExist(err) {
				continue
			}

			dieOnError("unable to remove extended attribute %q: %s", attr, err)
		}
	}

	fmt.Println("OK")
}

func clear(files []string) {
	for _, file := range files {
		attrs, err := xattr.List(file)
		if err != nil {
			dieOnError("unable to clear extended attributes: %s", err)
		}

		for _, attr := range attrs {
			if err := xattr.Remove(file, attr); err != nil {
				if xattr.IsNotExist(err) {
					continue
				}

				dieOnError("unable to clear extended attribute: %s", err)
			}
		}
	}

	fmt.Println("OK")
}

func printUsage(output io.Writer) {
	fmt.Fprintf(output, "Usage: %s [options] <command> [arg] file [file...]", path.Base(os.Args[0]))

	fmt.Fprintf(output, "\n\nOptions:\n")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(output, "  -%s  %s (default: %q)\n", f.Name, f.Usage, f.DefValue)
	})

	fmt.Fprintf(output, `
Commands:
  clear             clear files(s) all extended attributes
  get <a>           get files(s) extended attribute <a> value
  list|ls           list file(s) extended attributes
  remove|rm <a>	    remove file(s) extended attribute <a>
  set <a> <v>       set files(s) extended attribute <a> value to <v>
`)

	os.Exit(2)
}
