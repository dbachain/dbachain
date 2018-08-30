package vote

import (
	"reflect"

	project "github.com/dbachain/dbachain/x/project"

	"github.com/dbachain/dbachain/common/log"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// NewHandler returns a handler for "vote" type messages.
func NewHandler(vk VoteKeeper, bank bank.CoinKeeper, project project.ProjectKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case VoteMsg:
			return handleVoteMsg(ctx, vk, bank, project, msg)
		default:
			errMsg := "Unrecognized vote Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleVoteMsg(ctx sdk.Context, vk VoteKeeper, bank bank.CoinKeeper, pk project.ProjectKeeper, msg VoteMsg) sdk.Result {
	if project.StatusNormal != pk.GetProjectStatus(ctx, msg.ProjectID) {
		err := ErrCodeProjectStatusUnnormal("project status unnormal")
		log.Error(err.Error())
		return err.Result()
	}

	accm, err := vk.Vote(ctx, msg, bank)
	if nil != err {
		return err.Result()
	}

	bin, err2 := json.MarshalIndent(accm, "", " ")
	if nil != err2 {
		return sdk.ErrInternal(err2.Error()).Result()
	}

	// TODO: add some tags so we can search it!
	return sdk.Result{Data: bin}
}
