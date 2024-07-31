package seed

import (
	"context"
	"time"

	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Seeder struct {
	conn *pgxpool.Pool
}

func NewSeeder(conn *pgxpool.Pool) *Seeder {
	return &Seeder{conn: conn}
}

func (s *Seeder) Seed() {
	// check is table users empty
	var count int
	err := s.conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot count users")
	}

	if count == 0 {
		userUUID1 := uuid.New()
		userUUID2 := uuid.New()
		userUUID3 := uuid.New()
		hashedPassword, err := util.MakePasswordBcrypt("P4ssw0rd!")
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot hash password")
		}
		arg := []map[string]interface{}{
			{
				"username":                "superadmin",
				"user_uuid":               userUUID1,
				"hashed_password":         hashedPassword,
				"full_name":               "superadmin",
				"email":                   "superadmin@gmail.com",
				"role":                    "superadmin",
				"created_at":              time.Now(),
				"verified_email_at":       time.Now(),
				"verification_email_code": uuid.New(),
			},
			{
				"username":                "admin",
				"user_uuid":               userUUID2,
				"hashed_password":         hashedPassword,
				"full_name":               "admin",
				"email":                   "admin@gmail.com",
				"role":                    "admin",
				"created_at":              time.Now(),
				"verified_email_at":       time.Now(),
				"verification_email_code": uuid.New(),
			}, {
				"username":                "customer",
				"user_uuid":               userUUID3,
				"hashed_password":         hashedPassword,
				"full_name":               "customer",
				"email":                   "customer@gmail.com",
				"role":                    "customer",
				"created_at":              time.Now(),
				"verified_email_at":       time.Now(),
				"verification_email_code": uuid.New(),
			},
		}

		// insert to table users
		for _, v := range arg {
			_, err := s.conn.Exec(context.Background(), "INSERT INTO users (username, user_uuid, hashed_password, full_name, email, role, created_at, verified_email_at, verification_email_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", v["username"], v["user_uuid"], v["hashed_password"], v["full_name"], v["email"], v["role"], v["created_at"], v["verified_email_at"], v["verification_email_code"])
			if err != nil {
				log.Fatal().Err(err).Msg("Cannot insert users")
			}
		}

		argAccount := []map[string]interface{}{
			{
				"user_uuid":    userUUID1,
				"owner":        "superadmin",
				"currency":     "IDR",
				"balance":      1000000000,
				"status":       1,
				"created_at":   time.Now(),
				"account_uuid": uuid.New(),
			},
			{
				"user_uuid":    userUUID2,
				"owner":        "admin",
				"currency":     "IDR",
				"balance":      1000000000,
				"status":       1,
				"created_at":   time.Now(),
				"account_uuid": uuid.New(),
			},
			{
				"user_uuid":    userUUID3,
				"owner":        "customer",
				"currency":     "IDR",
				"balance":      1000000000,
				"status":       1,
				"created_at":   time.Now(),
				"account_uuid": uuid.New(),
			},
		}
		for _, v := range argAccount {
			_, err := s.conn.Exec(context.Background(), "INSERT INTO accounts (user_uuid, owner, currency, balance, status, created_at, account_uuid) VALUES ($1, $2, $3, $4, $5, $6, $7)", v["user_uuid"], v["owner"], v["currency"], v["balance"], v["status"], v["created_at"], v["account_uuid"])
			if err != nil {
				log.Fatal().Err(err).Msg("Cannot insert accounts")
			}
		}
	}
}
