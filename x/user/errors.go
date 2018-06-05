package user

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//custom ABCI Response Codes
const (
	CodeInvalidMessageType sdk.CodeType = 4000
	CodeInvalidUserID      sdk.CodeType = 4001
	CodeInvalidUserStatus  sdk.CodeType = 4002
	CodeUserAlreadyExisted sdk.CodeType = 4003
	CodeUserNotExist       sdk.CodeType = 4004
	CodeSaveUserFailed     sdk.CodeType = 4005
	CodeUpdateUserFailed   sdk.CodeType = 4006
)

func ErrInvalidMessageType(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidMessageType, msg)
}

func ErrInvalidUserID(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidUserID, msg)
}

func ErrInvalidUserStatus(msg string) sdk.Error {
	return sdk.NewError(CodeInvalidUserStatus, msg)
}

func ErrUserAlreadyExisted(msg string) sdk.Error {
	return sdk.NewError(CodeUserAlreadyExisted, msg)
}

func ErrUserNotExist(msg string) sdk.Error {
	return sdk.NewError(CodeUserNotExist, msg)
}

func ErrSaveUserFailed(msg string) sdk.Error {
	return sdk.NewError(CodeSaveUserFailed, msg)
}

func ErrUpdateUserFailed(msg string) sdk.Error {
	return sdk.NewError(CodeUpdateUserFailed, msg)
}
