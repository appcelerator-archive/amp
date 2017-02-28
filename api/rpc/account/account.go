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

func convertError(err error) error {
	switch err {
	case schema.InvalidName:
	case schema.InvalidEmail:
	case schema.InvalidToken:
	case schema.PasswordTooWeak:
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	case schema.WrongPassword:
	case schema.UserNotVerified:
	case schema.AtLeastOneOwner:
		return grpc.Errorf(codes.FailedPrecondition, err.Error())
	case schema.UserAlreadyExists:
	case schema.OrganizationAlreadyExists:
	case schema.TeamAlreadyExists:
		return grpc.Errorf(codes.AlreadyExists, err.Error())
	case schema.UserNotFound:
	case schema.OrganizationNotFound:
	case schema.TeamNotFound:
		return grpc.Errorf(codes.NotFound, err.Error())
	case schema.NotAuthorized:
		return grpc.Errorf(codes.PermissionDenied, err.Error())
	}
	return grpc.Errorf(codes.Internal, err.Error())
}

// Users

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (*pb.Empty, error) {
	// Create user
	user, err := s.accounts.CreateUser(ctx, in.Name, in.Email, in.Password)
	if err != nil {
		return nil, convertError(err)
	}
	// Create a verification token valid for an hour
	token, err := auth.CreateToken(user.Name, auth.TokenTypeVerify, time.Hour)
	if err != nil {
		return nil, convertError(err)
	}
	// Send the verification email
	if err := mail.SendAccountVerificationEmail(user.Email, user.Name, token); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully created user", user.Name)
	return &pb.Empty{}, nil
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*VerificationReply, error) {
	user, err := s.accounts.VerifyUser(ctx, in.Token)
	if err != nil {
		return nil, convertError(err)
	}
	if err := mail.SendAccountCreatedEmail(user.Email, user.Name); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully verified user", user.Name)
	return &VerificationReply{Reply: fmt.Sprintf("Account %s is ready", user.Name)}, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*pb.Empty, error) {
	// Check password
	if err := s.accounts.CheckUserPassword(ctx, in.Name, in.Password); err != nil {
		return nil, convertError(err)
	}
	// Create an authentication token valid for a day
	token, err := auth.CreateToken(in.Name, auth.TokenTypeLogin, 24*time.Hour)
	if err != nil {
		return nil, convertError(err)
	}
	// Send the authN token to the client
	md := metadata.Pairs(auth.TokenKey, token)
	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully logged user in", in.Name)
	return &pb.Empty{}, nil
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (*pb.Empty, error) {
	// Get the user
	user, err := s.accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found: %s", in.Name)
	}
	// Create a password reset token valid for an hour
	token, err := auth.CreateToken(user.Name, auth.TokenTypePassword, time.Hour)
	if err != nil {
		return nil, convertError(err)
	}
	// Send the password reset email
	if err := mail.SendAccountResetPasswordEmail(user.Email, user.Name, token); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully reset password for user", user.Name)
	return &pb.Empty{}, nil
}

// PasswordSet implements account.PasswordSet
func (s *Server) PasswordSet(ctx context.Context, in *PasswordSetRequest) (*pb.Empty, error) {
	// Validate token
	claims, err := auth.ValidateToken(in.Token, auth.TokenTypePassword)
	if err != nil {
		return nil, convertError(err)
	}
	// Sets the new password
	if err := s.accounts.SetUserPassword(ctx, claims.AccountName, in.Password); err != nil {
		return &pb.Empty{}, convertError(err)
	}
	log.Println("Successfully set new password for user", claims.AccountName)
	return &pb.Empty{}, nil
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (*pb.Empty, error) {
	// Get requester
	requester, err := auth.GetRequesterName(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	// Check the existing password password
	if err := s.accounts.CheckUserPassword(ctx, requester, in.ExistingPassword); err != nil {
		return nil, convertError(err)
	}
	// Sets the new password
	if err := s.accounts.SetUserPassword(ctx, requester, in.NewPassword); err != nil {
		return &pb.Empty{}, convertError(err)
	}
	log.Println("Successfully updated the password for user", requester)
	return &pb.Empty{}, nil
}

// ForgotLogin implements account.PasswordChange
func (s *Server) ForgotLogin(ctx context.Context, in *ForgotLoginRequest) (*pb.Empty, error) {
	// Get the user
	user, err := s.accounts.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found: %s", in.Email)
	}
	// Send the account name reminder email
	if err := mail.SendAccountNameReminderEmail(user.Email, user.Name); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully processed forgot login request for user", user.Name)
	return &pb.Empty{}, nil
}

// GetUser implements account.GetUser
func (s *Server) GetUser(ctx context.Context, in *GetUserRequest) (*GetUserReply, error) {
	// Get the user
	user, err := s.accounts.GetUser(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found: %s", in.Name)
	}
	log.Println("Successfully retrieved user", user.Name)
	return &GetUserReply{User: user}, nil
}

// ListUsers implements account.ListUsers
func (s *Server) ListUsers(ctx context.Context, in *ListUsersRequest) (*ListUsersReply, error) {
	users, err := s.accounts.ListUsers(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully list users")
	return &ListUsersReply{Users: users}, nil
}

// Organizations

// CreateOrganization implements account.CreateOrganization
func (s *Server) CreateOrganization(ctx context.Context, in *CreateOrganizationRequest) (*pb.Empty, error) {
	if err := s.accounts.CreateOrganization(ctx, in.Name, in.Email); err != nil {
		return nil, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendOrganizationCreatedEmail(email, in.Name); err != nil {
			return nil, convertError(err)
		}
	}
	log.Println("Successfully created organization", in.Name)
	return &pb.Empty{}, nil
}

// AddUserToOrganization implements account.AddOrganizationMember
func (s *Server) AddUserToOrganization(ctx context.Context, in *AddUserToOrganizationRequest) (*pb.Empty, error) {
	if err := s.accounts.AddUserToOrganization(ctx, in.OrganizationName, in.UserName); err != nil {
		return &pb.Empty{}, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendUserAddedInOrganizationEmail(email, in.OrganizationName, in.UserName); err != nil {
			return nil, convertError(err)
		}
	}
	log.Printf("Successfully added member %s to organization %s\n", in.UserName, in.OrganizationName)
	return &pb.Empty{}, nil
}

// RemoveUserFromOrganization implements account.RemoveOrganizationMember
func (s *Server) RemoveUserFromOrganization(ctx context.Context, in *RemoveUserFromOrganizationRequest) (*pb.Empty, error) {
	if err := s.accounts.RemoveUserFromOrganization(ctx, in.OrganizationName, in.UserName); err != nil {
		return &pb.Empty{}, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendUserRemovedFromOrganizationEmail(email, in.OrganizationName, in.UserName); err != nil {
			return nil, convertError(err)
		}
	}
	log.Printf("Successfully removed user %s from organization %s\n", in.UserName, in.OrganizationName)
	return &pb.Empty{}, nil
}

// DeleteOrganization implements account.DeleteOrganization
func (s *Server) DeleteOrganization(ctx context.Context, in *DeleteOrganizationRequest) (*pb.Empty, error) {
	if err := s.accounts.DeleteOrganization(ctx, in.Name); err != nil {
		return nil, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendOrganizationRemovedEmail(email, in.Name); err != nil {
			return nil, convertError(err)
		}
	}
	log.Println("Successfully deleted organization", in.Name)
	return &pb.Empty{}, nil
}

// GetOrganization implements account.GetOrganization
func (s *Server) GetOrganization(ctx context.Context, in *GetOrganizationRequest) (*GetOrganizationReply, error) {
	organization, err := s.accounts.GetOrganization(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if organization == nil {
		return nil, grpc.Errorf(codes.NotFound, "organization not found: %s", in.Name)
	}
	log.Println("Successfully retrieved organization", organization.Name)
	return &GetOrganizationReply{Organization: organization}, nil
}

// ListOrganizations implements account.ListOrganizations
func (s *Server) ListOrganizations(ctx context.Context, in *ListOrganizationsRequest) (*ListOrganizationsReply, error) {
	organizations, err := s.accounts.ListOrganizations(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully list organizations")
	return &ListOrganizationsReply{Organizations: organizations}, nil
}

// Teams

// CreateTeam implements account.CreateTeam
func (s *Server) CreateTeam(ctx context.Context, in *CreateTeamRequest) (*pb.Empty, error) {
	if err := s.accounts.CreateTeam(ctx, in.OrganizationName, in.TeamName); err != nil {
		return nil, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendTeamCreatedEmail(email, in.TeamName); err != nil {
			return nil, convertError(err)
		}
	}
	log.Printf("Successfully created team %s in organization %s\n", in.TeamName, in.OrganizationName)
	return &pb.Empty{}, nil
}

// AddUserToTeam implements account.AddUserToTeam
func (s *Server) AddUserToTeam(ctx context.Context, in *AddUserToTeamRequest) (*pb.Empty, error) {
	if err := s.accounts.AddUserToTeam(ctx, in.OrganizationName, in.TeamName, in.UserName); err != nil {
		return &pb.Empty{}, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendUserAddedInTeamEmail(email, in.TeamName, in.UserName); err != nil {
			return nil, convertError(err)
		}
	}
	log.Printf("Successfully added member %s to team %s in organization %s\n", in.UserName, in.TeamName, in.OrganizationName)
	return &pb.Empty{}, nil
}

// RemoveUserFromTeam implements account.RemoveUserFromTeam
func (s *Server) RemoveUserFromTeam(ctx context.Context, in *RemoveUserFromTeamRequest) (*pb.Empty, error) {
	if err := s.accounts.RemoveUserFromTeam(ctx, in.OrganizationName, in.TeamName, in.UserName); err != nil {
		return &pb.Empty{}, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendUserRemovedFromTeamEmail(email, in.TeamName, in.UserName); err != nil {
			return nil, convertError(err)
		}
	}
	log.Printf("Successfully removed user %s from teams %s in organization %s\n", in.UserName, in.TeamName, in.OrganizationName)
	return &pb.Empty{}, nil
}

// DeleteTeam implements account.DeleteTeam
func (s *Server) DeleteTeam(ctx context.Context, in *DeleteTeamRequest) (*pb.Empty, error) {
	if err := s.accounts.DeleteTeam(ctx, in.OrganizationName, in.TeamName); err != nil {
		return nil, convertError(err)
	}
	// Send confirmation email
	if email := s.getUserEmail(ctx); email != "" {
		if err := mail.SendTeamRemovedEmail(email, in.TeamName); err != nil {
			return nil, convertError(err)
		}
	}
	log.Printf("Successfully deleted team %s from organization %s\n", in.TeamName, in.OrganizationName)
	return &pb.Empty{}, nil
}

// GetTeam implements account.GetTeam
func (s *Server) GetTeam(ctx context.Context, in *GetTeamRequest) (*GetTeamReply, error) {
	team, err := s.accounts.GetTeam(ctx, in.OrganizationName, in.TeamName)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully retrieved team", team.Name)
	return &GetTeamReply{Team: team}, nil
}

// ListTeams implements account.ListTeams
func (s *Server) ListTeams(ctx context.Context, in *ListTeamsRequest) (*ListTeamsReply, error) {
	teams, err := s.accounts.ListTeams(ctx, in.OrganizationName)
	if err != nil {
		return nil, err
	}
	return &ListTeamsReply{Teams: teams}, nil
}

func (s *Server) getUserEmail(ctx context.Context) string {
	name, err := auth.GetRequesterName(ctx)
	if err != nil {
		return ""
	}
	user, errm := s.accounts.GetUser(ctx, name)
	if errm != nil {
		return ""
	}
	return user.Email
}
