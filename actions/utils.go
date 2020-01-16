package actions

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// BindJSON unmarshals the Request Body into an struct
func BindJSON(v interface{}, c *buffalo.Context) error {
	body, err := ioutil.ReadAll((*c).Request().Body)
	if err != nil {
		return errors.Wrap(err, "Error while reading request body")
	}

	err = json.Unmarshal([]byte(body), &v)
	if err != nil {
		return errors.Wrap(err, "Error while unmarshaling request body")
	}

	return nil
}
