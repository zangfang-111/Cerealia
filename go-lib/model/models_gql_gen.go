// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// AdminUser; User with special Admin information
type AdminUser struct {
	User      *User            `json:"user"`
	Approvals []AccessApproval `json:"approvals"`
}

// AuthUser; after user login, backend sends user token
type AuthUser struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Password change data
type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// New trade fields
type NewStageInput struct {
	Tid         string `json:"tid"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
}

// Trade creation fields
type NewTradeInput struct {
	TemplateID   string  `json:"templateID"`
	Name         string  `json:"name"`
	SellerID     string  `json:"sellerID"`
	BuyerID      string  `json:"buyerID"`
	Description  *string `json:"description"`
	TradeOfferID *string `json:"tradeOfferID"`
}

// New user input data
type NewUserInput struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Avatar    *string `json:"avatar"`
	PublicKey string  `json:"publicKey"`
	Biography *string `json:"biography"`
	OrgID     string  `json:"orgID"`
	OrgRole   string  `json:"orgRole"`
}

// Organization input data
type OrgInput struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

// StellarNet; stellar network infomation, we provide it to frontend
type StellarNet struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Passphrase string `json:"passphrase"`
}

// TradeActorWallet; data about an actor in a trade
type TradeActorWallet struct {
	PubKey   string `json:"pubKey"`
	KeyPath  string `json:"keyPath"`
	WalletID string `json:"walletID"`
}

// trade offer input data
type TradeOfferInput struct {
	Price       float64        `json:"price"`
	PriceType   OfferPriceType `json:"priceType"`
	IsSell      bool           `json:"isSell"`
	Currency    Currency       `json:"currency"`
	ExpiresAt   *time.Time     `json:"expiresAt"`
	Commodity   string         `json:"commodity"`
	ComType     []string       `json:"comType"`
	Quality     string         `json:"quality"`
	OrgID       string         `json:"orgID"`
	Origin      string         `json:"origin"`
	IsAnonymous bool           `json:"isAnonymous"`
	Incoterm    Incoterm       `json:"incoterm"`
	MarketLoc   string         `json:"marketLoc"`
	Vol         int            `json:"vol"`
	Shipment    []time.Time    `json:"shipment"`
	Note        string         `json:"note"`
	DocID       *string        `json:"docID"`
}

// Context of a document
type TradeStageDocPath struct {
	Tid          string `json:"tid"`
	StageIdx     uint   `json:"stageIdx"`
	StageDocIdx  uint   `json:"stageDocIdx"`
	StageDocHash string `json:"stageDocHash"`
}

// Context of a trade stage
type TradeStagePath struct {
	Tid      string `json:"tid"`
	StageIdx uint   `json:"stageIdx"`
}

// User data for logging in
type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User organization with role
type UserOrgMap struct {
	Org  Organization `json:"org"`
	Role string       `json:"role"`
}

// User orgMap input
type UserOrgMapInput struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

// User Profile update data
type UserProfileInput struct {
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	OrgMap    []UserOrgMapInput `json:"orgMap"`
	Biography string            `json:"biography"`
}

// ApprovalStatusOfATradeOrAStage
type Approval string

const (
	ApprovalNil       Approval = "nil"
	ApprovalPending   Approval = "pending"
	ApprovalExpired   Approval = "expired"
	ApprovalRejected  Approval = "rejected"
	ApprovalApproved  Approval = "approved"
	ApprovalSubmitted Approval = "submitted"
)

var AllApproval = []Approval{
	ApprovalNil,
	ApprovalPending,
	ApprovalExpired,
	ApprovalRejected,
	ApprovalApproved,
	ApprovalSubmitted,
}

func (e Approval) IsValid() bool {
	switch e {
	case ApprovalNil, ApprovalPending, ApprovalExpired, ApprovalRejected, ApprovalApproved, ApprovalSubmitted:
		return true
	}
	return false
}

func (e Approval) String() string {
	return string(e)
}

func (e *Approval) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Approval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Approval", str)
	}
	return nil
}

func (e Approval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Currency string

const (
	CurrencyUsd Currency = "USD"
	CurrencyRub Currency = "RUB"
	CurrencyEur Currency = "EUR"
	CurrencyTry Currency = "TRY"
	CurrencyUah Currency = "UAH"
	CurrencyCny Currency = "CNY"
)

var AllCurrency = []Currency{
	CurrencyUsd,
	CurrencyRub,
	CurrencyEur,
	CurrencyTry,
	CurrencyUah,
	CurrencyCny,
}

func (e Currency) IsValid() bool {
	switch e {
	case CurrencyUsd, CurrencyRub, CurrencyEur, CurrencyTry, CurrencyUah, CurrencyCny:
		return true
	}
	return false
}

func (e Currency) String() string {
	return string(e)
}

func (e *Currency) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Currency(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Currency", str)
	}
	return nil
}

func (e Currency) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ModeratingStatusOfATradeOrAStage
type DoneStatus string

const (
	DoneStatusNil   DoneStatus = "nil"
	DoneStatusDoing DoneStatus = "doing"
	DoneStatusDone  DoneStatus = "done"
)

var AllDoneStatus = []DoneStatus{
	DoneStatusNil,
	DoneStatusDoing,
	DoneStatusDone,
}

func (e DoneStatus) IsValid() bool {
	switch e {
	case DoneStatusNil, DoneStatusDoing, DoneStatusDone:
		return true
	}
	return false
}

func (e DoneStatus) String() string {
	return string(e)
}

func (e *DoneStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DoneStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DoneStatus", str)
	}
	return nil
}

func (e DoneStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// TradingIncoterm
type Incoterm string

const (
	IncotermCaf   Incoterm = "CAF"
	IncotermCfr   Incoterm = "CFR"
	IncotermCfrc  Incoterm = "CFRC"
	IncotermCif   Incoterm = "CIF"
	IncotermCiffo Incoterm = "CIFFO"
	IncotermCnf   Incoterm = "CNF"
	IncotermCnffo Incoterm = "CNFFO"
	IncotermDel   Incoterm = "DEL"
	IncotermFob   Incoterm = "FOB"
)

var AllIncoterm = []Incoterm{
	IncotermCaf,
	IncotermCfr,
	IncotermCfrc,
	IncotermCif,
	IncotermCiffo,
	IncotermCnf,
	IncotermCnffo,
	IncotermDel,
	IncotermFob,
}

func (e Incoterm) IsValid() bool {
	switch e {
	case IncotermCaf, IncotermCfr, IncotermCfrc, IncotermCif, IncotermCiffo, IncotermCnf, IncotermCnffo, IncotermDel, IncotermFob:
		return true
	}
	return false
}

func (e Incoterm) String() string {
	return string(e)
}

func (e *Incoterm) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Incoterm(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Incoterm", str)
	}
	return nil
}

func (e Incoterm) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type NotifType string

const (
	// triggered when a user perform some action and this triggers a notification for other users
	NotifTypeAction NotifType = "action"
	// triggered by a system when there is something to be do (eg documenet approval is close to expire)
	NotifTypeAlert NotifType = "alert"
)

var AllNotifType = []NotifType{
	NotifTypeAction,
	NotifTypeAlert,
}

func (e NotifType) IsValid() bool {
	switch e {
	case NotifTypeAction, NotifTypeAlert:
		return true
	}
	return false
}

func (e NotifType) String() string {
	return string(e)
}

func (e *NotifType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NotifType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NotifType", str)
	}
	return nil
}

func (e NotifType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IsItAFirmOfferOrJustAQuote
type OfferPriceType string

const (
	OfferPriceTypeFirm  OfferPriceType = "firm"
	OfferPriceTypeQuote OfferPriceType = "quote"
)

var AllOfferPriceType = []OfferPriceType{
	OfferPriceTypeFirm,
	OfferPriceTypeQuote,
}

func (e OfferPriceType) IsValid() bool {
	switch e {
	case OfferPriceTypeFirm, OfferPriceTypeQuote:
		return true
	}
	return false
}

func (e OfferPriceType) String() string {
	return string(e)
}

func (e *OfferPriceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OfferPriceType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OfferPriceType", str)
	}
	return nil
}

func (e OfferPriceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// SimpleApprovalIsABasicStatusForApprovals
type SimpleApproval string

const (
	SimpleApprovalRejected SimpleApproval = "rejected"
	SimpleApprovalApproved SimpleApproval = "approved"
)

var AllSimpleApproval = []SimpleApproval{
	SimpleApprovalRejected,
	SimpleApprovalApproved,
}

func (e SimpleApproval) IsValid() bool {
	switch e {
	case SimpleApprovalRejected, SimpleApprovalApproved:
		return true
	}
	return false
}

func (e SimpleApproval) String() string {
	return string(e)
}

func (e *SimpleApproval) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SimpleApproval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SimpleApproval", str)
	}
	return nil
}

func (e SimpleApproval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// TypeOfUserThatParticipatesInTrade
type TradeActor string

const (
	// No Owner
	TradeActorN TradeActor = "n"
	// Buyer
	TradeActorB TradeActor = "b"
	// Seller
	TradeActorS TradeActor = "s"
	// Moderator
	TradeActorM TradeActor = "m"
)

var AllTradeActor = []TradeActor{
	TradeActorN,
	TradeActorB,
	TradeActorS,
	TradeActorM,
}

func (e TradeActor) IsValid() bool {
	switch e {
	case TradeActorN, TradeActorB, TradeActorS, TradeActorM:
		return true
	}
	return false
}

func (e TradeActor) String() string {
	return string(e)
}

func (e *TradeActor) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TradeActor(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TradeActor", str)
	}
	return nil
}

func (e TradeActor) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// TxTradeEntityEnumsForTransactionBuild
type TxTradeEntity string

const (
	TxTradeEntityStageDoc       TxTradeEntity = "stage_doc"
	TxTradeEntityStageCloseReqs TxTradeEntity = "stage_closeReqs"
	TxTradeEntityStageAdd       TxTradeEntity = "stage_add"
	TxTradeEntityTradeCloseReqs TxTradeEntity = "trade_closeReqs"
)

var AllTxTradeEntity = []TxTradeEntity{
	TxTradeEntityStageDoc,
	TxTradeEntityStageCloseReqs,
	TxTradeEntityStageAdd,
	TxTradeEntityTradeCloseReqs,
}

func (e TxTradeEntity) IsValid() bool {
	switch e {
	case TxTradeEntityStageDoc, TxTradeEntityStageCloseReqs, TxTradeEntityStageAdd, TxTradeEntityTradeCloseReqs:
		return true
	}
	return false
}

func (e TxTradeEntity) String() string {
	return string(e)
}

func (e *TxTradeEntity) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TxTradeEntity(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TxTradeEntity", str)
	}
	return nil
}

func (e TxTradeEntity) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AvailableUserRoleOptionsForAuthorization
type UserRole string

const (
	// Trader is the participant of a trade. Can act on trade.
	UserRoleTrader UserRole = "trader"
	// Moderator has supervisor privileges for trade
	UserRoleModerator UserRole = "moderator"
)

var AllUserRole = []UserRole{
	UserRoleTrader,
	UserRoleModerator,
}

func (e UserRole) IsValid() bool {
	switch e {
	case UserRoleTrader, UserRoleModerator:
		return true
	}
	return false
}

func (e UserRole) String() string {
	return string(e)
}

func (e *UserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserRole", str)
	}
	return nil
}

func (e UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}