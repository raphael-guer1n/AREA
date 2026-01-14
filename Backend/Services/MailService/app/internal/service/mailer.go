package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/MailService/internal/config"
)

type Mailer struct {
	cfg config.Config
}

func NewMailer(cfg config.Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) Send(recipients []string, subject, body string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided")
	}
	if subject == "" {
		return fmt.Errorf("subject is required")
	}
	if body == "" {
		return fmt.Errorf("body is required")
	}

	smtpCfg := m.cfg.SMTP
	if smtpCfg.Host == "" || smtpCfg.Port == 0 || smtpCfg.Username == "" || smtpCfg.Password == "" || smtpCfg.From == "" {
		return fmt.Errorf("smtp configuration is incomplete")
	}

	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)
	security := strings.ToLower(strings.TrimSpace(smtpCfg.Security))

	var client *smtp.Client
	var err error
	if security == "ssl" {
		tlsConfig := &tls.Config{ServerName: smtpCfg.Host}
		conn, dialErr := tls.Dial("tcp", addr, tlsConfig)
		if dialErr != nil {
			return fmt.Errorf("smtp tls dial failed: %w", dialErr)
		}
		client, err = smtp.NewClient(conn, smtpCfg.Host)
	} else {
		client, err = smtp.Dial(addr)
		if err == nil && security == "starttls" {
			tlsConfig := &tls.Config{ServerName: smtpCfg.Host}
			if startErr := client.StartTLS(tlsConfig); startErr != nil {
				_ = client.Close()
				return fmt.Errorf("smtp starttls failed: %w", startErr)
			}
		}
	}
	if err != nil {
		return fmt.Errorf("smtp connection failed: %w", err)
	}
	defer func() { _ = client.Quit() }()

	auth := smtp.PlainAuth("", smtpCfg.Username, smtpCfg.Password, smtpCfg.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}
	if err = client.Mail(smtpCfg.From); err != nil {
		return fmt.Errorf("smtp from failed: %w", err)
	}
	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt failed: %w", err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data failed: %w", err)
	}
	defer func() { _ = writer.Close() }()

	subject = sanitizeHeader(subject)
	htmlBody, err := buildHTMLBody(subject, body)
	if err != nil {
		return err
	}

	message := buildMultipartMessage(smtpCfg, recipients, subject, body, htmlBody)
	if _, err = writer.Write([]byte(message)); err != nil {
		return fmt.Errorf("smtp write failed: %w", err)
	}
	return nil
}

func buildMultipartMessage(cfg config.SMTPConfig, recipients []string, subject, plainBody, htmlBody string) string {
	boundary := fmt.Sprintf("area_mail_%d", time.Now().UnixNano())
	fromHeader := cfg.From
	if cfg.FromName != "" {
		fromHeader = (&mail.Address{Name: cfg.FromName, Address: cfg.From}).String()
	}

	headers := []string{
		"From: " + fromHeader,
		"To: " + strings.Join(recipients, ", "),
		"Subject: " + subject,
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		"Content-Type: multipart/alternative; boundary=\"" + boundary + "\"",
	}

	plain := normalizeLineEndings(plainBody)
	html := normalizeLineEndings(htmlBody)

	var buffer bytes.Buffer
	buffer.WriteString(strings.Join(headers, "\r\n"))
	buffer.WriteString("\r\n\r\n")
	buffer.WriteString("--" + boundary + "\r\n")
	buffer.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buffer.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	buffer.WriteString(plain)
	buffer.WriteString("\r\n--" + boundary + "\r\n")
	buffer.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buffer.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	buffer.WriteString(html)
	buffer.WriteString("\r\n--" + boundary + "--\r\n")
	return buffer.String()
}

func buildHTMLBody(subject, body string) (string, error) {
	escaped := template.HTMLEscapeString(body)
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")

	data := struct {
		Subject string
		Body    template.HTML
	}{
		Subject: subject,
		Body:    template.HTML(escaped),
	}

	tmpl, err := template.ParseFiles("templates/email.html")
	if err != nil {
		return "", fmt.Errorf("email template parse failed: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("email template execute failed: %w", err)
	}
	return buf.String(), nil
}

func sanitizeHeader(value string) string {
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}

func normalizeLineEndings(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	return strings.ReplaceAll(value, "\n", "\r\n")
}
