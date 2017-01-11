package account

import (
	"net/mail"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const codeLength = 8
const minPasswordLength = 8

func checkName(name string) error {
	if name == "" {
		return grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	return nil
}

func checkUserName(name string) error {
	if name == "" {
		return grpc.Errorf(codes.InvalidArgument, "user name is mandatory")
	}
	return nil
}

func checkOrganizationName(name string) error {
	if name == "" {
		return grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil
}

func checkEmailAddress(email string) (processedEmail string, err error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	return address.Address, nil
}

func checkPasswordStrength(password string) error {
	if len(password) < minPasswordLength {
		return grpc.Errorf(codes.InvalidArgument, "password too weak")
	}
	return nil
}

func checkVerificationCodeFormat(code string) error {
	if len(code) != codeLength {
		return grpc.Errorf(codes.InvalidArgument, "invalid verification code")
	}
	return nil
}

// Validate validates SignUpRequest
func (r *SignUpRequest) Validate() (err error) {
	err = checkUserName(r.Name)
	if err != nil {
		return
	}
	email, err := checkEmailAddress(r.Email)
	if err != nil {
		return
	}
	r.Email = email
	err = checkPasswordStrength(r.Password)
	if err != nil {
		return
	}
	return
}

// Validate validates VerificationRequest
func (r *VerificationRequest) Validate() error {
	return checkVerificationCodeFormat(r.Code)
}

// Validate validates OrganizationRequest
func (r *OrganizationRequest) Validate() (err error) {
	err = checkOrganizationName(r.Name)
	if err != nil {
		return
	}
	email, err := checkEmailAddress(r.Email)
	if err != nil {
		return
	}
	r.Email = email
	return
}

// Validate validates LogInRequest
func (r *LogInRequest) Validate() error {
	return checkUserName(r.Name)
}

// Validate validates EditRequest
func (r *EditRequest) Validate() (err error) {
	err = checkName(r.Name)
	if err != nil {
		return
	}
	if r.Email != "" {
		var email string
		email, err = checkEmailAddress(r.Email)
		if err != nil {
			return
		}
		r.Email = email
	}
	if r.NewPassword != "" {
		err = checkPasswordStrength(r.NewPassword)
		if err != nil {
			return
		}
	}
	return
}
