package stellar

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/txlog/txlogi"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
)

func logTxAction(txb64 string, e *xdr.TransactionEnvelope, txErr error, txlogger txlogi.Logger) errstack.E {
	if txErr == nil {
		// Transaction passed
		return errstack.WrapAsInf(
			txlogger.LogTxStatus(txb64, e, model.TxStatusOk),
			"Transaction passed, can't add TxLog entry",
		)
	}
	logErr := txlogger.LogTxStatus(txb64, e, model.TxStatusFailed)
	if logErr == nil {
		// Transaction failed, logging is fine
		return errstack.WrapAsDomain(
			txErr,
			"Transaction failed, can't add TxLog entry",
		)
	}
	// Transaction failed and logging failed too. Let's at least return all errors.
	errb := errstack.NewBuilder()
	errb.Put("tx-err", txErr)
	errb.Put("log-err", logErr)
	return errstack.WrapAsInf(errb.ToReqErr(), "Tx failed and tx-log failed too.")
}

// Driver holds master secret and network parameters
type Driver struct {
	Network Network
	Client  Client
}

// NewDriver creates a new StellarDriver
func NewDriver(netName string) (*Driver, error) {
	var client Client

	net, client, errs := getNetworkAndClient(netName)
	logger.Info("Created Stellar Driver", "network", net)
	return &Driver{
		Client:  client,
		Network: net,
	}, errs
}

// SignEnvelope signs the envelope
func (c *Driver) SignEnvelope(txEnvelope *build.TransactionEnvelopeBuilder, key keypair.Full) (*build.TransactionEnvelopeBuilder, errstack.E) {
	err := txEnvelope.MutateTX(c.Network.Passphrase)
	if err != nil {
		return nil, errstack.WrapAsDomain(err, "Can't sign the envelope (Can't set network passphrase)")
	}
	signatureMutator := build.Sign{Seed: key.Seed()}
	err = signatureMutator.MutateTransactionEnvelope(txEnvelope)
	return txEnvelope, errstack.WrapAsDomain(err, "Can't sign the envelope (Can't add signature)")
}

// WithTxLogger wraps a driver with a tx logger
func (c Driver) WithTxLogger(l txlogi.Logger, isAccLockAcquiredFn func() error) *WrappedDriver {
	return &WrappedDriver{
		Driver:              c,
		txlogger:            l,
		isAccLockAcquiredFn: isAccLockAcquiredFn,
	}
}

// WrappedDriver combines a Driver operations with a Tx Logger
type WrappedDriver struct {
	Driver
	txlogger            txlogi.Logger
	isAccLockAcquiredFn func() error
}

// Send sends tx to stellar network
func (c *WrappedDriver) Send(signedTx *build.TransactionEnvelopeBuilder) (*hProtocol.TransactionSuccess, errstack.E) {
	err := c.isAccLockAcquiredFn()
	if err != nil {
		return nil, errstack.WrapAsDomain(err, "Refusing to send transacton. Lock is not acquired.")
	}
	txb64, err := signedTx.Base64()
	if err != nil {
		return nil, errstack.WrapAsDomain(err, "Can't convert transaction to base64")
	}
	err = c.txlogger.LogTxStatus(txb64, signedTx.E, model.TxStatusPending)
	if err != nil {
		return nil, errstack.WrapAsInf(err, "Transaction is pending, can't add TxLog entry")
	}
	response, err := c.Client.SubmitTransaction(txb64)
	return &response, logTxAction(txb64, signedTx.E, err, c.txlogger)
}

func convertKeypairsToStrings(signSecrets []keypair.Full) []string {
	var a []string
	for _, secret := range signSecrets {
		a = append(a, secret.Seed())
	}
	return a
}

// SignAndSend signs the given tx with all given secrets and sends it to the network
func (c *WrappedDriver) SignAndSend(tx build.TransactionBuilder, signKeyPairs ...keypair.Full) (*hProtocol.TransactionSuccess, errstack.E) {
	signedTx, err := tx.Sign(convertKeypairsToStrings(signKeyPairs)...)
	if err != nil {
		return nil, errstack.WrapAsDomain(err, "Can't sign transacton. Probably the key is malformed.")
	}
	return c.Send(&signedTx)
}

// SignAndSendEnvelope adds an additional signature and sends it to stellar
func (c *WrappedDriver) SignAndSendEnvelope(txEnvelope *build.TransactionEnvelopeBuilder, keypairs ...keypair.Full) (*hProtocol.TransactionSuccess, errstack.E) {
	signedEnvelope := txEnvelope
	var err errstack.E
	for _, key := range keypairs {
		signedEnvelope, err = c.SignEnvelope(txEnvelope, key)
		if err != nil {
			return nil, err
		}
	}
	return c.Send(signedEnvelope)
}

// SignAndSendEnvelopeSource is a convenience method for source account type
func (c *WrappedDriver) SignAndSendEnvelopeSource(txEnvelope *build.TransactionEnvelopeBuilder, sourceAccs *txsource.SourceAccs) (*hProtocol.TransactionSuccess, errstack.E) {
	if sourceAccs.PoolAcc.KeyPair == sourceAccs.TradeKeyPair {
		return c.SignAndSendEnvelope(txEnvelope, sourceAccs.TradeKeyPair)
	}
	return c.SignAndSendEnvelope(txEnvelope, sourceAccs.PoolAcc.KeyPair, sourceAccs.TradeKeyPair)
}
