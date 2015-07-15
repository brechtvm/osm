
osm
===

The osm package is used to parse OpenStreetMap files and manipulate the data

## Installation

```bash
$ go get github.com/vetinari/osm
```

## Usage

Yet another pbf to osm xml converter:
```Go
    package main

    import (
        "fmt"
        "os"
        "github.com/vetinari/osm"
        "github.com/vetinari/osm/pbf"
        "github.com/vetinari/osm/xml"
    )

    func main() {
        if len(os.Args) != 2 {
            fmt.Fprintf(os.Stderr, "%s: Usage: %s PBF_FILE\n", os.Args[0], os.Args[0])
            os.Exit(1)
        }

        fh, err := os.Open(os.Args[1])
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s: failed to open %s: %s\n", os.Args[1], err)
            os.Exit(1)
        }
        defer fh.Close()

        o, err := osm.New(pbf.Parser(fh))
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s: failed to parse %s: %s\n", os.Args[1], err)
            os.Exit(1)
        }

        xml.Dump(os.Stdout, o)
    }
```

## Documentation

http://godoc.org/github.com/vetinari/osm

## Help

pull requests welcome :-)

## To Do

- More documentation ;-)
- tests
- import "log" & debug level / debug
- don't write/dump deleted items unless explicitly requested
- download / parse from osm.org API
- merge 2 or more *osm.OSM into one

