package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands to interact with
// local account storage.
func Commands(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Add or view local user",
		Long:  `User allows you to manage your local user store.`,
	}
	cmd.AddCommand(
		client.PostCommands(
			addUserCommand(cdc),
			updateUserCommand(cdc),
			listUsersCommand(cdc),
		)...)

	return cmd
}
