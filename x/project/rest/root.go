package rest

import (
	"github.com/dbachain/dbachain/x/user"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, cdc *wire.Codec, ub user.Userbase) {
	r.HandleFunc("/project", ProjectRequestHandler(cdc, ub)).Methods("POST")
}
