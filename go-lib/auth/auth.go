// Package auth contains authentication functions
package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/setup"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/robert-zaremba/errstack"
)

// ExpireTime means the duration of expires time, its unit is minute.
// default expires time is 60min
const ExpireTime = 200

var signKey *rsa.PrivateKey

// AppClaims provides custom claim for JWT
type AppClaims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

// InitKeys Read the key files before starting http handlers
func InitKeys() errstack.E {
	signBytes, err := ioutil.ReadFile(setup.RootDir + setup.RsaKeyPath)
	if err != nil {
		return errstack.WrapAsInf(err, "Can not read private key from config data")
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	return errstack.WrapAsInf(err, "Can not parse private key")
}

// CreateJWT generates a new JWT token
func CreateJWT(userID string) (string, errstack.E) {
	errs := InitKeys()
	if errs != nil {
		return "", errs
	}
	// Create the Claims
	claims := AppClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * ExpireTime).Unix(),
			Issuer:    "admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(signKey)
	return ss, errstack.WrapAsDomain(err, "Can not generate JWT token")
}

// Authorize Middleware for validating JWT tokens
func Authorize(tokenStr string) (string, errstack.E) {
	// init config
	errs := InitKeys()
	if errs != nil {
		return "", errs
	}
	token, err := jwt.ParseWithClaims(tokenStr,
		&AppClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return &signKey.PublicKey, nil
		})
	if err != nil {
		switch vErr := err.(type) {
		case *jwt.ValidationError:
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return "", errstack.WrapAsReq(err, "JWT Access Token is expired, get a new Token")
			default:
				return "", errstack.WrapAsReq(err, "JWT Validation Error")
			}

		default:
			return "", errstack.WrapAsReq(err, "Error while parsing the JWT Access Token!")
		}

	}
	if !token.Valid {
		return "", errstack.NewReq("Invalid JWT token!")
	}
	// Set user name to HTTP context
	// context.Set(c.Request, "user", token.Claims.(*AppClaims).UserID)
	return token.Claims.(*AppClaims).UserID, nil
}

// TokenFromAuthHeader is a "TokenExtractor" that takes a given request and extracts
// the JWT token from the Authorization header.
func TokenFromAuthHeader(r *http.Request) (string, errstack.E) {
	// Look for an Authorization header
	if ah := r.Header.Get("Authorization"); ah != "" {
		// Should be a bearer token
		if len(ah) > 6 && strings.ToUpper(ah[:6]) == "BEARER" {
			return ah[7:], nil
		}
	}
	return "", errstack.NewReq("No token in the HTTP request")
}
