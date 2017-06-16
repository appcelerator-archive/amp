package accounts

import (
	"log"
	"os"
	"testing"

	. "github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/data/accounts"
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

func TestUserSignUpConflictWithOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create an organization
	h.CreateOrganization(t, &testOrg, &testUser)

	// SignUp user with organization name
	conflictSignUp := testUser
	conflictSignUp.Name = testOrg.Name
	_, err := h.Accounts().SignUp(ctx, &conflictSignUp)
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

func TestUserDeleteUserNotOwnerOfOrganizationShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create an organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create a member
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err := h.Accounts().DeleteUser(memberCtx, &DeleteUserRequest{Name: testMember.Name})
	assert.NoError(t, err)
}

// Organizations

func TestOrganizationCreate(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// CreateOrganization
	_, err := h.Accounts().CreateOrganization(ownerCtx, &testOrg)
	assert.NoError(t, err)
}

func TestOrganizationCreateInvalidNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Name = "this is not a valid name"
	_, err := h.Accounts().CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateInvalidEmailShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Email = "this is not a valid email"
	_, err := h.Accounts().CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateAlreadyExistsShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// CreateOrganization again
	_, err := h.Accounts().CreateOrganization(ownerCtx, &testOrg)
	assert.Error(t, err)
}

func TestOrganizationCreateConflictsWithUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create user
	ownerCtx := h.CreateUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Name = testUser.Name
	_, err := h.Accounts().CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationAddUser(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create owner
	ownerCtx := h.CreateUser(t, &testUser)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserNotOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(memberCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddSameUserTwiceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToOrganization
	_, err = h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUser(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create user
	ownerCtx := h.CreateUser(t, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// RemoveUserFromOrganization
	_, err := h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserNotOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(memberCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testUser.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// RemoveUserFromOrganization
	_, err := h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveSameUserTwiceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveAllOwnersShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testUser.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationChangeUserRole(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create a member
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	_, err := h.Accounts().ChangeOrganizationMemberRole(ownerCtx, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.NoError(t, err)
}

func TestOrganizationChangeUserRoleNotOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create a member
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	_, err := h.Accounts().ChangeOrganizationMemberRole(memberCtx, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.Error(t, err)
}

func TestOrganizationChangeUserRoleNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create user
	h.CreateUser(t, &testMember)

	_, err := h.Accounts().ChangeOrganizationMemberRole(ownerCtx, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.Error(t, err)
}

func TestOrganizationGet(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a organization
	h.CreateOrganization(t, &testOrg, &testUser)

	// Get
	getReply, err := h.Accounts().GetOrganization(ctx, &GetOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, getReply.Organization.Name, testOrg.Name)
	assert.Equal(t, getReply.Organization.Email, testOrg.Email)
	assert.NotEmpty(t, getReply.Organization.CreateDt)
}

func TestOrganizationGetMalformedOrganizationShouldFail(t *testing.T) {
	// Get
	_, err := h.Accounts().GetOrganization(ctx, &GetOrganizationRequest{
		Name: "this organization is malformed",
	})
	assert.Error(t, err)
}

func TestOrganizationList(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	h.CreateOrganization(t, &testOrg, &testUser)

	// List
	listReply, err := h.Accounts().ListOrganizations(ctx, &ListOrganizationsRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	found := false
	for _, org := range listReply.Organizations {
		if org.Name == testOrg.Name {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestOrganizationDelete(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Delete
	_, err := h.Accounts().DeleteOrganization(ownerCtx, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationDeleteNotOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create a member
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err := h.Accounts().DeleteOrganization(memberCtx, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationDeleteNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// Delete
	_, err := h.Accounts().DeleteOrganization(ownerCtx, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.Error(t, err)
}

// Teams

func TestTeamCreate(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// CreateTeam
	_, err := h.Accounts().CreateTeam(ownerCtx, &testTeam)
	assert.NoError(t, err)
}

func TestTeamCreateInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "this is not a valid name"
	_, err := h.Accounts().CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateInvalidTeamNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.TeamName = "this is not a valid name"
	_, err := h.Accounts().CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateUser(t, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "non-existing-org"
	_, err := h.Accounts().CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateNotOrgOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create organization
	h.CreateOrganization(t, &testOrg, &testUser)

	// Create a user not part of the organization
	notOrgOwnerCtx := h.CreateUser(t, &testMember)

	// CreateTeam
	_, err := h.Accounts().CreateTeam(notOrgOwnerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamCreateAlreadyExistsShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// CreateTeam again
	_, err := h.Accounts().CreateTeam(ownerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamAddUser(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidTeamNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateUser(t, &testUser)
	h.CreateUser(t, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddNonExistingUserToTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserNotOrganizationOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(memberCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddSameUserTwiceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToTeam again
	_, err = h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamChangeName(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Change team name
	_, err := h.Accounts().ChangeTeamName(ownerCtx, &ChangeTeamNameRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		NewName:          "newteamname",
	})
	assert.NoError(t, err)
}

func TestTeamChangeNameToAlreadyExistingTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	anotherTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create another team
	_, err := h.Accounts().CreateTeam(ownerCtx, &anotherTeam)
	assert.NoError(t, err)

	// Change team name
	_, err = h.Accounts().ChangeTeamName(ownerCtx, &ChangeTeamNameRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		NewName:          anotherTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUser(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserInvalidTeamNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	ownerCtx := h.CreateUser(t, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserFromNonExistingTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = h.Accounts().RemoveUserFromTeam(memberCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveNonExistingUserShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// RemoveUserFromTeam
	_, err := h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotPartOfTheTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	h.CreateUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamGet(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Get
	getReply, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.NotEmpty(t, getReply.Team)
	assert.Equal(t, getReply.Team.Name, testTeam.TeamName)
	assert.NotEmpty(t, getReply.Team.CreateDt)
	assert.NotEmpty(t, getReply.Team.Members)
	assert.Equal(t, getReply.Team.Members[0], testUser.Name)
}

func TestTeamGetNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	ownerCtx := h.CreateUser(t, &testUser)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetNonExistingTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create org
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetMalformedOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: "this is not a valid team name",
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetMalformedTeamShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         "this is not a valid team name",
	})
	assert.Error(t, err)
}

func TestTeamList(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a team
	h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// List
	listReply, err := h.Accounts().ListTeams(ctx, &ListTeamsRequest{
		OrganizationName: testOrg.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	assert.Len(t, listReply.Teams, 1)
	assert.Equal(t, listReply.Teams[0].Name, testTeam.TeamName)
	assert.NotEmpty(t, listReply.Teams[0].CreateDt)
	assert.NotEmpty(t, listReply.Teams[0].Members)
	assert.Equal(t, listReply.Teams[0].Members[0], testUser.Name)
}

func TestTeamListInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// List
	_, err := h.Accounts().ListTeams(ctx, &ListTeamsRequest{
		OrganizationName: "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamListNonExistingOrganizationNameShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	h.CreateUser(t, &testUser)

	// List
	_, err := h.Accounts().ListTeams(ctx, &ListTeamsRequest{
		OrganizationName: testTeam.OrganizationName,
	})
	assert.Error(t, err)
}

func TestTeamDelete(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Delete
	_, err := h.Accounts().DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}

func TestTeamDeleteNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create org
	ownerCtx := h.CreateUser(t, &testUser)

	// Delete
	_, err := h.Accounts().DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamDeleteNonExistingTeamShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create org
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Delete
	_, err := h.Accounts().DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}

func TestTeamDeleteNotOrgOwnerShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err := h.Accounts().DeleteTeam(memberCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
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

func TestSuperUserNotOwnerOfOrganizationAddUserShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create organization
	h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err = h.Accounts().AddUserToOrganization(su, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOwnerOfOrganizationRemoveUserShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create member
	h.CreateUser(t, &testMember)

	// AddUserToOrganization
	_, err = h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(su, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOwnerOfOrganizationChangeUserRoleShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create organization
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create a member
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	_, err = h.Accounts().ChangeOrganizationMemberRole(su, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOwnerOfOrganizationDeleteShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create a user
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Create a member
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err = h.Accounts().DeleteOrganization(su, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOrgOwnerTeamCreateShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create organization
	h.CreateOrganization(t, &testOrg, &testUser)

	// CreateTeam
	_, err = h.Accounts().CreateTeam(su, &testTeam)
	assert.NoError(t, err)
}

func TestSuperUserNotOrgOwnerTeamAddUserShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err = h.Accounts().AddUserToTeam(su, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOrgOwnerTeamRemoveUserShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err = h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = h.Accounts().RemoveUserFromTeam(su, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOrgOwnerTeamDeleteShouldSucceed(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)
	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err = h.Accounts().DeleteTeam(su, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}
