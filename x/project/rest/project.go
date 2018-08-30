package rest

import (
	"net/http"

	"github.com/dbachain/dbachain/common/utils"
	"github.com/dbachain/dbachain/x/user"

	"github.com/cosmos/cosmos-sdk/wire"
)

type projectBody struct {
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
}

func ProjectRequestHandler(cdc *wire.Codec, ub user.Userbase) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var project projectBody
		err := utils.JSONUnmarshal(r.Body, &project)
		if nil != err {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
