package user

import (
	"errors"
	"fmt"
	"strings"

	mycrypto "github.com/dbachain/dbachain/x/user/crypto"

	pkgerr "github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	wire "github.com/tendermint/go-wire"
	dbm "github.com/tendermint/tmlibs/db"
)

type dbAccountbase struct {
	db    dbm.DB
	codec words.Codec
}

func New(db dbm.DB, codec words.Codec) dbAccountbase {
	return dbAccountbase{
		db:    db,
		codec: codec,
	}
}

var _ Userbase = dbAccountbase{}

func (ab dbAccountbase) Create(name, passphrase string, algo keys.CryptoAlgo, info AccountInfo) (AccountInfo, string, error) {
	// NOTE: secret is SHA256 hashed by secp256k1 and ed25519.
	// 16 byte secret corresponds to 12 BIP39 words.
	// XXX: Ledgers use 24 words now - should we ?
	secret := crypto.CRandBytes(16)
	priv, err := generate(algo, secret)
	if err != nil {
		return AccountInfo{}, "", err
	}

	// encrypt and persist the key
	accountinfo := ab.writeKey(priv, name, passphrase, info)

	// we append the type byte to the serialized secret to help with
	// recovery
	// ie [secret] = [type] + [secret]
	typ := mycrypto.CryptoAlgoToByte(algo)
	secret = append([]byte{typ}, secret...)

	// return the mnemonic phrase
	words, err := ab.codec.BytesToWords(secret)
	seed := strings.Join(words, " ")
	return accountinfo, seed, err
}

// Recover converts a seedphrase to a private key and persists it,
// encrypted with the given passphrase.  Functions like Create, but
// seedphrase is input not output.
func (ab dbAccountbase) Recover(name, passphrase, seedphrase string, info AccountInfo) (AccountInfo, error) {
	words := strings.Split(strings.TrimSpace(seedphrase), " ")
	secret, err := ab.codec.WordsToBytes(words)
	if err != nil {
		return AccountInfo{}, err
	}

	// secret is comprised of the actual secret with the type
	// appended.
	// ie [secret] = [type] + [secret]
	typ, secret := secret[0], secret[1:]
	algo := mycrypto.ByteToCryptoAlgo(typ)
	priv, err := generate(algo, secret)
	if err != nil {
		return AccountInfo{}, err
	}

	// encrypt and persist key.
	public := ab.writeKey(priv, name, passphrase, info)
	return public, err
}

// List returns the keys from storage in alphabetical order.
func (ab dbAccountbase) List() ([]AccountInfo, error) {
	var res []AccountInfo
	iter := ab.db.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		info, err := readInfo(iter.Value())
		if err != nil {
			return nil, err
		}
		res = append(res, info)
	}
	return res, nil
}

// Get returns the public information about one key.
func (ab dbAccountbase) Get(name string) (AccountInfo, error) {
	bs := ab.db.Get(infoKey(name))
	return readInfo(bs)
}

// Sign signs the msg with the named key.
// It returns an error if the key doesn't exist or the decryption fails.
func (ab dbAccountbase) Sign(name, passphrase string, msg []byte) (sig crypto.Signature, pub crypto.PubKey, err error) {
	info, err := ab.Get(name)
	if err != nil {
		return
	}
	priv, err := mycrypto.UnarmorDecryptPrivKey(info.Info.PrivKeyArmor, passphrase)
	if err != nil {
		return
	}
	sig = priv.Sign(msg)
	pub = priv.PubKey()
	return
}

func (ab dbAccountbase) Export(name string) (armor string, err error) {
	info, err := ab.Get(name)
	if err != nil {
		return
	}
	bz, err := wire.MarshalBinary(info.Info)
	if err != nil {
		return
	}
	return mycrypto.ArmorInfoBytes(bz), nil
}

func (ab dbAccountbase) Import(name string, armor string) (err error) {
	bz := ab.db.Get(infoKey(name))
	if len(bz) > 0 {
		return errors.New("Cannot overwrite data for name " + name)
	}
	infoBytes, err := mycrypto.UnarmorInfoBytes(armor)
	if err != nil {
		return
	}
	accountinfo := AccountInfo{}
	err = wire.UnmarshalBinary(infoBytes, &accountinfo.Info)
	if err != nil {
		return
	}
	accountinfobytes, err := wire.MarshalBinary(accountinfo)
	if err != nil {
		return
	}
	ab.db.Set(infoKey(name), accountinfobytes)
	return nil
}

// Delete removes key forever, but we must present the
// proper passphrase before deleting it (for security).
func (ab dbAccountbase) Delete(name, passphrase string) error {
	//verify we have the proper password before deleting
	info, err := ab.Get(name)
	if err != nil {
		return err
	}
	_, err = mycrypto.UnarmorDecryptPrivKey(info.PrivKeyArmor, passphrase)
	if err != nil {
		return err
	}
	ab.db.DeleteSync(infoKey(name))
	return nil
}

//
func (ab dbAccountbase) Update(name, pass string, accountinfo AccountInfo) error {
	info, err := ab.Get(name)
	if err != nil {
		return err
	}
	key, err := mycrypto.UnarmorDecryptPrivKey(info.Info.PrivKeyArmor, pass)
	if err != nil {
		return err
	}

	ab.writeKey(key, name, pass, accountinfo)
	return nil
}

func (ab dbAccountbase) GetUserAddress(name string) crypto.Address {
	info, err := ab.Get(name)
	if err != nil {
		return crypto.Address{}
	}

	return info.Info.PubKey.Address()
}

func (ab dbAccountbase) writeKey(priv crypto.PrivKey, name, passphrase string, accountinfo AccountInfo) AccountInfo {
	// generate the encrypted privkey
	privArmor := mycrypto.EncryptArmorPrivKey(priv, passphrase)

	accountinfo.Info = keys.Info{name, priv.PubKey(), privArmor}

	bz, err := wire.MarshalBinary(accountinfo)
	if err != nil {
		panic(err)
	}

	// write them both
	ab.db.SetSync(infoKey(name), bz)
	return accountinfo
}

func generate(algo keys.CryptoAlgo, secret []byte) (crypto.PrivKey, error) {
	switch algo {
	case keys.AlgoEd25519:
		return crypto.GenPrivKeyEd25519FromSecret(secret).Wrap(), nil
	case keys.AlgoSecp256k1:
		return crypto.GenPrivKeySecp256k1FromSecret(secret).Wrap(), nil
	default:
		err := pkgerr.Errorf("Cannot generate keys for algorithm: %s", algo)
		return crypto.PrivKey{}, err
	}
}

func infoKey(name string) []byte {
	return []byte(fmt.Sprintf("%s.info", name))
}
