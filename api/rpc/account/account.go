package account

import (
	"fmt"
	"log"
	"time"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/mail"
	pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/hlandau/passlib"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// Server is used to implement account.UserServer
type Server struct {
	accounts account.Interface
}

// NewServer instantiates account.Server
func NewServer(store storage.Interface) *Server {
	return &Server{accounts: account.NewStore(store)}
}

// Users

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Check if user already exists
	alreadyExists, err := s.accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if alreadyExists != nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "user already exists")
	}

	// Hash password
	passwordHash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Create the new user
	user := &schema.User{
		Email:        in.Email,
		Name:         in.Name,
		IsVerified:   false,
		PasswordHash: passwordHash,
	}
	if err := s.accounts.CreateUser(ctx, user); err != nil {
		return nil, grpc.Errorf(codes.Internal, "storage error")
	}

	// Create a verification token valid for an hour
	token, err := auth.CreateUserToken(user.Name, time.Hour)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Send the verification email
	if err := mail.SendAccountVerificationEmail(user.Email, user.Name, token); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully created user", in.Name)

	return &pb.Empty{}, nil
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*VerificationReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	claims, err := auth.ValidateUserToken(in.Token)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "VadidateUsetToken: %v", err.Error())
	}

	// Get the user
	user, err := s.accounts.GetUser(ctx, claims.AccountName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}

	// Activate the user
	user.IsVerified = true
	if err := s.accounts.UpdateUser(ctx, user); err != nil {
		return &VerificationReply{}, grpc.Errorf(codes.Internal, err.Error())
	}
	// TODO: We probably need to send an email ...
	log.Println("Successfully verified user", user.Name)

	return &VerificationReply{
		Reply: fmt.Sprintf("Account %s is ready", user.Name),
	}, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Check password
	_, err = passlib.Verify(in.Password, user.PasswordHash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	// Create an authentication token valid for a day
	token, err := auth.CreateUserToken(user.Name, 24*time.Hour)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Send the authN token to the client
	md := metadata.Pairs(auth.TokenKey, token)
	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully logged user in", user.Name)

	return &pb.Empty{}, nil
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Create a password reset token valid for an hour
	token, err := auth.CreateUserToken(user.Name, time.Hour)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Send the password reset email
	if err := mail.SendAccountResetPasswordEmail(user.Email, user.Name, token); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully reset password for user", user.Name)

	return &pb.Empty{}, nil
}

// PasswordSet implements account.PasswordSet
func (s *Server) PasswordSet(ctx context.Context, in *PasswordSetRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	claims, err := auth.ValidateUserToken(in.Token)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Get the user
	user, err := s.accounts.GetUser(ctx, claims.AccountName)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Sets the new password
	passwordHash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	user.PasswordHash = passwordHash
	if err := s.accounts.UpdateUser(ctx, user); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully set new password for user", user.Name)

	return &pb.Empty{}, nil
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the requester
	requesterName, err := auth.GetRequesterName(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	requester, err := s.accounts.GetUser(ctx, requesterName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if requester == nil {
		return nil, grpc.Errorf(codes.NotFound, "requester not found")
	}
	if !requester.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "requester not verified")
	}

	// Check the existing password password
	_, err = passlib.Verify(in.ExistingPassword, requester.PasswordHash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	// Sets the new password
	newPasswordHash, err := passlib.Hash(in.NewPassword)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	requester.PasswordHash = newPasswordHash
	if err := s.accounts.UpdateUser(ctx, requester); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully updated the password for user", requester.Name)

	return &pb.Empty{}, nil
}

// ForgotLogin implements account.PasswordChange
func (s *Server) ForgotLogin(ctx context.Context, in *ForgotLoginRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.accounts.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}

	// Send the account name reminder email
	if err := mail.SendAccountNameReminderEmail(user.Email, user.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully processed forgot login request for user", user.Name)

	return &pb.Empty{}, nil
}

// GetUser implements account.GetUser
func (s *Server) GetUser(ctx context.Context, in *GetUserRequest) (*GetUserReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	log.Println("Successfully retrieved user", user.Name)

	return &GetUserReply{User: FromSchema(user)}, nil
}

// ListUsers implements account.ListUsers
func (s *Server) ListUsers(ctx context.Context, in *ListUsersRequest) (*ListUsersReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// List users
	users, err := s.accounts.ListUsers(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	reply := &ListUsersReply{}
	for _, user := range users {
		reply.Users = append(reply.Users, FromSchema(user))
	}
	log.Println("Successfully list users")

	return reply, nil
}

// Organizations

// CreateOrganization implements account.CreateOrganization
func (s *Server) CreateOrganization(ctx context.Context, in *CreateOrganizationRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the requester
	requesterName, err := auth.GetRequesterName(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	requester, err := s.accounts.GetUser(ctx, requesterName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if requester == nil {
		return nil, grpc.Errorf(codes.NotFound, "requester not found")
	}
	if !requester.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "requester not verified")
	}

	// Check if organization already exists
	alreadyExists, err := s.accounts.GetOrganization(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if alreadyExists != nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "organization already exists")
	}

	// Create the new organization
	organization := &schema.Organization{
		Email: in.Email,
		Name:  in.Name,
		Members: []*schema.OrganizationMember{
			{
				Name: requester.Name,
				Role: schema.OrganizationRole_OWNER,
			},
		},
	}
	if err := s.accounts.CreateOrganization(ctx, organization); err != nil {
		return nil, grpc.Errorf(codes.Internal, "storage error")
	}
	// TODO: We probably need to send an email ...
	log.Println("Successfully created organization", in.Name)

	return &pb.Empty{}, nil
}

// AddUserToOrganization implements account.AddOrganizationMember
func (s *Server) AddUserToOrganization(ctx context.Context, in *AddUserToOrganizationRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the requester
	requesterName, err := auth.GetRequesterName(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	requester, err := s.accounts.GetUser(ctx, requesterName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if requester == nil {
		return nil, grpc.Errorf(codes.NotFound, "requester not found")
	}
	if !requester.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "requester not verified")
	}

	// Check if organization exists
	organization, err := s.accounts.GetOrganization(ctx, in.OrganizationName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if organization == nil {
		return nil, grpc.Errorf(codes.NotFound, "organization not found")
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  requester.Name,
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"owners": organization.GetOwners(),
		},
	}); err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
	}

	// Check if the user exists
	user, err := s.accounts.GetUser(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Add the new member
	if err := s.accounts.AddUserToOrganization(ctx, organization, user); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	// TODO: We probably need to send an email ...
	log.Printf("Successfully added member %s to organization %s\n", in.UserName, organization.Name)

	return &pb.Empty{}, nil
}

// RemoveUserFromOrganization implements account.RemoveOrganizationMember
func (s *Server) RemoveUserFromOrganization(ctx context.Context, in *RemoveUserFromOrganizationRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the requester
	requesterName, err := auth.GetRequesterName(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	requester, err := s.accounts.GetUser(ctx, requesterName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if requester == nil {
		return nil, grpc.Errorf(codes.NotFound, "requester not found")
	}
	if !requester.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "requester not verified")
	}

	// Check if organization exists
	organization, err := s.accounts.GetOrganization(ctx, in.OrganizationName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if organization == nil {
		return nil, grpc.Errorf(codes.NotFound, "organization not found")
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  requester.Name,
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"owners": organization.GetOwners(),
		},
	}); err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
	}

	// Check if the user exists
	user, err := s.accounts.GetUser(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Remove the existing member
	if err := s.accounts.RemoveUserFromOrganization(ctx, organization, user); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	// TODO: We probably need to send an email ...
	log.Printf("Successfully removed user %s from organization %s\n", in.UserName, organization.Name)

	return &pb.Empty{}, nil
}

// DeleteOrganization implements account.DeleteOrganization
func (s *Server) DeleteOrganization(ctx context.Context, in *DeleteOrganizationRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the requester
	requesterName, err := auth.GetRequesterName(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	requester, err := s.accounts.GetUser(ctx, requesterName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if requester == nil {
		return nil, grpc.Errorf(codes.NotFound, "requester not found")
	}
	if !requester.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "requester not verified")
	}

	// Check if organization exists
	organization, err := s.accounts.GetOrganization(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if organization == nil {
		return nil, grpc.Errorf(codes.NotFound, "organization not found")
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  requester.Name,
		Action:   auth.DeleteAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"owners": organization.GetOwners(),
		},
	}); err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
	}

	if err := s.accounts.DeleteOrganization(ctx, organization.Id); err != nil {
		return nil, grpc.Errorf(codes.Internal, "storage error")
	}
	// TODO: We probably need to send an email ...
	log.Println("Successfully deleted organization", in.Name)

	return &pb.Empty{}, nil
}
