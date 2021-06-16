// Package txvalidation contains all transaction validation functions
package txvalidation

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	bat "github.com/robert-zaremba/go-bat"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/davecgh/go-spew/spew"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

// tx field for front-end errors
const validationFieldTX = "tx"

// translation keys for front-end
const txUnparsable = "validation.stellar.parse-error"
const txSignature = "validation.stellar.signature-error"
const oneSignatureExpected = "validation.stellar.one-signature-expected"
const badData = "validation.stellar.bad-data"
const badMemoErr = "validation.stellar.bad-memo"

// Data field names
const dataKeyEntity = "entity"
const dataKeyIdx = "idx"
const dataKeyOperation = "operation"
const dataKeyExpireTime = "expireTime"

func idxJoin(indexes ...uint) string {
	res := make([]string, len(indexes))
	for i := 0; i < len(indexes); i++ {
		res[i] = strconv.Itoa(int(indexes[i]))
	}
	return strings.Join(res, ":")
}

func parseUserKeypair(t *model.Trade, u *model.User) (keypair.KP, error) {
	var err error
	keyStr, err := t.FindPubKey(u)
	if err != nil {
		return nil, err
	}
	return keypair.Parse(string(keyStr))
}

//This function verifies that the signer and the submitter have the same keys
//
//	The transaction can be signed by other party too, so this check removes the possibility of this kind of attack.
//
//	The attack would be like this:
//    I steal the private secret key of my Victim
//    I open the trade
//    I send any transaction I want using that stolen key
//
//  This could work the opposite way too:
//    I steal the cookie of my Victim
//    I create the transaction as Me, because I don’t have the private key of the Victim
//    I send the transaction to back-end and use his cookie to do that.
//
//	This kind of check makes sure that tx signer and submitter are the same user
func validateSignature(userAddr keypair.KP, se *SimplifiedEnvelope, eb *build.TransactionEnvelopeBuilder) error {
	if len(eb.E.Signatures) != 1 {
		return errstack.NewReq(oneSignatureExpected)
	}
	return userAddr.Verify(se.TxHash[:], eb.E.Signatures[0].Signature)
}

// Validates that transaction is of data submission type.
// Validation cases:
// * tx can be decoded
// * user address from tx and submitter's main key matches
// * data keys from tx are only from a permitted list (values are not verified here)
// * the data keys are only set to the trade account and no other accounts
func prevalidateTradeDataTx(t *model.Trade, u *model.User, signedTx string) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, *validation.Builder) {
	vb := validation.Builder{}
	simplified, eBuilder, err := Simplify(signedTx)
	if err != nil {
		logger.Error("tx parsing", err, "user", *u, "trade", *t)
		vb.Append(validationFieldTX, txUnparsable)
		return nil, nil, &vb
	}
	userAddr, err := parseUserKeypair(t, u)
	if err != nil {
		logger.Error("bad user address", err, "user", *u, "trade", *t)
		vb.Append(validationFieldTX, txSignature)
		return nil, nil, &vb
	}
	err = validateSignature(userAddr, simplified, eBuilder)
	if err != nil {
		logger.Error("bad user's signature", err, "user", u, "trade", *t)
		vb.Append(validationFieldTX, txSignature)
	}
	if !vb.IsEmpty() {
		return nil, nil, &vb
	}
	return eBuilder, simplified, &vb
}

func compareData(vb *validation.Builder, obtained, expected dataMap) bool {
	isEqual := reflect.DeepEqual(obtained, expected)
	if !isEqual {
		logger.Error(spew.Sprintf("User data validation. Expected: %s; Actual: %s", expected, obtained))
		vb.Append(validationFieldTX, badData)
	}
	return isEqual
}

// validateManageData does the actual validation of the data values
func validateManageData(vb *validation.Builder, obtained dataMap, source model.SCAddr, idx string, entity model.TxTradeEntity, op model.Approval) bool {
	return compareData(vb, obtained, dataMap{
		string(source): {
			dataKeyIdx:       idx,
			dataKeyEntity:    entity.String(),
			dataKeyOperation: op.String(),
		}})
}

func validateManageDataWithExpireTime(vb *validation.Builder, obtained dataMap, source model.SCAddr, idx string,
	entity model.TxTradeEntity, op model.Approval, expireTime time.Time) bool {
	return compareData(vb, obtained, dataMap{
		string(source): {
			dataKeyIdx:        idx,
			dataKeyEntity:     entity.String(),
			dataKeyOperation:  op.String(),
			dataKeyExpireTime: bat.I64toa(expireTime.Unix()),
		}})
}

func validateMemo(vb *validation.Builder, obtained, expected string) {
	if obtained != expected {
		logger.Warn("invalid tx doc hash", "tx-hash", obtained, "expected", expected)
		vb.Append(validationFieldTX, badMemoErr)
	}
}
