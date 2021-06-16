package main

import (
	"bitbucket.org/cerealia/apps/go-lib/setup"
	"github.com/robert-zaremba/flag"
)

// CleanFlags is stellar account cleaner config type
type CleanFlags struct {
	setup.SrvFlags
	DeleteAccountOfTrade *string
	FundDestinationAddr  *string
	AdditionalSigners    arrayFlags
	DismantleAllTrades   *bool
}

type arrayFlags []string

const delTradeAccConfigKey = "delete-trade-account"
const fundDestinationAddr = "fund-destination-addr"
const signer = "signer"
const cleanAllTrades = "clean-all-trades"

func init() {
	flag.Var(&F.AdditionalSigners, signer, "Additional signers for trade account removal. Multiple of them are accepted: --signer <key1> --signer <key2>")
}

// String implements github.com/namsral/flag Value interface
func (i *arrayFlags) String() string {
	return "arrayFlags"
}

// Set implements github.com/namsral/flag Value interface
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var isParsed = false

// F stores command line flags
var F = CleanFlags{
	SrvFlags:             setup.NewSrvFlags(),
	DeleteAccountOfTrade: flag.String(delTradeAccConfigKey, "", "Trade ID to delete stellar account of and reclaim it's XLMs from (reclamation account can be specified using 'stellar-master-secret' key)"),
	FundDestinationAddr: flag.String(fundDestinationAddr, "Fund destination address is not defined",
		"Destination address for funds of dismantled accounts"),
	AdditionalSigners:  arrayFlags{}, // refer to init method
	DismantleAllTrades: flag.Bool(cleanAllTrades, false, "Will query all trades and remove their SCs one by one"),
}
