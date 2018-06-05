package project

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"encoding/json"
)

type ProjectMsg struct {
	Address sdk.Address `json:"address"`
	MsgType string      `json:"msg_type"`
	ID      string      `json:"id"`
	Status  int         `json:"status"`
	CTime   int64       `json:"ctime"`
}

func NewProjectMsg(msgType string, address sdk.Address, projectID string, projectStatus int, ctime int64) ProjectMsg {
	return ProjectMsg{
		MsgType: msgType,
		Address: address,
		ID:      projectID,
		Status:  projectStatus,
		CTime:   ctime,
	}
}

func (msg ProjectMsg) ValidateBasic() sdk.Error {
	if strings.TrimSpace(msg.MsgType) == "" {
		return ErrInvalidMessageType("ProjectMsg type cannot be empty.")
	}

	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidProjectID("ProjectID cannot be empty.")
	}

	return nil
}

func (msg ProjectMsg) Type() string {
	return "project"
}

func (msg ProjectMsg) String() string {
	return fmt.Sprintf("ProjectMsg{MsgType:%v,ID:%v,Status:%v}", msg.MsgType, msg.ID, msg.Status)
}

func (msg ProjectMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg ProjectMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

func (msg ProjectMsg) Get(key interface{}) (value interface{}) {
	return nil
}
