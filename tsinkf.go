package main

import (
	"flag"
	"os"
	"strings"
	"text/template"
)

var Version string

var (
	showHelp = flag.Bool("h", false, "show this help")
	root     = flag.String("root", ".tsinkf", "directory where state files are created")
)


var cmdList = map[string]*Cmd{}

var usageTmpl = `usage: tsinkf [globals] command [arguments]

Globals:
  -root     storage directory ({{.Globals.root}})

Commands:{{range $k, $v := .Commands}}
  {{$k | printf "%-20s"}} {{$v.Desc}}{{end}}
`

func usage() {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(usageTmpl))
	data := struct {
		Commands map[string]*Cmd
		Globals  map[string]string
	}{
		cmdList,
		map[string]string{"root": *root},
	}

	if err := t.Execute(os.Stderr, data); err != nil {
		panic(err.Error())
	}

	os.Exit(1)
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) == 0 || *showHelp {
		flag.Usage()
		return
	}

	c, found := cmdList[args[0]]
	if !found {
		flag.Usage()
		return
	}

	os.Exit(c.Run(args[1:]))
}
