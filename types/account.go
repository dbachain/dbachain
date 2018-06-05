package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	crypto "github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
)

var _ sdk.Account = (*AppAccount)(nil)

// Custom extensions for this application.  This is just an example of
// extending auth.BaseAccount with custom fields.
//
// This is compatible with the stock auth.AccountStore, since
// auth.AccountStore uses the flexible go-wire library.
type AppAccount struct {
	auth.BaseAccount
	Name        string       `json:"name"`
	UserHash    cmn.HexBytes `json:"user_hash"`
	UserStatus  int          `json:"user_status"`
	Nationality string       `json:"nationality"`
}

// Implements sdk.Account.
func (acc AppAccount) GetAddress() sdk.Address {
	return acc.BaseAccount.Address
}

// Implements sdk.Account.
func (acc *AppAccount) SetAddress(addr sdk.Address) error {
	return acc.BaseAccount.SetAddress(addr)
}

// Implements sdk.Account.
func (acc AppAccount) GetPubKey() crypto.PubKey {
	return acc.BaseAccount.PubKey
}

// Implements sdk.Account.
func (acc *AppAccount) SetPubKey(pubKey crypto.PubKey) error {
	return acc.BaseAccount.SetPubKey(pubKey)
}

// Implements sdk.Account.
func (acc AppAccount) GetCoins() sdk.Coins {
	return acc.BaseAccount.GetCoins()
}

// Implements sdk.Account.
func (acc *AppAccount) SetCoins(coins sdk.Coins) error {
	return acc.BaseAccount.SetCoins(coins)
}

// Implements sdk.Account.
func (acc AppAccount) GetSequence() int64 {
	return acc.BaseAccount.GetSequence()
}

// Implements sdk.Account.
func (acc *AppAccount) SetSequence(seq int64) error {
	return acc.BaseAccount.SetSequence(seq)
}

func (acc AppAccount) GetName() string                    { return acc.Name }
func (acc *AppAccount) SetName(name string)               { acc.Name = name }
func (acc AppAccount) GetUserStatus() int                 { return acc.UserStatus }
func (acc *AppAccount) SetUserStatus(userstatus int)      { acc.UserStatus = userstatus }
func (acc AppAccount) GetNationality() string             { return acc.Nationality }
func (acc *AppAccount) SetNationality(nationality string) { acc.Nationality = nationality }

// Get the AccountDecoder function for the custom AppAccount
func GetAccountDecoder(cdc *wire.Codec) sdk.AccountDecoder {
	return func(accBytes []byte) (res sdk.Account, err error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}
		acct := new(AppAccount)
		err = cdc.UnmarshalBinary(accBytes, &acct)
		if err != nil {
			panic(err)
		}
		return acct, err
	}
}

// State to Unmarshal
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name    string      `json:"name"`
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

func NewGenesisAccount(aa *AppAccount) *GenesisAccount {
	return &GenesisAccount{
		Name:    aa.Name,
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// convert GenesisAccount to AppAccount
func (ga *GenesisAccount) ToAppAccount() (acc *AppAccount, err error) {
	baseAcc := auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
	return &AppAccount{
		BaseAccount: baseAcc,
		Name:        ga.Name,
	}, nil
}
