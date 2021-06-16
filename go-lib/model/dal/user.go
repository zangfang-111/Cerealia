package dal

import (
	"bytes"
	"context"
	"crypto/rand"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/google/uuid"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
	"golang.org/x/crypto/argon2"
)

var logger = log15.Root()

const saltLength = 16

// Input types
type (
	// AuthUserModel Model for authorized user with access token
	AuthUserModel struct {
		User  model.User `json:"user"`
		Token string     `json:"token"`
	}
)

// Generate cryptographically-sound random salt via crypto/rand
func createSalt() ([]byte, errstack.E) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt[:saltLength])
	return salt, errstack.WrapAsInf(err, "Can't create random salt")
}

// Calculate srgon2 scrypt salted hash key from password and salt
func createPwdHash(password string, salt []byte) []byte {
	key := argon2.Key([]byte(password), salt, 3, saltLength*1024, 2, 32)
	return key
}

// InsertUser create new user
func InsertUser(ctx context.Context, db driver.Database, nu *model.NewUserInput) (*model.User, errstack.E) {
	errs := assertEmailNotExist(ctx, db, nu.Email)
	if errs != nil {
		return nil, errs
	}
	if errs = assertOrgExists(ctx, db, nu.OrgID); errs != nil {
		return nil, errs
	}
	u := model.User{
		CreatedAt:     time.Now().UTC(),
		Emails:        []string{nu.Email},
		FirstName:     nu.FirstName,
		LastName:      nu.LastName,
		Biography:     *nu.Biography,
		Avatar:        *nu.Avatar,
		Organizations: map[string]string{nu.OrgID: nu.OrgRole},
		// currently we set trader role for all new users
		Roles:     []model.UserRole{model.UserRoleTrader},
		Approvals: []model.AccessApproval{},
	}
	if u.Salt, errs = createSalt(); errs != nil {
		return nil, errs
	}
	u.Password = createPwdHash(nu.Password, u.Salt)
	walletID, err := uuid.NewRandom()
	if err != nil {
		return nil, errstack.WrapAsInf(err, "Can't generate wallet ID")
	}
	walletIDStr := walletID.String()
	u.StaticWallets = map[string]model.StaticWallet{
		walletIDStr: model.StaticWallet{
			PubKey: nu.PublicKey,
			Wallet: model.Wallet{
				Note: "Single-key wallet created at registration",
			},
		},
	}
	u.DefaultWalletID = walletIDStr
	_, errs = insertHasID(ctx, dbconst.ColUsers, &u, db)
	return &u, errs
}

// UserLogin authenticates and signs in the user
func UserLogin(ctx context.Context, db driver.Database, ul model.UserLoginInput) (*model.User, errstack.E) {
	u, errs := GetUserByEmail(ctx, db, ul.Email)
	var notFound = IsNotFound(errs)
	if errs != nil && !notFound {
		return u, errs
	}
	if len(u.Approvals) == 0 || u.Approvals[len(u.Approvals)-1].Status == model.SimpleApprovalRejected {
		return nil, errstack.NewDomain("Failed to login, Your account should be accepted by Cerealia team")
	}
	if notFound || !bytes.Equal(u.Password, createPwdHash(ul.Password, u.Salt)) {
		return u, errstack.NewReq("Invalid credentials!")
	}
	return u, nil
}

// CreateOrganization creates new Organization and returns new organization ID.
func CreateOrganization(ctx context.Context, db driver.Database,
	org *model.Organization) (*model.Organization, errstack.E) {
	_, errs := insertHasID(ctx, dbconst.ColOrganizations, org, db)
	return org, errs
}

// GetUser fetches User from DB by ID
func GetUser(ctx context.Context, db driver.Database, userID string) (*model.User, errstack.E) {
	var u model.User
	return &u, DBGetOneFromColl(ctx, &u, userID, dbconst.ColUsers, db)
}

// Get2Users calls GetUser twice. Convenient for buyer+seller
// TODO: get many documents at once
func Get2Users(ctx context.Context, db driver.Database,
	userID1 string, userID2 string) (*model.User, *model.User, errstack.E) {
	u1, err := GetUser(ctx, db, userID1)
	if err != nil {
		return nil, nil, err
	}
	u2, err := GetUser(ctx, db, userID2)
	return u1, u2, err
}

// GetApprovedUsers fetches all approved users from db.
func GetApprovedUsers(ctx context.Context, db driver.Database) ([]model.User, errstack.E) {
	q := "FOR d IN users FILTER LAST(d.approvals).status=='approved' RETURN d"
	var us []model.User
	return us, DBQueryMany(ctx, &us, q, nil, db)
}

// GetAdminUsers fetches all users from DB.
func GetAdminUsers(ctx context.Context, db driver.Database) ([]model.AdminUser, errstack.E) {
	q := "FOR d IN users RETURN d"
	var us []model.User
	var aus []model.AdminUser
	errs := DBQueryMany(ctx, &us, q, nil, db)
	if errs != nil {
		return nil, errs
	}
	for _, v := range us {
		aus = append(aus, model.AdminUser{
			User:      &v,
			Approvals: v.Approvals,
		})
	}
	return aus, nil
}

// GetUserByEmail fetches user from DB
func GetUserByEmail(ctx context.Context, db driver.Database, email string) (*model.User, errstack.E) {
	var u = new(model.User)
	query := "FOR d IN users FILTER @userEmail IN d.emails RETURN d"
	bindVars := map[string]interface{}{
		"userEmail": email,
	}
	return u, DBQueryOne(ctx, u, query, bindVars, db)
}

// GetOrganization fetches Organization from DB by its ID
func GetOrganization(ctx context.Context, db driver.Database, id string) (*model.Organization, errstack.E) {
	var o model.Organization
	return &o, DBGetOneFromColl(ctx, &o, id, dbconst.ColOrganizations, db)
}

// GetAllOrganizations fetches all organizations
func GetAllOrganizations(ctx context.Context, db driver.Database) ([]model.Organization, errstack.E) {
	q := "FOR d IN organizations sort d.name RETURN d"
	var orgs []model.Organization
	return orgs, DBQueryMany(ctx, &orgs, q, nil, db)
}

// GetOrgMap builds new OrgMap array of a user from his orgmap(map[string]string)
func GetOrgMap(ctx context.Context, db driver.Database, userID string) ([]model.UserOrgMap, errstack.E) {
	var orgMaps []model.UserOrgMap
	query := `let orgs = (for d in users filter d._key==@userID return d.organizations)[0]
		let keys = attributes(orgs, false, true)
		for key in keys for d in organizations filter d._key==key return {'Org':d,'Role':orgs[key]}`
	bindVars := map[string]interface{}{
		"userID": userID,
	}
	err := DBQueryMany(ctx, &orgMaps, query, bindVars, db)
	return orgMaps, err
}

// DeleteUser removes the user by its email
func DeleteUser(ctx context.Context, db driver.Database, email string) errstack.E {
	query := "FOR d IN users FILTER @userEmail IN d.emails REMOVE d IN users"
	bindVars := map[string]interface{}{
		"userEmail": email,
	}
	c, err := db.Query(ctx, query, bindVars)
	errstack.CallAndLog(logger, c.Close)
	return errstack.WrapAsInf(err, "Failed to remove user")
}

// DeleteOrg removes the Organization by its id
func DeleteOrg(ctx context.Context, db driver.Database, id string) errstack.E {
	return deleteDoc(ctx, db, dbconst.ColOrganizations, id)
}

// UpdateAvatar updates user avatar
func UpdateAvatar(ctx context.Context, db driver.Database, uid string, avatar string) errstack.E {
	query := "FOR u IN users UPDATE { _key: @uid, avatar: @avatar } In users"
	bindVars := map[string]interface{}{
		"uid":    uid,
		"avatar": avatar,
	}
	c, err := db.Query(ctx, query, bindVars)
	errstack.CallAndLog(logger, c.Close)
	return errstack.WrapAsInf(err, "Failed to update user avatar")
}

// ChangePassword updates the original user password to new password
func ChangePassword(ctx context.Context, db driver.Database, u *model.User,
	input model.ChangePasswordInput) errstack.E {
	if !bytes.Equal(u.Password, createPwdHash(input.OldPassword, u.Salt)) {
		return errstack.NewReq("Your old password is invalid")
	}
	u.Password = createPwdHash(input.NewPassword, u.Salt)
	_, err := UpdateDoc(ctx, db, dbconst.ColUsers, u.ID, u)
	return errstack.WrapAsInf(err, "Failed to update password")
}

// ChangeEmail updates the user emails
func ChangeEmail(ctx context.Context, db driver.Database, u *model.User, emails []string) errstack.E {
	// We don't need to check if email exists because it will be handled by DB.
	// Arango fails if we insert a user with same emails so let's remove them. This is a bug in ArangoDB: https://trello.com/c/Ms2SbXK7/667-bring-back-the-emails-unique-index

	q := `RETURN flatten(FOR u in users FILTER u._key != @uid RETURN intersection(u.emails, @emails))`
	bindVars := map[string]interface{}{"emails": emails, "uid": u.ID}
	var dupEmails = []string{}
	if err := DBQueryMany(ctx, &dupEmails, q, bindVars, db); err != nil {
		return err
	}
	if len(dupEmails) > 0 {
		return errstack.NewReqF("%v  emails are already used", dupEmails)
	}
	diff := map[string][]string{"emails": emails}
	_, err := UpdateDoc(ctx, db, dbconst.ColUsers, u.ID, diff)
	if err == nil {
		u.Emails = emails
	}
	return err
}

// UpdateUserProfile updates user profile
func UpdateUserProfile(ctx context.Context, db driver.Database, u *model.User, input model.UserProfileInput) (*model.User, errstack.E) {
	var errs errstack.E
	u.FirstName = input.FirstName
	u.LastName = input.LastName
	m := make(map[string]string)
	for _, v := range input.OrgMap {
		if errs = assertOrgExists(ctx, db, v.ID); errs != nil {
			return u, errs
		}
		m[v.ID] = v.Role
	}
	u.Organizations = m
	u.Biography = input.Biography
	return u, ReplaceUser(ctx, db, u)
}

// ReplaceUser replaces user object in DB
func ReplaceUser(ctx context.Context, db driver.Database, u *model.User) errstack.E {
	_, err := replaceDoc(ctx, db, dbconst.ColUsers, u.ID, u)
	return err
}

// assertEmailNotExist returns error if an email exists.
// It will ignore emails emails from `excludedUser` if it is provided (not empty).
func assertEmailNotExist(ctx context.Context, db driver.Database, email string) errstack.E {
	bindVars := map[string]interface{}{"email": email}
	query := existsQuery("FOR d IN users FILTER @email IN d.emails")
	var exists bool
	err := DBQueryOne(ctx, &exists, query, bindVars, db)
	if err == nil && exists {
		err = errstack.NewReq("User with email [" + email + "] already exists")
	}
	return err
}

func assertOrgExists(ctx context.Context, db driver.Database, orgID string) errstack.E {
	var exists bool
	query := existsQuery(`FOR d IN organizations FILTER d._key==@orgID`)
	bindVars := map[string]interface{}{
		"orgID": orgID,
	}
	err := DBQueryOne(ctx, &exists, query, bindVars, db)
	if err == nil && !exists {
		err = errstack.NewReq("Organization with id [" + orgID + "] doesn't exists")
	}
	return err
}
