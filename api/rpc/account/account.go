package account

import (
	"fmt"
	"strconv"

	context "golang.org/x/net/context"

	"github.com/appcelerator/amp/data/schema"
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
	account := &schema.Account{
		Id:         strconv.Itoa(len(accounts)),
		Name:       in.Name,
		Type:       schema.AccountType_USER,
		Email:      in.Email,
		PwHashcode: hash,
		IsVerified: false,
	}
	accounts = append(accounts, account)
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
	for _, account := range accounts {
		if account.Name == in.Name {
			account.IsVerified = true
		}
	}
	fmt.Println(in.Code)
	return
}

// CreateOrganization implements account.CreateOranization
func (s *Server) CreateOrganization(ctx context.Context, in *OrganizationRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	organizationID := strconv.Itoa(len(accounts))
	account := &schema.Account{
		Id:         organizationID,
		Name:       in.Name,
		Type:       schema.AccountType_ORGANIZATION,
		Email:      in.Email,
		IsVerified: false,
	}
	accounts = append(accounts, account)
	team := &schema.Team{
		Id:           strconv.Itoa(len(teams)),
		OrgAccountId: organizationID,
		Name:         "owners",
		Desc:         in.Name + " owners team",
	}
	teams = append(teams, team)
	// create organization membership
	// create team membership
	return
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (out *SessionReply, err error) {
	out = &SessionReply{}
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	var account *schema.Account
	for _, a := range accounts {
		if a.Name == in.Name {
			account = a
		}
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	_, err = passlib.Verify(in.Password, account.PwHashcode)
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
	for _, a := range accounts {
		if a.Name == in.Organization {
			return
		}
	}
	return nil, grpc.Errorf(codes.NotFound, "organization not found")
}

// ListAccounts implements account.ListAccounts
func (s *Server) ListAccounts(ctx context.Context, in *AccountsRequest) (out *ListReply, err error) {
	out = &ListReply{}
	if in.Type != "individual" && in.Type != "organization" {
		return nil, grpc.Errorf(codes.InvalidArgument, "account type is mandatory")
	}
	for _, account := range accounts {
		if in.Type == "individual" && account.Type == schema.AccountType_USER {
			out.Accounts = append(out.Accounts, &Account{
				Name: account.Name,
			})
		}
		if in.Type == "organization" && account.Type == schema.AccountType_ORGANIZATION {
			out.Accounts = append(out.Accounts, &Account{
				Name: account.Name,
			})
		}
	}
	return
}

// GetAccountDetails implements account.GetAccountDetails
func (s *Server) GetAccountDetails(ctx context.Context, in *AccountRequest) (out *AccountReply, err error) {
	out = &AccountReply{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	for _, account := range accounts {
		if account.Name == in.Name {
			out.Account.Name = account.Name
			out.Account.Email = account.Email
			return
		}
	}
	return nil, grpc.Errorf(codes.NotFound, "account not found")
}

// EditAccount implements account.EditAccount
func (s *Server) EditAccount(ctx context.Context, in *EditRequest) (out *AccountReply, err error) {
	out = &AccountReply{}
	err = in.Validate()
	if err != nil {
		return
	}
	var account *schema.Account
	for _, a := range accounts {
		if a.Name == in.Name {
			account = a
		}
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "account not found")
	}
	if in.NewPassword != "" {
		_, err := passlib.Verify(in.Password, account.PwHashcode)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
		}
		hash, err := passlib.Hash(in.NewPassword)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "hashing error")
		}
		account.PwHashcode = hash
	}
	if in.Email != "" {
		account.Email = in.Email
	}
	out.Account.Name = account.Name
	out.Account.Email = account.Email
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
