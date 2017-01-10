package account

import (
	"fmt"
	"net/mail"

	context "golang.org/x/net/context"

	google_protobuf1 "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/hlandau/passlib.v1"
)

const codeLength = 8

// Server is used to implement account.AccountServer
type Server struct{}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (*SessionReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "user name is mandatory")
	}
	address, err := mail.ParseAddress(in.Email)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	if len(in.Password) < 8 {
		return nil, grpc.Errorf(codes.InvalidArgument, "password too weak")
	}
	hash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "hashing error")
	}
	fmt.Println(in.Name, address, hash)
	return &SessionReply{
		SessionKey: in.Name,
	}, nil
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*google_protobuf1.Empty, error) {
	if len(in.Code) != codeLength {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid verification code")
	}
	fmt.Println(in.Code)
	return nil, nil
}

// CreateOrganization implements account.CreateOranization
func (s *Server) CreateOrganization(ctx context.Context, in *OrganizationRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	address, err := mail.ParseAddress(in.Email)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	fmt.Println(in.Name, address)
	return nil, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*SessionReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	_, err := passlib.Verify(in.Password, "hash")
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}
	return &SessionReply{
		SessionKey: in.Name,
	}, nil
}

// Switch implements account.Switch
func (s *Server) Switch(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// ListAccounts implements account.ListAccounts
func (s *Server) ListAccounts(ctx context.Context, in *AccountsRequest) (*ListReply, error) {
	if in.Type != "individual" && in.Type != "organization" {
		return nil, grpc.Errorf(codes.InvalidArgument, "account type is mandatory")
	}
	return nil, nil
}

// GetAccountDetails implements account.GetAccountDetails
func (s *Server) GetAccountDetails(ctx context.Context, in *AccountRequest) (*AccountReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return nil, nil
}

// EditAccount implements account.EditAccount
func (s *Server) EditAccount(ctx context.Context, in *EditRequest) (*AccountReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	if in.Email != "" {
		_, err := mail.ParseAddress(in.Email)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}
	if in.NewPassword != "" {
		_, err := passlib.Verify(in.Password, "hash")
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
		}
		if len(in.NewPassword) < 8 {
			return nil, grpc.Errorf(codes.InvalidArgument, "password too weak")
		}
		_, err = passlib.Hash(in.NewPassword)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "hashing error")
		}
	}
	return nil, nil
}

// DeleteAccount implements account.DeleteAccount
func (s *Server) DeleteAccount(ctx context.Context, in *AccountRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return nil, nil
}

// AddOrganizationMemberships implements account.AddOrganizationMemberships
func (s *Server) AddOrganizationMemberships(ctx context.Context, in *OrganizationMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return nil, nil
}

// DeleteOrganizationMemberships implements account.DeleteOrganizationMemberships
func (s *Server) DeleteOrganizationMemberships(ctx context.Context, in *OrganizationMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return nil, nil
}

// CreateTeam implements account.CreateTeam
func (s *Server) CreateTeam(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// ListTeams implements account.ListTeams
func (s *Server) ListTeams(ctx context.Context, in *TeamRequest) (*ListReply, error) {
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// EditTeam implements account.EditTeam
func (s *Server) EditTeam(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// GetTeamDetails implements account.GetTeamDetails
func (s *Server) GetTeamDetails(ctx context.Context, in *TeamRequest) (*TeamReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// DeleteTeam implements account.DeleteTeam
func (s *Server) DeleteTeam(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// AddTeamMemberships implements account.AddTeamMemberships
func (s *Server) AddTeamMemberships(ctx context.Context, in *TeamMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return nil, nil
}

// DeleteTeamMemberships implements account.DeleteTeamMemberships
func (s *Server) DeleteTeamMemberships(ctx context.Context, in *TeamMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return nil, nil
}

// GrantPermission implements account.GrantPermission
func (s *Server) GrantPermission(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Team == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.Level == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "permission level is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	return nil, nil
}

// ListPermissions implements account.ListPermissions
func (s *Server) ListPermissions(ctx context.Context, in *PermissionRequest) (*ListReply, error) {
	if in.Team != "" && in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil, nil
}

// EditPermission implements account.EditPermission
func (s *Server) EditPermission(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Team == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.Level == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "permission level is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	return nil, nil
}

// RevokePermission implements account.RevokePermission
func (s *Server) RevokePermission(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Team == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	return nil, nil
}

// TransferOwnership implements account.TransferOwnership
func (s *Server) TransferOwnership(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Team == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	return nil, nil
}
