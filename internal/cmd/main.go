package main

import (
	"flag"
	helpers "github.com/point-c/generator-helpers"
	"github.com/point-c/wgevents/internal/pkg/app"
	"github.com/point-c/wgevents/internal/pkg/templates"
)

var args = struct {
	config      string
	outfile     string
	testOutfile string
	pkg         string
}{
	config:      "events.yml",
	outfile:     helpers.OutputFilename(helpers.EnvGoFile()),
	testOutfile: helpers.TestOutputFilename(helpers.OutputFilename(helpers.EnvGoFile())),
	pkg:         helpers.EnvGoPackage(),
}

func init() {
	flag.StringVar(&args.config, "config", args.config, "events config file")
	flag.StringVar(&args.outfile, "out", args.outfile, "output filename")
	flag.StringVar(&args.testOutfile, "tout", args.testOutfile, "test output filename")
	flag.StringVar(&args.pkg, "pkg", args.pkg, "gopackage of output")
}

func main() {
	flag.Parse()
	// Read config
	cfg := helpers.Must(helpers.UnmarshalYAML[app.Config](args.config))
	// Convert to template context
	dot := app.Cfg2Dot(&cfg, args.pkg)
	// Generate files
	templates.Generate(dot, templates.GetTemplate(false), args.outfile)
	templates.Generate(dot, templates.GetTemplate(true), args.testOutfile)
}
