package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/dappledger/ann-contracts/evm/client/commands"
)

func main() {
	tool := cli.NewApp()

	tool.Name = "evm"

	tool.Commands = []cli.Command{
		commands.EVMCommands,
	}

	_ = tool.Run(os.Args)
}
