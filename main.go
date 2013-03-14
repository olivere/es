package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	Version    = "0.1.1"
	DefaultUrl = "http://localhost:9200"
)

var verbose bool
var force bool
var flush bool
var refresh bool

type Command struct {
	Run  func(cmd *Command, args []string)
	Flag flag.FlagSet

	Usage string
	Short string
	Long  string

	ApiUrl string
}

func (c *Command) printUsage() {
	fmt.Printf("Usage: es %s\n\n", c.Usage)
	fmt.Println(strings.TrimSpace(c.Long))
}

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Running es on the command line will print these commands in order.
var commands = []*Command{
	cmdSearch,

	cmdIndices,
	cmdCreateIndex,
	cmdDeleteIndex,
	cmdSettings,
	cmdStatus,
	cmdStats,
	cmdRefresh,
	cmdOptimize,
	cmdFlush,
	cmdFlushDisable,
	cmdFlushEnable,

	cmdAliases,
	cmdIndexAliases,

	cmdMapping,
	cmdPutMapping,

	cmdTemplates,
	cmdTemplate,
	cmdCreateTemplate,
	cmdDeleteTemplate,

	cmdBulk,
	cmdReindex,

	cmdClusterHealth,
	cmdClusterState,
	cmdClusterNodes,

	cmdVersion,
	cmdApi,
	cmdHelp,
}

var (
	esUrl string
)

func main() {
	log.SetFlags(0)

	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	esUrl = DefaultUrl
	if s := os.Getenv("ES_URL"); s != "" {
		esUrl = strings.TrimRight(s, "/")
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = usage
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	usage()
}
