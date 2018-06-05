package commands

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/client"
)

// Commands registers a sub-tree of commands to interact with
// local project storage.
func Commands(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create or update a particular project",
		Long:  `Project allows you to manage your local project store.`,
	}

	cmd.AddCommand(
		client.PostCommands(
			createProjectCommand(cdc),
			updateProjectCommand(cdc),
		)...)

	return cmd
}
