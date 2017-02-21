package server

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
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

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
		tokens := md["amp.token"]
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
