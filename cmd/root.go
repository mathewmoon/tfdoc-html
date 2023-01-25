/*
Copyright © 2023 Mathew Moon <me@mathewmoon.net>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mathewmoon/tfdoc-html/formatter"
	"github.com/mathewmoon/tfdoc-html/writers"

	"github.com/spf13/cobra"
	"github.com/terraform-docs/terraform-docs/format"
	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

type config struct {
	includeOutputs bool   // Will attempt to generate outputs from Terraform and include in the html
	markdownOnly   bool   // Don't generate HTML, just markdown
	s3Uri          string // Fully qualified S3 URI
	cssFile        string // Path to a file of CSS to inject into the <head> of the html
	noStdout       bool   // Don't print to stdout
	file           string // File path to write output to
	sourcePath     string // Directory containing your Terraform
	header         string // An extra header to add to HTML docs
}

/*
Validate inputs and return a `Config`
*/
func parseConfig(cmd *cobra.Command, args []string) config {
	include_outputs, err := cmd.Flags().GetBool("outputs")
	if err != nil {
		exitWithError(err, 1)
	}

	markdown_only, err := cmd.Flags().GetBool("markdown")
	if err != nil {
		exitWithError(err, 1)
	}

	s3_uri, err := cmd.Flags().GetString("s3-uri")
	if err != nil {
		exitWithError(err, 1)
	}

	css_file, err := cmd.Flags().GetString("css-file")
	if err != nil {
		exitWithError(err, 1)
	}

	no_stdout, err := cmd.Flags().GetBool("no-stdout")
	if err != nil {
		exitWithError(err, 1)
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		exitWithError(err, 1)
	}

	header, err := cmd.Flags().GetString("header")
	if err != nil {
		exitWithError(err, 1)
	}

	if len(args) == 0 {
		exitWithError(errors.New("must provide [PATH] as first argument"), 1)
	}

	return config{
		includeOutputs: include_outputs,
		markdownOnly:   markdown_only,
		s3Uri:          s3_uri,
		cssFile:        css_file,
		noStdout:       no_stdout,
		file:           file,
		sourcePath:     args[0],
		header:         header,
	}
}

var rootCmd = &cobra.Command{
	Use:   "tfdoc-html [PATH]",
	Short: "Generate Terraform Docs in HTML, optionally uploading to S3 or writing to a file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		config := parseConfig(cmd, args)

		tfConfig := print.DefaultConfig()
		tfConfig.OutputValues.Enabled = config.includeOutputs
		tfConfig.ModuleRoot = config.sourcePath

		module, err := terraform.LoadWithOptions(tfConfig)

		if err != nil {
			exitWithError(err, 2)
		}

		table := format.NewMarkdownTable(tfConfig)

		if err := table.Generate(module); err != nil {
			exitWithError(err, 1)
		}

		doc := table.Content()

		if !config.markdownOnly {
			doc, err = formatter.GenerateHtml(doc, config.cssFile, config.header)
			if err != nil {
				exitWithError(err, 1)
			}
		}

		if config.file != "" {
			if err := writers.WriteToFile(config.file, doc); err != nil {
				exitWithError(err, 2)
			}
		}

		if config.s3Uri != "" {
			_, err := writers.S3Upload(config.s3Uri, doc)
			if err != nil {
				exitWithError(err, 2)
			}
		}

		if !config.noStdout {
			fmt.Println(doc)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		exitWithError(err, 3)
	}
}

func exitWithError(msg error, code int) {
	fmt.Printf("Error: %s ", msg)
	os.Exit(code)
}

func init() {
	rootCmd.Flags().BoolP("outputs", "o", false, "Inject outputs from state file on the fly. This requires having access to the state file declared in your backend config.")
	rootCmd.Flags().StringP("file", "f", "", "Write output to file.")
	rootCmd.Flags().Bool("no-stdout", false, "Don't write to stdout.")
	rootCmd.Flags().BoolP("markdown", "m", false, "Output MarkDown instead of HTML.")
	rootCmd.Flags().StringP("s3-uri", "s", "", "A full S3 uri that, if provided, the generated output will be uploaded to")
	rootCmd.Flags().StringP("css-file", "C", "", "A file containing CSS that will be used to override the default styling")
	rootCmd.Flags().StringP("header", "H", "", "A string to add as a header to html docs")

}
