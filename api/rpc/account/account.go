package account

import (
	"context"
	"fmt"
	"net/mail"

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
	return nil, &SessionReply{
		SessionKey: in.Name,
	}
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*google_protobuf1.Empty, error) {
	if len(in.Code) != codeLength {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid verification code")
	}
	fmt.Println(in.Code)
	return
}

// CreateOrganization implements account.CreateOranization
func CreateOrganization(ctx context.Context, in *OrganizationRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	address, err := mail.ParseAddress(in.Email)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	fmt.Println(in.Name, address)
	return
}

// Login implements account.Login
func Login(ctx context.Context, in *LogInRequest) (*SessionReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	_, err := passlib.Verify(in.Password, "hash")
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}
	return nil, &SessionReply{
		SessionKey: in.Name,
	}
}

// Switch implements account.Switch
func Switch(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// ListAccounts implements account.ListAccounts
func ListAccounts(ctx context.Context, in *AccountsRequest) (*ListReply, error) {
	if in.Type != "individual" && in.Type != "organization" {
		return nil, grpc.Errorf(codes.InvalidArgument, "account type is mandatory")
	}
	return
}

// GetAccountDetails implements account.GetAccountDetails
func GetAccountDetails(ctx context.Context, in *AccountRequest) (*AccountReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return
}

// EditAccount implements account.EditAccount
func EditAccount(ctx context.Context, in *EditRequest) (*AccountReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	if in.Email != "" {
		address, err := mail.ParseAddress(in.Email)
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
		hash, err := passlib.Hash(in.NewPassword)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "hashing error")
		}
	}
	return
}

// DeleteAccount implements account.DeleteAccount
func DeleteAccount(ctx context.Context, in *AccountRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return
}

// AddOrganizationMemberships implements account.AddOrganizationMemberships
func AddOrganizationMemberships(ctx context.Context, in *OrganizationMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return
}

// DeleteOrganizationMemberships implements account.DeleteOrganizationMemberships
func DeleteOrganizationMemberships(ctx context.Context, in *OrganizationMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return
}

// CreateTeam implements account.CreateTeam
func CreateTeam(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// ListTeams implements account.ListTeams
func ListTeams(ctx context.Context, in *TeamRequest) (*ListReply, error) {
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// EditTeam implements account.EditTeam
func EditTeam(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// GetTeamDetails implements account.GetTeamDetails
func GetTeamDetails(ctx context.Context, in *TeamRequest) (*TeamReply, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// DeleteTeam implements account.DeleteTeam
func DeleteTeam(ctx context.Context, in *TeamRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// AddTeamMemberships implements account.AddTeamMemberships
func AddTeamMemberships(ctx context.Context, in *TeamMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return
}

// DeleteTeamMemberships implements account.DeleteTeamMemberships
func DeleteTeamMemberships(ctx context.Context, in *TeamMembershipsRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return
}

// GrantPermission implements account.GrantPermission
func GrantPermission(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
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
	return
}

// ListPermissions implements account.ListPermissions
func ListPermissions(ctx context.Context, in *PermissionRequest) (*ListReply, error) {
	if in.Name != "" && in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// EditPermission implements account.EditPermission
func EditPermission(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
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
	return
}

// RevokePermission implements account.RevokePermission
func RevokePermission(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	return
}

// TransferOwnership implements account.TransferOwnership
func TransferOwnership(ctx context.Context, in *PermissionRequest) (*google_protobuf1.Empty, error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	return
}
