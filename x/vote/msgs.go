package vote

import (
	"dbachain/common/log"

	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

type VoteMsg struct {
	UserAddr  sdk.Address `json:"user_addr"`
	ProjectID string      `json:"project_id"`
	Coin      string      `json:"coin"`
	Round     int64       `json:"round"`
	CTime     int64       `json:"ctime"`
}

func NewVoteMsg(addr sdk.Address, prjID, coin string, round int64, ctime int64) VoteMsg {
	return VoteMsg{
		UserAddr:  addr,
		ProjectID: prjID,
		Coin:      coin,
		Round:     round,
		CTime:     ctime,
	}
}

func (m VoteMsg) Type() string {
	return "vote"
}

func (m VoteMsg) ValidateBasic() sdk.Error {
	if m.Round <= 0 {
		return ErrInvalidVoteRound("vote round <= 0")
	}

	if m.Coin == "" {
		return ErrInvalidVoteCoin("vote coin cannot be empty")
	}

	if strings.TrimSpace(m.ProjectID) == "" {
		return ErrInvalidProjectID("projectID cannot be empty")
	}

	return nil
}

func (m VoteMsg) GetSignBytes() []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(m)
	if err != nil {
		log.Fatal(err.Error()) //TODO better proccess?
	}

	return bz
}

func (m VoteMsg) GetSigners() []sdk.Address {
	return []sdk.Address{m.UserAddr}
}

func (m VoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}
