package schema

type Error string

func (e Error) Error() string { return string(e) }

const InvalidName = Error("username is invalid")
const InvalidEmail = Error("email is invalid")
const PasswordTooWeak = Error("password is too weak")
const WrongPassword = Error("password is wrong")
const InvalidToken = Error("token is invalid")
const UserAlreadyExists = Error("user already exists")
const UserNotFound = Error("user not found")
const UserNotVerified = Error("user not verified")
const OrganizationAlreadyExists = Error("organization already exists")
const OrganizationNotFound = Error("organization not found")
const TeamAlreadyExists = Error("team already exists")
const TeamNotFound = Error("team not found")
const AtLeastOneOwner = Error("organization must have at least one owner")
const NotAuthorized = Error("not authorized")
