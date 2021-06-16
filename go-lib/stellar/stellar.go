// Package stellar defines all stellar functions and builders
package stellar

import (
	"net/http"

	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
)

// Network holds network info for stellar platform
type Network struct {
	Name       string
	Passphrase build.Network
	URL        string
}

type networksT map[string]Network

// Networks are list of networks
var Networks = networksT{}

func init() {
	ns := []Network{
		// Using build.TestNetwork network passphrase for testnet because tests load it both: from config and from noop
		Network{fakeNetName, build.TestNetwork, fakeNetName},
		Network{"horizon-test", build.TestNetwork, "https://horizon-testnet.stellar.org"},
		Network{"horizon-main", build.PublicNetwork, "https://horizon.stellar.org"},
	}
	for i := range ns {
		Networks[ns[i].Name] = ns[i]
	}
}

func (ns networksT) Keys() []string {
	var keys []string
	for k := range ns {
		keys = append(keys, k)
	}
	return keys
}

func getNetworkAndClient(netName string) (Network, Client, errstack.E) {
	n, ok := Networks[netName]
	if !ok {
		return n, nil, errstack.NewReq("Unknown Stellar network name: " + netName)
	}
	if netName == fakeNetName {
		return n, &NoopClient{}, nil
	}
	return n, &horizon.Client{
		HTTP: http.DefaultClient,
		URL:  n.URL}, nil
}

// Client is interface for stellar clients that have only one function
type Client interface {
	SubmitTransaction(tx string) (hProtocol.TransactionSuccess, error)
	SequenceForAccount(accountID string) (xdr.SequenceNumber, error)
}
