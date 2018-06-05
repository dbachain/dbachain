package user

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	wire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
)

// Accountbase
type Userbase interface {
	// Sign some bytes
	Sign(name, passphrase string, msg []byte) (crypto.Signature, crypto.PubKey, error)
	// Create a new keypair
	Create(name, passphrase string, algo keys.CryptoAlgo, info AccountInfo) (accountinfo AccountInfo, seed string, err error)
	// Recover takes a seedphrase and loads in the key
	Recover(name, passphrase, seedphrase string, info AccountInfo) (accountinfo AccountInfo, erro error)
	List() ([]AccountInfo, error)
	Get(name string) (AccountInfo, error)
	Update(name, newpass string, info AccountInfo) error
	Delete(name, passphrase string) error

	Import(name string, armor string) (err error)
	Export(name string) (armor string, err error)
	GetUserAddress(name string) crypto.Address
}

// Info is the public information about a key
type AccountInfo struct {
	keys.Info
	UserHash    cmn.HexBytes
	UserStatus  int
	Nationality string
	// TO DO
}

func newAccountInfo(name string, pub crypto.PubKey, privArmor string) AccountInfo {
	return AccountInfo{
		Info: keys.Info{
			Name:         name,
			PubKey:       pub,
			PrivKeyArmor: privArmor,
		},
		UserHash:    nil,
		UserStatus:  -1,
		Nationality: "",
		// TO DO
	}
}

// Address is a helper function to calculate the address from the pubkey
func (i AccountInfo) AccountAddress() []byte {
	return i.Info.PubKey.Address()
}

func (i AccountInfo) bytes() []byte {
	bz, err := wire.MarshalBinary(i)
	if err != nil {
		panic(err)
	}
	return bz
}

func readInfo(bz []byte) (info AccountInfo, err error) {
	err = wire.UnmarshalBinary(bz, &info)
	return
}
