package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command, cdc *wire.Codec) {
	cmdSub := voteTxCmd(cdc)
	cmdSub.AddCommand(
		client.GetCommands(voteStatusTxCmd(cdc))...,
	)

	cmd.AddCommand(client.PostCommands(cmdSub)...)
}
