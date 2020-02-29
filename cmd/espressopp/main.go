/**
 * @begin 16-Feb-2020
 * @author <a href="mailto:giuseppe.greco@skeeterhealth.com">Giuseppe Greco</a>
 * @copyright 2020 <a href="skeeterhealth.com">Skeeter</a>
 */

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
	interpreter := &espressopp.Espressopp{}
	codeGenerator := &espressopp.SqlCodeGenerator{}
	err := interpreter.Accept(codeGenerator, r, w)

	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		fmt.Sprintf("Error generating sql from %s\n", buf.String())
	} else {
		fmt.Println(w.String())
	}
}
