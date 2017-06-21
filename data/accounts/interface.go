package accounts

import "golang.org/x/net/context"

// Error type
type Error string

func (e Error) Error() string {
	return string(e)
}

// Errors
const (
	InvalidName               = Error("username is invalid")
	InvalidEmail              = Error("email is invalid")
	PasswordTooWeak           = Error("password is too weak")
	WrongPassword             = Error("password is wrong")
	InvalidToken              = Error("token is invalid")
	UserAlreadyExists         = Error("user already exists")
	EmailAlreadyUsed          = Error("email is already in use")
	UserNotFound              = Error("user not found")
	UserNotVerified           = Error("user not verified")
	OrganizationAlreadyExists = Error("organization already exists")
	OrganizationNotFound      = Error("organization not found")
	TeamAlreadyExists         = Error("team already exists")
	TeamNotFound              = Error("team not found")
	AtLeastOneOwner           = Error("organization must have at least one owner")
	NotAuthorized             = Error("user not authorized")
	NotPartOfOrganization     = Error("user is not part of the organization")
	InvalidResourceID         = Error("invalid resource ID")
	ResourceNotFound          = Error("resource not found")
	ResourceAlreadyExists     = Error("resource already exists")
	TokenAlreadyUsed          = Error("token has already been used")
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user with given password
	CreateUser(ctx context.Context, name string, email string, password string) (user *User, err error)

	// CheckUserPassword checks the given user password
	CheckUserPassword(ctx context.Context, name string, password string) (err error)

	// SetUserPassword sets the given user password
	SetUserPassword(ctx context.Context, name string, password string) (err error)

	// GetUser fetches a user by name
	GetUser(ctx context.Context, name string) (user *User, err error)

	// GetUserByEmail fetches a user by email
	GetUserByEmail(ctx context.Context, email string) (user *User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*User, err error)

	// VerifyUser verifies a user account
	VerifyUser(ctx context.Context, name string) (err error)

	// DeleteNotVerifedUser deletes a not verified user by-passing the authorization check
	DeleteNotVerifiedUser(ctx context.Context, name string) (err error)

	// DeleteUser deletes a user by name
	DeleteUser(ctx context.Context, name string) (err error)

	// IsAuthorized returns whether the requesting user is authorized to perform the given action on given resource
	IsAuthorized(ctx context.Context, owner *Account, action string, resource string, resourceId string) bool

	// Reset resets the user storage
	Reset(ctx context.Context)
}
