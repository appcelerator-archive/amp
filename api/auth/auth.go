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
		"/account.Account/GetUser",
		"/account.Account/ListUsers",
		"/account.Account/GetOrganization",
		"/account.Account/ListOrganizations",
		"/account.Account/GetTeam",
		"/account.Account/ListTeams",
		"/cluster.Cluster/GetRegistrationStatus",
		"/version.Version/Get",
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
	if !isAnonymous(info.FullMethod) {
		if _, err := authorize(stream.Context()); err != nil {
			return err
		}
	}
	return handler(srv, stream)
}

// Interceptor is an interceptor checking for authentication tokens
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (i interface{}, err error) {
	if !isAnonymous(info.FullMethod) {
		if ctx, err = authorize(ctx); err != nil {
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

func authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, CredentialsRequired)
	}
	authorizations := md[AuthorizationHeader]
	if len(authorizations) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, CredentialsRequired)
	}
	authorization := authorizations[0]
	scheme, token := parseAuthorizationHeader(authorization)
	if scheme != AuthorizationScheme || token == "" {
		return nil, status.Errorf(codes.Unauthenticated, CredentialsRequired)
	}
	claims, err := ValidateToken(token, TokenTypeLogin)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials. Please log in again.")
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
