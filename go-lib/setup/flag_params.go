package setup

import (
	"fmt"
	"net/url"
	"os"

	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/flag"
	bat "github.com/robert-zaremba/go-bat"
)

const cfgNameStellarNetwork = "stellar-network"
const cfgNameSCLockDuration = "tx-source-acc-lock-duration"

// RsaKeyPath is the file path of rsa private key to sign jwt-token
const RsaKeyPath = "/config/app.rsa"

// URLFlag is a structure containing network connection details to a service
// The value should have the following structure:
//     //username:password@host:port/directory
type URLFlag struct {
	url.URL
}

// Set implements github.com/namsral/flag Value interface
func (a *URLFlag) Set(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return errstack.WrapAsReq(err, "Can't parse network address config")
	}
	a.URL = *u
	return nil
}

// PathFlag represents a file in a filesystem
type PathFlag struct {
	Path string
}

// Set implements github.com/namsral/flag Value interface
func (a *PathFlag) Set(filePath string) error {
	a.Path = filePath
	return a.Check()
}

// String implements github.com/namsral/flag Value interface
func (a *PathFlag) String() string {
	return a.Path
}

// Check returns an error if it can't find the file
func (a PathFlag) Check() error {
	if a.Path == "" {
		return errstack.NewReq("File path can't be empty")
	}
	_, err := os.Stat(a.Path)
	return err

}

var rollbarFlag = flag.String("rollbar", "", "rollbar token [required in production env]")
var flagLogColored = flag.Bool("log-colored", true, "Use colored log (good for terminal output)")

func init() {
	var config string
	RootProjectDir, err := bat.FindRoot()
	if err != nil {
		logger.Fatal("Can't set a Root directory", err)
	}
	var configFilename = RootProjectDir + "/config/config.ini"
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		logger.Fatal("config.ini file does not exist", err)
	}
	var rsaKeyFilename = RootProjectDir + RsaKeyPath
	if _, err := os.Stat(rsaKeyFilename); os.IsNotExist(err) {
		logger.Fatal("app.rsa file does not exist", err)
	}
	flag.StringVar(&config, flag.DefaultConfigFlagname, configFilename, "config file")
}

// SrvFlags represents common server flags
type SrvFlags struct {
	Production         *bool
	Port               *string
	StellarNetwork     *string
	SCAddrLockDuration *uint
}

// NewSrvFlags setups common server flags
func NewSrvFlags() SrvFlags {
	return SrvFlags{
		flag.Bool("production", false, "Run in production mode"),
		flag.String("port", "8000", "The HTTP listening port"),
		flag.String(cfgNameStellarNetwork, "", "Stellar network name. Must be one of "+fmt.Sprint(stellar.Networks.Keys())),
		flag.Uint(cfgNameSCLockDuration, 4, "Smart contract address lock time."),
	}
}

// Check validates the flags. It may panic!
func (f SrvFlags) Check() error {
	errb := errstack.NewBuilder()
	validation.NotEmpty(*f.StellarNetwork, errb.Putter(cfgNameStellarNetwork))
	validation.NotEmpty(*f.StellarNetwork, errb.Putter(cfgNameSCLockDuration))
	validation.Positive(*f.SCAddrLockDuration, errb.Putter(cfgNameSCLockDuration))
	return errb.ToReqErr()
}
