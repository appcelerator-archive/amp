package account

import (
	"net/mail"

	"github.com/holys/safe"
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

func CheckUserName(name string) error {
	if name == "" {
		return grpc.Errorf(codes.InvalidArgument, "user name is mandatory")
	}
	return nil
}

func CheckPassword(password string) error {
	if password == "" {
		return grpc.Errorf(codes.InvalidArgument, "password is mandatory")
	}
	return nil
}

func checkOrganizationName(name string) error {
	if name == "" {
		return grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return nil
}

func CheckEmailAddress(email string) (processedEmail string, err error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	return address.Address, nil
}

func CheckPasswordStrength(password string) error {
	s := safe.New(8, 0, 0, safe.Simple)
	err := s.SetWords("./vendor/github.com/holys/safe/words.dat")
	if err != nil {
		return grpc.Errorf(codes.NotFound, "cannot find common password data")
	}
	str := s.Check(password)
	if str < 2 {
		return grpc.Errorf(codes.InvalidArgument, "password too weak")
	}
	return nil
}

func CheckVerificationCodeFormat(code string) error {
	if len(code) != codeLength {
		return grpc.Errorf(codes.InvalidArgument, "invalid verification code")
	}
	return nil
}

// Validate validates SignUpRequest
func (r *SignUpRequest) Validate() (err error) {
	err = CheckUserName(r.Name)
	if err != nil {
		return
	}
	email, err := CheckEmailAddress(r.Email)
	if err != nil {
		return
	}
	r.Email = email
	err = CheckPasswordStrength(r.Password)
	if err != nil {
		return
	}
	return
}

// Validate validates VerificationRequest
func (r *VerificationRequest) Validate() error {
	return CheckVerificationCodeFormat(r.Code)
}

// Validate validates OrganizationRequest
func (r *OrganizationRequest) Validate() (err error) {
	err = checkOrganizationName(r.Name)
	if err != nil {
		return
	}
	email, err := CheckEmailAddress(r.Email)
	if err != nil {
		return
	}
	r.Email = email
	return
}

// Validate validates LogInRequest
func (r *LogInRequest) Validate() error {
	return CheckUserName(r.Name)
}

// Validate validates EditRequest
func (r *EditRequest) Validate() (err error) {
	err = checkName(r.Name)
	if err != nil {
		return
	}
	if r.Email != "" {
		var email string
		email, err = CheckEmailAddress(r.Email)
		if err != nil {
			return
		}
		r.Email = email
	}
	if r.NewPassword != "" {
		err = CheckPasswordStrength(r.NewPassword)
		if err != nil {
			return
		}
	}
	return
}

// Validate validates PasswordResetRequest
func (r *PasswordResetRequest) Validate() (err error) {
	err = CheckUserName(r.Username)
	if err != nil {
		return
	}
	email, err := CheckEmailAddress(r.Email)
	if err != nil {
		return
	}
	r.Email = email
	return
}
