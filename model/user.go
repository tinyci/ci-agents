package model

import (
	"encoding/json"
	"time"

	gh "github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/types"
	"golang.org/x/oauth2"
)

// Capability is a type of access gating mechanism. If present on the user
// account access is granted, otherwise not.
type Capability string

const (
	// CapabilityModifyCI is required for modifying CI properties such as adding or removing a repo.
	CapabilityModifyCI Capability = "modify:ci"
	// CapabilityModifyUser allows you to modify users; including caps.
	CapabilityModifyUser Capability = "modify:user"
	// CapabilitySubmit allows manual submissions
	CapabilitySubmit Capability = "submit"
	// CapabilityCancel allows cancels
	CapabilityCancel Capability = "cancel"
)

var (
	// AllCapabilities comprises the superuser account's list of capabilities.
	AllCapabilities = []Capability{CapabilityModifyCI, CapabilityModifyUser, CapabilitySubmit, CapabilityCancel}

	// TokenCryptKey is the standard token crypt key.
	// NOTE: the default is only used by tests; it is overwritten on service boot; see config/auth.go.
	TokenCryptKey = []byte{1, 2, 3, 4, 5, 6, 7, 8}
)

// UserError is the encapsulation of many *errors.Errors that need to be presented to
// the user.
type UserError struct {
	ID     int64  `gorm:"primary key" json:"id"`
	UserID int64  `json:"-"`
	Error  string `json:"error"`
}

// NewUserErrorFromProto converts a user error from proto to canonical representation
func NewUserErrorFromProto(ue *types.UserError) *UserError {
	return &UserError{
		ID:     ue.Id,
		UserID: ue.UserID,
		Error:  ue.Error,
	}
}

// ToProto converts the usererror to its protocol buffer representation.
func (ue UserError) ToProto() *types.UserError {
	return &types.UserError{
		Id:     ue.ID,
		UserID: ue.UserID,
		Error:  ue.Error,
	}
}

// User is a user record.
type User struct {
	ID               int64         `gorm:"primary_key" json:"id"`
	Username         string        `gorm:"unique" json:"username"`
	LastScannedRepos *time.Time    `json:"last_scanned_repos,omitempty"`
	Errors           []UserError   `json:"errors,omitempty"`
	Subscribed       []*Repository `gorm:"many2many:subscriptions;preload:false" json:"subscribed,omitempty"`
	LoginToken       []byte        `json:"-"`

	TokenJSON []byte        `gorm:"column:token;not null" json:"-"`
	Token     *oauth2.Token `json:"token,omitempty"`
}

// SetToken sets the token's byte stream, and encrypts it.
func (u *User) SetToken() *errors.Error {
	var err *errors.Error
	u.TokenJSON, err = encryptToken(TokenCryptKey, u.Token)
	return err
}

// FetchToken retrieves the token from the db, decrypting it if necessary.
func (u *User) FetchToken() *errors.Error {
	if u.Token != nil {
		return nil
	}

	var err *errors.Error
	u.Token, err = decryptToken(TokenCryptKey, u.TokenJSON)
	return err
}

func encryptToken(key []byte, tok *oauth2.Token) ([]byte, *errors.Error) {
	str, err := securecookie.EncodeMulti("token", tok, securecookie.CodecsFromPairs(key)...)
	if err != nil {
		return nil, errors.New(err)
	}

	return []byte(str), nil
}

func decryptToken(key, tokenBytes []byte) (*oauth2.Token, *errors.Error) {
	tok := oauth2.Token{}

	if len(tokenBytes) == 0 {
		return &tok, nil
	}

	err := securecookie.DecodeMulti("token", string(tokenBytes), &tok, securecookie.CodecsFromPairs(key)...)
	if err != nil {
		return nil, errors.New(err)
	}

	return &tok, nil
}

// NewUserFromProto converts a proto user to a real user.
func NewUserFromProto(u *types.User) (*User, *errors.Error) {
	errs := []UserError{}

	if len(u.Errors) != 0 {
		for _, e := range u.Errors {
			errs = append(errs, *NewUserErrorFromProto(e))
		}
	}

	token := &oauth2.Token{}

	if u.TokenJSON != nil {
		if err := json.Unmarshal(u.TokenJSON, token); err != nil {
			return nil, errors.New(err)
		}
	}

	return &User{
		ID:               u.Id,
		Username:         u.Username,
		LastScannedRepos: MakeTime(u.LastScannedRepos, true),
		Errors:           errs,
		TokenJSON:        u.TokenJSON,
		Token:            token,
	}, nil
}

// ToProto converts the user struct to a protobuf capable one
func (u *User) ToProto() *types.User {
	errors := []*types.UserError{}

	for _, e := range u.Errors {
		errors = append(errors, e.ToProto())
	}

	if u.Token != nil {
		u.TokenJSON, _ = json.Marshal(u.Token)
	}

	return &types.User{
		Id:               u.ID,
		Username:         u.Username,
		LastScannedRepos: MakeTimestamp(u.LastScannedRepos),
		Errors:           errors,
		TokenJSON:        u.TokenJSON,
	}
}

// CreateUser initializes a user struct and writes it to the db.
func (m *Model) CreateUser(username string, token *oauth2.Token) (*User, *errors.Error) {
	u := &User{Username: username, Token: token}
	return u, m.WrapError(m.Create(u), "creating user")
}

// FindUserByID finds the user by integer ID.
func (m *Model) FindUserByID(id int64) (*User, *errors.Error) {
	u := &User{}
	return u, m.WrapError(m.Where("id = ?", id).First(u), "finding user by id")
}

// FindUserByName finds a user by unique key username.
func (m *Model) FindUserByName(username string) (*User, *errors.Error) {
	u := &User{}
	return u, m.WrapError(m.Where("username = ?", username).First(u), "finding user by name")
}

// FindUserByNameWithSubscriptions finds a user by unique key username. It also fetches the subscriptions for the user.
func (m *Model) FindUserByNameWithSubscriptions(username string) (*User, *errors.Error) {
	u := &User{}
	return u, m.WrapError(m.Preload("Subscribed").Where("username = ?", username).First(u), "preloading subscriptions with user")
}

// DeleteError deletes a given *errors.Error for a user.
func (m *Model) DeleteError(u *User, id int64) *errors.Error {
	return m.WrapError(m.Where("id = ?", id).Delete(u.Errors), "deleting errors for user")
}

// AddSubscriptionsForUser adds the repositories to the subscriptions table. Access is
// validated at the API level, not here.
func (m *Model) AddSubscriptionsForUser(u *User, repos []*Repository) *errors.Error {
	return errors.New(m.Model(u).Association("Subscribed").Append(repos).Error)
}

// RemoveSubscriptionForUser removes an item from the subscriptions table.
func (m *Model) RemoveSubscriptionForUser(u *User, repo *Repository) *errors.Error {
	return errors.New(m.Model(u).Association("Subscribed").Delete(repo).Error)
}

// AddError adds an error to the error list.
func (u *User) AddError(err *errors.Error) {
	u.Errors = append(u.Errors, UserError{Error: err.Error()})
}

// AfterFind is a gorm hook to unmarshal the Token JSON after finding the record.
func (u *User) AfterFind(tx *gorm.DB) error {
	if err := u.FetchToken(); err != nil {
		return err
	}

	if err := u.Validate(); err != nil {
		return errors.New(err).Wrapf("reading user id %d (%q)", u.ID, u.Username)
	}

	return nil
}

// BeforeCreate just calls BeforeSave.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.BeforeSave(tx)
}

// BeforeSave is a gorm hook to marshal the Token JSON before saving the record
func (u *User) BeforeSave(tx *gorm.DB) error {
	if err := u.ValidateWrite(); err != nil {
		return errors.New(err).Wrapf("saving user %q", u.Username)
	}

	if err := u.SetToken(); err != nil {
		return err
	}

	return nil
}

// ValidateWrite is for write-only validations.
func (u *User) ValidateWrite() *errors.Error {
	if u.Token == nil || u.Token.AccessToken == "" || !u.Token.Valid() {
		return errors.New("cannot be written because the oauth credentials are not valid")
	}

	return u.Validate()
}

// Validate validates the user record to ensure it can be written.
func (u *User) Validate() *errors.Error {
	if u.Username == "" {
		return errors.New("username is empty")
	}

	return nil
}

func (m *Model) mkRepositoryFromGithub(repo *gh.Repository, owner *User, autoCreated bool) *Repository {
	return &Repository{
		Name:        repo.GetFullName(),
		Private:     repo.GetPrivate(),
		Disabled:    true, // created repos are disabled by default
		Github:      repo,
		Owner:       owner,
		AutoCreated: autoCreated,
	}
}

// SaveRepositories saves github repositories; it sets the *User provided to
// the owner of it.
func (m *Model) SaveRepositories(repos []*gh.Repository, username string, autoCreated bool) *errors.Error {
	owner, err := m.FindUserByName(username)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		_, err := m.GetRepositoryByName(repo.GetFullName())
		if err != nil {
			localRepo := m.mkRepositoryFromGithub(repo, owner, autoCreated)
			if err := m.WrapError(m.Create(localRepo), "creating repository"); err != nil {
				return err.Wrapf("could not create repository %q", repo.GetFullName())
			}
		}
	}

	t := time.Now()

	owner.LastScannedRepos = &t
	return m.WrapError(m.Save(owner), "saving owner cache data")
}

// ListSubscribedTasksForUser lists all tasks related to the subscribed repositories
// for the user.
func (m *Model) ListSubscribedTasksForUser(userID, page, perPage int64) ([]*Task, *errors.Error) {
	tasks := []*Task{}
	call := m.Limit(perPage).Offset(page*perPage).Joins(
		"inner join subscriptions on subscriptions.repository_id = tasks.parent_id",
	).Where("subscriptions.user_id = ?", userID).Find(&tasks)

	return tasks, m.WrapError(call, "locating user's subscribed tasks")
}

// AddCapabilityToUser adds a capability to a user account.
func (m *Model) AddCapabilityToUser(u *User, cap Capability) *errors.Error {
	return m.WrapError(m.Exec("insert into user_capabilities (user_id, name) values (?, ?)", u.ID, cap), "adding capability for user")
}

// RemoveCapabilityFromUser removes a capability from a user account.
func (m *Model) RemoveCapabilityFromUser(u *User, cap Capability) *errors.Error {
	return m.WrapError(m.Exec("delete from user_capabilities where user_id = ? and name = ?", u.ID, cap), "removing capability from user")
}

// HasCapability returns true if the user is capable of performing the operation.
func (m *Model) HasCapability(u *User, cap Capability) (bool, *errors.Error) {
	var slice []int64
	err := m.WrapError(m.Raw("select 1 from user_capabilities where user_id = ? and name = ?", u.ID, cap).Find(&slice), "checking capabilities for user")
	return len(slice) > 0, err
}
