package auth

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Keys used in context metadata
const (
	AuthorizationHeader   = "authorization" // HTTP2 authorization header name, cf. https://http2.github.io/http2-spec/compression.html#static.table.definition
	AuthorizationScheme   = "amp"
	TokenKey              = "amp.token"
	UserKey               = "amp.user"
	ActiveOrganizationKey = "amp.organization"
	CredentialsRequired   = "credentials required"
)

var (
	anonymousAllowed = []string{
		"/account.Account/SignUp",
		"/account.Account/Verify",
		"/account.Account/Login",
		"/account.Account/PasswordReset",
		"/account.Account/PasswordSet",
		"/account.Account/ForgotLogin",
		"/account.Account/ResendVerify",
		"/account.Account/GetUser",
		"/account.Account/ListUsers",
		"/account.Account/GetOrganization",
		"/account.Account/ListOrganizations",
		"/account.Account/GetTeam",
		"/account.Account/ListTeams",
		"/version.Version/VersionGet",
	}
)

type ValidateUser func(string) bool

type Interceptors struct {
	Tokens      *Tokens
	IsUserValid ValidateUser
}

func isAnonymous(elem string) bool {
	for _, e := range anonymousAllowed {
		if e == elem {
			return true
		}
	}
	return false
}

// StreamInterceptor is an interceptor checking for authentication tokens
func (i *Interceptors) StreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if _, err := i.authorize(stream.Context()); err != nil {
		if !isAnonymous(info.FullMethod) {
			return err
		}
	}
	return handler(srv, stream)
}

// Interceptor is an interceptor checking for authentication tokens
func (i *Interceptors) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (h interface{}, err error) {
	if ctx, err = i.authorize(ctx); err != nil {
		if !isAnonymous(info.FullMethod) {
			return nil, err
		}
	}
	return handler(ctx, req)
}

// Authorization header is formatted like this: "Authorization: <scheme> <token>", cf. https://tools.ietf.org/html/rfc7235#section-4.2
func parseAuthorizationHeader(header string) (scheme string, token string) {
	fields := strings.Fields(header)
	if len(fields) != 2 {
		return "", ""
	}
	scheme = fields[0]
	token = fields[1]
	return scheme, token
}

// ForgeAuthorizationHeader forges an amp authorization header
func ForgeAuthorizationHeader(token string) string {
	return AuthorizationScheme + " " + token
}

func (i *Interceptors) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, CredentialsRequired)
	}
	authorizations := md[AuthorizationHeader]
	if len(authorizations) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, CredentialsRequired)
	}
	authorization := authorizations[0]
	scheme, token := parseAuthorizationHeader(authorization)
	if scheme != AuthorizationScheme || token == "" {
		return ctx, status.Errorf(codes.Unauthenticated, CredentialsRequired)
	}
	claims, err := i.Tokens.ValidateToken(token, TokenTypeLogin)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid credentials. Please log in again.")
	}
	if !i.IsUserValid(claims.AccountName) {
		return ctx, status.Errorf(codes.Unauthenticated, "user not found. Please sign up.")
	}
	// Enrich the context
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(UserKey, claims.AccountName, ActiveOrganizationKey, claims.ActiveOrganization))
	return ctx, nil
}

// GetUser gets the user from context metadata
func GetUser(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	users := md[UserKey]
	if len(users) == 0 {
		return ""
	}
	user := users[0]
	return user
}

// GetActiveOrganization gets the active organization from context metadata
func GetActiveOrganization(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	activeOrganizations := md[ActiveOrganizationKey]
	if len(activeOrganizations) == 0 {
		return ""
	}
	activeOrganization := activeOrganizations[0]
	return activeOrganization
}
