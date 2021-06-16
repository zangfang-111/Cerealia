package txvalidation

import (
	"encoding/base64"

	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

func base64ToXDR(b64str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64str)
}

func readEnvelope(envelopeBase64 string) (*xdr.TransactionEnvelope, error) {
	binary, err := base64ToXDR(envelopeBase64)
	if err != nil {
		return nil, err
	}
	envelope := xdr.TransactionEnvelope{}
	err = envelope.UnmarshalBinary(binary)
	return &envelope, err
}

// ReadEnvelopeBuilder reads base64 signed envelope into typed obj
func ReadEnvelopeBuilder(envelopeBase64 string) (*build.TransactionEnvelopeBuilder, error) {
	txEnvelope, err := readEnvelope(envelopeBase64)
	if err != nil {
		return nil, err
	}
	builder := build.TransactionEnvelopeBuilder{E: txEnvelope}
	return &builder, err
}
