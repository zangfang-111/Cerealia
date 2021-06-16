package trades

import (
	"encoding/json"

	routing "github.com/go-ozzo/ozzo-routing"
)

func respondWithJSON(c *routing.Context, marshallable interface{}) error {
	bytes, err := json.Marshal(marshallable)
	if err != nil {
		return err
	}
	return c.Write(string(bytes))
}
