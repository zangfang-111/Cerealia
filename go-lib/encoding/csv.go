// Package encoding contains csv encoding functions
package encoding

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/robert-zaremba/errstack"
)

// NewCSVReader creates csv.Reader. It assumes that
// * fields are separated with `,`
// * `#` is used for comments
// * accepts trailing commas
// Furthermore it reads and skips given number of records
func NewCSVReader(r io.Reader, skip int) (*csv.Reader, errstack.E) {
	var reader = csv.NewReader(r)
	reader.Comment = '#'
	reader.TrailingComma = true
	for i := 0; i < skip; i++ {
		if _, err := reader.Read(); err != nil {
			return nil, errstack.WrapAsReq(err, "Can't read the CSV header")
		}
	}
	return reader, nil
}

// NewCSVFileReader opens the `fname` file and calls NewCSVReader.
// It also returns file close function as the second field.
func NewCSVFileReader(fname string, skip int) (*csv.Reader, func() error, errstack.E) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, errstack.WrapAsReq(err, "Can't open csv file: "+fname)
	}
	r, err2 := NewCSVReader(f, skip)
	return r, f.Close, err2
}
