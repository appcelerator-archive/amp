package accounts

import (
	"log"
	"os"
	"testing"

	. "github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	ctx context.Context
	h   *helpers.Helper
)

func setup() (err error) {
	h, err = helpers.New()
	if err != nil {
		return err
	}
	ctx = context.Background()
	return nil
}

func tearDown() {
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	tearDown()
	os.Exit(code)
}

// Users

func TestUserShouldSignUpAndVerify(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Create a token
	token, err := h.Tokens().CreateVerificationToken(testUser.Name)
	assert.NoError(t, err)

	// Verify
	_, err = h.Accounts().Verify(ctx, &VerificationRequest{Token: token})
	assert.NoError(t, err)
}

func TestUserSignUpInvalidNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Name = "UpperCaseIsNotAllowed"
	_, err := h.Accounts().SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Email = "this is not an email"
	_, err := h.Accounts().SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Password = ""
	_, err := h.Accounts().SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// SignUp
	_, err = h.Accounts().SignUp(ctx, &testUser)
	assert.Error(t, err)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Verify
	_, err = h.Accounts().Verify(ctx, &VerificationRequest{Token: "this is not a token"})
	assert.Error(t, err)
}

func TestUserVerifyNonExistingUserShouldFail(t *testing.T) {
	// Create a verify token
	token, err := h.Tokens().CreateVerificationToken("nonexistinguser")
	assert.NoError(t, err)

	// Verify
	_, err = h.Accounts().Verify(ctx, &VerificationRequest{Token: token})
	assert.Error(t, err)
}

func TestUserLogin(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Login
	_, err := h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.NoError(t, err)
}

func TestUserLoginNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Login
	_, err := h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Login
	_, err := h.Accounts().Login(ctx, &LogInRequest{
		Name:     "not the right user name",
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Login
	_, err := h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "not the right password",
	})
	assert.Error(t, err)
}

func TestUserPasswordReset(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Reset
	_, err := h.Accounts().PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)
}

func TestUserPasswordResetMalformedRequestShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Reset
	_, err := h.Accounts().PasswordReset(ctx, &PasswordResetRequest{Name: "this is not a valid user name"})
	assert.Error(t, err)
}

func TestUserPasswordResetNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Reset
	_, err := h.Accounts().PasswordReset(ctx, &PasswordResetRequest{Name: "nonexistinguser"})
	assert.Error(t, err)
}

func TestUserPasswordSet(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Set
	token, _ := h.Tokens().CreatePasswordToken(testUser.Name)
	_, err := h.Accounts().PasswordSet(ctx, &PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.NoError(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.NoError(t, err)
}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Reset
	_, err := h.Accounts().PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	_, err = h.Accounts().PasswordSet(ctx, &PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Reset
	_, err := h.Accounts().PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := h.Tokens().CreatePasswordToken("nonexistinguser")
	_, err = h.Accounts().PasswordSet(ctx, &PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Password Reset
	_, err := h.Accounts().PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := h.Tokens().CreatePasswordToken(testUser.Name)
	_, err = h.Accounts().PasswordSet(ctx, &PasswordSetRequest{
		Token:    token,
		Password: "",
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "",
	})
	assert.Error(t, err)
}

func TestUserPasswordChange(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.NoError(t, err)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: "this is not the right password",
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeEmptyNewPassword(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Password Change
	newPassword := ""
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Password Change
	newPassword := "aze"
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserForgotLogin(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = h.Accounts().ForgotLogin(ctx, &ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.NoError(t, err)
}

func TestUserForgotLoginMalformedEmailShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = h.Accounts().ForgotLogin(ctx, &ForgotLoginRequest{
		Email: "this is not a valid email",
	})
	assert.Error(t, err)
}

func TestUserForgotLoginNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// ForgotLogin
	_, err := h.Accounts().ForgotLogin(ctx, &ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.Error(t, err)
}

func TestUserGet(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Get
	getReply, err := h.Accounts().GetUser(ctx, &GetUserRequest{
		Name: testUser.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, getReply.User.Name, testUser.Name)
	assert.Equal(t, getReply.User.Email, testUser.Email)
	assert.NotEmpty(t, getReply.User.CreateDt)
}

func TestUserGetMalformedUserShouldFail(t *testing.T) {
	// Get
	_, err := h.Accounts().GetUser(ctx, &GetUserRequest{
		Name: "this user is malformed",
	})
	assert.Error(t, err)
}

func TestUserGetNonExistingUserShouldFail(t *testing.T) {
	// Get
	_, err := h.Accounts().GetUser(ctx, &GetUserRequest{
		Name: "nonexistinguser",
	})
	assert.Error(t, err)
}

func TestUserList(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// List
	listReply, err := h.Accounts().ListUsers(ctx, &ListUsersRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	found := false
	for _, user := range listReply.Users {
		if user.Name == testUser.Name {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestUserDelete(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Delete
	_, err := h.Accounts().DeleteUser(ownerCtx, &DeleteUserRequest{Name: testUser.Name})
	assert.NoError(t, err)
}

func TestUserDeleteSomeoneElseAccountShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Create another user
	h.CreateUser(t, &testMember)

	// Delete
	_, err := h.Accounts().DeleteUser(ownerCtx, &DeleteUserRequest{Name: testMember.Name})
	assert.Error(t, err)
}

// Super Accounts

func TestSuperUserLogin(t *testing.T) {
	superUser := h.SuperUser()

	// Login
	_, err := h.Accounts().Login(ctx, &LogInRequest{
		Name:     superUser.Name,
		Password: superUser.Password,
	})
	assert.NoError(t, err)
}

func TestSuperUserDeleteSomeoneElseAccountShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testUser)

	// Super user
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Delete
	_, err = h.Accounts().DeleteUser(su, &DeleteUserRequest{Name: testUser.Name})
	assert.NoError(t, err)
}
