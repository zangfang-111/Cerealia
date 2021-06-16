package model

import (
	"fmt"

	"github.com/robert-zaremba/errstack"
)

const hdCerealiaWalletFormat = "m/44'/148'/%v'"

// KeyWallet is an interface to produce user's public keys
// It may mutate the user object and return it
type KeyWallet interface {
	// DeriveKey produces same key for same keyDerivationPath
	// This method mutates the wallet
	DeriveNewKey() (pubKey string, keyDerivationPath string, error errstack.E)
}

// DeriveNewKey implements interface KeyWallet
func (w *StaticWallet) DeriveNewKey() (pubKey string, keyDerivationPath string, error errstack.E) {
	return w.PubKey, "", nil
}

// DeriveNewKey implements interface KeyWallet
// TODO: add pubkey derivation lib
func (w *HDCerealiaWallet) DeriveNewKey() (pubKey string, keyDerivationPath string, error errstack.E) {
	format := fmt.Sprintf(hdCerealiaWalletFormat, w.DerivationIndex)
	w.DerivationIndex++
	// TODO: Derive the key here
	return "GD3EPS4EBOK6ZELDEN466I6EU4LW7TK6UL6INZRC6OKLZKYXXESS64VE", format, nil
}
