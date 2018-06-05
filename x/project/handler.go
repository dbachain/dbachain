package project

import (
	"dbachain/common/log"

	"reflect"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
)

func NewHandler(k ProjectKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ProjectMsg:
			return handleProjectMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized project Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleProjectMsg(ctx sdk.Context, k ProjectKeeper, msg ProjectMsg) sdk.Result {
	msgType := msg.MsgType
	if msgType == "create" {
		return saveProject(ctx, k, msg)
	} else if msgType == "update" {
		return updateProject(ctx, k, msg)
	} else {
		return sdk.Result{Code: sdk.CodeInternal, Log: fmt.Sprintf("illegal argument ProjectMsg.MsgType:%s", msgType),}
	}

	return sdk.Result{Code: sdk.CodeOK, Log: fmt.Sprintf("handle `ProjectMsg` successfully. ProjectMsg:%s", msg.String()),}
}

func saveProject(ctx sdk.Context, k ProjectKeeper, msg ProjectMsg) sdk.Result {
	existed := k.IsProjectExisted(ctx, msg.ID)
	if existed {
		log.Errorf("Project already existed, project id: %s", msg.ID)
		return ErrProjectAlreadyExisted("Project already existed.").Result()
	}

	if ctx.IsCheckTx() {
		return sdk.Result{} // TODO
	}

	//p := &Project{msg.ID, msg.Status}
	err := k.saveProject(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func updateProject(ctx sdk.Context, k ProjectKeeper, msg ProjectMsg) sdk.Result {
	existed := k.IsProjectExisted(ctx, msg.ID)
	if !existed {
		log.Errorf("Project not exist, project id: %s", msg.ID)
		return ErrProjectNotExist("Project not exist.").Result()
	}

	status := k.GetProjectStatus(ctx, msg.ID)
	if msg.Status == status {
		log.Errorf("Project state is consistent, project id: %s; status: %s", msg.ID, msg.Status)
		return ErrUpdateProjectFailed("Project state is consistent.").Result()
	}

	if ctx.IsCheckTx() {
		return sdk.Result{} // TODO
	}

	//p := &Project{msg.ID, msg.Status}
	err := k.updateProject(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
