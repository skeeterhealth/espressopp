# Espresso++

<!-- TOC -->
1. [Introduction](#introduction)
2. [Building Instructions](#building-instructions)
3. [Getting Started](#getting-started)
<!-- /TOC -->

<a id="markdown-introduction" name="introduction"></a>
## Introduction

**Espresso++** is a natural query language that abstracts away the native query language
of any data management system.

This distro consists of both a library to be used by any client application that
needs to support filtering and a command-line utility to help developers debug
filters written in the Espresso++ language.

Further information about Espresso++ is available in the following documents:
* [The Espresso++ Language Specification](docs/espressopp-spec.adoc)
* [Espresso++: Software Design Document](docs/espressopp-sdd.adoc)

<a id="markdown-building-instructions" name="building-instructions"></a>
## Building Instructions

The `Makefile` that comes with this project supports the following targets:

* `all`:            Runs in sequence `check`, `test`, and `install`
* `get-tools`:      Retrieves and installs `goimports`
* `init`:           Creates `go.mod` if missing
* `build`:          Builds the `espressopp` utility in the `out` sub-directory
* `clean`:          Cleans the last build
* `test`:           Runs unit tests
* `install`:        Installs the `espressopp` utility into `$GOPATH/bin`
* `uninstall`:      Uninstalls the `espressopp` utility from `$GOPATH/bin`
* `docs`:           Runs in sequence `docs-html` and `docs-pdf`
* `docs-html`:      Generates HTML documentation in the `out/html` sub-directory
* `docs-pdf`:       Generates PDF documentation in the `out/pdf` sub-directory
* `fmt`:            Formats source code according to the Go best practices
* `simplify`:       Simplifies source code according to the Go best practices
* `check`:          Checks whether source code is well formatted

The Espesso++ library gets compiled together with the client application, therefore
this build procedure only applies to `espressopp`, a command-line utility that
takes an Espresso++ expression as an input and returns the resulting native query.
To build and install `espressopp` into `$GOPATH/bin`, issue the following command
&ndash; ensure you added `GOPATH/bin` to your `PATH`:

```sh
$ make install
```

Finally, to generate the HTML and PDF documentation in the `out/docs` sub-directory,
issue the following command  &ndash; ensure `asciidoctor`, `asciidoctor-pdf`, and
`asciidoctor-diagrams` are installed on your system:

```sh
$ make docs
```

To generate HTML only:

```sh
$ make docs-html
```

To generate PDF only:

```sh
$ make docs-pdf
```

<a id="markdown-getting-started" name="getting-started"></a>
## Getting Started

Currently Espresso++ suports just SQL, but the way a filter is create for data
management systems with other query languages does not change. Below is an example
of how to get an Espresso++ expression translated into SQL:

```go
package main

import (
	"bytes"
	"fmt"
	"strings"

	"gitlab.com/skeeterhealth/espressopp"
)

func main() {
	r := strings.NewReader("age gte 30")
	w := new(bytes.Buffer)
	interpreter := espressopp.NewEspressoppInterpreter()
	codeGenerator := espressopp.NewSqlCodeGenerator()
	err := interpreter.Accept(codeGenerator, r, w)

	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		msg := fmt.Errorf("Error generating sql from %v: %v", buf.String(), err)
		fmt.Println(msg)
	} else {
		fmt.Println(w.String())
	}
}
```

If Espresso++ supported MongoDB, the client code would be something like this: 

```go
package main

import (
	"bytes"
	"fmt"
	"strings"

	"gitlab.com/skeeterhealth/espressopp"
)

func main() {
	r := strings.NewReader("age gte 30")
	w := new(bytes.Buffer)
	interpreter := espressopp.NewEspressoppInterpreter()
	codeGenerator := espressopp.NewMongoCodeGenerator()
	err := interpreter.Accept(codeGenerator, r, w)

	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		msg := fmt.Errorf("Error generating sql from %v: %v", buf.String(), err)
		fmt.Println(msg)
	} else {
		fmt.Println(w.String())
	}
}
```

Last but not least, developers can debug their Espresso++ expressions with the
`espressopp`command-line utility:

```sh
$ espressopp --help

A utility that converts input Espresso++ expressions into native queries.

Usage:

espressopp --target <LANGUAGE> "<EXPRESSION>"

Supported targets:
1. sql
```

For example, let's translate the Espresso++ expression `age gte 30` into SQL:

```sh
$ espressopp --target sql "age gte 30"

age >= 30
```

Finally, the same Espresso++ expression translated into MongoDB query language:

 ```sh
$ espressopp --target mongo "age gte 30"

{ age: { $gte: 30 } }
 ```

---

*Copyright 2020 Skeeter Health*
