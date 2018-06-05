package commands

import (
	"dbachain/common/log"
	"dbachain/common/utils"
	"errors"

	"dbachain/x/project"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func updateProjectCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Update an existing project",
		Long:  `Update an existing project such as status to the project store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper()

			// get the from address from the name flag
			from, err := utils.GetUserFromCoreContext(ctx)
			if err != nil {
				return err
			}

			projectId := viper.GetString(flagProjectID)
			if projectId == "" {
				return errors.New("project-id is empty")
			}

			status := viper.GetInt(flagProjectStatus)
			if status != project.StatusInit && status != project.StatusNormal && status != project.StatusClosed {
				return errors.New("project status is undefined")
			}

			msg := project.NewProjectMsg("update", from, projectId, status, time.Now().Unix())

			// build and sign the transaction, then broadcast to Tendermint
			res, err := utils.SignBuildBroadcast(ctx, msg, cdc)
			if err != nil {
				return err
			}

			log.Infof("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())

			return nil
		},
	}

	cmd.Flags().String(flagProjectID, "", "Project id from the command line parameter.")
	cmd.Flags().String(flagProjectStatus, "", "Project status from the command line parameter.")

	return cmd
}
