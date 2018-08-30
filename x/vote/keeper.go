package vote

import (
	"github.com/dbachain/dbachain/common/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type VoteKeeper struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewVoteKeeper(key sdk.StoreKey, cdc *wire.Codec) VoteKeeper {
	return VoteKeeper{
		key: key,
		cdc: cdc,
	}
}

func (vk VoteKeeper) Vote(ctx sdk.Context, msg VoteMsg, bank bank.CoinKeeper) (*AccumulatedProjectVote, sdk.Error) {
	log.Debugf("VoteKeeper vote")

	//TODO checkout user status
	coins, err2 := sdk.ParseCoins(msg.Coin)
	if nil != err2 {
		log.Error(err2.Error())
		return nil, sdk.ErrInternal(err2.Error())
	}

	_, err := bank.SubtractCoins(ctx, msg.UserAddr, coins)
	if nil != err {
		log.Error(err.Error())
		return nil, err
	}

	bin, err2 := vk.cdc.MarshalJSON(msg)
	if nil != err2 {
		log.Error(err2.Error())
		return nil, sdk.ErrInternal(err2.Error())
	}
	ctx.KVStore(vk.key).Set(GetUserVoteKey(msg.UserAddr, msg.ProjectID, msg.CTime), bin)

	projectVoteState := ctx.KVStore(vk.key).Get(GetAccumulatedProjectVoteKey(msg.ProjectID))
	var accm AccumulatedProjectVote
	if nil != projectVoteState {
		err := vk.cdc.UnmarshalJSON(projectVoteState, &accm)
		if nil != err {
			log.Error(err.Error())
			return nil, sdk.ErrInternal(err.Error())
		}
		accm.Ammount += coins[0].Amount //TODO
	} else {
		accm.Ammount = coins[0].Amount //TODO
		accm.ProjectID = msg.ProjectID
	}

	bin, err2 = vk.cdc.MarshalJSON(accm)
	if nil != err2 {
		log.Error(err2.Error())
		return nil, sdk.ErrInternal(err2.Error())
	}
	ctx.KVStore(vk.key).Set(GetAccumulatedProjectVoteKey(accm.ProjectID), bin)

	return &accm, nil
}
