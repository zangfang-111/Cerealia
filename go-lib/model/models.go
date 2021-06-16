// Package model represents database object relations
// This file represents the domain model
package model

import (
	"time"

	"github.com/stellar/go/xdr"
)

// SCAddr represents a blockchain account / smart contract adddress
type SCAddr string

// SCSecret is a private key of smart contract address.
type SCSecret string

// TXSourceAccType represents a type of entity that is represented by smart contract address (trade/user/etc)
type TXSourceAccType string

// ApproveReq represents approve request action
type ApproveReq struct {
	Status       Approval   `json:"status"`
	ReqActor     TradeActor `json:"reqActor"`
	ReqBy        string     `json:"reqBy"`
	ReqAt        time.Time  `json:"reqAt"`
	ReqTx        string     `json:"reqTx"`
	ApprovedBy   string     `json:"approvedBy"`
	ApprovedAt   *time.Time `json:"approvedAt"`
	ApprovedTx   string     `json:"approvedTx"`
	ReqReason    string     `json:"reqReason"`
	RejectReason string     `json:"rejectReason"`
}

// Wallet type for generic wallet description
type Wallet struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

// StaticWallet is a wallet containing only one key
// Shouldn't be used in prod
type StaticWallet struct {
	Wallet
	PubKey string `json:"pubKey"`
}

// HDCerealiaWallet wallet stores metadata for HD public key derivation wallet
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#public-parent-key--public-child-key
// TODO: Currently it's dead code
type HDCerealiaWallet struct {
	Wallet
	MasterPubKey    string `json:"masterPubKey"`
	DerivationIndex int    `json:"derivationIndex"` // Next index for derivation
}

// StageModerator is a type for moderating user info
type StageModerator struct {
	UserID    string     `json:"userID"`
	CreatedAt *time.Time `json:"createdAt"`
}

// WalletType is a string key for wallet type
type WalletType string

// AccessApproval is a approval type for user management
type AccessApproval struct {
	Status    SimpleApproval `json:"status"`
	Approver  string         `json:"approver"`
	Reason    *string        `json:"reason"`
	CreatedAt time.Time      `json:"createdAt"`
}

// User type for user info
type User struct {
	ID                string                      `json:"_key,omitempty"`
	FirstName         string                      `json:"firstname"`
	LastName          string                      `json:"lastname"`
	Emails            []string                    `json:"emails"`
	Roles             []UserRole                  `json:"roles"`
	Avatar            string                      `json:"avatar"`
	Password          []byte                      `json:"password"`
	Organizations     map[string]string           `json:"organizations"`
	CreatedAt         time.Time                   `json:"createdAt"`
	Salt              []byte                      `json:"salt"`
	Biography         string                      `json:"biography"`
	DefaultWalletID   string                      `json:"defaultwalletID"`
	StaticWallets     map[string]StaticWallet     `json:"staticWallets"`
	HDCerealiaWallets map[string]HDCerealiaWallet `json:"hdCerealiaWallets"` // currently it's dead code
	Approvals         []AccessApproval            `json:"approvals"`
}

// TradeParticipant shows the user and his key in his wallet
type TradeParticipant struct {
	UserID            string `json:"userID"`
	KeyDerivationPath string `json:"keyPath,omitempty"` // should be used in HD wallets later, currently it's dead code
	WalletID          string `json:"walletID"`
	PubKey            string `json:"pubKey"`
}

// Trade type for trade info
type Trade struct {
	ID           string             `json:"_key,omitempty"`
	Name         string             `json:"name"`
	Description  *string            `json:"description"`
	TemplateID   string             `json:"templateID"`
	Buyer        TradeParticipant   `json:"buyer"`
	Seller       TradeParticipant   `json:"seller"`
	SCAddr       SCAddr             `json:"scAddr"`
	SCVersion    uint               `json:"scVersion"`
	Stages       []TradeStage       `json:"stages"`
	StageAddReqs []TradeStageAddReq `json:"stageAddReqs"`
	CloseReqs    []ApproveReq       `json:"closeReqs"`
	CreatedBy    string             `json:"createdBy"`
	CreatedAt    time.Time          `json:"createdAt"`
	TradeOfferID *string            `json:"tradeOffer,omitempty"`
	Moderating   DoneStatus         `json:"moderating"`
}

// TradeStageAddReq type for tradeStageAddReq info
type TradeStageAddReq struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Owner       TradeActor `json:"owner"`

	ApproveReq
}

// TradeStage type for stage info
type TradeStage struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	AddReqIdx   int             `json:"addReqIdx"`
	Owner       TradeActor      `json:"owner"`
	ExpiresAt   *time.Time      `json:"expiresAt"`
	Docs        []TradeStageDoc `json:"docs"`
	DelReqs     []ApproveReq    `json:"delReqs"`
	CloseReqs   []ApproveReq    `json:"closeReqs"`
	Moderator   StageModerator  `json:"moderator"`
}

// TradeStageDoc type for TradeStageDoc info
type TradeStageDoc struct {
	DocID        string     `json:"docID"`
	Status       Approval   `json:"status"`
	ReqTx        string     `json:"reqTx"`
	ApprovedTx   string     `json:"approvedTx"`
	ApprovedBy   string     `json:"approvedBy,omitempty"`
	ApprovedAt   *time.Time `json:"approvedAt,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt"`
	RejectReason string     `json:"rejectReason,omitempty"`
}

// TradeDocEdge represents graph edge between Doc and Trade
type TradeDocEdge struct {
	TradeID     string
	StageIdx    uint
	StageDocIdx uint
}

// TradeDocEdgeDO is a database object for document edge relation
type TradeDocEdgeDO struct {
	FullDocID   string `json:"_from"`
	FullTradeID string `json:"_to"`
	StageIdx    uint   `json:"stageIdx"`
	StageDocIdx uint   `json:"stageDocIdx"`
}

// TradeDocOfferEdgeDO is a database object for document-tradeOffer edge relation
type TradeDocOfferEdgeDO struct {
	FullDocID        string `json:"_from"`
	FullTradeOfferID string `json:"_to"`
}

// TradeTemplate type for collection model of stages for  one trade
type TradeTemplate struct {
	ID          string               `json:"_key,omitempty"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Stages      []TradeStageTemplate `json:"stages"`
}

// TradeStageTemplate type is for template stage model in trade template.
type TradeStageTemplate struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Owner       TradeActor `json:"owner"`
}

// TradeOffer object
type TradeOffer struct {
	ID          string         `json:"_key,omitempty"`
	Price       float64        `json:"price"`
	IsSell      bool           `json:"isSell"`
	PriceType   OfferPriceType `json:"priceType"`
	Currency    Currency       `json:"currency"`
	CreatedBy   string         `json:"createdBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	ExpiresAt   *time.Time     `json:"expiresAt"`
	ClosedAt    *time.Time     `json:"closedAt"`
	OrgID       string         `json:"orgID"`
	IsAnonymous bool           `json:"isAnonymous"`
	Commodity   string         `json:"commodity"`
	ComType     []string       `json:"comType"`
	Quality     string         `json:"quality"`
	Origin      string         `json:"origin"`
	Incoterm    Incoterm       `json:"incoterm"`
	MarketLoc   string         `json:"marketLoc"`
	Vol         int            `json:"vol"`
	Shipment    []time.Time    `json:"shipment"`
	Note        string         `json:"note"`
	DocID       *string        `json:"docID"`
}

// Doc type for trade document
type Doc struct {
	ID        string    `json:"_key,omitempty"`
	Hash      string    `json:"hash"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// Organization type for organization info
type Organization struct {
	ID        string `json:"_key,omitempty"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

// TxLog type for transaction logging into DB
type TxLog struct {
	ID        string             `json:"_key,omitempty"`
	TxStatus  TxStatusEnum       `json:"txStatus"`
	Ledger    LedgerEnum         `json:"ledger"` // Stellar or other
	RawTx     string             `json:"rawTx"`  // Can be any tx representation, XDR too
	CreatedBy string             `json:"createdBy"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Nonce     xdr.SequenceNumber `json:"nonce"`
	Notes     string             `json:"notes"`
	SourceAcc string             `json:"sourceAcc"`
}

// TxLogEdgeDTO represents graph edge between TxLogEntry and Trade
type TxLogEdgeDTO struct {
	FullTxLogID string `json:"_from"`       // LogEntryEdge entity
	FullTradeID string `json:"_to"`         // Trade entity
	StageIdx    *uint  `json:"stageIdx"`    // non-mandatory, omitempty does not work
	StageDocIdx *uint  `json:"stageDocIdx"` // non-mandatory, omitempty does not work
}

// TxLogEdge represents graph edge between TxLogEntry and Trade
// Set up with "naked" ids and call toDTO to get object with full IDs
type TxLogEdge struct {
	TradeID     string
	StageIdx    *uint // non-mandatory
	StageDocIdx *uint // non-mandatory
}

// FileInfo is for uploaded files infomation
type FileInfo struct {
	FileName string `json:"filename"`
	Hash     string `json:"hash"`
	URL      string `json:"url"`
}

// Notification object type
type Notification struct {
	ID          string    `json:"_key,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	TriggeredBy string    `json:"triggeredBy"`
	Receiver    []string  `json:"receiver"`
	Type        NotifType `json:"type"`
	Dismissed   []string  `json:"dismissed"`
	EntityID    string    `json:"entityID"`
	Msg         string    `json:"msg"`
	Action      Approval  `json:"approval"`
}

// TXSourceAcc stores unlock times and secret keys for smart contracts
type TXSourceAcc struct {
	PubKey         SCAddr          `json:"_key"`
	SCSecret       SCSecret        `json:"sKey"`
	Type           TXSourceAccType `json:"type"`
	LockTradeID    string          `json:"lockTradeID"`
	LockUserID     string          `json:"lockUserID"`
	LockExpiresAt  time.Time       `json:"lockExpiresAt"`
	LockUnlockedAt *time.Time      `json:"lockUnlockedAt"`
}
