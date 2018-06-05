package project

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//custom ABCI Response Codes
const (
	CodeInvalidMessageType    sdk.CodeType = 3000
	CodeInvalidProjectID      sdk.CodeType = 3001
	CodeInvalidProjectStatus  sdk.CodeType = 3002
	CodeProjectAlreadyExisted sdk.CodeType = 3003
	CodeProjectNotExist       sdk.CodeType = 3004
	CodeSaveProjectFailed     sdk.CodeType = 3005
	CodeUpdateProjectFailed   sdk.CodeType = 3006
)

func ErrInvalidMessageType(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidMessageType, msg)
}

func ErrInvalidProjectID(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidProjectID, msg)
}

func ErrInvalidProjectStatus(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidProjectStatus, msg)
}

func ErrProjectAlreadyExisted(msg string) sdk.Error {
	return sdk.NewError(CodeProjectAlreadyExisted, msg)
}

func ErrProjectNotExist(msg string) sdk.Error {
	return sdk.NewError(CodeProjectNotExist, msg)
}

func ErrSaveProjectFailed(msg string) sdk.Error {
	return sdk.NewError(CodeSaveProjectFailed, msg)
}

func ErrUpdateProjectFailed(msg string) sdk.Error {
	return sdk.NewError(CodeUpdateProjectFailed, msg)
}
