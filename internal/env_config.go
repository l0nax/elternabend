package internal

import (
	"strconv"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

func init() {
	// intialize global Config (env) Variables
	getGlobalConf()
}

// ======================================================
//		Global Configuration Variables
// ======================================================
var TEACHER_PASSWORD_LENGTH int

var PASSWORD_HASH_MEMORY int
var PASSWORD_HASH_ITERATIONS int
var PASSWORD_HASH_THREADS int
var PASSWORD_HASH_SALT_LEN int
var PASSWORD_HASH_KEY_LEN int

// getGlobalConf initializes all global (env) Variables.
// Panics if there was an error
func getGlobalConf() {
	intConv(&TEACHER_PASSWORD_LENGTH, "TEACHER_PASSWORD_LENGTH", "6",
		"Error while converting 'TEACHER_PASSWORD_LENGTH' to int")

	intConv(&PASSWORD_HASH_MEMORY, "PASSWORD_HASH_MEMORY", "512",
		"Error while converting 'PASSWORD_HASH_MEMORY' to int")

	intConv(&PASSWORD_HASH_ITERATIONS, "PASSWORD_HASH_ITERATIONS", "4",
		"Error while converting 'PASSWORD_HASH_ITERATIONS' to int")

	intConv(&PASSWORD_HASH_THREADS, "PASSWORD_HASH_THREADS", "4",
		"Error while converting 'PASSWORD_HASH_THREADS' to int")

	intConv(&PASSWORD_HASH_SALT_LEN, "PASSWORD_HASH_SALT_LEN", "16",
		"Error while converting 'PASSWORD_HASH_SALT_LEN' to int")

	intConv(&PASSWORD_HASH_KEY_LEN, "PASSWORD_HASH_KEY_LEN", "32",
		"Error while converting 'PASSWORD_HASH_KEY_LEN' to int")
}

// intConv converts an environment variable to an integer and puts it into @res.
// panics with error message @errMsg if conversion was not possible
func intConv(res *int, name string, defaultVal string, errMsg string) {
	var err error

	(*res), err = strconv.Atoi(envy.Get(name, defaultVal))
	if err != nil {
		panic(Wrap(err, errMsg))
	}
}

// Wrap wraps an error and a message AND adds the 'WithStack'
func Wrap(err error, msg string) error {
	return errors.WithStack(errors.Wrap(err, msg))
}
