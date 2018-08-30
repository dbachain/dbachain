package main

import (
	"os"
	"path/filepath"

	"github.com/dbachain/dbachain/app"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/server"
)

// rootCmd is the entry point for this binary
var (
	context = server.NewDefaultContext()
	rootCmd = &cobra.Command{
		Use:               "dbachaind",
		Short:             "DBAChaind Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(context),
	}
)

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dataDir := filepath.Join(rootDir, "data")
	dbMain, err := dbm.NewGoLevelDB("dbachain", dataDir)
	if err != nil {
		return nil, err
	}
	dbAcc, err := dbm.NewGoLevelDB("dbachain-acc", dataDir)
	if err != nil {
		return nil, err
	}
	dbIBC, err := dbm.NewGoLevelDB("dbachain-ibc", dataDir)
	if err != nil {
		return nil, err
	}
	dbStaking, err := dbm.NewGoLevelDB("dbachain-staking", dataDir)
	if err != nil {
		return nil, err
	}
	dbVote, err := dbm.NewGoLevelDB("dbachain-vote", dataDir)
	if err != nil {
		return nil, err
	}
	dbProject, err := dbm.NewGoLevelDB("dbachain-project", dataDir)
	if err != nil {
		return nil, err
	}

	dbs := map[string]dbm.DB{
		"main":    dbMain,
		"acc":     dbAcc,
		"ibc":     dbIBC,
		"staking": dbStaking,
		"vote":    dbVote,
		"project": dbProject,
	}
	bapp := app.NewDbaApp(logger, dbs)
	return bapp, nil
}

func main() {
	server.AddCommands(rootCmd, server.DefaultGenAppState, generateApp, context)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.dbachaind")
	executor := cli.PrepareBaseCmd(rootCmd, "DBAC", rootDir)
	executor.Execute()
}
