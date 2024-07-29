package validate

import (
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, helper.FieldViolation("username", err))
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, helper.FieldViolation("password", err))
	}

	if err := ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, helper.FieldViolation("currency", err))
	}

	if err := helper.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, helper.FieldViolation("email", err))
	}

	if err := ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, helper.FieldViolation("full_name", err))
	}

	return violations
}

func ValidateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUserUUID(req.GetUserUuid()); err != nil {
		violations = append(violations, helper.FieldViolation("user_uuid", err))
	}

	if err := ValidateUsernameNoRequired(req.GetUsername()); err != nil {
		violations = append(violations, helper.FieldViolation("username", err))
	}

	if err := ValidatePasswordNotRequired(req.GetPassword()); err != nil {
		violations = append(violations, helper.FieldViolation("password", err))
	}

	if req.GetEmail() != "" {
		if err := helper.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, helper.FieldViolation("email", err))
		}
	}

	if err := ValidateFullNameNotRequired(req.GetFullName()); err != nil {
		violations = append(violations, helper.FieldViolation("full_name", err))
	}

	return violations
}
