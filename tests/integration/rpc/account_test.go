package tests

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

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

	// Establish an anonymous connection
	conn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	assert.NoError(t, err)

	// Recreate the account client
	accountClient := account.NewAccountClient(conn)

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

	// Password Change
	md := metadata.Pairs(auth.TokenKey, token)
	ctx = metadata.NewContext(ctx, md)
	newPassword := "newPassword"
	_, passwordChangeErr := accountClient.PasswordChange(ctx, &account.PasswordChangeRequest{
		ExistingPassword: signUpRequest.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, passwordChangeErr)

	// Login
	_, newLoginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: newPassword,
	})
	assert.NoError(t, newLoginErr)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Establish an anonymous connection
	conn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	assert.NoError(t, err)

	// Recreate the account client
	accountClient := account.NewAccountClient(conn)

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

	// Password Change
	newPassword := "newPassword"
	md := metadata.Pairs(auth.TokenKey, token)
	ctx = metadata.NewContext(ctx, md)
	_, passwordChangeErr := accountClient.PasswordChange(ctx, &account.PasswordChangeRequest{
		ExistingPassword: "this is not a valid password",
		NewPassword:      newPassword,
	})
	assert.Error(t, passwordChangeErr)

	// Login
	_, newLoginErr := accountClient.Login(ctx, &account.LogInRequest{
		Name:     signUpRequest.Name,
		Password: newPassword,
	})
	assert.Error(t, newLoginErr)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Establish an anonymous connection
	conn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	assert.NoError(t, err)

	// Recreate the account client
	accountClient := account.NewAccountClient(conn)

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

	// Password Change
	newPassword := ""
	md := metadata.Pairs(auth.TokenKey, token)
	ctx = metadata.NewContext(ctx, md)
	_, passwordChangeErr := accountClient.PasswordChange(ctx, &account.PasswordChangeRequest{
		ExistingPassword: signUpRequest.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, passwordChangeErr)

	// Login
	_, newLoginErr := accountClient.Login(ctx, &account.LogInRequest{
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
