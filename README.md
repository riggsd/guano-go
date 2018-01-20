# GUANO-Go - GUANO bat acoustic metadata for Golang

GUANO is the Grand Unified Acoustic Notation Ontology, a metadata for bat
acoustic recordings. GUANO-Go is a package for reading GUANO metadata
from .WAV files using the Go (Golang) programming language.

**Status:** Experimental, don't use this in production yet!


## Installation

`$> go get github.com/riggsd/guano-go/guano`


## Usage

Import the package as `"github.com/riggsd/guano-go/guano"`.

The `guano.Guano` struct represents a file's GUANO metadata. Its `Fields`
member is a map of fieldname to value, where a fieldname includes the
optional namespace (eg. "Timestamp" or "GUANO|Version"), and its value is
the UTF-8 string representation (eg. "500" the string, *not* 500 the int).

Read the metadata from a .WAV file with `guano.Read(filename string)`, or
from any UTF-8 string with `guano.ParseGuanoString(s string)`.

```go
package main

import "github.com/riggsd/guano-go/guano"

func main() {
    filename := "test.wav"

    // read a named .WAV file
    gfile, err := guano.Read(filename)
    if err != nil {
        ...
    }

    // get a specific field
    version := gfile.Fields["GUANO|Version"]
    fmt.Printf("GUANO Version: %s", version)

    // iterate over all fields present
    for key, value := range gfile.Fields {
        fmt.Printf("%s:\t%s\n", key, value)
    }
}
```

## License

Copyright 2018, Myotisoft LLC

This project is Free / Open Source software under the MIT License.
