package vote

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// Keys for store prefixes
	UserVoteKey               = []byte{0x00} // prefix for each key to a user's vote info
	AccumulatedProjectVoteKey = []byte{0x01} // prefix for each key to a accumlated project vote
)

func GetUserVoteKey(addr sdk.Address, projID string, ctime int64) []byte {
	strCTime := time.Unix(ctime, 0).String()

	return []byte(addr.String() + projID + strCTime)
}

func GetUserVotePrefixKey(addr sdk.Address) []byte {
	return append(UserVoteKey, addr.Bytes()...)
}

func GetAccumulatedProjectVoteKey(projID string) []byte {
	return append(AccumulatedProjectVoteKey, []byte(projID)...)
}
