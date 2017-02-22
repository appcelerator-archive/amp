package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

// UserClaims represents user claims
type UserClaims struct {
	AccountName string `json:"AccountName"`
	jwt.StandardClaims
}

// LoginCredentials represents login credentials
type LoginCredentials struct {
	Token string
}

// GetRequestMetadata implements credentials.PerRPCCredentials
func (c *LoginCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		TokenKey: c.Token,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (c *LoginCredentials) RequireTransportSecurity() bool {
	return false
}

// Keys used in context metadata
const (
	TokenKey     = "amp.token"
	RequesterKey = "amp.requester"
)

var (
	// TODO: this MUST NOT be public
	// TODO: find a way to store this key secretly
	secretKey = []byte("&kv@l3go-f=@^*@ush0(o5*5utxe6932j9di+ume=$mkj%d&&9*%k53(bmpksf&!c2&zpw$z=8ndi6ib)&nxms0ia7rf*sj9g8r4")

	anonymousAllowed = []string{
		"/account.Account/SignUp",
		"/account.Account/Verify",
		"/account.Account/Login",
		"/account.Account/PasswordReset",
		"/account.Account/PasswordSet",
		"/account.Account/ForgotLogin",

		"/version.Version/List",
	}
)

func isAnonymous(elem string) bool {
	for _, e := range anonymousAllowed {
		if e == elem {
			return true
		}
	}
	return false
}

// StreamInterceptor is an interceptor checking for authentication tokens
func StreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if anonymous := isAnonymous(info.FullMethod); !anonymous {
		if _, err := authorize(stream.Context()); err != nil {
			return err
		}
	}
	return handler(srv, stream)
}

// Interceptor is an interceptor checking for authentication tokens
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (i interface{}, err error) {
	if anonymous := isAnonymous(info.FullMethod); !anonymous {
		if ctx, err = authorize(ctx); err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}

func authorize(ctx context.Context) (context.Context, error) {
	if md, ok := metadata.FromContext(ctx); ok {
		tokens := md[TokenKey]
		if len(tokens) == 0 {
			return nil, grpc.Errorf(codes.Unauthenticated, "credentials required")
		}
		token := tokens[0]
		if token == "" {
			return nil, grpc.Errorf(codes.Unauthenticated, "credentials required")
		}
		claims, err := ValidateUserToken(token)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		// Enrich the context with the requester
		md := metadata.Pairs(RequesterKey, claims.AccountName)
		ctx = metadata.NewContext(ctx, md)
		return ctx, nil
	}
	return nil, grpc.Errorf(codes.Unauthenticated, "credentials required")
}

// CreateUserToken creates a token for a given user name
func CreateUserToken(name string, validFor time.Duration) (string, error) {
	// Forge the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		name, // The token contains the user name to verify
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(validFor).Unix(),
			Issuer:    os.Args[0],
		},
	})
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("unable to issue verification token")
	}
	return ss, nil
}

// ValidateUserToken validates a user token and return its claims
func ValidateUserToken(signedString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(signedString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
