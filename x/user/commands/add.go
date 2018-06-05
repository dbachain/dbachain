package commands

import (
	"dbachain/common/log"
	"dbachain/x/user"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	keys "github.com/tendermint/go-crypto/keys"

	"dbachain/common/utils"
)

func addUserCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a new user account",
		Long:  `Add a user to the user store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var name, pass string

			buf := client.BufferStdin()

			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide a name for the user")
			}
			name = args[0]
			ab, err := utils.GetUserBase()
			if err != nil {
				return err
			}

			_, err = ab.Get(name)
			if err == nil {
				// account exists, ask for user confirmation
				if response, err := client.GetConfirmation(
					fmt.Sprintf("override the existing name %s", name), buf); err != nil || !response {
					return err
				}
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

			pass, err = client.GetCheckPassword(
				"Enter a passphrase for your key:",
				"Repeat the passphrase:", buf)
			if err != nil {
				return err
			}

			if viper.GetBool(utils.FlagRecover) {
				seed, err := client.GetSeed(
					"Enter your recovery seed phrase:", buf)
				if err != nil {
					return err
				}
				info, err := ab.Recover(name, pass, seed, info)
				if err != nil {
					return err
				}

				// print out results without the seed phrase
				viper.Set(utils.FlagNoBackup, true)
				utils.PrintCreate(info, "")
			} else {
				algo := keys.CryptoAlgo(viper.GetString(utils.FlagType))
				// TO DO : add more detail information in local db
				info, seed, err := ab.Create(name, pass, algo, info)
				if err != nil {
					return err
				}
				utils.PrintCreate(info, seed)
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
				msg := user.NewUserMsg("add", from, useraddress, name, userstatus)

				// build and sign the transaction, then broadcast to Tendermint

				res, err := utils.SignBuildBroadcast(ctx, msg, cdc)
				if err != nil {
					log.Infof("Failed to commit at current block\n")
					if ab.Delete(name, pass) == nil {
						log.Infof("Local db was clean, done\n")
					} else {
						log.Infof("Local db failed to delete\n")
					}

					return err
				}
				log.Infof("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())

			}

			return nil
		},
	}
	cmd.Flags().StringP(utils.FlagType, "t", "ed25519", "Type of private key (ed25519|secp256k1|ledger)")
	cmd.Flags().Bool(utils.FlagRecover, false, "Provide seed phrase to recover existing key instead of creating")
	cmd.Flags().Bool(utils.FlagNoBackup, false, "Don't print out seed phrase (if others are watching the terminal)")
	cmd.Flags().Bool(utils.FlagDryRun, false, "Perform action, but don't add key to local keystore")
	cmd.Flags().Int(utils.FlagUserStatus, -1, "User status from the command line parameter.")
	cmd.Flags().String(utils.FlagNationality, "", "User nationality from the command line parameter.")
	cmd.Flags().Bool(utils.FlagOnchain, false, "whether broadcast the tx from the command line parameter.")

	return cmd
}
