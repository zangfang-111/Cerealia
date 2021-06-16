package stellar

import (
	"fmt"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
	b "github.com/stellar/go/build"
)

// This fee should incorporate:
// + 1XLM for account creation
// + other set up fees for later
// https://www.stellar.org/developers/guides/concepts/fees.html
// (2 + 3 signers + 2 * 3 data-entries) * 0.5 (base fee) = 11 * 0.5 = 5.5
// Can't remove original account's owner from acc in one tx. Increase by 0.5
const initialNewAccountFunds = "6.5"

// CreateTradeAccount creates a tx with CreateTradeAccount operation
// Also it sets the buyer and seller
//
// There are two kinds of signers:
// 1. validation signer
// 2. trade-party signer
//
// Our back-end acts as a validation signer and validates transactions
// before transmitting them to blockchain
// Trade-party signer produces transactions and submits them to our back-end
//
// Weights are calculated in a way that neither we alone,
// nor trade party signers together could add any transactions to the blockchain.
//
// This means that all calls to it are done via our service
// and therefore are validated using it as a validation mechanism.
//
func CreateTradeAccount(d *WrappedDriver, pks model.TradeParticipants, sources *txsource.SourceAccs) errstack.E {
	var tradeParty uint32 = 1                      // weight of any trade-party key
	tradeParticipants := tradeParty * 2            // 2 parties: buyer & seller
	validator := tradeParticipants                 // validator == trade participants
	signatureRequirement := validator + tradeParty // validator + any trade party
	t, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: string(sources.PoolAcc.PubKey)},
		b.AutoSequence{SequenceProvider: d.Client},
		d.Network.Passphrase,
		b.CreateAccount(
			b.Destination{AddressOrSeed: sources.TradeKeyPair.Seed()},
			b.NativeAmount{Amount: initialNewAccountFunds},
		),
		b.SetOptions(
			b.SourceAccount{AddressOrSeed: sources.TradeKeyPair.Seed()},
			// Can't add two signers using one SetOptions operation
			b.AddSigner(pks.Seller.PubKey, tradeParty),
		),
		b.SetOptions(
			b.SourceAccount{AddressOrSeed: sources.TradeKeyPair.Seed()},
			b.AddSigner(pks.Buyer.PubKey, tradeParty),
			b.MasterWeight(validator),
			b.SetThresholds(
				signatureRequirement,
				signatureRequirement,
				validator+tradeParticipants), // All have to agree
		),
	)
	if err != nil {
		return errstack.WrapAsDomain(err, "Can't construct a 'create account' transaction")
	}
	_, errs := d.SignAndSend(*t, sources.PoolAcc.KeyPair, sources.TradeKeyPair)
	return errs
}

// MkTradeDocApprovalTx makes tx for stage document approval.
// input:
// d:        stellar driver
// sources:  trade account sources with pool and trade account addresses
// id:       Document and stage identifier
// entity:   stage_doc
// op:       pending/approved/rejected
func MkTradeDocApprovalTx(d *Driver, sources *txsource.SourceAccs, id model.TradeStageDocPath, entity model.TxTradeEntity, op model.Approval) (string, error) {
	return mkDataMemoTx(d, sources, docPathFormat(id), entity, op, id.StageDocHash)
}

// MkTradeDocApprovalExpireTx makes tx for stage document approval.
// input:
// d:         stellar driver
// sources:   trade account sources with pool and trade account addresses
// id:        Document and stage identifier
// entity:    stage_doc
// op:        pending/approved/rejected
// expiresAt: expiration time
func MkTradeDocApprovalExpireTx(d *Driver, sources *txsource.SourceAccs, id model.TradeStageDocPath, entity model.TxTradeEntity, op model.Approval, expiresAt time.Time) (string, error) {
	return mkDataMemoTxExpire(d, sources, docPathFormat(id), entity, op, id.StageDocHash, expiresAt)
}

// MkTradeStageOperationTx makes tx for stage operation.
// input:
// d:       stellar driver
// sources: trade account sources with pool and trade account addresses
// stageID: 2
// entity:  stage_closeReqs/stage_add
// op:      pending/approved/rejected
func MkTradeStageOperationTx(d *Driver, sources *txsource.SourceAccs, stageID uint, entity model.TxTradeEntity, op model.Approval) (string, error) {
	return mkDataTx(d, sources, fmt.Sprintf("%d", stageID), entity, op)
}

// MkTradeCloseTx makes tx for trade completion.
// input:
// d:        stellar driver
// sources:  trade account sources with pool and trade account addresses
// entityID: 2 (tradeID)
// entity:	 trade_closeReqs
// op:	     pending/approved/rejected
func MkTradeCloseTx(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval) (string, error) {
	return mkDataTx(d, sources, entityID, entity, op)
}

// mkDataMemoTxExpire makes tx for document action with memo.
// input:
// d:          stellar driver
// sources:    trade account sources with pool and trade account addresses
// entityID:   2:4 (stageIdx:stageDocIdx)
// entity:     stage_doc
// op:         pending/approved/rejected
// memoText:   document hash
func mkDataMemoTx(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval, memoText string) (string, error) {
	return makeTx(mkDataMutationsWithHash(d, sources, entityID, entity, op, memoText)...)
}

// mkDataMemoTxExpire makes tx for document action with memo and expiration time.
// input:
// d:          stellar driver
// sources:    trade account sources with pool and trade account addresses
// entityID:   2:4 (stageIdx:stageDocIdx)
// entity:     stage_doc
// op:         pending/approved/rejected
// memoText:   document hash
// expireTime: time attribute
func mkDataMemoTxExpire(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval, memoText string, expireTime time.Time) (string, error) {
	return makeTx(mkDataMutationsWithHashAndExpiration(d, sources, entityID, entity, op, memoText, expireTime)...)
}

// mkDataTx makes tx for stage operation.
// input:
// driver: 			stellar driver
// scAddr: 			trade account address
// entityID:		2 (stageIdx)
// entity:		stage_closeReqs/stage_add
// op:	pending/approved/rejected
func mkDataTx(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval) (string, error) {
	return makeTx(mkDataMutations(d, sources, entityID, entity, op)...)
}

// mkDataMutations produces basic data fields for tx
func mkDataMutations(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval) []b.TransactionMutator {
	return []b.TransactionMutator{
		b.SourceAccount{AddressOrSeed: string(sources.PoolAcc.PubKey)},
		b.AutoSequence{SequenceProvider: d.Client},
		d.Network.Passphrase,
		b.SetData("entity", []byte(entity), b.SourceAccount{AddressOrSeed: sources.TradeKeyPair.Seed()}),
		b.SetData("idx", []byte(entityID), b.SourceAccount{AddressOrSeed: sources.TradeKeyPair.Seed()}),
		b.SetData("operation", []byte(op), b.SourceAccount{AddressOrSeed: sources.TradeKeyPair.Seed()}),
	}
}

// mkDataMutationsWithHash produces basic data fields with hash
func mkDataMutationsWithHash(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval, memoText string) []b.TransactionMutator {
	return append(
		mkDataMutations(d, sources, entityID, entity, op),
		b.MemoHash{Value: memoTextToHash(memoText)},
	)
}

// mkDataMutationsWithHashAndExpiration produces basic data fields with hash and expireTime
func mkDataMutationsWithHashAndExpiration(d *Driver, sources *txsource.SourceAccs, entityID string, entity model.TxTradeEntity, op model.Approval, memoText string, expireTime time.Time) []b.TransactionMutator {
	return append(
		mkDataMutationsWithHash(d, sources, entityID, entity, op, memoText),
		b.SetData("expireTime", []byte(bat.I64toa(expireTime.Unix())), b.SourceAccount{AddressOrSeed: sources.TradeKeyPair.Seed()}),
	)
}
