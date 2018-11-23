# xattr

Package xattr provides a simple interface to user extended attributes on Linux and OSX.

Install it: `go get github.com/falzm/go-xattr`

Documentation is available on [godoc.org](http://godoc.org/github.com/falzm/go-xattr).

## CLI

A convenience `xattr` command is available in the `cmd/xattr` sub-package. To install it:

```console
go get -u github.com/falzm/go-xattr/cmd/xattr
```

Usage:

```console
$ $GOPATH/bin/xattr -h
Usage: xattr [options] <command> [arg] file [file...]

Options:
  -h  display this help and exit (default: "false")

Commands:
  clear             clear files(s) all extended attributes
  get <a>           get files(s) extended attribute <a> value
  list|ls           list file(s) extended attributes
  remove|rm <a>	    remove file(s) extended attribute <a>
  set <a> <v>       set files(s) extended attribute <a> value to <v>
```

## License

Simplified BSD License (see LICENSE).
