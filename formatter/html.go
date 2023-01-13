package formatter

import (
	"fmt"
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/yosssi/gohtml"
)

/*
	Converts a string to html
	Args:
		`data`: The document body to inject.
		`css_file`. If `css_file` is anything other than an empty string then we will attempt to read the file
								 and inject the contents into the style tag
*/
func GenerateHtml(data, css_file string) (string, error) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	md := []byte(data)
	body := string(markdown.ToHTML(md, parser, nil))

	style := `
  	td, tr, th, table {
			border: 1px solid black;
			border-collapse: collapse;
		}
		td, th {
			padding-left: 10px;
			padding-right: 10px;
			padding-top: 5px;
			padding-bottom: 5px;
		}
		th {
			background-color: #d6d5d2;
		}
		body {
			margin: 50px;
		}
		a {
			text-decoration: none;
		}
    `

	if css_file != "" {
		buf, err := ioutil.ReadFile(css_file)

		if err != nil {
			return "", err
		}
		style = string(buf)

	}
	html_template := `
	<html>
		<head>
			<style>
				%s
			</style>
		</head>
		<body>
			%s
		</body>
	</html>	
	`
	html := fmt.Sprintf(html_template, style, body)

	html = gohtml.Format(html)
	return html, nil
}
