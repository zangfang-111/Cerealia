package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	time "time"

	utils "bitbucket.org/cerealia/apps/go-lib/utils"
	"github.com/99designs/gqlgen/graphql"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
)

func wrapStr(s string) string {
	return bat.StrJoin(s, "\"", "\"")
}

// MarshalID marshals the ID to string.
func MarshalID(b string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bat.UnsafeStrToByteArray(wrapStr(b)))
	})
}

// UnmarshalID check the ID type with Regexp.
func UnmarshalID(v interface{}) (string, error) {
	IDStr, ok := v.(string)
	if !ok || IDStr == "" {
		return "", errstack.NewReq("ID must be string and not empty")
	}
	if utils.ReArangoID.MatchString(IDStr) {
		return IDStr, nil
	}
	return IDStr, errstack.NewReq("The Id type is not correct!")
}

// MarshalEmail marshals the email to string
func MarshalEmail(b string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bat.UnsafeStrToByteArray(wrapStr(b)))
	})
}

// UnmarshalEmail check the Email type with Regexp.
func UnmarshalEmail(v interface{}) (string, error) {
	email, ok := v.(string)
	if !ok || email == "" {
		return "", errstack.NewReq("Email must be string and not empty")
	}
	email = strings.TrimSpace(email)
	if utils.ReEmail.MatchString(email) {
		return email, nil
	}
	return email, errstack.NewReq("Malformed email!")
}

// MarshalTime marshals the time to int
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bat.UnsafeStrToByteArray(wrapStr(t.Format(time.RFC3339))))
	})
}

// UnmarshalTime check the time type
func UnmarshalTime(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		return time.Parse(time.RFC3339, tmpStr)
	}
	return time.Time{}, errstack.NewReq("time should be correct format. correct time format: 2006-01-02T15:04:05Z07:00")
}

// MarshalTelephone marshals the telephone to string
func MarshalTelephone(b string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bat.UnsafeStrToByteArray(wrapStr(b)))
	})
}

// UnmarshalTelephone check the telephone type with Regexp.
func UnmarshalTelephone(v interface{}) (string, error) {
	telStr, ok := v.(string)
	if !ok || telStr == "" {
		return "", errstack.NewReq("phone number should be string and not empty")
	}
	telStr = utils.ReSpace.ReplaceAllString(telStr, "") //remove all spaces
	if utils.RePhone.MatchString(telStr) {
		return telStr, nil
	}
	return telStr, errstack.NewReq("The Telephone type is not correct!")
}

// MarshalHash marshals the hash to string
func MarshalHash(b string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bat.UnsafeStrToByteArray(wrapStr(b)))
	})
}

// UnmarshalHash check the hash type with Regexp.
func UnmarshalHash(v interface{}) (string, error) {
	hashStr, ok := v.(string)
	if !ok || hashStr == "" {
		return "", errstack.NewReq("hash should be string and not empty")
	}
	if utils.ReHex.MatchString(hashStr) {
		return hashStr, nil
	}
	return hashStr, errstack.NewReq("The Hash type is not correct!")
}

// MarshalUint marshals the unit to int
func MarshalUint(b uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, b)
	})
}

// UnmarshalUint check the uint type with Regexp.
func UnmarshalUint(v interface{}) (uint, error) {
	i, ok := v.(json.Number)
	if !ok {
		return 0, errstack.NewReq("Uint should be uint type and not empty")
	}
	u, err := bat.Atoui64(string(i))
	return uint(u), err
}
