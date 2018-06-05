package rest

import (
	"net/http"

	"dbachain/common/utils"
	"dbachain/x/user"

	"github.com/cosmos/cosmos-sdk/wire"
	cmn "github.com/tendermint/tmlibs/common"
)

type accountBody struct {
	UserName   string       `json:"user_name"`
	UserHash   cmn.HexBytes `json:"user_hash"`
	UserStatus bool         `json:"user_status"`
}

// JSON RPC

func AccountRequestHandler(cdc *wire.Codec, ub user.Userbase) func(http.ResponseWriter, *http.Request) {
	//	c := commander{cdc}
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var acc accountBody
		err := utils.JSONUnmarshal(r.Body, &acc)
		if nil != err {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// TO DO
	}
}
