// Package espressopp provides primitives for parsing Espresso++ expressions and
// convert them into native queries.
//
// For example, here below is the code to convert an Espresso++ expression into
// native SQL:
//
//  package main
//
//  import (
//      "bytes"
//      "fmt"
//      "strings"
//
//      "gitlab.com/skeeterhealth/espressopp"
//  )
//
//  func main() {
//      r := strings.NewReader("age gte 30")
//      w := new(bytes.Buffer)
//      interpreter := &Espressopp{}
//      codeGenerator := &SqlCodeGenerator{}
//      interpreter.Accept(codeGenerator, r, w)
//      fmt.Println(w.String()) // prints "age >= 30"
//  }
package espressopp