package utils

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/tendermint/go-crypto/keys/words"
	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client/core"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/dbachain/dbachain/x/user"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	dbm "github.com/tendermint/tmlibs/db"
)

const (
	FlagType        = "type"
	FlagRecover     = "recover"
	FlagNoBackup    = "no-backup"
	FlagDryRun      = "dry-run"
	FlagUserStatus  = "status"
	FlagNationality = "nationality"
	FlagOnchain     = "onchain"
)

const AccountDBName = "accounts"

var (
	// keybase is used to make GetUserBase a singleton
	userbase user.Userbase
)

func GetUserBase() (user.Userbase, error) {
	if userbase == nil {
		rootDir := viper.GetString(cli.HomeFlag)
		db, err := dbm.NewGoLevelDB(AccountDBName, filepath.Join(rootDir, "accounts"))
		if err != nil {
			return nil, err
		}
		userbase = user.New(
			db,
			words.MustLoadCodec("english"),
		)
	}
	return userbase, nil
}

type addOutput struct {
	Key  user.AccountInfo `json:"key"`
	Seed string           `json:"seed"`
}

func PrintAccountInfo(info user.AccountInfo) {
	switch viper.Get(cli.OutputFlag) {
	case "text":
		addr := info.Info.PubKey.Address().String()
		sep := "\t\t"
		if len(info.Info.Name) > 7 {
			sep = "\t"
		}

		fmt.Printf("%s%s%s\n", info.Info.Name, sep, addr)

	case "json":
		json, err := user.MarshalJSON(info.Info)
		if err != nil {
			panic(err) // really shouldn't happen...
		}
		fmt.Println(string(json))
	}
}

func PrintAccountInfos(infos []user.AccountInfo) {
	switch viper.Get(cli.OutputFlag) {
	case "text":
		fmt.Println("All accounts:")
		for _, i := range infos {
			PrintAccountInfo(i)
		}
	case "json":
		json, err := user.MarshalJSON(infos)
		if err != nil {
			panic(err) // really shouldn't happen...
		}
		fmt.Println(string(json))
	}
}

// Get the from address from the name flag
func GetUserFromCoreContext(ctx core.CoreContext) (from sdk.Address, err error) {

	userbase, err := GetUserBase()
	if err != nil {
		return nil, err
	}

	name := ctx.FromAddressName
	if name == "" {
		return nil, errors.Errorf("must provide a from address name")
	}

	info, err := userbase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", name)
	}

	return info.Info.PubKey.Address(), nil
}

// sign and build the transaction from the msg
func SignBuildBroadcast(ctx core.CoreContext, msg sdk.Msg, cdc *wire.Codec) (*ctypes.ResultBroadcastTxCommit, error) {
	name := ctx.FromAddressName
	passphrase, err := ctx.GetPassphraseFromStdin(name)
	if err != nil {
		return nil, err
	}

	txBytes, err := signAndBuild(ctx, name, passphrase, msg, cdc)
	if err != nil {
		return nil, err
	}

	return ctx.BroadcastTx(txBytes)
}

// sign and build the transaction from the msg
func signAndBuild(ctx core.CoreContext, name, passphrase string, msg sdk.Msg, cdc *wire.Codec) ([]byte, error) {

	// build the Sign Messsage from the Standard Message
	chainID := ctx.ChainID
	sequence := ctx.Sequence
	signMsg := sdk.StdSignMsg{
		ChainID:   chainID,
		Sequences: []int64{sequence},
		Msg:       msg,
	}

	userbase, err := GetUserBase()
	if err != nil {

		return nil, err
	}

	// sign and build
	bz := signMsg.Bytes()

	sig, pubkey, err := userbase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []sdk.StdSignature{{
		PubKey:    pubkey,
		Signature: sig,
		Sequence:  sequence,
	}}

	// marshal bytes
	tx := sdk.NewStdTx(signMsg.Msg, signMsg.Fee, sigs)

	return cdc.MarshalBinary(tx)
}

func PrintCreate(info user.AccountInfo, seed string) {
	output := viper.Get(cli.OutputFlag)
	switch output {
	case "text":
		PrintAccountInfo(info)
		// print seed unless requested not to.
		if !viper.GetBool(FlagNoBackup) {
			fmt.Println("**Important** write this seed phrase in a safe place.")
			fmt.Println("It is the only way to recover your user if you ever forget your password.")
			fmt.Println()
			fmt.Println(seed)
		}
	case "json":
		out := addOutput{Key: info}
		if !viper.GetBool(FlagNoBackup) {
			out.Seed = seed
		}
		json, err := user.MarshalJSON(out)
		if err != nil {
			panic(err) // really shouldn't happen...
		}
		fmt.Println(string(json))
	default:
		panic(fmt.Sprintf("I can't speak: %s", output))
	}
}
