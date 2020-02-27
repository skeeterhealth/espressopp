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

This distro consists of both a library to be used by any client application and a
command-line utility to help developers debug filters written in the Espresso++ language.

Further information about Espresso++ is available in the following documents:
* [The Espresso++ Language Specification](docs/espressopp-spec.adoc)
* [Espresso++: Software Design Document](docs/espressopp-sdd.adoc)

<a id="markdown-building-instructions" name="building-instructions"></a>
## Building Instructions

The `Makefile` that comes with the project supports the following targets:

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

`espressopp` is a command-line utility that takes an Espresso++ expression as an
input and returns the resulting native query. To build and install `espressopp`
into `$GOPATH/bin`, issue the following command -- ensure you added `GOPATH/bin`
to your `PATH`:

```sh
$ make install
```

Finally, to generate the HTML and PDF documentation in the `out/docs` sub-directory,
issue the following command -- ensure `asciidoctor`, `asciidoctor-pdf`, and
`asciidoctor-diagrams` are installed on your system:

```sh
$ make docs
```

<a id="markdown-getting-started" name="getting-started"></a>
## Getting Started

Currently Espresso++ only suports SQL, but the way a filter is create for data
management systems with other query languages does not change.

---

*Copyright 2020 Skeeter Health*
