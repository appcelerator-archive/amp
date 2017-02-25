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
	testUser = account.SignUpRequest{
		Name:     "user",
		Password: "userPassword",
		Email:    "user@amp.io",
	}
)

func createUser(t *testing.T, user *account.SignUpRequest) context.Context {
	// SignUp
	_, err := accountClient.SignUp(ctx, user)
	assert.NoError(t, err)

	// Create a verify token
	token, err := auth.CreateToken(user.Name, auth.TokenTypeVerify, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, err)

	// Create a login token
	token, err = auth.CreateToken(user.Name, auth.TokenTypeLogin, time.Hour)
	return metadata.NewContext(ctx, metadata.Pairs(auth.TokenKey, token))
}

func TestUserSignUpInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Name = "UpperCaseIsNotAllowed"
	_, err := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Email = "this is not an email"
	_, err := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Password = ""
	_, err := accountClient.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserShouldSignUpAndVerify(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Create a token
	token, err := auth.CreateToken(testUser.Name, auth.TokenTypeVerify, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, err)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err1 := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err1)

	// SignUp
	_, err2 := accountClient.SignUp(ctx, &testUser)
	assert.Error(t, err2)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: "this is not a token"})
	assert.Error(t, err)
}

func TestUserVerifyNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a verify token
	token, err := auth.CreateToken("nonexistinguser", auth.TokenTypeVerify, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.Error(t, err)
}

// TODO: Check expired token

func TestUserLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.NoError(t, err)
}

func TestUserLoginNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginNonVerifiedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     "not the right user name",
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "not the right password",
	})
	assert.Error(t, err)
}

func TestUserPasswordReset(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)
}

func TestUserPasswordResetMalformedRequestShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: "this is not a valid user name"})
	assert.Error(t, err)
}

func TestUserPasswordResetNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: "nonexistinguser"})
	assert.Error(t, err)
}

func TestUserPasswordSet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Set
	token, _ := auth.CreateToken(testUser.Name, auth.TokenTypePassword, time.Hour)
	_, err := accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.NoError(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.NoError(t, err)
}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	_, err = accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := auth.CreateToken("nonexistinguser", auth.TokenTypePassword, time.Hour)
	_, err = accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := accountClient.PasswordReset(ctx, &account.PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := auth.CreateToken(testUser.Name, auth.TokenTypePassword, time.Hour)
	_, err = accountClient.PasswordSet(ctx, &account.PasswordSetRequest{
		Token:    token,
		Password: "",
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: "",
	})
	assert.Error(t, err)
}

func TestUserPasswordChange(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.NoError(t, err)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: "this is not the right password",
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeEmptyNewPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := ""
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "aze"
	_, err := accountClient.PasswordChange(ownerCtx, &account.PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = accountClient.Login(ctx, &account.LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserForgotLogin(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.NoError(t, err)
}

func TestUserForgotLoginMalformedEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// SignUp
	_, err := accountClient.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{
		Email: "this is not a valid email",
	})
	assert.Error(t, err)
}

func TestUserForgotLoginNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// ForgotLogin
	_, err := accountClient.ForgotLogin(ctx, &account.ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.Error(t, err)
}

func TestUserGet(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// Get
	getReply, err := accountClient.GetUser(ctx, &account.GetUserRequest{
		Name: testUser.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, getReply.User.Name, testUser.Name)
	assert.Equal(t, getReply.User.Email, testUser.Email)
	assert.NotEmpty(t, getReply.User.CreateDt)
}

func TestUserGetMalformedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Get
	_, err := accountClient.GetUser(ctx, &account.GetUserRequest{
		Name: "this user is malformed",
	})
	assert.Error(t, err)
}

func TestUserGetNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Get
	_, err := accountClient.GetUser(ctx, &account.GetUserRequest{
		Name: testUser.Name,
	})
	assert.Error(t, err)
}

func TestUserList(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	createUser(t, &testUser)

	// List
	listReply, err := accountClient.ListUsers(ctx, &account.ListUsersRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	assert.Len(t, listReply.Users, 1)
	assert.Equal(t, listReply.Users[0].Name, testUser.Name)
	assert.Equal(t, listReply.Users[0].Email, testUser.Email)
	assert.NotEmpty(t, listReply.Users[0].CreateDt)
}

// Organizations

var (
	testOrg = account.CreateOrganizationRequest{
		Name:  "organization",
		Email: "organization@amp.io",
	}
	testMember = account.SignUpRequest{
		Name:     "organization-member",
		Password: "organizationMemberPassword",
		Email:    "organization.member@amp.io",
	}
)

func createOrganization(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest) context.Context {
	// Create a user
	ownerCtx := createUser(t, owner)

	// CreateOrganization
	_, err := accountClient.CreateOrganization(ownerCtx, org)
	assert.NoError(t, err)

	return ownerCtx
}

func TestOrganizationCreate(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	_, err := accountClient.CreateOrganization(ownerCtx, &testOrg)
	assert.NoError(t, err)
}

func TestOrganizationCreateInvalidNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Name = "this is not a valid name"
	_, err := accountClient.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateInvalidEmailShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Email = "this is not a valid email"
	_, err := accountClient.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateAlreadyExistsShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateOrganization again
	_, err := accountClient.CreateOrganization(ownerCtx, &testOrg)
	assert.Error(t, err)
}

func TestOrganizationAddUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create owner
	ownerCtx := createUser(t, &testUser)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	createOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(memberCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddNonValidatedUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// SignUp member
	_, err := accountClient.SignUp(ctx, &testMember)
	assert.NoError(t, err)

	// AddUserToOrganization
	_, err = accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToOrganization
	_, err = accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUser(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create user
	ownerCtx := createUser(t, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromOrganization
	_, err := accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserNotOwnerShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(memberCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveNonExistingUserShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// RemoveUserFromOrganization
	_, err := accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveSameUserTwiceShouldSucceed(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveAllOwnersShouldFail(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = accountClient.RemoveUserFromOrganization(ownerCtx, &account.RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testUser.Name,
	})
	assert.Error(t, err)
}
