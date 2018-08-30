package project

import (
	"github.com/dbachain/dbachain/common/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	StatusNormal = iota
	StatusInit
	StatusClosed
)

type ProjectKeeper struct {
	cdc      *wire.Codec
	storeKey sdk.StoreKey
}

func NewProjectKeeper(cdc *wire.Codec, key sdk.StoreKey) ProjectKeeper {
	return ProjectKeeper{
		cdc:      cdc,
		storeKey: key,
	}
}

//type Project struct {
//	id     string
//	status int
//}

func (k ProjectKeeper) saveProject(ctx sdk.Context, msg ProjectMsg) sdk.Error {
	bz, err := k.cdc.MarshalBinary(msg)
	if err != nil {
		log.Errorf("Save project failed, project id: %s, status: %s", msg.ID, msg.Status)
		return ErrSaveProjectFailed("Save project failed.")
	}

	k.setProject(ctx, []byte(msg.ID), bz)
	return nil
}

func (k ProjectKeeper) updateProject(ctx sdk.Context, msg ProjectMsg) sdk.Error {
	bz, err := k.cdc.MarshalBinary(msg)
	if err != nil {
		log.Errorf("Update project failed, project id: %s, status: %s", msg.ID, msg.Status)
		return ErrUpdateProjectFailed("Update project failed.")
	}

	k.setProject(ctx, []byte(msg.ID), bz)
	return nil
}

func (k ProjectKeeper) getProject(ctx sdk.Context, id []byte) (msg ProjectMsg) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(id)
	if bz == nil {
		return ProjectMsg{}
	}

	err := k.cdc.UnmarshalBinary(bz, &msg)
	if err != nil {
		panic(err)
	}

	return
}

func (k ProjectKeeper) setProject(ctx sdk.Context, id []byte, bz []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(id, bz)
}

func (k ProjectKeeper) IsProjectExisted(ctx sdk.Context, projectId string) bool {
	p := k.getProject(ctx, []byte(projectId))
	if p.ID != "" {
		return true
	}
	return false
}

func (k ProjectKeeper) GetProjectStatus(ctx sdk.Context, projectId string) int {
	return k.getProject(ctx, []byte(projectId)).Status
}
