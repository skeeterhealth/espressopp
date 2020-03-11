/**
 * @begin 2020-02-27
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

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
//      interpreter := espressopp.NewEspressoppInterpreter()
//      codeGenerator := espressopp.NewSqlCodeGenerator()
//      err := interpreter.Accept(codeGenerator, r, w)
//
//      if err != nil {
//          buf := new(bytes.Buffer)
//          buf.ReadFrom(r)
//          msg := fmt.Errorf("Error generating sql from %v: %v", buf.String(), err)
//          fmt.Println(msg)
//      } else {
//          fmt.Println(w.String())
//      }
// }
package espressopp
