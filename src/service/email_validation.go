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
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
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

func (s *EmailValidationService) sendEmail(to string, code string, template string) error {
	if len(s.smtpEmail) == 0 || len(s.smtpPassword) == 0 || len(s.smtpHost) == 0 || len(s.smtpPort) == 0 {
		log.Println("SMTP configuration is incomplete. Cannot send email.")
		return nil
	}

	subject := "Verify your email for Desconectapp"
	body := fmt.Sprintf(template, code)

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

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

	err = s.sendEmail(email, code, EMAIL_VALIDATION_TEMPLATE)
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

	if *code.Code == strings.ToUpper(codeToValidate) {
		err = s.queries.VerifyEmail(s.ctx, userId)
		if err != nil {
			log.Println("Error verifying email:", err)
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid verification code")
}

func (s *EmailValidationService) ResendEmail(userId int32) error {
	newCode, err := RandomCode(6)
	if err != nil {
		return err
	}

	user, err := s.queries.GetUserById(s.ctx, userId)
	if err != nil {
		log.Println("Error fetching user:", err)
		return err
	}

	if user.EmailValidated {
		return fmt.Errorf("email is already verified")
	}

	var args repository.UpdateVerificationCodeParams
	args.Code = &newCode
	args.UserID = userId

	err = s.queries.UpdateVerificationCode(s.ctx, args)
	if err != nil {
		log.Println("Error creating email verification token:", err)
		return err
	}

	if user.Email == "" {
		return fmt.Errorf("user does not have a valid email")
	}

	return s.sendEmail(user.Email, newCode, EMAIL_VALIDATION_TEMPLATE)
}

func (s *EmailValidationService) StartForgotPasswordFlow(email string) (*int32, error) {
	newCode, err := RandomCode(6)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUserByEmail(s.ctx, email)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	if user.EmailValidated == false {
		return nil, fmt.Errorf("email is not verified")
	}

	err = s.queries.CreateEmailVerificationToken(s.ctx, repository.CreateEmailVerificationTokenParams{
		UserID: user.ID,
		Code:   &newCode,
	})
	if err != nil {
		log.Println("Error creating email verification token:", err)
		return nil, err
	}

	err = s.sendEmail(user.Email, newCode, FORGOT_PASSWORD_TEMPLATE)
	if err != nil {
		log.Println("Error sending forgot password email:", err)
		return nil, err
	}

	return &user.ID, nil
}

func (s *EmailValidationService) UpdatePassword(userId int32, newPassword string, codeToValidate string) error {
	if newPassword == "" || len(newPassword) < 8 {
		return fmt.Errorf("Invalid new password")
	}

	code, err := s.queries.GetVerificationCode(s.ctx, userId)
	if err != nil {
		log.Println("Error fetching verification code:", err)
		return err
	}

	if *code.Code == strings.ToUpper(codeToValidate) {
		err = s.queries.VerifyEmail(s.ctx, userId)
		if err != nil {
			log.Println("Invalidating verification code:", err)
			return err
		}
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	err = s.queries.UpdateUserPassword(s.ctx, repository.UpdateUserPasswordParams{
		UserID:   userId,
		Password: string(hashedPassword),
	})
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}

const EMAIL_VALIDATION_TEMPLATE = `
	<html>
		<body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color:#f4f4f4;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4; padding:40px 0;">
				<tr>
					<td align="center">
						<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:12px; overflow:hidden; box-shadow:0 4px 12px rgba(0,0,0,0.1);">
							<tr>
								<td align="center" style="background:#4A90E2; padding:20px;">
									<img src="https://img.icons8.com/color/96/000000/verified-badge.png" alt="Verify Icon" width="80" style="display:block; margin:0 auto;">
									<h1 style="color:#ffffff; margin:10px 0 0; font-size:24px;">Email Verification</h1>
								</td>
							</tr>
							<tr>
								<td style="padding:30px; color:#333333;">
									<p style="font-size:16px; margin-bottom:20px;">
										Hello,
									</p>
									<p style="font-size:16px; margin-bottom:20px;">
										Thank you for signing up for <strong>Desconectapp</strong>! To complete your registration, please verify your email address.
									</p>
									<div style="text-align:center; margin:30px 0;">
										<p style="font-size:18px; margin-bottom:10px;">Your verification code:</p>
										<p style="font-size:28px; font-weight:bold; color:#4A90E2; letter-spacing:2px;">%s</p>
									</div>
									<p style="font-size:14px; color:#777777;">
										⚠️ Do not share this code with anyone. If you did not request this verification, you can safely ignore this email.
									</p>
								</td>
							</tr>
							<tr>
								<td align="center" style="background:#f4f4f4; padding:20px; font-size:12px; color:#999999;">
									<p>&copy; 2025 Desconectapp. All rights reserved.</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
	</html>
`

const FORGOT_PASSWORD_TEMPLATE = `
	<html>
		<body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color:#f4f4f4;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4; padding:40px 0;">
				<tr>
					<td align="center">
						<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:12px; overflow:hidden; box-shadow:0 4px 12px rgba(0,0,0,0.1);">
							<tr>
								<td align="center" style="background:#E94E77; padding:20px;">
									<img src="https://img.icons8.com/color/96/000000/forgot-password.png" alt="Forgot Password Icon" width="80" style="display:block; margin:0 auto;">
									<h1 style="color:#ffffff; margin:10px 0 0; font-size:24px;">Reset Your Password</h1>
								</td>
							</tr>
							<tr>
								<td style="padding:30px; color:#333333;">
									<p style="font-size:16px; margin-bottom:20px;">
										Hello,
									</p>
									<p style="font-size:16px; margin-bottom:20px;">
										We received a request to reset the password for your <strong>Desconectapp</strong> account. Use the code below to proceed with creating a new password.
									</p>
									<div style="text-align:center; margin:30px 0;">
										<p style="font-size:18px; margin-bottom:10px;">Your reset code:</p>
										<p style="font-size:28px; font-weight:bold; color:#E94E77; letter-spacing:2px;">%s</p>
									</div>
									<p style="font-size:14px; color:#777777; margin-top:20px;">
										⚠️ If you did not request a password reset, please ignore this email or contact support.
									</p>
								</td>
							</tr>
							<tr>
								<td align="center" style="background:#f4f4f4; padding:20px; font-size:12px; color:#999999;">
									<p>&copy; 2025 Desconectapp. All rights reserved.</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
	</html>
`
