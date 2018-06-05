package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//custom ABCI Response Codes
const (
	CodeInvalidVoteRound      sdk.CodeType = 5000
	CodeInvalidVoteCoin       sdk.CodeType = 5001
	CodeInvalidProjectID      sdk.CodeType = 5002
	CodeProjectStatusUnnormal sdk.CodeType = 5003
)

func ErrInvalidVoteRound(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidVoteRound, msg)
}

func ErrInvalidVoteCoin(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidVoteCoin, msg)
}

func ErrInvalidProjectID(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidProjectID, msg)
}

func ErrCodeProjectStatusUnnormal(msg string) sdk.Error {
	return sdk.NewError(CodeProjectStatusUnnormal, msg)
}
