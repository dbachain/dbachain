package commands

import (
	"github.com/dbachain/dbachain/common/utils"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

// CMD
func listUsersCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Long:  `Return a list of all users stored by local client.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ab, err := utils.GetUserBase()
			if err != nil {
				return err
			}

			infos, err := ab.List()
			if err == nil {
				utils.PrintAccountInfos(infos)
			}
			return err
		},
	}
	return cmd
}
