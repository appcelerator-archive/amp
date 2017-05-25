package accounts

import (
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/auth"
	. "github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/tests"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	as     accounts.Interface
	ctx    context.Context
	client AccountClient
)

func setup() error {
	as = helpers.NewAccountsStore()
	ctx = context.Background()
	conn, err := helpers.AmplifierConnection()
	if err != nil {
		return err
	}
	client = NewAccountClient(conn)
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

//func TestUserShouldSignUpAndVerify(t *testing.T) {
//	testUser := randomUser()
//
//	// SignUp
//	_, err := client.SignUp(ctx, &testUser)
//	assert.NoError(t, err)
//
//	// Create a token
//	token, err := auth.CreateVerificationToken(testUser.Name)
//	assert.NoError(t, err)
//
//	// Verify
//	_, err = client.Verify(ctx, &VerificationRequest{Token: token})
//	assert.NoError(t, err)
//}

func TestUserSignUpInvalidNameShouldFail(t *testing.T) {
	testUser := randomUser()

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Name = "UpperCaseIsNotAllowed"
	_, err := client.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	testUser := randomUser()

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Email = "this is not an email"
	_, err := client.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	testUser := randomUser()

	// SignUp
	invalidSignUp := testUser
	invalidSignUp.Password = ""
	_, err := client.SignUp(ctx, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	testUser := randomUser()

	// SignUp
	_, err := client.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// SignUp
	_, err = client.SignUp(ctx, &testUser)
	assert.Error(t, err)
}

func TestUserSignUpConflictWithOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create an organization
	createOrganization(t, &testOrg, &testUser)

	// SignUp user with organization name
	conflictSignUp := testUser
	conflictSignUp.Name = testOrg.Name
	_, err := client.SignUp(ctx, &conflictSignUp)
	assert.Error(t, err)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	testUser := randomUser()

	// SignUp
	_, err := client.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// Verify
	_, err = client.Verify(ctx, &VerificationRequest{Token: "this is not a token"})
	assert.Error(t, err)
}

//func TestUserVerifyNonExistingUserShouldFail(t *testing.T) {
//	// Create a verify token
//	token, err := auth.CreateVerificationToken("nonexistinguser")
//	assert.NoError(t, err)
//
//	// Verify
//	_, err = client.Verify(ctx, &VerificationRequest{Token: token})
//	assert.Error(t, err)
//}

func TestUserLogin(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.NoError(t, err)
}

func TestUserLoginNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()

	// Login
	_, err := client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidNameShouldFail(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := client.Login(ctx, &LogInRequest{
		Name:     "not the right user name",
		Password: testUser.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Login
	_, err := client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "not the right password",
	})
	assert.Error(t, err)
}

func TestUserPasswordReset(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := client.PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)
}

func TestUserPasswordResetMalformedRequestShouldFail(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := client.PasswordReset(ctx, &PasswordResetRequest{Name: "this is not a valid user name"})
	assert.Error(t, err)
}

func TestUserPasswordResetNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := client.PasswordReset(ctx, &PasswordResetRequest{Name: "nonexistinguser"})
	assert.Error(t, err)
}

//func TestUserPasswordSet(t *testing.T) {
//	testUser := randomUser()
//
//	// Create a user
//	createUser(t, &testUser)
//
//	// Password Set
//	token, _ := auth.CreatePasswordToken(testUser.Name)
//	_, err := client.PasswordSet(ctx, &PasswordSetRequest{
//		Token:    token,
//		Password: "newPassword",
//	})
//	assert.NoError(t, err)
//
//	// Login
//	_, err = client.Login(ctx, &LogInRequest{
//		Name:     testUser.Name,
//		Password: "newPassword",
//	})
//	assert.NoError(t, err)
//}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Password Reset
	_, err := client.PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
	assert.NoError(t, err)

	// Password Set
	_, err = client.PasswordSet(ctx, &PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

//func TestUserPasswordSetNonExistingUserShouldFail(t *testing.T) {
//	testUser := randomUser()
//
//	// Create a user
//	createUser(t, &testUser)
//
//	// Password Reset
//	_, err := client.PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
//	assert.NoError(t, err)
//
//	// Password Set
//	token, _ := auth.CreatePasswordToken("nonexistinguser")
//	_, err = client.PasswordSet(ctx, &PasswordSetRequest{
//		Token:    token,
//		Password: "newPassword",
//	})
//	assert.Error(t, err)
//
//	// Login
//	_, err = client.Login(ctx, &LogInRequest{
//		Name:     testUser.Name,
//		Password: "newPassword",
//	})
//	assert.Error(t, err)
//}

//func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
//	testUser := randomUser()
//
//	// Create a user
//	createUser(t, &testUser)
//
//	// Password Reset
//	_, err := client.PasswordReset(ctx, &PasswordResetRequest{Name: testUser.Name})
//	assert.NoError(t, err)
//
//	// Password Set
//	token, _ := auth.CreatePasswordToken(testUser.Name)
//	_, err = client.PasswordSet(ctx, &PasswordSetRequest{
//		Token:    token,
//		Password: "",
//	})
//	assert.Error(t, err)
//
//	// Login
//	_, err = client.Login(ctx, &LogInRequest{
//		Name:     testUser.Name,
//		Password: "",
//	})
//	assert.Error(t, err)
//}

func TestUserPasswordChange(t *testing.T) {
	testUser := randomUser()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := client.PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, err)

	// Login
	_, err = client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.NoError(t, err)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	testUser := randomUser()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "newPassword"
	_, err := client.PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: "this is not the right password",
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeEmptyNewPassword(t *testing.T) {
	testUser := randomUser()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := ""
	_, err := client.PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	testUser := randomUser()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Password Change
	newPassword := "aze"
	_, err := client.PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testUser.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = client.Login(ctx, &LogInRequest{
		Name:     testUser.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserForgotLogin(t *testing.T) {
	testUser := randomUser()

	// SignUp
	_, err := client.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = client.ForgotLogin(ctx, &ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.NoError(t, err)
}

func TestUserForgotLoginMalformedEmailShouldFail(t *testing.T) {
	testUser := randomUser()

	// SignUp
	_, err := client.SignUp(ctx, &testUser)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = client.ForgotLogin(ctx, &ForgotLoginRequest{
		Email: "this is not a valid email",
	})
	assert.Error(t, err)
}

func TestUserForgotLoginNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()

	// ForgotLogin
	_, err := client.ForgotLogin(ctx, &ForgotLoginRequest{
		Email: testUser.Email,
	})
	assert.Error(t, err)
}

func TestUserGet(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// Get
	getReply, err := client.GetUser(ctx, &GetUserRequest{
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
	_, err := client.GetUser(ctx, &GetUserRequest{
		Name: "this user is malformed",
	})
	assert.Error(t, err)
}

func TestUserGetNonExistingUserShouldFail(t *testing.T) {
	// Get
	_, err := client.GetUser(ctx, &GetUserRequest{
		Name: "nonexistinguser",
	})
	assert.Error(t, err)
}

func TestUserList(t *testing.T) {
	testUser := randomUser()

	// Create a user
	createUser(t, &testUser)

	// List
	listReply, err := client.ListUsers(ctx, &ListUsersRequest{})
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
	testUser := randomUser()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Delete
	_, err := client.DeleteUser(ownerCtx, &DeleteUserRequest{Name: testUser.Name})
	assert.NoError(t, err)
}

func TestUserDeleteSomeoneElseAccountShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Create another user
	createUser(t, &testMember)

	// Delete
	_, err := client.DeleteUser(ownerCtx, &DeleteUserRequest{Name: testMember.Name})
	assert.Error(t, err)
}

func TestUserDeleteUserOnlyOwnerOfOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create an organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Delete
	_, err := client.DeleteUser(ownerCtx, &DeleteUserRequest{Name: testUser.Name})
	assert.Error(t, err)
}

func TestUserDeleteUserNotOwnerOfOrganizationShouldSucceed(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create an organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create a member
	memberCtx := createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err := client.DeleteUser(memberCtx, &DeleteUserRequest{Name: testMember.Name})
	assert.NoError(t, err)
}

// Organizations

func TestOrganizationCreate(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	_, err := client.CreateOrganization(ownerCtx, &testOrg)
	assert.NoError(t, err)
}

func TestOrganizationCreateInvalidNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Name = "this is not a valid name"
	_, err := client.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateInvalidEmailShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Email = "this is not a valid email"
	_, err := client.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationCreateAlreadyExistsShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateOrganization again
	_, err := client.CreateOrganization(ownerCtx, &testOrg)
	assert.Error(t, err)
}

func TestOrganizationCreateConflictsWithUserShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create user
	ownerCtx := createUser(t, &testUser)

	// CreateOrganization
	invalidRequest := testOrg
	invalidRequest.Name = testUser.Name
	_, err := client.CreateOrganization(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestOrganizationAddUser(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create owner
	ownerCtx := createUser(t, &testUser)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddUserNotOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	createOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(memberCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationAddSameUserTwiceShouldSucceed(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToOrganization
	_, err = client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUser(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create user
	ownerCtx := createUser(t, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromOrganization
	_, err := client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveUserNotOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	memberCtx := createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(memberCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testUser.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// RemoveUserFromOrganization
	_, err := client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationRemoveSameUserTwiceShouldSucceed(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationRemoveAllOwnersShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = client.RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testUser.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationChangeUserRole(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create a member
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	_, err := client.ChangeOrganizationMemberRole(ownerCtx, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.NoError(t, err)
}

func TestOrganizationChangeUserRoleNotOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create a member
	memberCtx := createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	_, err := client.ChangeOrganizationMemberRole(memberCtx, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.Error(t, err)
}

func TestOrganizationChangeUserRoleNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create organization
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create user
	createUser(t, &testMember)

	_, err := client.ChangeOrganizationMemberRole(ownerCtx, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.Error(t, err)
}

func TestOrganizationGet(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a organization
	createOrganization(t, &testOrg, &testUser)

	// Get
	getReply, err := client.GetOrganization(ctx, &GetOrganizationRequest{
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
	_, err := client.GetOrganization(ctx, &GetOrganizationRequest{
		Name: "this organization is malformed",
	})
	assert.Error(t, err)
}

func TestOrganizationList(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a user
	createOrganization(t, &testOrg, &testUser)

	// List
	listReply, err := client.ListOrganizations(ctx, &ListOrganizationsRequest{})
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
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Delete
	_, err := client.DeleteOrganization(ownerCtx, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.NoError(t, err)
}

func TestOrganizationDeleteNotOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create a member
	memberCtx := createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// Delete
	_, err := client.DeleteOrganization(memberCtx, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.Error(t, err)
}

func TestOrganizationDeleteNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// Delete
	_, err := client.DeleteOrganization(ownerCtx, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.Error(t, err)
}

// Teams

func TestTeamCreate(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateTeam
	_, err := client.CreateTeam(ownerCtx, &testTeam)
	assert.NoError(t, err)
}

func TestTeamCreateInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "this is not a valid name"
	_, err := client.CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateInvalidTeamNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create a user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.TeamName = "this is not a valid name"
	_, err := client.CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create a user
	ownerCtx := createUser(t, &testUser)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "non-existing-org"
	_, err := client.CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateNotOrgOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create organization
	createOrganization(t, &testOrg, &testUser)

	// Create a user not part of the organization
	notOrgOwnerCtx := createUser(t, &testMember)

	// CreateTeam
	_, err := client.CreateTeam(notOrgOwnerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamCreateAlreadyExistsShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// CreateTeam again
	_, err := client.CreateTeam(ownerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamAddUser(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidTeamNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createUser(t, &testUser)
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingTeamShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createOrganization(t, &testOrg, &testUser)
	createUser(t, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddNonExistingUserToTeamShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserNotOrganizationOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	createTeam(t, &testOrg, &testUser, &testTeam)
	memberCtx := createUser(t, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(memberCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddNonValidatedUserShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// SignUp member
	_, err := client.SignUp(ctx, &testMember)
	assert.NoError(t, err)

	// AddUserToTeam
	_, err = client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddSameUserTwiceShouldSucceed(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// AddUserToTeam again
	_, err = client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamRemoveUser(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserInvalidTeamNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserInvalidUserNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create user
	ownerCtx := createUser(t, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserFromNonExistingTeamShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create user
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotOwnerShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	/// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	memberCtx := createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)

	// AddUserToTeam
	_, err := client.AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// RemoveUserFromTeam
	_, err = client.RemoveUserFromTeam(memberCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveNonExistingUserShouldFail(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// RemoveUserFromTeam
	_, err := client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotPartOfTheTeamShouldSucceed(t *testing.T) {
	testUser := randomUser()
	testMember := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	/// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Create member
	createUser(t, &testMember)

	// RemoveUserFromTeam
	_, err := client.RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamGet(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Get
	getReply, err := client.GetTeam(ownerCtx, &GetTeamRequest{
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
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create user
	ownerCtx := createUser(t, &testUser)

	// Get
	_, err := client.GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetNonExistingTeamShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create org
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Get
	_, err := client.GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetMalformedOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Get
	_, err := client.GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: "this is not a valid team name",
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetMalformedTeamShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Get
	_, err := client.GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         "this is not a valid team name",
	})
	assert.Error(t, err)
}

func TestTeamList(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create a team
	createTeam(t, &testOrg, &testUser, &testTeam)

	// List
	listReply, err := client.ListTeams(ctx, &ListTeamsRequest{
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
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	createTeam(t, &testOrg, &testUser, &testTeam)

	// List
	_, err := client.ListTeams(ctx, &ListTeamsRequest{
		OrganizationName: "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamListNonExistingOrganizationNameShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create user
	createUser(t, &testUser)

	// List
	_, err := client.ListTeams(ctx, &ListTeamsRequest{
		OrganizationName: testTeam.OrganizationName,
	})
	assert.Error(t, err)
}

func TestTeamDelete(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create team
	ownerCtx := createTeam(t, &testOrg, &testUser, &testTeam)

	// Delete
	_, err := client.DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}

func TestTeamDeleteNonExistingOrganizationShouldFail(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create org
	ownerCtx := createUser(t, &testUser)

	// Delete
	_, err := client.DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamDeleteNonExistingTeamShouldSucceed(t *testing.T) {
	testUser := randomUser()
	testOrg := randomOrg()
	testTeam := randomTeam(testOrg.Name)

	// Create org
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Delete
	_, err := client.DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}

// Helpers

func createUser(t *testing.T, user *SignUpRequest) context.Context {
	// SignUp
	_, err := client.SignUp(ctx, user)
	assert.NoError(t, err)

	// Login
	header := metadata.MD{}
	_, err = client.Login(ctx, &LogInRequest{Name: user.Name, Password: user.Password}, grpc.Header(&header))
	assert.NoError(t, err)

	// Extract token from header
	tokens := header[auth.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

	return metadata.NewContext(ctx, metadata.Pairs(auth.AuthorizationHeader, auth.ForgeAuthorizationHeader(token)))
}

func createOrganization(t *testing.T, org *CreateOrganizationRequest, owner *SignUpRequest) context.Context {
	// Create a user
	ownerCtx := createUser(t, owner)

	// CreateOrganization
	_, err := client.CreateOrganization(ownerCtx, org)
	assert.NoError(t, err)

	return ownerCtx
}

func createAndAddUserToOrganization(ownerCtx context.Context, t *testing.T, org *CreateOrganizationRequest, user *SignUpRequest) context.Context {
	// Create a user
	userCtx := createUser(t, user)

	// AddUserToOrganization
	_, err := client.AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
		OrganizationName: org.Name,
		UserName:         user.Name,
	})
	assert.NoError(t, err)
	return userCtx
}

func createTeam(t *testing.T, org *CreateOrganizationRequest, owner *SignUpRequest, team *CreateTeamRequest) context.Context {
	// Create a user
	ownerCtx := createOrganization(t, org, owner)

	// CreateTeam
	_, err := client.CreateTeam(ownerCtx, team)
	assert.NoError(t, err)

	return ownerCtx
}

func randomUser() SignUpRequest {
	id := stringid.GenerateNonCryptoID()
	return SignUpRequest{
		Name:     id,
		Password: "userPassword",
		Email:    id + "@user.email",
	}
}

func randomOrg() CreateOrganizationRequest {
	id := stringid.GenerateNonCryptoID()
	return CreateOrganizationRequest{
		Name:  id,
		Email: id + "@org.email",
	}
}

func randomTeam(org string) CreateTeamRequest {
	id := stringid.GenerateNonCryptoID()
	return CreateTeamRequest{
		OrganizationName: org,
		TeamName:         id,
	}
}

//func switchAccount(userCtx context.Context, t *testing.T, accountName string) context.Context {
//	header := metadata.MD{}
//	_, err := client.Switch(userCtx, &SwitchRequest{Account: accountName}, grpc.Header(&header))
//	assert.NoError(t, err)
//
//	// Extract token from header
//	tokens := header[auth.AuthorizationHeader]
//	assert.NotEmpty(t, tokens)
//	token := tokens[0]
//	assert.NotEmpty(t, token)
//
//	return metadata.NewContext(ctx, metadata.Pairs(auth.AuthorizationHeader, token))
//}
//
//func changeOrganizationMemberRole(userCtx context.Context, t *testing.T, org *CreateOrganizationRequest, user *SignUpRequest, role accounts.OrganizationRole) {
//	_, err := client.ChangeOrganizationMemberRole(userCtx, &ChangeOrganizationMemberRoleRequest{OrganizationName: org.Name, UserName: user.Name, Role: role})
//	assert.NoError(t, err)
//}
