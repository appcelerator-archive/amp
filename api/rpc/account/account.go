package account

import (
	"fmt"
	"log"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/pkg/mail"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Server is used to implement account.AccountServer
type Server struct {
	Accounts accounts.Interface
	Mailer   *mail.Mailer
	Config   *configuration.Configuration
	Tokens   *auth.Tokens
}

func convertError(err error) error {
	switch err {
	case accounts.InvalidName,
		accounts.InvalidEmail,
		accounts.InvalidToken,
		accounts.PasswordTooWeak:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case accounts.WrongPassword:
		return status.Errorf(codes.Unauthenticated, err.Error())
	case accounts.UserNotVerified,
		accounts.AtLeastOneOwner,
		accounts.TokenAlreadyUsed:
		return status.Errorf(codes.FailedPrecondition, err.Error())
	case accounts.UserAlreadyExists,
		accounts.EmailAlreadyUsed,
		accounts.OrganizationAlreadyExists,
		accounts.TeamAlreadyExists,
		accounts.ResourceAlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case accounts.UserNotFound,
		accounts.OrganizationNotFound,
		accounts.TeamNotFound,
		accounts.ResourceNotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case accounts.NotAuthorized:
		return status.Errorf(codes.PermissionDenied, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}

func getServerAddress(ctx context.Context) string {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return ""
	}
	authorities := md[":authority"]
	if len(authorities) == 0 {
		return ""
	}
	return authorities[0]
}

func (s *Server) getRequesterEmail(ctx context.Context) string {
	user, err := s.Accounts.GetUser(ctx, auth.GetUser(ctx))
	if err != nil {
		return ""
	}
	if user == nil {
		return ""
	}
	return user.Email
}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (*empty.Empty, error) {
	// Create user
	user, err := s.Accounts.CreateUser(ctx, in.Name, in.Email, in.Password)
	if err != nil {
		if errd := s.Accounts.DeleteNotVerifiedUser(ctx, in.Name); errd != nil {
			return nil, convertError(fmt.Errorf("Delete user error [%v] comming after CreateUser error [%v]", errd, err))
		}
		return nil, convertError(err)
	}

	switch s.Config.Registration {
	case configuration.RegistrationEmail:
		// Create a verification token valid for an hour
		token, err := s.Tokens.CreateVerificationToken(user.Name)
		if err != nil {
			if errd := s.Accounts.DeleteNotVerifiedUser(ctx, user.Name); errd != nil {
				return nil, convertError(fmt.Errorf("Delete user error [%v] comming after to VerificationToken error [%v]", errd, err))
			}
			return nil, convertError(err)
		}

		// Send the verification email
		if err := s.Mailer.SendAccountVerificationEmail(user.Email, user.Name, token, in.Url); err != nil {
			if errd := s.Accounts.DeleteNotVerifiedUser(ctx, user.Name); errd != nil {
				return nil, convertError(fmt.Errorf("Delete user error [%v] comming after SendAccountVerificationEmail error [%v]", errd, err))
			}
			return nil, err
		}
	}
	log.Println("Successfully created user", user.Name)
	return &empty.Empty{}, nil
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*empty.Empty, error) {
	// Validate the token
	log.Printf("verify token=%s\n", in.Token)
	claims, err := s.Tokens.ValidateToken(in.Token, auth.TokenTypeVerification)
	if err != nil {
		return nil, accounts.InvalidToken
	}
	if err := s.Accounts.VerifyUser(ctx, claims.AccountName); err != nil {
		return nil, convertError(err)
	}
	user, err := s.Accounts.GetUser(ctx, claims.AccountName)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, accounts.UserNotFound
	}
	if err := s.Mailer.SendAccountVerifiedEmail(user.Email, user.Name); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully verified user", user.Name)
	return &empty.Empty{}, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*LogInReply, error) {
	// Check password
	if err := s.Accounts.CheckUserPassword(ctx, in.Name, in.Password); err != nil {
		return nil, convertError(err)
	}
	// Create an authentication token valid for a day
	token, err := s.Tokens.CreateLoginToken(in.Name, "")
	if err != nil {
		return nil, convertError(err)
	}
	// Send the auth token to the client
	md := metadata.Pairs(auth.TokenKey, token)
	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully logged user in", in.Name)
	return &LogInReply{Auth: token}, nil //Angular issue with custom header
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (*empty.Empty, error) {
	// Get the user
	user, err := s.Accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", in.Name)
	}
	// Create a password reset token valid for an hour
	token, err := s.Tokens.CreatePasswordToken(user.Name)
	if err != nil {
		return nil, convertError(err)
	}
	// Send the password reset email
	if err := s.Mailer.SendAccountResetPasswordEmail(user.Email, user.Name, token, getServerAddress(ctx)); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully reset password for user", user.Name)
	return &empty.Empty{}, nil
}

// PasswordSet implements account.PasswordSet
func (s *Server) PasswordSet(ctx context.Context, in *PasswordSetRequest) (*empty.Empty, error) {
	// Validate token
	claims, err := s.Tokens.ValidateToken(in.Token, auth.TokenTypePassword)
	if err != nil {
		return nil, convertError(err)
	}
	// Sets the new password
	if err := s.Accounts.SetUserPassword(ctx, claims.AccountName, in.Password); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Println("Successfully set new password for user", claims.AccountName)
	return &empty.Empty{}, nil
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (*empty.Empty, error) {
	// Get requesting user
	requester := auth.GetUser(ctx)

	// Check the existing password
	if err := s.Accounts.CheckUserPassword(ctx, requester, in.ExistingPassword); err != nil {
		return nil, convertError(err)
	}
	// Sets the new password
	if err := s.Accounts.SetUserPassword(ctx, requester, in.NewPassword); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Println("Successfully updated the password for user", requester)
	return &empty.Empty{}, nil
}

// ForgotLogin implements account.PasswordChange
func (s *Server) ForgotLogin(ctx context.Context, in *ForgotLoginRequest) (*empty.Empty, error) {
	// Get the user
	user, err := s.Accounts.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", in.Email)
	}
	// Send the account name reminder email
	if err := s.Mailer.SendAccountNameReminderEmail(user.Email, user.Name); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully processed forgot login request for user", user.Name)
	return &empty.Empty{}, nil
}

// GetUser implements account.GetUser
func (s *Server) GetUser(ctx context.Context, in *GetUserRequest) (*GetUserReply, error) {
	// Get the user
	user, err := s.Accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", in.Name)
	}
	log.Println("Successfully retrieved user", user.Name)
	return &GetUserReply{User: user}, nil
}

// ListUsers implements account.ListUsers
func (s *Server) ListUsers(ctx context.Context, in *ListUsersRequest) (*ListUsersReply, error) {
	users, err := s.Accounts.ListUsers(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully list users")
	return &ListUsersReply{Users: users}, nil
}

// DeleteUser implements account.DeleteUser
func (s *Server) DeleteUser(ctx context.Context, in *DeleteUserRequest) (*empty.Empty, error) {
	// Get requesting user
	user, err := s.Accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", in.Name)
	}

	if err := s.Accounts.DeleteUser(ctx, in.Name); err != nil {
		return nil, convertError(err)
	}
	if err := s.Mailer.SendAccountRemovedEmail(user.Email, user.Name); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully deleted user", in.Name)
	return &empty.Empty{}, nil
}
