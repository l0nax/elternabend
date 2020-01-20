package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	suuid "github.com/google/uuid"
	"github.com/l0nax/elternabend/internal"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id" rw:"r"`
	Username     string     `json:"username" db:"username"`
	PasswordHash string     `json:"-" db:"password_hash"`
	ClassTeacher bool       `json:"-" db:"class_teacher"` // check if only rw: r
	Roles        string     `json:"-" db:"roles"`         // a comma seperated list of all roles
	SessionID    suuid.UUID `json:"session_id" db:"session_id"`
	Password     string     `json:"password" db:"-"` // this field is for (eg) the Users/Edit or Users/Create site
	Email        string     `json:"email" db:"email"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Destroy deletes a User from the Database
// NOTE: This function does NOT check if any relations exists!
func (u *User) Destroy(tx *pop.Connection) error {
	// reject deleting Admin Users
	if len(u.Roles) != 0 {
		// check if one of the roles is 'admin'
		roles := strings.Split(u.Roles, ",")
		isAdmin := false

		for _, v := range roles {
			if v == "admin" {
				isAdmin = true
				break
			}
		}

		if isAdmin {
			log.Printf("[ERROR] Can not delete Users which have 'admin' role! (%s)\n",
				u.Username)
			return errors.New("Can not delete Administrators!")
		}
	}

	log.Printf("[INFO] Deleting User '%s'\n", u.Username)

	return tx.Destroy(u)
}

// Create creates a new User and adds it to the Database
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	// If the mail is not providet than use the Username as mail-address
	if len(u.Email) == 0 {
		u.Email = u.Username
	}

	// check set'ed roles
	if len(u.Roles) == 0 {
		// set – as default – 'class_teacher' as role
		u.Roles = "class_teacher"
	}

	log.Printf("Hashing Password")

	// hash Password
	paramTmp := &hashParams{}
	pwdHash, err := hashPassword(u.Password, &paramTmp)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(errors.Wrap(err, "Error while hashing Password"))
	}

	log.Printf("Password hashed: '%s'", string(pwdHash))

	u.PasswordHash = pwdHash
	return tx.ValidateAndCreate(u)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		// &UsernameNotTaken{Name: "Username", Field: u.Username, tx: tx},
		// &EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
	// &UsernameNotTaken{Name: "Username", Field: u.Username, tx: tx},
	// &EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ==========> Custom Validators <==========

type UsernameNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

func (v *UsernameNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("username = ?", v.Field)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil {
		// found a user with same username
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("The username %s is not available.", v.Field))
	}
}

type EmailNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

// IsValid performs the validation check for unique emails
func (v *EmailNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("email = ?", v.Field)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil {
		// found a user with the same email
		errors.Add(validators.GenerateKey(v.Name), "An account with that email already exists.")
	}
}

// ==========> Internal Functions <==========

type hashParams struct {
	memory     uint32
	iterations uint32
	threads    uint8
	saltLength uint32
	keyLength  uint32
}

func generateRandomBytes(n uint32) ([]byte, error) {
	ret := make([]byte, n)
	_, err := rand.Read(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// hashPassword hashes a password with the given Params
// set @p to nil to use the Env Configuration (set by the user), ONLY set @p
// if want to compare two hashes.
func hashPassword(password string, p **hashParams) (string, error) {
	var param *hashParams

	// generate cryptographically secure salt
	salt, err := generateRandomBytes(uint32(internal.PASSWORD_HASH_SALT_LEN))
	if err != nil {
		return "", err
	}

	if (hashParams{}) == **p {
		param = &hashParams{}

		(*param).iterations = uint32(internal.PASSWORD_HASH_ITERATIONS)
		(*param).memory = uint32(internal.PASSWORD_HASH_MEMORY)
		(*param).threads = uint8(internal.PASSWORD_HASH_THREADS)
		(*param).keyLength = uint32(internal.PASSWORD_HASH_KEY_LEN)
		(*param).saltLength = uint32(internal.PASSWORD_HASH_SALT_LEN)
	} else {
		param = *p
	}

	// hash the Password
	hash := argon2.IDKey([]byte(password), salt, (*param).iterations,
		(*param).memory, (*param).threads,
		(*param).keyLength)

	// encode salt and hash since they are in binary format
	_salt := base64.RawStdEncoding.EncodeToString(salt)
	_hash := base64.RawStdEncoding.EncodeToString(hash)

	// "encode" all needed data into a ASCII string
	// which is:
	//	$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	//
	// * $argon2id				— the variant of Argon2 being used.
	// * $v=19				— the version of Argon2 being used.
	// * $m=65536,t=3,p=2			— the memory (m),
	//					  iterations (t) and
	//					  parallelism (p) parameters being used.
	// * $c29tZXNhbHQ			— the base64-encoded salt,
	//					  using standard base64-encoding
	//					  and no padding.
	// * $RdescudvJCsgt3ub+b+dWRWJTmaaJObG	— the base64-encoded hashed
	//					  password (derived key),
	//					  using standard base64-encoding
	//					  and no padding.
	hashedPwd := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, (*param).memory,
		(*param).iterations, (*param).threads, _salt,
		_hash)

	return hashedPwd, nil
}

// decodeArgonHash decodes an Argon2 hash and returns all relevant/ needed data
// to generate an identical hash.
func decodeArgonHash(hash string) (p *hashParams, salt, rawHash []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &hashParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations,
		&p.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	rawHash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, rawHash, nil
}

// comparePasswordAndHash compares an raw Password and an Argon2 hash
func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeArgonHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory,
		p.threads, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}
