package actions

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// BindJSON unmarshals the Request Body into an struct
func BindJSON(v interface{}, c buffalo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.Wrap(err, "Error while reading request body")
	}

	err = json.Unmarshal([]byte(body), &v)
	if err != nil {
		return errors.Wrap(err, "Error while unmarshaling request body")
	}

	return nil
}

// Wrap wraps an error and a message AND adds the 'WithStack'
func Wrap(err error, msg string) error {
	return errors.WithStack(errors.Wrap(err, msg))
}

// RandomString generates a random string
func RandomString(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789"
	charsLen := len(chars)

	randString := make([]byte, length)

	// add Chars
	for i := 0; i < length; i++ {
		randString[i] = chars[rand.Intn(charsLen)]
	}

	return string(randString)
}
