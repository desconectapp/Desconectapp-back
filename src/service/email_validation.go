package service

import (
	"context"
	"crypto/rand"
	"fmt"
	repository "gin/db/generated"
	"log"
	"math/big"
	"net/smtp"
	"os"

	"github.com/jackc/pgx/v5"
)

type EmailValidationService struct {
	smtpHost     string
	smtpPort     string
	smtpEmail    string
	smtpPassword string
	queries      *repository.Queries
	ctx          context.Context
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
		queries:      queries,
		ctx:          ctx,
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpEmail:    smtpEmail,
		smtpPassword: smtpPassword,
	}
}

func (s *EmailValidationService) sendEmail(to string, code string) error {
	if len(s.smtpEmail) == 0 || len(s.smtpPassword) == 0 || len(s.smtpHost) == 0 || len(s.smtpPort) == 0 {
		log.Println("SMTP configuration is incomplete. Cannot send email.")
		return nil
	}

	subject := "Verify your email"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h1>Email Verification For Desconectapp</h1>
				<p>Your verification code is: <strong>%s</strong></p>
				<p>Please enter this code in the app to verify your email address.</p>
				<p>Do not share this code with anyone.</p>
				<p>If you did not request this, please ignore this email.</p>
			</body>
		</html>
	`, code)

	msg := []byte("From: " + s.smtpEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", s.smtpEmail, s.smtpPassword, s.smtpHost)

	host := s.smtpHost + ":" + s.smtpPort
	return smtp.SendMail(host, auth, s.smtpEmail, []string{to}, msg)
}

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomCode(n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

func (s *EmailValidationService) StartEmailVerification(userId int32, email string) error {
	code, err := RandomCode(6)
	if err != nil {
		return err
	}

	var args repository.CreateEmailVerificationTokenParams
	args.Code = &code
	args.UserID = userId

	err = s.queries.CreateEmailVerificationToken(s.ctx, args)
	if err != nil {
		log.Println("Error creating email verification token:", err)
		return err
	}

	err = s.sendEmail(email, code)
	if err != nil {
		log.Println("Error sending verification email:", err)
		return err
	}

	log.Println("Verification email sent to", email)

	return nil
}

func (s *EmailValidationService) ValidateEmailCode(userId int32, codeToValidate string) error {
	code, err := s.queries.GetVerificationCode(s.ctx, userId)
	if err != nil {
		log.Println("Error fetching verification code:", err)
		return err
	}

	if *code.Code == "" {
		return fmt.Errorf("no verification code found for user")
	}

	if *code.Code == codeToValidate {
		err = s.queries.VerifyEmail(s.ctx, userId)
		if err != nil {
			log.Println("Error verifying email:", err)
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid verification code")
}
