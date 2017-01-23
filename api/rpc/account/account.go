package account

import (
	"fmt"

	context "golang.org/x/net/context"

	pb "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/hlandau/passlib.v1"
)

const hash = "$s2$16384$8$1$42JtddBgSqrJMwc3YuTNW+R+$ISfEF3jkvYQYk4AK/UFAxdqnmNFVeUw2gUVXEMBDAng=" // password

// Server is used to implement account.AccountServer
type Server struct{}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (out *SessionReply, err error) {
	out = &SessionReply{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	_, err = passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "hashing error")
	}
	out.SessionKey = in.Name
	return
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	_, err = passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "hashing error")
	}
	fmt.Println(in.Code)
	return
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	return
}

// CreateOrganization implements account.CreateOranization
func (s *Server) CreateOrganization(ctx context.Context, in *OrganizationRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	return
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (out *SessionReply, err error) {
	out = &SessionReply{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	_, err = passlib.Verify(in.Password, hash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}
	out.SessionKey = in.Name
	return

}

// Switch implements account.Switch
func (s *Server) Switch(ctx context.Context, in *TeamRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// ListAccounts implements account.ListAccounts
func (s *Server) ListAccounts(ctx context.Context, in *AccountsRequest) (out *ListReply, err error) {
	out = &ListReply{}
	if in.Type != "individual" && in.Type != "organization" {
		return nil, grpc.Errorf(codes.InvalidArgument, "account type is mandatory")
	}
	return
}

// GetAccountDetails implements account.GetAccountDetails
func (s *Server) GetAccountDetails(ctx context.Context, in *AccountRequest) (out *AccountReply, err error) {
	out = &AccountReply{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return
}

// EditAccount implements account.EditAccount
func (s *Server) EditAccount(ctx context.Context, in *EditRequest) (out *AccountReply, err error) {
	out = &AccountReply{}
	err = in.Validate()
	if err != nil {
		return
	}
	if in.NewPassword != "" {
		_, err := passlib.Verify(in.Password, hash)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
		}
		_, err = passlib.Hash(in.NewPassword)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "hashing error")
		}
	}
	return
}

// DeleteAccount implements account.DeleteAccount
func (s *Server) DeleteAccount(ctx context.Context, in *AccountRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return
}

// AddOrganizationMemberships implements account.AddOrganizationMemberships
func (s *Server) AddOrganizationMemberships(ctx context.Context, in *OrganizationMembershipsRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return
}

// DeleteOrganizationMemberships implements account.DeleteOrganizationMemberships
func (s *Server) DeleteOrganizationMemberships(ctx context.Context, in *OrganizationMembershipsRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	return
}

// CreateTeam implements account.CreateTeam
func (s *Server) CreateTeam(ctx context.Context, in *TeamRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// ListTeams implements account.ListTeams
func (s *Server) ListTeams(ctx context.Context, in *TeamRequest) (out *ListReply, err error) {
	out = &ListReply{}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// EditTeam implements account.EditTeam
func (s *Server) EditTeam(ctx context.Context, in *TeamRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// GetTeamDetails implements account.GetTeamDetails
func (s *Server) GetTeamDetails(ctx context.Context, in *TeamRequest) (out *TeamReply, err error) {
	out = &TeamReply{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// DeleteTeam implements account.DeleteTeam
func (s *Server) DeleteTeam(ctx context.Context, in *TeamRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// AddTeamMemberships implements account.AddTeamMemberships
func (s *Server) AddTeamMemberships(ctx context.Context, in *TeamMembershipsRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
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
func (s *Server) DeleteTeamMemberships(ctx context.Context, in *TeamMembershipsRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
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
func (s *Server) GrantPermission(ctx context.Context, in *PermissionRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
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
	return
}

// ListPermissions implements account.ListPermissions
func (s *Server) ListPermissions(ctx context.Context, in *PermissionRequest) (out *ListReply, err error) {
	out = &ListReply{}
	if in.Team != "" && in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// EditPermission implements account.EditPermission
func (s *Server) EditPermission(ctx context.Context, in *PermissionRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
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
	return
}

// RevokePermission implements account.RevokePermission
func (s *Server) RevokePermission(ctx context.Context, in *PermissionRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Team == "" {
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
func (s *Server) TransferOwnership(ctx context.Context, in *PermissionRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Team == "" {
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
