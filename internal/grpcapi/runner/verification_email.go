package runner

import (
	"context"
	"errors"

	"time"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"gopkg.in/gomail.v2"
)

const (
	maxRetries        = 5
	baseBackoff       = time.Second
	maxBackoff        = 30 * time.Second
	emailsPerMinute   = 10 // Maximum number of emails to send per minute
	rateLimitInterval = time.Minute
)

func sendEmailWithRetry(m *gomail.Message, d *gomail.Dialer) error {
	var err error
	backoff := baseBackoff

	for i := 0; i < maxRetries; i++ {
		if err = d.DialAndSend(m); err == nil {
			return nil
		}

		log.Error().Err(err).Msg("Failed to send email")
		time.Sleep(backoff)

		// Exponential backoff with a maximum limit
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	return errors.New("failed to send email after maximum retries")
}

// SendVerificationEmails scans for verification email keys and sends emails at a controlled rate.
func SendVerificationEmails(ctx context.Context, rdb *redis.Client, config util.Config) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processEmails(ctx, rdb, config)
		case <-ctx.Done():
			log.Info().Msg("Context canceled, stopping email sender")
			return
		}
	}
}

func processEmails(ctx context.Context, rdb *redis.Client, config util.Config) {
	cursor := uint64(0)
	prefix := "verification_email:"
	limiter := rate.NewLimiter(rate.Every(rateLimitInterval/emailsPerMinute), emailsPerMinute) // Rate limiter
	for {
		// Scan for keys with the specified prefix
		keys, nextCursor, err := rdb.Scan(ctx, cursor, prefix+"*", 0).Result()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to scan keys")
		}

		// Process each key
		for _, key := range keys {
			if err := limiter.Wait(ctx); err != nil {
				log.Error().Err(err).Msg("Rate limiter error")
				continue
			}
			processKey(ctx, key, rdb, config)
		}

		// If nextCursor is 0, we have finished scanning
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
}

func processKey(ctx context.Context, key string, rdb *redis.Client, config util.Config) {
	// Retrieve the hash data
	fields, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get hash fields for key %s", key)
		return
	}

	// Log the retrieved fields

	// Extract specific fields
	email, emailExists := fields["email"]
	verificationCode, codeExists := fields["verification_code"]

	if !emailExists {
		log.Error().Msgf("Email for key %s not found", key)
	}

	if !codeExists {
		log.Error().Msgf("Verification code for key %s not found", key)
	}

	// todo sendemail
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "simplebank@fajaramaulanadev.com")
	mailer.SetHeader("To", email)
	// mailer.SetAddressHeader("Cc", "tralalala@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Verification Code Simplebank")
	mailer.SetBody("text/html", "Hello, this is your verification code: "+verificationCode)

	dialer := gomail.NewDialer(
		config.MailHost, config.MailPort, config.MailUser, config.MailPassword,
	)
	if err = sendEmailWithRetry(mailer, dialer); err != nil {
		log.Error().Err(err).Msg("Failed to send email")
	}

	log.Info().Msgf("Sending email to %s with verification code %s", email, verificationCode)

	// Delete the key after processing
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Error().Err(err).Msgf("Failed to delete key %s", key)
	}
}
