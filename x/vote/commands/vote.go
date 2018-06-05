package commands

import (
	"dbachain/common/log"
	"dbachain/common/utils"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dbachain/x/vote"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	flagProjectID = "project-id"
	flagCoin      = "coin"
	flagRound     = "round"
)

func voteTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Vote or view vote info",
		Long:  `Vote allows you to vote or query vote info.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper()

			// get the from address
			from, err := utils.GetUserFromCoreContext(ctx)
			if err != nil {
				return err
			}

			projectID := viper.GetString(flagProjectID)
			if "" == projectID {
				return errors.New("project-id is empty")
			}

			coin := viper.GetString(flagCoin)
			if "" == coin {
				return errors.New("coin is empty")
			}

			round := viper.GetInt64(flagRound)
			if round <= 0 {
				return errors.New("round must be > 0")
			}

			// build message
			msg := vote.NewVoteMsg(from, projectID, coin, round, time.Now().Unix())

			// build and sign the transaction, then broadcast to Tendermint
			res, err := utils.SignBuildBroadcast(ctx, msg, cdc)
			if err != nil {
				return err
			}

			log.Infof("Committed at block %d. Hash: %s,data:%s\n", res.Height, res.Hash.String(), res.DeliverTx.GetData())
			return nil
		},
	}

	cmd.Flags().String(flagProjectID, "", "Project to vote")
	cmd.Flags().String(flagCoin, "", "Amount of coins to vote")
	cmd.Flags().Int64(flagRound, 0, "the round of vote")
	return cmd
}
