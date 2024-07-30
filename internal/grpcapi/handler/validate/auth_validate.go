package validate

import (
	"github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(req.GetUsername()); err != nil {
		log.Error().Err(err).Msg("Invalid username")
		violations = append(violations, helper.FieldViolation("username", err))
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		log.Error().Err(err).Msg("Invalid password")
		violations = append(violations, helper.FieldViolation("password", err))
	}

	return violations
}

func ValidateVerifyEmailUserRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := helper.ValidateRequired(req.GetVerificationCode()); err != nil {
		log.Error().Err(err).Msg("Invalid verification code")
		violations = append(violations, helper.FieldViolation("verification_code", err))
	}

	if err := helper.ValidateUUID(req.GetVerificationCode()); err != nil {
		log.Error().Err(err).Msg("Invalid user id")
		violations = append(violations, helper.FieldViolation("verficiation_code", err))
	}

	return violations
}
