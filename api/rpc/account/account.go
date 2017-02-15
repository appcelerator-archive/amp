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

type accountClaims struct {
	AccountID string `json:"AccountID"`
	jwt.StandardClaims
}

// Server is used to implement account.AccountServer
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

	// Check if account already exists
	alreadyExists, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if alreadyExists != nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "account already exists")
	}

	// Hash password
	passwordHash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Create the new account
	account := &schema.Account{
		Email:        in.Email,
		UserName:     in.UserName,
		Type:         in.AccountType,
		IsVerified:   false,
		PasswordHash: passwordHash,
	}
	id, err := s.accounts.CreateAccount(ctx, account)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "storage error")
	}
	log.Println("Successfully created account", in.UserName)

	// Forge the verification token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accountClaims{
		id, // The token contains the account id to verify
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
	token, err := jwt.ParseWithClaims(in.Token, &accountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if !token.Valid {
		return &pb.Empty{}, grpc.Errorf(codes.InvalidArgument, "invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*accountClaims)
	if !ok {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, "invalid claims")
	}

	// Activate the account
	account, err := s.accounts.GetAccount(ctx, claims.AccountID)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	account.IsVerified = true
	if err := s.accounts.UpdateAccount(ctx, account); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully verified account", account.UserName)

	return &pb.Empty{}, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*LogInReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the account
	account, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "account not found")
	}
	if !account.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "account not verified")
	}

	// Check password
	_, err = passlib.Verify(in.Password, account.PasswordHash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	// Forge the authentication token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accountClaims{
		account.Id, // The token contains the account id
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
	log.Println("Successfully login for account", account.UserName)

	return &LogInReply{Token: ss}, nil
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (*PasswordResetReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the account
	account, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "account not found")
	}
	// TODO: Do we need the account to be verified?
	if !account.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "account not verified")
	}

	// Forge the password reset token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accountClaims{
		account.Id, // The token contains the account id to reset
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
	log.Println("Successfully reset password for account", account.UserName)

	// TODO: send password reset email with token

	return &PasswordResetReply{Token: ss}, nil
}

// PasswordSet implements account.PasswordSet
func (s *Server) PasswordSet(ctx context.Context, in *PasswordSetRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	token, err := jwt.ParseWithClaims(in.Token, &accountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if !token.Valid {
		return &pb.Empty{}, grpc.Errorf(codes.InvalidArgument, "invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*accountClaims)
	if !ok {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, "invalid claims")
	}

	// Get the account
	account, err := s.accounts.GetAccount(ctx, claims.AccountID)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}

	// Sets the new password
	passwordHash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	account.PasswordHash = passwordHash
	if err := s.accounts.UpdateAccount(ctx, account); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully set new password for account", account.UserName)

	return &pb.Empty{}, nil
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the account
	account, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "account not found")
	}
	if !account.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "account not verified")
	}

	// Check the existing password password
	_, err = passlib.Verify(in.ExistingPassword, account.PasswordHash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	// Sets the new password
	newPasswordHash, err := passlib.Hash(in.NewPassword)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	account.PasswordHash = newPasswordHash
	if err := s.accounts.UpdateAccount(ctx, account); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully updated the password for account", account.UserName)

	return &pb.Empty{}, nil
}
