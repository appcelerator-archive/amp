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
	anonymous context.Context
	h         *helpers.Helper
)

func setup() (err error) {
	h, err = helpers.New()
	if err != nil {
		return err
	}
	anonymous = context.Background()
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
	testOwner := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(anonymous, &testOwner)
	assert.NoError(t, err)

	// Create a token
	token, err := h.Tokens().CreateVerificationToken(testOwner.Name)
	assert.NoError(t, err)

	// Verify
	_, err = h.Accounts().Verify(anonymous, &VerificationRequest{Token: token})
	assert.NoError(t, err)
}

func TestUserSignUpInvalidNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	invalidSignUp := testOwner
	invalidSignUp.Name = "UpperCaseIsNotAllowed"
	_, err := h.Accounts().SignUp(anonymous, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidEmailShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	invalidSignUp := testOwner
	invalidSignUp.Email = "this is not an email"
	_, err := h.Accounts().SignUp(anonymous, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpInvalidPasswordShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	invalidSignUp := testOwner
	invalidSignUp.Password = ""
	_, err := h.Accounts().SignUp(anonymous, &invalidSignUp)
	assert.Error(t, err)
}

func TestUserSignUpAlreadyExistsShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(anonymous, &testOwner)
	assert.NoError(t, err)

	// SignUp
	_, err = h.Accounts().SignUp(anonymous, &testOwner)
	assert.Error(t, err)
}

func TestUserSignUpConflictWithOrganizationShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create an organization
	h.CreateOrganization(t, &testOrg, &testOwner)

	// SignUp user with organization name
	conflictSignUp := testOwner
	conflictSignUp.Name = testOrg.Name
	_, err := h.Accounts().SignUp(anonymous, &conflictSignUp)
	assert.Error(t, err)
}

func TestUserVerifyNotATokenShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(anonymous, &testOwner)
	assert.NoError(t, err)

	// Verify
	_, err = h.Accounts().Verify(anonymous, &VerificationRequest{Token: "this is not a token"})
	assert.Error(t, err)
}

func TestUserVerifyNonExistingUserShouldFail(t *testing.T) {
	// Create a verify token
	token, err := h.Tokens().CreateVerificationToken("nonexistinguser")
	assert.NoError(t, err)

	// Verify
	_, err = h.Accounts().Verify(anonymous, &VerificationRequest{Token: token})
	assert.Error(t, err)
}

func TestUserLogin(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Login
	_, err := h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: testOwner.Password,
	})
	assert.NoError(t, err)
}

func TestUserLoginNonExistingUserShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Login
	_, err := h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: testOwner.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Login
	_, err := h.Accounts().Login(anonymous, &LogInRequest{
		Name:     "not the right user name",
		Password: testOwner.Password,
	})
	assert.Error(t, err)
}

func TestUserLoginInvalidPasswordShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Login
	_, err := h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: "not the right password",
	})
	assert.Error(t, err)
}

func TestUserPasswordReset(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Reset
	_, err := h.Accounts().PasswordReset(anonymous, &PasswordResetRequest{Name: testOwner.Name})
	assert.NoError(t, err)
}

func TestUserPasswordResetMalformedRequestShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Reset
	_, err := h.Accounts().PasswordReset(anonymous, &PasswordResetRequest{Name: "this is not a valid user name"})
	assert.Error(t, err)
}

func TestUserPasswordResetNonExistingUserShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Reset
	_, err := h.Accounts().PasswordReset(anonymous, &PasswordResetRequest{Name: "nonexistinguser"})
	assert.Error(t, err)
}

func TestUserPasswordSet(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Set
	token, _ := h.Tokens().CreatePasswordToken(testOwner.Name)
	_, err := h.Accounts().PasswordSet(anonymous, &PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.NoError(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: "newPassword",
	})
	assert.NoError(t, err)
}

func TestUserPasswordSetInvalidTokenShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Reset
	_, err := h.Accounts().PasswordReset(anonymous, &PasswordResetRequest{Name: testOwner.Name})
	assert.NoError(t, err)

	// Password Set
	_, err = h.Accounts().PasswordSet(anonymous, &PasswordSetRequest{
		Token:    "this is an invalid token",
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetNonExistingUserShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Reset
	_, err := h.Accounts().PasswordReset(anonymous, &PasswordResetRequest{Name: testOwner.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := h.Tokens().CreatePasswordToken("nonexistinguser")
	_, err = h.Accounts().PasswordSet(anonymous, &PasswordSetRequest{
		Token:    token,
		Password: "newPassword",
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: "newPassword",
	})
	assert.Error(t, err)
}

func TestUserPasswordSetInvalidPasswordShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Password Reset
	_, err := h.Accounts().PasswordReset(anonymous, &PasswordResetRequest{Name: testOwner.Name})
	assert.NoError(t, err)

	// Password Set
	token, _ := h.Tokens().CreatePasswordToken(testOwner.Name)
	_, err = h.Accounts().PasswordSet(anonymous, &PasswordSetRequest{
		Token:    token,
		Password: "",
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: "",
	})
	assert.Error(t, err)
}

func TestUserPasswordChange(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Password Change
	newPassword := "newPassword"
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testOwner.Password,
		NewPassword:      newPassword,
	})
	assert.NoError(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: newPassword,
	})
	assert.NoError(t, err)
}

func TestUserPasswordChangeInvalidExistingPassword(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Password Change
	newPassword := "newPassword"
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: "this is not the right password",
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeEmptyNewPassword(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Password Change
	newPassword := ""
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testOwner.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserPasswordChangeInvalidNewPassword(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Password Change
	newPassword := "aze"
	_, err := h.Accounts().PasswordChange(ownerCtx, &PasswordChangeRequest{
		ExistingPassword: testOwner.Password,
		NewPassword:      newPassword,
	})
	assert.Error(t, err)

	// Login
	_, err = h.Accounts().Login(anonymous, &LogInRequest{
		Name:     testOwner.Name,
		Password: newPassword,
	})
	assert.Error(t, err)
}

func TestUserForgotLogin(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(anonymous, &testOwner)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = h.Accounts().ForgotLogin(anonymous, &ForgotLoginRequest{
		Email: testOwner.Email,
	})
	assert.NoError(t, err)
}

func TestUserForgotLoginMalformedEmailShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// SignUp
	_, err := h.Accounts().SignUp(anonymous, &testOwner)
	assert.NoError(t, err)

	// ForgotLogin
	_, err = h.Accounts().ForgotLogin(anonymous, &ForgotLoginRequest{
		Email: "this is not a valid email",
	})
	assert.Error(t, err)
}

func TestUserForgotLoginNonExistingUserShouldFail(t *testing.T) {
	testOwner := h.RandomUser()

	// ForgotLogin
	_, err := h.Accounts().ForgotLogin(anonymous, &ForgotLoginRequest{
		Email: testOwner.Email,
	})
	assert.Error(t, err)
}

func TestUserGet(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	userCtx := h.CreateUser(t, &testOwner)

	// Get
	getReply, err := h.Accounts().GetUser(userCtx, &GetUserRequest{
		Name: testOwner.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, testOwner.Name, getReply.User.Name)
	assert.Equal(t, testOwner.Email, getReply.User.Email)
	assert.NotEmpty(t, getReply.User.CreateDt)
}

func TestAnonymousUserGetShouldNotReturnUserEmail(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Get
	getReply, err := h.Accounts().GetUser(anonymous, &GetUserRequest{
		Name: testOwner.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getReply)
	assert.Equal(t, testOwner.Name, getReply.User.Name)
	assert.Empty(t, getReply.User.Email)
	assert.NotEmpty(t, getReply.User.CreateDt)
}

func TestUserGetMalformedUserShouldFail(t *testing.T) {
	// Get
	_, err := h.Accounts().GetUser(anonymous, &GetUserRequest{
		Name: "this user is malformed",
	})
	assert.Error(t, err)
}

func TestUserGetNonExistingUserShouldFail(t *testing.T) {
	// Get
	_, err := h.Accounts().GetUser(anonymous, &GetUserRequest{
		Name: "nonexistinguser",
	})
	assert.Error(t, err)
}

func TestUserList(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// List
	listReply, err := h.Accounts().ListUsers(anonymous, &ListUsersRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	found := false
	for _, user := range listReply.Users {
		if user.Name == testOwner.Name {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestUserDelete(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Delete
	_, err := h.Accounts().DeleteUser(ownerCtx, &DeleteUserRequest{Name: testOwner.Name})
	assert.NoError(t, err)
}

func TestUserDeleteSomeoneElseAccountShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Create another user
	h.CreateUser(t, &testMember)

	// Delete
	_, err := h.Accounts().DeleteUser(ownerCtx, &DeleteUserRequest{Name: testMember.Name})
	assert.Error(t, err)
}

//func testOwnerDeleteUserNotOwnerOfOrganizationShouldSucceed(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create an organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create a member
//	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)
//
//	// Delete
//	_, err := h.Accounts().DeleteUser(memberCtx, &DeleteUserRequest{Name: testMember.Name})
//	assert.NoError(t, err)
//}
//
//// Organizations
//
//func TestOrganizationCreate(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// CreateOrganization
//	_, err := h.Accounts().CreateOrganization(ownerCtx, &testOrg)
//	assert.NoError(t, err)
//}
//
//func TestOrganizationCreateInvalidNameShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// CreateOrganization
//	invalidRequest := testOrg
//	invalidRequest.Name = "this is not a valid name"
//	_, err := h.Accounts().CreateOrganization(ownerCtx, &invalidRequest)
//	assert.Error(t, err)
//}
//
//func TestOrganizationCreateInvalidEmailShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// CreateOrganization
//	invalidRequest := testOrg
//	invalidRequest.Email = "this is not a valid email"
//	_, err := h.Accounts().CreateOrganization(ownerCtx, &invalidRequest)
//	assert.Error(t, err)
//}
//
//func TestOrganizationCreateAlreadyExistsShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// CreateOrganization again
//	_, err := h.Accounts().CreateOrganization(ownerCtx, &testOrg)
//	assert.Error(t, err)
//}
//
//func TestOrganizationCreateConflictsWithUserShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// CreateOrganization
//	invalidRequest := testOrg
//	invalidRequest.Name = testOwner.Name
//	_, err := h.Accounts().CreateOrganization(ownerCtx, &invalidRequest)
//	assert.Error(t, err)
//}
//
//func TestOrganizationAddUser(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//}
//
//func TestOrganizationAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: "this is not a valid name",
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationAddUserInvalidUserNameShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         "this is not a valid name",
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create owner
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationAddUserNotOwnerShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	memberCtx := h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(memberCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationAddNonExistingUserShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationAddSameUserTwiceShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// AddUserToOrganization
//	_, err = h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveUser(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//}
//
//func TestOrganizationRemoveUserInvalidOrganizationNameShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: "this is not a valid name",
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveUserInvalidUserNameShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         "this is not a valid name",
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveUserFromNonExistingOrganizationShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// RemoveUserFromOrganization
//	_, err := h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveUserNotOwnerShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	memberCtx := h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(memberCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testOwner.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveNonExistingUserShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// RemoveUserFromOrganization
//	_, err := h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveSameUserTwiceShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationRemoveAllOwnersShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err := h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//
//	// RemoveUserFromOrganization
//	_, err = h.Accounts().RemoveUserFromOrganization(ownerCtx, &RemoveUserFromOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testOwner.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationChangeUserRole(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create a member
//	h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)
//
//	_, err := h.Accounts().ChangeOrganizationMemberRole(ownerCtx, &ChangeOrganizationMemberRoleRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
//	})
//	assert.NoError(t, err)
//}
//
//func TestOrganizationChangeUserRoleNotOwnerShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create a member
//	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)
//
//	_, err := h.Accounts().ChangeOrganizationMemberRole(memberCtx, &ChangeOrganizationMemberRoleRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationChangeUserRoleNonExistingUserShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create organization
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create user
//	h.CreateUser(t, &testMember)
//
//	_, err := h.Accounts().ChangeOrganizationMemberRole(ownerCtx, &ChangeOrganizationMemberRoleRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationGet(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a organization
//	h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Get
//	getReply, err := h.Accounts().GetOrganization(anonymous, &GetOrganizationRequest{
//		Name: testOrg.Name,
//	})
//	assert.NoError(t, err)
//	assert.NotEmpty(t, getReply)
//	assert.Equal(t, getReply.Organization.Name, testOrg.Name)
//	assert.Equal(t, getReply.Organization.Email, testOrg.Email)
//	assert.NotEmpty(t, getReply.Organization.CreateDt)
//}
//
//func TestOrganizationGetMalformedOrganizationShouldFail(t *testing.T) {
//	// Get
//	_, err := h.Accounts().GetOrganization(anonymous, &GetOrganizationRequest{
//		Name: "this organization is malformed",
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationList(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// List
//	listReply, err := h.Accounts().ListOrganizations(anonymous, &ListOrganizationsRequest{})
//	assert.NoError(t, err)
//	assert.NotEmpty(t, listReply)
//	found := false
//	for _, org := range listReply.Organizations {
//		if org.Name == testOrg.Name {
//			found = true
//			break
//		}
//	}
//	assert.True(t, found)
//}
//
//func TestOrganizationDelete(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Delete
//	_, err := h.Accounts().DeleteOrganization(ownerCtx, &DeleteOrganizationRequest{
//		Name: testOrg.Name,
//	})
//	assert.NoError(t, err)
//}
//
//func TestOrganizationDeleteNotOwnerShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create a member
//	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)
//
//	// Delete
//	_, err := h.Accounts().DeleteOrganization(memberCtx, &DeleteOrganizationRequest{
//		Name: testOrg.Name,
//	})
//	assert.Error(t, err)
//}
//
//func TestOrganizationDeleteNonExistingOrganizationShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.RandomOrg()
//
//	// Create a user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// Delete
//	_, err := h.Accounts().DeleteOrganization(ownerCtx, &DeleteOrganizationRequest{
//		Name: testOrg.Name,
//	})
//	assert.Error(t, err)
//}

// Teams

func TestTeamCreate(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// CreateTeam
	_, err := h.Accounts().CreateTeam(ownerCtx, &testTeam)
	assert.NoError(t, err)
}

func TestTeamCreateInvalidOrganizationNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.OrganizationName = "this is not a valid name"
	_, err := h.Accounts().CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

func TestTeamCreateInvalidTeamNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a user
	ownerCtx := h.CreateUser(t, &testOwner)

	// CreateTeam
	invalidRequest := testTeam
	invalidRequest.TeamName = "this is not a valid name"
	_, err := h.Accounts().CreateTeam(ownerCtx, &invalidRequest)
	assert.Error(t, err)
}

//func TestTeamCreateNonExistingOrganizationShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testOrg := h.DefaultOrg()
//	testTeam := h.RandomTeam(testOrg.Name)
//
//	// Create a user
//	ownerCtx := h.CreateUser(t, &testOwner)
//
//	// CreateTeam
//	invalidRequest := testTeam
//	invalidRequest.OrganizationName = "non-existing-org"
//	_, err := h.Accounts().CreateTeam(ownerCtx, &invalidRequest)
//	assert.Error(t, err)
//}
//s
//func TestTeamCreateNotOrgOwnerShouldFail(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.DefaultOrg()
//	testOrg := h.RandomOrg()
//	testTeam := h.RandomTeam(testOrg.Name)
//
//	// Create organization
//	h.CreateOrganization(t, &testOrg, &testOwner)
//
//	// Create a user not part of the organization
//	notOrgOwnerCtx := h.CreateUser(t, &testMember)
//
//	// CreateTeam
//	_, err := h.Accounts().CreateTeam(notOrgOwnerCtx, &testTeam)
//	assert.Error(t, err)
//}

func TestTeamCreateAlreadyExistsShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// CreateTeam again
	_, err := h.Accounts().CreateTeam(ownerCtx, &testTeam)
	assert.Error(t, err)
}

func TestTeamCreateByOrgMemberShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create organization and add a member
	h.CreateUser(t, &testOwner)
	memberCtx := h.CreateUser(t, &testMember)

	// CreateTeam
	_, err := h.Accounts().CreateTeam(memberCtx, &testTeam)
	assert.NoError(t, err)
}

func TestTeamAddUser(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)
	h.CreateUser(t, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestTeamAddUserInvalidOrganizationNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: "this is not a valid name",
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidTeamNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         "this is not a valid name",
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserInvalidUserNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamAddUserToNonExistingOrganizationShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateUser(t, &testOwner)
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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateUser(t, &testOwner)
	h.CreateUser(t, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddNonExistingUserToTeamShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddUserNotOrganizationOwnerShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	h.CreateTeam(t, &testOrg, &testOwner, &testTeam)
	memberCtx := h.CreateUser(t, &testMember)

	// AddUserToTeam
	_, err := h.Accounts().AddUserToTeam(memberCtx, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamAddSameUserTwiceShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)
	h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Change team name
	_, err := h.Accounts().ChangeTeamName(ownerCtx, &ChangeTeamNameRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		NewName:          "new" + testTeam.TeamName,
	})
	assert.NoError(t, err)
}

func TestTeamChangeNameToAlreadyExistingTeamShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	anotherTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	ownerCtx := h.CreateUser(t, &testOwner)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	ownerCtx := h.CreateUser(t, &testOwner)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	memberCtx := h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// RemoveUserFromTeam
	_, err := h.Accounts().RemoveUserFromTeam(ownerCtx, &RemoveUserFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.Error(t, err)
}

func TestTeamRemoveUserNotPartOfTheTeamShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

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
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

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
	assert.Equal(t, getReply.Team.Members[0], testOwner.Name)
}

func TestTeamGetNonExistingOrganizationShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	ownerCtx := h.CreateUser(t, &testOwner)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetNonExistingTeamShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create org
	ownerCtx := h.CreateUser(t, &testOwner)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetMalformedOrganizationShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: "this is not a valid team name",
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamGetMalformedTeamShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Get
	_, err := h.Accounts().GetTeam(ownerCtx, &GetTeamRequest{
		OrganizationName: testOrg.Name,
		TeamName:         "this is not a valid team name",
	})
	assert.Error(t, err)
}

func TestTeamList(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create a team
	h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// List
	listReply, err := h.Accounts().ListTeams(anonymous, &ListTeamsRequest{
		OrganizationName: testOrg.Name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, listReply)
	assert.NotEmpty(t, listReply.Teams)
	var team *accounts.Team
	for _, t := range listReply.Teams {
		if t.Name == testTeam.TeamName {
			team = t
			break
		}
	}
	assert.NotEmpty(t, team)
	assert.Equal(t, team.Name, testTeam.TeamName)
	assert.NotEmpty(t, team.CreateDt)
	assert.NotEmpty(t, team.Members)
	assert.Equal(t, team.Members[0], testOwner.Name)
}

func TestTeamListInvalidOrganizationNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// List
	_, err := h.Accounts().ListTeams(anonymous, &ListTeamsRequest{
		OrganizationName: "this is not a valid name",
	})
	assert.Error(t, err)
}

func TestTeamListNonExistingOrganizationNameShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user
	h.CreateUser(t, &testOwner)

	// List
	_, err := h.Accounts().ListTeams(anonymous, &ListTeamsRequest{
		OrganizationName: testTeam.OrganizationName,
	})
	assert.Error(t, err)
}

func TestTeamDelete(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Delete
	_, err := h.Accounts().DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}

func TestTeamDeleteNonExistingOrganizationShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create org
	ownerCtx := h.CreateUser(t, &testOwner)

	// Delete
	_, err := h.Accounts().DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamDeleteNonExistingTeamShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create org
	ownerCtx := h.CreateUser(t, &testOwner)

	// Delete
	_, err := h.Accounts().DeleteTeam(ownerCtx, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.Error(t, err)
}

func TestTeamDeleteNotOrgOwnerShouldFail(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	/// Create team
	h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	memberCtx := h.CreateUser(t, &testMember)

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
	_, err := h.Accounts().Login(anonymous, &LogInRequest{
		Name:     superUser.Name,
		Password: superUser.Password,
	})
	assert.NoError(t, err)
}

func TestSuperUserDeleteSomeoneElseAccountShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()

	// Create a user
	h.CreateUser(t, &testOwner)

	// Super user
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Delete
	_, err = h.Accounts().DeleteUser(su, &DeleteUserRequest{Name: testOwner.Name})
	assert.NoError(t, err)
}

//func TestSuperUserNotOwnerOfOrganizationAddUserShouldSucceed(t *testing.T) {
//	testOwner := h.RandomUser()
//	testMember := h.RandomUser()
//	testOrg := h.DefaultOrg()
//	su, err := h.SuperLogin()
//	assert.NoError(t, err)
//
//	// Create owner
//	h.CreateUser(t, &testOwner)
//
//	// Create member
//	h.CreateUser(t, &testMember)
//
//	// AddUserToOrganization
//	_, err = h.Accounts().AddUserToOrganization(su, &AddUserToOrganizationRequest{
//		OrganizationName: testOrg.Name,
//		UserName:         testMember.Name,
//	})
//	assert.NoError(t, err)
//}

func TestSuperUserNotOwnerOfOrganizationRemoveUserShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create owner
	h.CreateUser(t, &testOwner)

	// Create member
	h.CreateUser(t, &testMember)

	//// AddUserToOrganization
	//_, err = h.Accounts().AddUserToOrganization(ownerCtx, &AddUserToOrganizationRequest{
	//	OrganizationName: testOrg.Name,
	//	UserName:         testMember.Name,
	//})
	//assert.NoError(t, err)

	// RemoveUserFromOrganization
	_, err = h.Accounts().RemoveUserFromOrganization(su, &RemoveUserFromOrganizationRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOwnerOfOrganizationChangeUserRoleShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create owner
	h.CreateUser(t, &testOwner)

	// Create a member
	h.CreateUser(t, &testMember)

	_, err = h.Accounts().ChangeOrganizationMemberRole(su, &ChangeOrganizationMemberRoleRequest{
		OrganizationName: testOrg.Name,
		UserName:         testMember.Name,
		Role:             accounts.OrganizationRole_ORGANIZATION_OWNER,
	})
	assert.NoError(t, err)
}

func TestSuperDeleteDefaultOrganizationShouldFail(t *testing.T) {
	testOrg := h.DefaultOrg()
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Delete
	_, err = h.Accounts().DeleteOrganization(su, &DeleteOrganizationRequest{
		Name: testOrg.Name,
	})
	assert.Error(t, err)
}

func TestSuperUserNotOrgOwnerTeamCreateShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create owner
	h.CreateUser(t, &testOwner)

	// CreateTeam
	_, err = h.Accounts().CreateTeam(su, &testTeam)
	assert.NoError(t, err)
}

func TestSuperUserNotOrgOwnerTeamAddUserShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	// Create team
	h.CreateTeam(t, &testOrg, &testOwner, &testTeam)
	h.CreateUser(t, &testMember)

	// AddUserToTeam
	_, err = h.Accounts().AddUserToTeam(su, &AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)
}

func TestSuperUserNotOrgOwnerTeamRemoveUserShouldSucceed(t *testing.T) {
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testOwner, &testTeam)

	// Create member in org
	h.CreateUser(t, &testMember)

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
	testOwner := h.RandomUser()
	testMember := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	/// Create team
	h.CreateTeam(t, &testOrg, &testOwner, &testTeam)
	h.CreateUser(t, &testMember)

	// Delete
	_, err = h.Accounts().DeleteTeam(su, &DeleteTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
}
