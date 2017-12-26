package main

import (
	"fmt"
	"os"
	"os/user"
	"errors"
	"strings"

	"github.com/urfave/cli"
	"github.com/GeertJohan/go.linenoise"
)

var cmds = []cli.Command{
	CdCommand,
	GetCommand,
	SetCommand,
	DelCommand,
	CreateCommand,
	LsCommand,
}

const (
	HISTORY_FILE = "/.z_history"
)

var cmdmap = map[string]cli.Command{}

func init() {
	for _, cmd := range cmds {
		cmdmap[cmd.Name] = cmd
	}
}

func showHelp() {
	fmt.Println("commands:")
	for _, cmd := range cmds {
		fmt.Printf("%-3s%-14s  -  %-50s\n", "   ", cmd.Name, cmd.Usage)
	}
}

func validateHost(hosts []string) error {
	if len(hosts) == 0 {
		return errors.New("no host specified")
	}

	for _, host := range hosts {
		seg := strings.Split(host, ":")
		if len(seg) != 2 {
			return fmt.Errorf("invalidated host: %s", host)
		}
	}
	return nil
}

func main() {
	if (len(os.Args) == 2 && (string(os.Args[1]) == "-h" || string(os.Args[1]) == "--help") || (len(os.Args) == 1)) {
		fmt.Println("Usage: z <zkHost0>,<zkHost1> [<Command>] [options], -h for more details")
		os.Exit(1)
	}

	hosts := strings.Split(os.Args[1], ",")
	if err := validateHost(hosts); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	Meta.Hosts = hosts
	if err := Meta.Init(); err != nil {
		fmt.Println("meta init failed")
		os.Exit(-1)
	}

	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	historyFile := user.HomeDir + HISTORY_FILE
	_, err = os.Stat(historyFile)
	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(historyFile)
		if err != nil {
			fmt.Println("historty: " + historyFile + " create failed")
		}
	}

	if len(os.Args) == 2 {
		if err := linenoise.LoadHistory(historyFile); err != nil {
			fmt.Println(err)
		}

		for {
			prefix := "ZK:" + Meta.Path + "> "
			str, err := linenoise.Line(prefix)
			if err != nil {
				if err == linenoise.KillSignalError {
					os.Exit(1)
				}
				fmt.Printf("unexpected error: %s\n", err)
				os.Exit(1)
			}
			fields := strings.Fields(str)
			linenoise.AddHistory(str)
			if err := linenoise.SaveHistory(historyFile); err != nil {
				fmt.Println(err)
			}

			if len(fields) == 0 {
				continue
			}
			if fields[0] == "quit" {
				os.Exit(0)
			}

			cmd, ok := cmdmap[fields[0]]
			if !ok {
				showHelp()
				continue
			}
			app := cli.NewApp()
			app.Name = cmd.Name
			app.Commands = []cli.Command(cmds)
			app.Run(append(os.Args[:1], fields...))
		}
	}

	if (len(os.Args) > 2) {
		app := cli.NewApp()
		app.Name = "z"
		app.Usage = "sample client for zookeeper"
		app.Commands = cmds
		app.Run(append(os.Args[:1], os.Args[2:]...))
	}
}
