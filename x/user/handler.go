package user

import (
	"fmt"
	"reflect"

	"dbachain/common/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k UserKeeper, am sdk.AccountMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case UserMsg:
			return handleUserMsg(ctx, k, msg, am)
		default:
			errMsg := "Unrecognized project Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleUserMsg(ctx sdk.Context, k UserKeeper, msg UserMsg, am sdk.AccountMapper) sdk.Result {
	msgType := msg.MsgType
	if msgType == "add" {
		return saveUser(ctx, k, msg, am)
	} else if msgType == "update" {
		return updateUser(ctx, k, msg, am)
	} else {
		return sdk.Result{Code: sdk.CodeInternal, Log: fmt.Sprintf("illegal argument UserMsg.MsgType:%s", msgType)}
	}

}

func saveUser(ctx sdk.Context, k UserKeeper, msg UserMsg, am sdk.AccountMapper) sdk.Result {
	existed := k.IsUserExisted(ctx, msg, am)
	if existed {
		log.Errorf("User already existed, user id: %s", msg.ID)
		return ErrUserAlreadyExisted("Project already existed.").Result()
	}

	// TO DO
	// if ctx.IsCheckTx() {
	// 	return sdk.Result{}
	// }

	err := k.saveUser(ctx, msg, am)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func updateUser(ctx sdk.Context, k UserKeeper, msg UserMsg, am sdk.AccountMapper) sdk.Result {
	existed := k.IsUserExisted(ctx, msg, am)
	if !existed {
		log.Errorf("User not exist, user id: %s", msg.ID)
		return ErrUserNotExist("User not exist.").Result()
	}

	// TODO
	// if ctx.IsCheckTx() {
	// 	return sdk.Result{}
	// }

	err := k.updateUser(ctx, msg, am)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
