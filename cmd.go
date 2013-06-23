package main

import (
	"flag"
	"fmt"
	"os"
)

type CmdFn func(*Cmd, []string) int

type Cmd struct {
	Name      string
	Desc      string
	UsageLine string
	Flags     *flag.FlagSet
	fn        CmdFn
}

func NewCmd(name, desc, usage string, fn CmdFn) *Cmd {
	cmd := &Cmd{
		Name:      name,
		Desc:      desc,
		UsageLine: usage,
		Flags:     flag.NewFlagSet(name, flag.ExitOnError),
		fn:        fn,
	}

	return cmd
}

func (c *Cmd) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
	os.Exit(1)
}

func (c *Cmd) Run(args []string) int {
	return c.fn(c, args)
}
