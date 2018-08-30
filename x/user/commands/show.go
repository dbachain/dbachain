package commands

import (
	"github.com/dbachain/dbachain/common/utils"
	"github.com/dbachain/dbachain/x/user"
	"errors"

	"github.com/spf13/cobra"
)

var showUserCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show user info for the given name",
	Long:  `Return public details of one account.`,
	RunE:  runShowCmd,
}

func getKey(name string) (user.AccountInfo, error) {
	ab, err := utils.GetUserBase()
	if err != nil {
		return user.AccountInfo{}, err
	}

	return ab.Get(name)
}

func runShowCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a name for the key")
	}
	name := args[0]

	info, err := getKey(name)
	if err == nil {
		utils.PrintAccountInfo(info)
	}
	return err
}
