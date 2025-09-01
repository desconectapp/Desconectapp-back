package service

import (
	"context"
	repository "gin/db/generated"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type EmailValidationService struct {
	smtpHost  string
	smtpPort  string
	smtpEmail string
	smtpPassword  string
	queries   *repository.Queries
	ctx       context.Context
}

func NewEmailValidationService(conn *pgx.Conn) *EmailValidationService {
	queries := repository.New(conn)
	ctx := context.Background()

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if len(smtpHost) == 0 || len(smtpPort) == 0 || len(smtpEmail) == 0 || len(smtpPassword) == 0 {
		log.Println("SMTP environment variables are not set.")
	}

	return &EmailValidationService{
		queries: queries,
		ctx:     ctx,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpEmail: smtpEmail,
		smtpPassword: smtpPassword,
	}
}

func (s *EmailValidationService) StartEmailVerification(userId int32) error {
	var args repository.CreateEmailVerificationTokenParams

	code := "ABC123"
	args.Code = &code
	args.UserID = userId

	err := s.queries.CreateEmailVerificationToken(s.ctx, args)
	if err != nil {
		log.Println("Error creating email verification token:", err)
		return err
	}

	return nil
}
