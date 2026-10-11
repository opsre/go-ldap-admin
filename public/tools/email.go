package tools

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/patrickmn/go-cache"

	"strconv"

	"gopkg.in/gomail.v2"
)

// Store the verification code in the cache
var VerificationCodeCache = cache.New(24*time.Hour, 48*time.Hour)

func email(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": config.Conf.Email.User,
		"pass": config.Conf.Email.Pass,
		"host": config.Conf.Email.Host,
		"port": config.Conf.Email.Port,
	}
	port, _ := strconv.Atoi(mailConn["port"]) //Convert the port to an int

	newmail := gomail.NewMessage()

	newmail.SetHeader("From", newmail.FormatAddress(mailConn["user"], config.Conf.Email.From))
	newmail.SetHeader("To", mailTo...)    //Send to multiple users
	newmail.SetHeader("Subject", subject) //Set the email subject
	newmail.SetBody("text/html", body)    //Set the email body

	do := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	return do.DialAndSend(newmail)
}

func SendMail(sendto []string, pass string) error {
	return SendMailI18n(sendto, pass, "")
}

func SendMailI18n(sendto []string, pass string, locale string) error {
	subject := i18n.T(locale, "email.password_reset_subject", nil)
	body := fmt.Sprintf(i18n.T(locale, "email.password_reset_body", nil), pass)
	return email(sendto, subject, body)
}

// SendCode sends a verification code
func SendCode(sendto []string) error {
	return SendCodeI18n(sendto, "")
}

func SendCodeI18n(sendto []string, locale string) error {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	// Store the verification code in the cache so it can be retrieved for validation
	VerificationCodeCache.Set(sendto[0], vcode, time.Minute*5)
	subject := i18n.T(locale, "email.verification_code_subject", nil)
	body := fmt.Sprintf(i18n.T(locale, "email.verification_code_body", nil), vcode)
	return email(sendto, subject, body)
}

// SendUserCreationNotification sends the user creation success notification email
func SendUserCreationNotification(username, nickname, mail, password string) error {
	return SendUserCreationNotificationI18n(username, nickname, mail, password, "")
}

func SendUserCreationNotificationI18n(username, nickname, mail, password string, locale string) error {
	subject := i18n.T(locale, "email.user_creation_subject", nil)
	body := fmt.Sprintf(i18n.T(locale, "email.user_creation_body", nil), nickname, username, nickname, password)
	return email([]string{mail}, subject, body)
}

// SendPasswordResetNotification sends the password reset success notification email
func SendPasswordResetNotification(username, nickname, mail, newPassword string) error {
	return SendPasswordResetNotificationI18n(username, nickname, mail, newPassword, "")
}

func SendPasswordResetNotificationI18n(username, nickname, mail, newPassword string, locale string) error {
	subject := i18n.T(locale, "email.admin_password_reset_subject", nil)
	body := fmt.Sprintf(i18n.T(locale, "email.admin_password_reset_body", nil), nickname, username, newPassword)
	return email([]string{mail}, subject, body)
}
