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
//
//      interpreter := espressopp.NewEspressoppInterpreter()
//      codeGenerator := espressopp.NewSqlCodeGenerator()
//      codeGenerator.RenderingOptions.EnableNamedParams()
//      err := interpreter.Accept(codeGenerator, r, w)
//
//      if err != nil {
//          msg := fmt.Errorf("Error generating sql: %v", err)
//          fmt.Println(err)
//      } else {
//          fmt.Println(w.String())
//          namedParamValues, _ := codeGenerator.RenderingOptions.GetNamedParamValues()
//          for k, v := range namedParamValues {
//              fmt.Printf("%s: %s\n", k, v)
//          }
//      }
// }
package espressopp
