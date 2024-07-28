package validate

import (
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, helper.FieldViolation("username", err))
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, helper.FieldViolation("password", err))
	}

	return violations
}
