package account

import (
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/holys/safe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
)

func isEmpty(s string) bool {
	return s == "" || strings.TrimSpace(s) == ""
}

// CheckPassword checks password
func CheckPassword(password string) error {
	if isEmpty(password) {
		return grpc.Errorf(codes.InvalidArgument, "password is mandatory")
	}
	safety := safe.New(8, 0, 0, safe.Simple)
	if passwordStrength := safety.Check(password); passwordStrength <= safe.Simple {
		return grpc.Errorf(codes.InvalidArgument, "password too weak")
	}
	return nil
}

// CheckVerificationCode checks verification code
func CheckVerificationCode(code string) error {
	if isEmpty(code) {
		return grpc.Errorf(codes.InvalidArgument, "invalid verification code")
	}
	return nil
}

// Validate validates SignUpRequest
func (r *SignUpRequest) Validate() (err error) {
	if err = schema.CheckName(r.Name); err != nil {
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	if r.Email, err = schema.CheckEmailAddress(r.Email); err != nil {
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	if err = CheckPassword(r.Password); err != nil {
		return err
	}

	return nil
}

// Validate validates VerificationRequest
func (r *VerificationRequest) Validate() error {
	return nil
}

// Validate validates LogInRequest
func (r *LogInRequest) Validate() error {
	if err := schema.CheckName(r.Name); err != nil {
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := CheckPassword(r.Password); err != nil {
		return err
	}
	return nil
}

// Validate validates PasswordResetRequest
func (r *PasswordResetRequest) Validate() error {
	if err := schema.CheckName(r.Name); err != nil {
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

// Validate validates PasswordSetRequest
func (r *PasswordSetRequest) Validate() error {
	if err := CheckPassword(r.Password); err != nil {
		return err
	}
	return nil
}

// Validate validates PasswordChangeRequest
func (r *PasswordChangeRequest) Validate() error {
	if err := CheckPassword(r.NewPassword); err != nil {
		return err
	}
	return nil
}

// Validate validates ForgotLoginRequest
func (r *ForgotLoginRequest) Validate() (err error) {
	if r.Email, err = schema.CheckEmailAddress(r.Email); err != nil {
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

// Validate validates GetUserRequest
func (r *GetUserRequest) Validate() error {
	if err := schema.CheckName(r.Name); err != nil {
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

// Validate validates ListUsersRequest
func (r *ListUsersRequest) Validate() error {
	return nil
}
