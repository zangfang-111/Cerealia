// Package dbconst contains all constant collection names of database and related functions
package dbconst

// Col is a type for collection name
type Col string

// FullID creates a collection document full id
func (d Col) FullID(key string) string {
	return string(d) + "/" + key
}

// List of ArangoDB collections
const (
	GraphDoc              Col = "graph_doc_to_trade"
	GraphTxLog            Col = "graph_txlog_to_trade"
	GraphDocTradeOffer    Col = "graph_doc_to_tradeoffer"
	ColTrades             Col = "trades"
	ColUsers              Col = "users"
	ColDocs               Col = "docs"
	ColDocEdges           Col = "doc_edges"
	ColTradeTemplates     Col = "trade_templates"
	ColOrganizations      Col = "organizations"
	ColTxEntryLog         Col = "tx_entry_log"
	ColTxEntryLogEdges    Col = "tx_entry_log_edges"
	ColTradeOffers        Col = "trade_offers"
	ColDocTradeOfferEdges Col = "doc_tradeoffer_edges"
	ColNotifications      Col = "notifications"
	ColTxSourceAccs       Col = "tx_source_accounts"
)
