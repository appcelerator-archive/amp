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

const (
	TokenKey = "amp.token"
)

var (
	// TODO: this MUST NOT be public
	// TODO: find a way to store this key secretly
	secretKey = []byte("&kv@l3go-f=@^*@ush0(o5*5utxe6932j9di+ume=$mkj%d&&9*%k53(bmpksf&!c2&zpw$z=8ndi6ib)&nxms0ia7rf*sj9g8r4")

	// TODO: there is probably a better way of achieving this
	anonymousAllowed = []string{
		// TODO: Temporarily allow access to everything
		"/account.Account/SignUp",
		"/account.Account/Verify",
		"/account.Account/Login",
		"/account.Account/PasswordChange",
		"/account.Account/PasswordReset",
		"/account.Account/PasswordSet",
		"/account.Account/ForgotLogin",

		"/function.Function/Create",
		"/function.Function/List",
		"/function.Function/Delete",

		"/logs.Logs/Get",
		"/logs.Logs/GetStream",

		"/oauth.Github/Create",

		"/service.Service/Create",
		"/service.Service/Remove",

		"/stack.StackService/Up",
		"/stack.StackService/Create",
		"/stack.StackService/Start",
		"/stack.StackService/Stop",
		"/stack.StackService/Remove",
		"/stack.StackService/Get",
		"/stack.StackService/List",
		"/stack.StackService/Tasks",

		"/stats.Stats/StatsQuery",

		"/storage.Storage/Put",
		"/storage.Storage/Get",
		"/storage.Storage/Delete",
		"/storage.Storage/List",

		"/storage.Storage/List",

		"/topic.Topic/Create",
		"/topic.Topic/List",
		"/topic.Topic/Delete",

		"/version.Version/List",
	}
)

// AuthStreamInterceptor is an interceptor checking for authentication tokens
func AuthStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	anonymous := false
	for _, method := range anonymousAllowed {
		if method == info.FullMethod {
			anonymous = true
			break
		}
	}
	if !anonymous {
		if err := authorize(stream.Context()); err != nil {
			return err
		}
	}
	return handler(srv, stream)
}

// AuthInterceptor is an interceptor checking for authentication tokens
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	anonymous := false
	for _, method := range anonymousAllowed {
		if method == info.FullMethod {
			anonymous = true
			break
		}
	}
	if !anonymous {
		if err := authorize(ctx); err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}

func authorize(ctx context.Context) error {
	if md, ok := metadata.FromContext(ctx); ok {
		tokens := md[TokenKey]
		if len(tokens) == 0 {
			return grpc.Errorf(codes.Unauthenticated, "credentials required")
		}
		token := tokens[0]
		if token == "" {
			return grpc.Errorf(codes.Unauthenticated, "credentials required")
		}
		fmt.Println("token", token)
		return nil
	}
	return grpc.Errorf(codes.Internal, "empty metadata")
}

// CreateUserToken creates a token for a given user name
func CreateUserToken(name string, validFor time.Duration) (string, error) {
	// Forge the verification token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		name, // The token contains the user name to verify
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(validFor).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
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

	// Get the claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
