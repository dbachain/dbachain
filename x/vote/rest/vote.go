package rest

import (
	"net/http"

	"github.com/dbachain/dbachain/common/utils"
	"github.com/dbachain/dbachain/x/user"

	"github.com/cosmos/cosmos-sdk/wire"
)

type voteBody struct {
	UserName  string `json:"user_name"`
	ProjectID string `json:"project_id"`
	Weight    int    `json:"weight"`
	Round     int    `json:"round"`
}

// JSON RPC

func VoteRequestHandler(cdc *wire.Codec, ub user.Userbase) func(http.ResponseWriter, *http.Request) {
	//	c := commander{cdc}
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var vote voteBody
		err := utils.JSONUnmarshal(r.Body, &vote)
		if nil != err {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		//		output, err := c.queryTx(hashHexStr, trustNode)
		//		if err != nil {
		//			w.WriteHeader(500)
		//			w.Write([]byte(err.Error()))
		//			return
		//		}
		//		w.Write(output)
	}
}
