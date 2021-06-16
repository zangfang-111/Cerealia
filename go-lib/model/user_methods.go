package model

import (
	"regexp"
	"strings"

	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/robert-zaremba/errstack"
)

const shortLengthMessage = "Password length must be at least 8 characters"
const characterPwdMessage = "Password must include at least one small and one capital letter, and one sign"

// SetID implements dal.HasID interface
func (u *User) SetID(id string) {
	u.ID = id
}

// FindWallet returns a wallet from the user if it exists
func (u *User) FindWallet(walletID string) (KeyWallet, errstack.E) {
	hd, ok := u.HDCerealiaWallets[walletID]
	if ok {
		return &hd, nil
	}
	static, ok := u.StaticWallets[walletID]
	if ok {
		return &static, nil
	}
	return nil, errstack.NewDomainF("User '%s' does not have a wallet with ID: '%s'", u.ID, walletID)
}

// DefaultWallet returns a default public key wallet of the user
func (u *User) DefaultWallet() (KeyWallet, errstack.E) {
	return u.FindWallet(u.DefaultWalletID)
}

// IsModerator returns true if user has moderator role.
func (u *User) IsModerator() bool {
	for _, role := range u.Roles {
		if role == UserRoleModerator {
			return true
		}
	}
	return false
}

// CleanAndValidate validates the new userinput data and changes the nil value
func (nu *NewUserInput) CleanAndValidate() errstack.E {
	assignEmptyStringAndTrim(&nu.Avatar)
	assignEmptyStringAndTrim(&nu.Biography)
	nu.FirstName = strings.TrimSpace(nu.FirstName)
	nu.LastName = strings.TrimSpace(nu.LastName)
	nu.Email = strings.TrimSpace(nu.Email)

	vb := validation.Builder{}
	vb.Required("FirstName", nu.FirstName)
	vb.Required("LastName", nu.LastName)
	vb.Required("Email", nu.Email)
	vb.Required("Password", nu.Password)
	vb.Required("OrgID", nu.OrgID)
	vb.Required("OrgRole", nu.OrgRole)
	errb := vb.ToErrstackBuilder()
	if errb != nil {
		return errb.ToReqErr()
	}
	return checkPassword(nu.Password)
}

func assignEmptyStringAndTrim(s **string) {
	if *s == nil {
		*s = new(string)
	} else {
		**s = strings.TrimSpace(**s)
	}
}

// SetID implements dal.HasID interface
func (o *Organization) SetID(id string) {
	o.ID = id
}

func checkPassword(password string) errstack.E {
	if len(password) < 8 {
		return errstack.NewReq(shortLengthMessage)
	}
	lowerCond := regexp.MustCompile(`[a-z]`)
	upperCond := regexp.MustCompile(`[A-Z]`)
	specialCond := regexp.MustCompile(`[!@#$%^&*,.;'"()]`)
	if !lowerCond.MatchString(password) ||
		!upperCond.MatchString(password) ||
		!specialCond.MatchString(password) {
		return errstack.NewReq(characterPwdMessage)
	}
	return nil
}

// ValidateNewOrg validates new organization input
func (o *Organization) ValidateNewOrg() errstack.E {
	vb := validation.Builder{}
	vb.Required("Name", o.Name)
	vb.Required("Address", o.Address)
	vb.Required("Email", o.Email)
	vb.Required("Telephone", o.Telephone)
	errb := vb.ToErrstackBuilder()
	if len(o.Name) < 2 {
		errb.Put("Name length", errstack.NewReq("Organization name must be at least 2 characters long"))
	}
	if errb != nil {
		return errb.ToReqErr()
	}
	return nil
}
