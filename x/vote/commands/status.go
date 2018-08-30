package commands

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"

	"github.com/dbachain/dbachain/x/vote"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/core"
	"github.com/spf13/viper"
)

func voteStatusTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query vote status",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectID := viper.GetString(flagProjectID)
			if "" == projectID {
				return errors.New("project-id is empty")
			}

			queryVoteStatus(cdc, context.NewCoreContextFromViper(), projectID)

			return nil
		},
	}

	cmd.Flags().String(flagProjectID, "", "ID of Project")

	return cmd
}

func queryVoteStatus(cdc *wire.Codec, ctx core.CoreContext, projID string) {
	projectID := vote.GetAccumulatedProjectVoteKey(projID)
	// get the node
	node, err := ctx.GetNode()
	if err != nil {
		fmt.Println(err)

		return
	}

	res, err := node.ABCIQuery("/vote/key", projectID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", res.Response.GetValue())
}
