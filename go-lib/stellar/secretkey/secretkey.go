package secretkey

import (
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/keypair"
)

// Parse parsed keypair and returns a full keypair
func Parse(secret string) (*keypair.Full, errstack.E) {
	kp, err := keypair.Parse(secret)
	if err != nil {
		return nil, errstack.WrapAsReq(err, "Can't parse the secret")
	}
	full, ok := kp.(*keypair.Full)
	if !ok {
		return nil, errstack.WrapAsReq(err, "Provided secret value is not a secret")
	}
	return full, nil
}
