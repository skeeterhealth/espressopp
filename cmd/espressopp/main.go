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
	interpreter := &Espressopp{}
	codeGenerator := &SqlCodeGenerator{}
	interpreter.Accept(codeGenerator, r, w)
	fmt.Println(w.String())
}
