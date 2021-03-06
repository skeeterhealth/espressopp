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

	"github.com/alecthomas/kong"
	"gitlab.com/skeeterhealth/espressopp"
)

var cli struct {
	Generate struct {
		Target            string            `arg name:"target" help:"Target query language." name:"target"`
		Expression        string            `arg name:"expression" help:"Source expression." name:"expression"`
		FieldMap          map[string]string `arg optional name:"fieldmap" help:"Mapping to native column names." type:"string:string"`
		EnableNamedParams bool              `help:"Enable named parameters." short:"e"`
	} `cmd help:"Generate target native query."`
}

// emitSql renders SQL from e applying m.
func emitSql(e string, m map[string]string, b bool) {
	r := strings.NewReader(e)
	w := new(bytes.Buffer)

	interpreter := espressopp.NewEspressoppInterpreter()
	codeGenerator := espressopp.NewSqlCodeGenerator()
	codeGenerator.RenderingOptions.FieldsWithDefault(m)

	if b {
		codeGenerator.RenderingOptions.EnableNamedParams()
	}

	if err := interpreter.Accept(codeGenerator, r, w); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(w.String())

	if b {
		values, _ := codeGenerator.RenderingOptions.GetNamedParamValues()
		if len(values) > 0 {
			fmt.Println()
			fmt.Println("Named Parameters")
			fmt.Println("================")
			for k, v := range values {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	}
}

// main is the program's entry point.
func main() {
	ctx := kong.Parse(&cli,
		kong.Name("espressopp"),
		kong.Description("A utility that converts input Espresso++ expressions into native queries."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	switch ctx.Command() {
	case "generate <target> <expression>", "generate <target> <expression> <fieldmap>":
		switch strings.ToLower(cli.Generate.Target) {
		case "sql":
			emitSql(cli.Generate.Expression, cli.Generate.FieldMap, cli.Generate.EnableNamedParams)
		default:
			fmt.Println(fmt.Errorf("Target '%v' not supported.", cli.Generate.Target))
		}
	}
}
