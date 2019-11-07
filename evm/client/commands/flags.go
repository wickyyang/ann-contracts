package commands

import "gopkg.in/urfave/cli.v1"

type ToolFlags struct {
	callf cli.Flag
	nonce cli.Flag
	abif  cli.Flag
}

var toolFlags = ToolFlags{
	callf: cli.StringFlag{
		Name:  "callf",
		Usage: "contract parameters, such as address, bytecode, method and parameters .",
	},
	nonce: cli.StringFlag{
		Name:  "nonce",
		Usage: "account's nonce",
	},
	abif: cli.StringFlag{
		Name:  "abif",
		Usage: "abi file",
	},
}
