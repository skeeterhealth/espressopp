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
		Target     string `arg name:"target" help:"Target query language." name:"target"`
		Expression string `arg name:"expression" help:"Source expression." name:"expression"`
	} `cmd help:"Generate target native query."`
}

// emitSql renders SQL from e.
func emitSql(e string) (error, string) {
	r := strings.NewReader(e)
	w := new(bytes.Buffer)

	interpreter := espressopp.NewEspressoppInterpreter()
	codeGenerator := espressopp.NewSqlCodeGenerator()
	if err := interpreter.Accept(codeGenerator, r, w); err != nil {
		return err, ""
	}

	return nil, w.String()
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

	if ctx.Command() == "generate <target> <expression>" {
		if strings.ToLower(cli.Generate.Target) == "sql" {
			err, sql := emitSql(cli.Generate.Expression)
			if err != nil {
				msg := fmt.Errorf("Error generating sql from %v: %v", cli.Generate.Expression, err)
				fmt.Println(msg)
			} else {
				fmt.Println(sql)
			}
		} else {
			msg := fmt.Errorf("Target '%v' not supported.", cli.Generate.Target)
			fmt.Println(msg)
		}
	}
}
