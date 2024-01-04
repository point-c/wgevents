package main

import (
	"flag"
	helpers "github.com/point-c/generator-helpers"
	"github.com/point-c/wgevents/internal/pkg/app"
	"github.com/point-c/wgevents/internal/pkg/templates"
	"log/slog"
	"os"
)

var args = struct {
	config  string
	outfile string
	pkg     string
}{
	config:  "events.yml",
	outfile: helpers.OutputFilename(helpers.EnvGoFile()),
	pkg:     helpers.EnvGoPackage(),
}

func init() {
	flag.StringVar(&args.config, "config", args.config, "events config file")
	flag.StringVar(&args.outfile, "out", args.outfile, "output filename")
	flag.StringVar(&args.pkg, "pkg", args.pkg, "gopackage of output")
	flag.Parse()
}

func main() {
	cfg := helpers.Must(helpers.UnmarshalYAML[app.Config](args.config))
	dot := Cfg2Dot(&cfg)
	b := ExecTemplate(dot)
	helpers.Check(os.WriteFile(args.outfile, b, os.ModePerm))
}

func ExecTemplate(dot *templates.Dot) []byte {
	b := helpers.Must(templates.Exec(dot))
	bf, err := helpers.GoFmt(b)
	if err != nil {
		slog.Error("failed to format code", "error", err)
		return b
	}
	return bf
}

func Cfg2Dot(cfg *app.Config) *templates.Dot {
	d := templates.Dot{
		Package: args.pkg,
		Imports: cfg.Imports,
	}

	for name, ev := range cfg.Events {
		ev := templates.Event{
			Name:   name,
			Type:   string(ev.Type),
			Level:  string(ev.Level),
			Nice:   helpers.IfStringEmptyUseDefault(ev.Nice, ev.Format),
			Format: ev.Format,
			Args: func(args []templates.Arg) []templates.Arg {
				for i, arg := range ev.Args {
					args[i].Name = arg.Name
					args[i].Type = arg.Type
				}
				return args
			}(make([]templates.Arg, len(ev.Args))),
			Custom: ev.Custom,
		}

		d.Events = append(d.Events, &ev)

		if ev.Custom {
			d.Custom = append(d.Custom, &ev)
		}

		switch app.LogType(ev.Type) {
		case app.LogTypeErrorf:
			d.Errorf = append(d.Errorf, &ev)
		case app.LogTypeVerbosef:
			d.Verbosef = append(d.Verbosef, &ev)
		}
	}
	return &d
}
