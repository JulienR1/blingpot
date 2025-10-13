package query

import (
	"cmp"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/dtos"
)

var ErrExpectedLessThan = errors.New("expected value to be smaller")

func Integer(r *http.Request, key string, out *int) (err error) {
	str := r.URL.Query().Get(key)
	*out, err = strconv.Atoi(str)
	return err
}

func UnixTime(r *http.Request, key string, out *dtos.UnixTime) error {
	assert.Assert(out != nil, "query.UnixTime: Invalid out pointer")

	var milliseconds int
	if err := Integer(r, key, &milliseconds); err != nil {
		return err
	}

	*out = dtos.NewUnixTime(int64(milliseconds))
	return nil
}

func Less[T cmp.Ordered](a, b T) error {
	if a < b {
		return nil
	}
	return ErrExpectedLessThan
}
