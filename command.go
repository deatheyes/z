package main

import (
	"github.com/urfave/cli"
)

var GetCommand = cli.Command {
	Name: "get",
	Usage: "get [<Path>]",
	Description: "get path value",
	Action: GetAction,
}

var LsCommand = cli.Command {
	Name: "ls",
	Usage: "ls [<Path>]",
	Description: "list dir",
	Action: LsAction,
}

var SetCommand = cli.Command {
	Name: "set",
	Usage: "set [<Path>] <Value>",
	Description: "set path value",
	Action: SetAction,
}

var CdCommand = cli.Command {
	Name: "cd",
	Usage: "cd <Path>",
	Description: "channge dir",
	Action: CdAction,
}

var CreateCommand = cli.Command {
	Name: "create",
	Usage: "create [<Path>] <Value>",
	Description: "create path",
	Action: CreateAction,
}

var DelCommand = cli.Command {
	Name: "del",
	Usage: "del <Path>",
	Description: "delete path",
	Action: DelAction,
}
