package account

import (
	"github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/dgrijalva/jwt-go"
	pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/hlandau/passlib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"os"
	"time"
)

// TODO: this MUST NOT be public
// TODO: find a way to store this key secretly
var secretKey = []byte("&kv@l3go-f=@^*@ush0(o5*5utxe6932j9di+ume=$mkj%d&&9*%k53(bmpksf&!c2&zpw$z=8ndi6ib)&nxms0ia7rf*sj9g8r4")

type userClaims struct {
	UserID string `json:"UserID"`
	jwt.StandardClaims
}

// Server is used to implement account.UserServer
type Server struct {
	accounts account.Interface
}

// NewServer instantiates account.Server
func NewServer(store storage.Interface) *Server {
	return &Server{accounts: account.NewStore(store)}
}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (*SignUpReply, error) {
	// Validate input
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Check if user already exists
	alreadyExists, err := s.accounts.GetUserByName(ctx, in.Name)
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
	id, err := s.accounts.CreateUser(ctx, user)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "storage error")
	}
	log.Println("Successfully created user", in.Name)

	// Forge the verification token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		id, // The token contains the user id to verify
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// TODO: send confirmation email with token

	return &SignUpReply{Token: ss}, nil
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	token, err := jwt.ParseWithClaims(in.Token, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if !token.Valid {
		return &pb.Empty{}, grpc.Errorf(codes.InvalidArgument, "invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*userClaims)
	if !ok {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, "invalid claims")
	}

	// Activate the user
	user, err := s.accounts.GetUser(ctx, claims.UserID)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	user.IsVerified = true
	if err := s.accounts.UpdateUser(ctx, user); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully verified user", user.Name)

	return &pb.Empty{}, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*LogInReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.accounts.GetUserByName(ctx, in.Name)
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

	// Forge the authentication token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		user.Id, // The token contains the user id
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully login for user", user.Name)

	return &LogInReply{Token: ss}, nil
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (*PasswordResetReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.accounts.GetUserByName(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	// TODO: Do we need the user to be verified?
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Forge the password reset token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		user.Id, // The token contains the user id to reset
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully reset password for user", user.Name)

	// TODO: send password reset email with token

	return &PasswordResetReply{Token: ss}, nil
}

// PasswordSet implements account.PasswordSet
func (s *Server) PasswordSet(ctx context.Context, in *PasswordSetRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	token, err := jwt.ParseWithClaims(in.Token, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if !token.Valid {
		return &pb.Empty{}, grpc.Errorf(codes.InvalidArgument, "invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*userClaims)
	if !ok {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, "invalid claims")
	}

	// Get the user
	user, err := s.accounts.GetUser(ctx, claims.UserID)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
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

	// Get the user
	user, err := s.accounts.GetUserByName(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	if !user.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user not verified")
	}

	// Check the existing password password
	_, err = passlib.Verify(in.ExistingPassword, user.PasswordHash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	// Sets the new password
	newPasswordHash, err := passlib.Hash(in.NewPassword)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	user.PasswordHash = newPasswordHash
	if err := s.accounts.UpdateUser(ctx, user); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully updated the password for user", user.Name)

	return &pb.Empty{}, nil
}

// CreateOrganization implements account.CreateOrganization
func (s *Server) CreateOrganization(ctx context.Context, in *CreateOrganizationRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
