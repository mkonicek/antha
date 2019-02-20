package align

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io"
	"strconv"
	"strings"
)

func makeAlignmentFile(results Result, algorithmName string) (anthafile wtype.File, err error) {
	htmlHeader := `<html>
<body>
<link rel="stylesheet" type="text/css" href="https://fonts.googleapis.com/css?family=Roboto+Mono:400,700|Roboto:100,400,300,300italic,400italic,500,500italic,700,700italic" crossorigin="anonymous">
<style>
.amber{color:rgb(255,195,16)}
.blue{color:rgb(16,192,215)}
.green{color:rgb(110,182,22)}
.red{color:red}
.justify{align:justify}
.bold{font-weight: bold}
.courier{font-family: Courier New, Courier, monospace}
body{font-family:Roboto}
</style>	
	`
	htmlFooter := `
</body>
</html>`

	htmlSpace := `<span style='mso-tab-count:1'>&nbsp; </span>`
	htmlNewLine := `<br /> `
	header := func(s string) string {
		return `<header> <h1>` + s + `</h1> </header>`
	}
	justify := func(s string) string {
		return `<span class="justify">` + s + `</span>`
	}
	bold := func(s string) string {
		return `<span class="bold">` + s + `</span>`
	}
	htmlParagraph := func(s string) string {
		return `<div>` + justify(s) + `</div>`
	}
	courier := func(s string) string {
		return `<span class="courier">` + s + `</span>`
	}
	writeParagraph := func(buf io.Writer, s ...interface{}) (int, error) {
		return fmt.Fprintln(buf, htmlParagraph(fmt.Sprint(s...)))
	}

	red := func(s string) string {
		return `<span class="red">` + s + `</span>`
	}
	green := func(s string) string {
		return `<span class="green">` + s + `</span>`
	}
	blue := func(s string) string {
		return `<span class="blue">` + s + `</span>`
	}
	amber := func(s string) string {
		return `<span class="amber">` + s + `</span>`
	}

	// colourFormat colour codes all mismatches (lower case bp) in a formatted alignment sequence
	colourFormat := func(alignment string) (formatted string) {
		var newstring []string
		for _, letter := range alignment {
			if string(letter) == "a" {
				newstring = append(newstring, red(string(letter)))
			} else if string(letter) == "t" {
				newstring = append(newstring, amber(string(letter)))
			} else if string(letter) == "c" {
				newstring = append(newstring, blue(string(letter)))
			} else if string(letter) == "g" {
				newstring = append(newstring, green(string(letter)))
			} else {
				newstring = append(newstring, (string(letter)))
			}
		}

		return strings.Join(newstring, "")
	}

	// tabSpacer creates a tab of the correct size based on the size of the number
	tabSpacer := func(number int) string {
		var tab string
		if number < 10 {
			tab = htmlSpace + htmlSpace + htmlSpace + htmlSpace + htmlSpace
		} else if number < 100 {
			tab = htmlSpace + htmlSpace + htmlSpace + htmlSpace
		} else if number < 1000 {
			tab = htmlSpace + htmlSpace + htmlSpace
		} else if number < 10000 {
			tab = htmlSpace + htmlSpace
		} else {
			tab = htmlSpace
		}
		return tab
	}

	var buf bytes.Buffer
	writeParagraph(&buf, htmlHeader)
	writeParagraph(&buf, header("Antha Sequencing  Alignment"))
	writeParagraph(&buf, "Template: ", results.Template.Name())
	writeParagraph(&buf, "Query: ", results.Query.Name())
	writeParagraph(&buf, "Algorithm: ", algorithmName)
	writeParagraph(&buf, htmlNewLine)
	writeParagraph(&buf, "Matches: ", results.Matches())
	writeParagraph(&buf, "Mismatches: ", results.Mismatches())
	writeParagraph(&buf, "Gaps: ", results.Gaps())

	writeParagraph(&buf, fmt.Sprintf("Positions matched of query in template: %+v", results.Positions().Positions))
	writeParagraph(&buf, htmlNewLine)
	writeParagraph(&buf, bold(amber("line1 ")), htmlSpace, results.Template.Name())
	writeParagraph(&buf, bold(blue("line2 ")), htmlSpace, results.Query.Name())
	writeParagraph(&buf, htmlNewLine)

	for i := 0; i < len(results.Alignment.TemplateResult); i = i + 60 {
		var endLine int

		if i+60 > len(results.Alignment.TemplateResult) {
			endLine = len(results.Alignment.TemplateResult)
		} else {
			endLine = i + 60
		}

		templatePosition := results.Alignment.TemplatePositions[i]
		queryPosition := results.Alignment.QueryPositions[i]

		writeParagraph(&buf, bold(amber(strconv.Itoa(results.Alignment.TemplatePositions[i]))), tabSpacer(templatePosition), courier(colourFormat(results.Alignment.TemplateResult[i:endLine])))
		writeParagraph(&buf, bold(blue(strconv.Itoa(results.Alignment.QueryPositions[i]))), tabSpacer(queryPosition), courier(colourFormat(results.Alignment.QueryResult[i:endLine])))
		writeParagraph(&buf, htmlNewLine)

	}

	writeParagraph(&buf, htmlFooter)

	anthafile.Name = "Alignment_" + results.Template.Name() + "_" + results.Query.Name() + ".html"

	err = anthafile.WriteAll(buf.Bytes())

	if err != nil {
		return
	}

	return
}
