package commands

import (
	"dbachain/common/log"
	"dbachain/common/utils"
	"dbachain/x/user"
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func updateUserCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Update user account info",
		Long:  `Update user in the user store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide a name for the user")
			}
			name := args[0]
			ab, err := utils.GetUserBase()
			if err != nil {
				return err
			}

			userstatus := viper.GetInt(utils.FlagUserStatus)
			if userstatus != user.UserStatusNormal && userstatus != user.UserStatusIllegal {
				return errors.New("user status is undefined")
			}

			nationality := viper.GetString(utils.FlagNationality)
			if nationality == "" {
				return errors.New("user nationality is empty")
			}

			info := user.AccountInfo{
				UserStatus:  userstatus,
				Nationality: nationality,
			}

			buf := client.BufferStdin()
			pass, err := client.GetPassword(
				"Enter the current passphrase:", buf)
			if err != nil {
				return err
			}

			err = ab.Update(name, pass, info)
			if err != nil {
				log.Infof("Failed to update local db\n")
				return err
			}

			onchain := viper.GetBool(utils.FlagOnchain)
			if onchain && !viper.GetBool(utils.FlagRecover) {
				ctx := context.NewCoreContextFromViper()

				// get the from address from the name flag
				from, err := utils.GetUserFromCoreContext(ctx)
				if err != nil {
					return err
				}

				useraddress := ab.GetUserAddress(name)
				msg := user.NewUserMsg("update", from, useraddress, name, userstatus)

				// build and sign the transaction, then broadcast to Tendermint

				res, err := utils.SignBuildBroadcast(ctx, msg, cdc)
				if err != nil {
					log.Infof("Failed to commit at current block\n")
					return err
				}
				log.Infof("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			}
			return nil
		},
	}

	cmd.Flags().StringP(utils.FlagType, "t", "ed25519", "Type of private key (ed25519|secp256k1|ledger)")
	cmd.Flags().Bool(utils.FlagDryRun, false, "Perform action, but don't add key to local keystore")
	cmd.Flags().Int(utils.FlagUserStatus, -1, "User status from the command line parameter.")
	cmd.Flags().String(utils.FlagNationality, "", "User nationality from the command line parameter.")
	cmd.Flags().Bool(utils.FlagOnchain, false, "whether broadcast the tx from the command line parameter.")
	return cmd
}
