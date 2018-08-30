package user

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/dbachain/dbachain/common/log"
	"github.com/dbachain/dbachain/types"
)

const (
	UserStatusNormal = iota
	UserStatusIllegal
)

type UserKeeper struct {
	cdc      *wire.Codec
	storeKey sdk.StoreKey
}

func NewUserKeeper(cdc *wire.Codec, key sdk.StoreKey) UserKeeper {
	return UserKeeper{
		cdc:      cdc,
		storeKey: key,
	}
}

func (k UserKeeper) saveUser(ctx sdk.Context, msg UserMsg, am sdk.AccountMapper) sdk.Error {

	baseAcc := auth.NewBaseAccountWithAddress(msg.Address)

	acc := types.AppAccount{
		BaseAccount: baseAcc,
		Name:        msg.ID,
		UserStatus:  msg.UserStatus,
	}

	am.SetAccount(ctx, &acc)

	return nil
}

func (k UserKeeper) updateUser(ctx sdk.Context, msg UserMsg, am sdk.AccountMapper) sdk.Error {

	baseAcc := auth.NewBaseAccountWithAddress(msg.Address)

	acc := types.AppAccount{
		BaseAccount: baseAcc,
		Name:        msg.ID,
		UserStatus:  msg.UserStatus,
	}

	lastAcc, ok := am.GetAccount(ctx, msg.Address).(*types.AppAccount)

	if ok == false {
		log.Errorf("AppAccount type cast failed")
		return nil
	}
	if appAccountCompare(acc, *lastAcc) {
		log.Warnf("No new changes in account info, skip update operations on chain...\n")
		return nil
	}
	am.SetAccount(ctx, &acc)

	return nil
}

func (k UserKeeper) IsUserExisted(ctx sdk.Context, msg UserMsg, am sdk.AccountMapper) bool {

	if am.GetAccount(ctx, msg.Address) == nil {
		return false
	}

	return true
}

// Currently we compare the following fields in app account:
// AppAccount.Name
// AppAccount.UserStatus
func appAccountCompare(a, b types.AppAccount) bool {
	return a.Name == b.Name &&
		a.UserStatus == b.UserStatus
}
