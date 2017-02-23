package tests

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

// Users

var (
	signUpRequest = account.SignUpRequest{
		Name:     "user",
		Password: "userPassword",
		Email:    "user@amp.io",
	}
)

func TestUserSignUpInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequest
	invalidSignUp.Name = "UpperCaseIsNotAllowed"
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequest
	invalidSignUp.Email = "this is not an email"
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := signUpRequest
	invalidSignUp.Password = ""
	_, signUpErr := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, signUpErr)
}

func TestUserShouldSignUpAndVerify(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err1 := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, err1)

	// SignUp
	_, err2 := accountClient.SignUp(ctx, &signUpRequest)
	assert.Error(t, err2)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: "this is not a token"})
	assert.Error(t, verifyErr)
}

// TODO: Check token with invalid signature
// TODO: Check token with non existing account id
// TODO: Check expired token

func TestUserLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: signUpRequest.Password,
	})
	assert.NoError(t, loginErr)
}

func TestUserLoginNonExistingAccountShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: signUpRequest.Password,
	})
	assert.Error(t, loginErr)
}

func TestUserLoginNonVerifiedAccountShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: signUpRequest.Password,
	})
	assert.Error(t, loginErr)
}

func TestUserLoginInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     "not the right user name",
		Password: signUpRequest.Password,
	})
	assert.Error(t, loginErr)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: "not the right password",
	})
	assert.Error(t, loginErr)
}

func TestUserPasswordReset(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, passwordResetErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: signUpRequest.Name})
	assert.NoError(t, passwordResetErr)
}

func TestUserPasswordResetNonExistingAccountShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, passwordResetErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: "This is not an existing user"})
	assert.Error(t, passwordResetErr)
}

func TestUserPasswordSet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, passwordResetErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: signUpRequest.Name})
	assert.NoError(t, passwordResetErr)

	// Password Set
	_, passwordSetErr := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.NoError(t, passwordSetErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: "newPassword",
	})
	assert.NoError(t, loginErr)
}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, passwordResetErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: signUpRequest.Name})
	assert.NoError(t, passwordResetErr)

	// Password Set
	_, passwordSetErr := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, passwordSetErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: "newPassword",
	})
	assert.Error(t, loginErr)
}

func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Password Reset
	_, passwordResetErr := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: signUpRequest.Name})
	assert.NoError(t, passwordResetErr)

	// Password Set
	_, passwordSetErr := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "",
	})
	assert.Error(t, passwordSetErr)

	// Login
	_, loginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: "",
	})
	assert.Error(t, loginErr)
}

func TestUserPasswordChange(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := anonymousAccountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: signUpRequest.Password,
	})
	assert.NoError(t, loginErr)

	// Password Change
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	ctx = metadata.NewContext(ctx, ownerRequester)
	newPassword := "newPassword"
	_, passwordChangeErr := anonymousAccountClient.PasswordChange(ctx, &account.PasswordChangeRequest{
		ExistingPassword: signUpRequest.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, passwordChangeErr)

	// Login
	_, newLoginErr := anonymousAccountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: newPassword,
	})
	assert.NoError(t, newLoginErr)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := anonymousAccountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: signUpRequest.Password,
	})
	assert.NoError(t, loginErr)

	// Password Change
	newPassword := "newPassword"
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	ctx = metadata.NewContext(ctx, ownerRequester)
	_, passwordChangeErr := anonymousAccountClient.PasswordChange(ctx, &account.PasswordChangeRequest{
		ExistingPassword: "this is not a valid password",
		NewPassword:      newPassword,
	})
	assert.Error(t, passwordChangeErr)

	// Login
	_, newLoginErr := anonymousAccountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: newPassword,
	})
	assert.Error(t, newLoginErr)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Login
	_, loginErr := anonymousAccountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: signUpRequest.Password,
	})
	assert.NoError(t, loginErr)

	// Password Change
	newPassword := ""
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	ctx = metadata.NewContext(ctx, ownerRequester)
	_, passwordChangeErr := anonymousAccountClient.PasswordChange(ctx, &account.PasswordChangeRequest{
		ExistingPassword: signUpRequest.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, passwordChangeErr)

	// Login
	_, newLoginErr := anonymousAccountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: newPassword,
	})
	assert.Error(t, newLoginErr)
}

func TestUserForgotLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// ForgotLogin
	_, forgotLoginErr := accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{Email: signUpRequest.Email})
	assert.NoError(t, forgotLoginErr)
}

func TestUserForgotLoginInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// ForgotLogin
	_, forgotLoginErr := accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{Email: "this is not a valid email"})
	assert.Error(t, forgotLoginErr)
}

func TestUserGet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// Get
	getReply, getErr := accountClient.GetUser(ctx, &account.GetUserRequest{Name: signUpRequest.Name})
	assert.NoError(t, getErr)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, getReply.User.Name, signUpRequest.Name)
	assert.Equal(t, getReply.User.Email, signUpRequest.Email)
	assert.NotEmpty(t, getReply.User.CreateDt)
}

func TestUserList(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := accountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// List
	listReply, listErr := accountClient.ListUsers(ctx, &account.ListUsersRequest{})
	assert.NoError(t, listErr)
	assert.NotEmpty(t, listReply)
	assert.Len(t, listReply.Users, 1)
	assert.Equal(t, listReply.Users[0].Name, signUpRequest.Name)
	assert.Equal(t, listReply.Users[0].Email, signUpRequest.Email)
	assert.NotEmpty(t, listReply.Users[0].CreateDt)
}

// Organizations

var (
	createOrganizationRequest = account.CreateOrganizationRequest{
		Name:  "organization",
		Email: "organization@amp.io",
	}
	orgMemberSignUpRequest = account.SignUpRequest{
		Name:     "organization-member",
		Password: "organizationMemberPassword",
		Email:    "organization.member@amp.io",
	}
)

func TestOrganizationCreate(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)
}

func TestOrganizationCreateNotVerifiedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.Error(t, createOrganizationErr)
}

func TestOrganizationCreateAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// CreateOrganization again
	_, createOrganizationAgainErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.Error(t, createOrganizationAgainErr)
}

func TestOrganizationAddUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)
}

func TestOrganizationAddUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	memberRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, memberRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.Error(t, addUserToOrganizationErr)
}

func TestOrganizationAddNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.Error(t, addUserToOrganizationErr)
}

func TestOrganizationAddNonValidatedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.Error(t, addUserToOrganizationErr)
}

func TestOrganizationAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.Error(t, addUserToOrganizationErr)
}

func TestOrganizationAddSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr = anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)
}

func TestOrganizationRemoveUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)

	// RemoveUserFromOrganization
	_, removeUserFromOrganizationErr := anonymousAccountClient.RemoveUserFromOrganization(metadata.NewContext(ctx, ownerRequester), &account.RemoveUserFromOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, removeUserFromOrganizationErr)
}

func TestOrganizationRemoveUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	memberRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)

	// RemoveUserFromOrganization
	_, removeUserFromOrganizationErr := anonymousAccountClient.RemoveUserFromOrganization(metadata.NewContext(ctx, memberRequester), &account.RemoveUserFromOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.Error(t, removeUserFromOrganizationErr)
}

func TestOrganizationRemoveNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// RemoveUserFromOrganization
	_, removeUserFromOrganizationErr := anonymousAccountClient.RemoveUserFromOrganization(metadata.NewContext(ctx, ownerRequester), &account.RemoveUserFromOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.Error(t, removeUserFromOrganizationErr)
}

func TestOrganizationRemoveSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)

	// RemoveUserFromOrganization
	_, removeUserFromOrganizationErr := anonymousAccountClient.RemoveUserFromOrganization(metadata.NewContext(ctx, ownerRequester), &account.RemoveUserFromOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, removeUserFromOrganizationErr)

	// RemoveUserFromOrganization
	_, removeUserFromOrganizationErr = anonymousAccountClient.RemoveUserFromOrganization(metadata.NewContext(ctx, ownerRequester), &account.RemoveUserFromOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, removeUserFromOrganizationErr)
}

func TestOrganizationRemoveAllOwnersShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, signUpErr := anonymousAccountClient.SignUp(ctx, &signUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr := auth.CreateUserToken(signUpRequest.Name, time.Hour)
	ownerRequester := metadata.Pairs(auth.TokenKey, token)
	assert.NoError(t, createTokenErr)

	// Verify
	_, verifyErr := anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// SignUp member
	_, signUpErr = anonymousAccountClient.SignUp(ctx, &orgMemberSignUpRequest)
	assert.NoError(t, signUpErr)

	// Create a token
	token, createTokenErr = auth.CreateUserToken(orgMemberSignUpRequest.Name, time.Hour)
	assert.NoError(t, createTokenErr)

	// Verify member
	_, verifyErr = anonymousAccountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, verifyErr)

	// CreateOrganization
	_, createOrganizationErr := anonymousAccountClient.CreateOrganization(metadata.NewContext(ctx, ownerRequester), &createOrganizationRequest)
	assert.NoError(t, createOrganizationErr)

	// AddUserToOrganization
	_, addUserToOrganizationErr := anonymousAccountClient.AddUserToOrganization(metadata.NewContext(ctx, ownerRequester), &account.AddUserToOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: orgMemberSignUpRequest.Name})
	assert.NoError(t, addUserToOrganizationErr)

	// RemoveUserFromOrganization
	_, removeUserFromOrganizationErr := anonymousAccountClient.RemoveUserFromOrganization(metadata.NewContext(ctx, ownerRequester), &account.RemoveUserFromOrganizationRequest{OrganizationName: createOrganizationRequest.Name, UserName: signUpRequest.Name})
	assert.Error(t, removeUserFromOrganizationErr)
}
