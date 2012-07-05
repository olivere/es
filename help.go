package main

import (
	"fmt"
	"log"
	"os"
)

var cmdVersion = &Command{
	Run:    runVersion,
	Usage:  "version",
	Short:  "show version",
	Long:   `Shows the version string.`,
	ApiUrl: "",
}

func runVersion(cmd *Command, args []string) {
	log.Println(Version)
}

var cmdHelp = &Command{
	Usage: "help [command]",
	Short: "show help",
	Long:  `Help shows usage for a command.`,
}

func init() {
	cmdHelp.Run = runHelp
}

func runHelp(cmd *Command, args []string) {
	if len(args) == 0 {
		printUsage()
		return
	}

	if len(args) != 1 {
		log.Fatal("too many arguments")
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.printUsage()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic: %s\n\n", args[0])
	os.Exit(2)
}

func printUsage() {
	fmt.Printf("Usage: es <command> [options] [arguments]\n\n")
	fmt.Printf("Supported commands are:\n\n")
	for _, cmd := range commands {
		if cmd.Short != "" {
			fmt.Printf("  %-16s   %s\n", cmd.Name(), cmd.Short)
		}
	}
	fmt.Println()
	fmt.Println("See 'es help [command]' for more information about a command.")
}

func usage() {
	printUsage()
	os.Exit(2)
}
