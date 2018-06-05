package user

import (
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type UserMsg struct {
	Sender     sdk.Address `json:"address"`
	Address    sdk.Address `json:"address"`
	MsgType    string      `json:"msg_type"`
	ID         string      `json:"id"`
	UserStatus int         `json:"user_status"`
}

func NewUserMsg(msgType string, sender sdk.Address, address sdk.Address, name string, userStatus int) UserMsg {
	return UserMsg{
		Sender:     sender,
		Address:    address,
		MsgType:    msgType,
		ID:         name,
		UserStatus: userStatus,
	}
}

func (msg UserMsg) Type() string { return "user" }

func (msg UserMsg) ValidateBasic() sdk.Error {

	//TO DO
	return nil
}

func (msg UserMsg) String() string {
	return fmt.Sprintf("UserMsg{MsgType:%v,UserStatus:%v}", msg.MsgType, msg.UserStatus)
}

func (msg UserMsg) GetSignBytes() []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		log.Fatal(err.Error()) //TODO better proccess?
	}

	return bz
}

func (msg UserMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Sender}
}

func (msg UserMsg) Get(key interface{}) (value interface{}) {
	// TO DO
	return nil
}
