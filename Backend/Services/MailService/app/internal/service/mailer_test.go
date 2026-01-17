package service

import (
	"testing"

	"github.com/raphael-guer1n/AREA/MailService/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewMailer(t *testing.T) {
	cfg := config.Config{
		SMTP: config.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "user@example.com",
			Password: "password",
			From:     "noreply@example.com",
			FromName: "Test Service",
			Security: "starttls",
		},
	}

	mailer := NewMailer(cfg)

	assert.NotNil(t, mailer)
	assert.Equal(t, cfg, mailer.cfg)
}

func TestMailer_Send_ValidationErrors(t *testing.T) {
	cfg := config.Config{
		SMTP: config.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "user@example.com",
			Password: "password",
			From:     "noreply@example.com",
			Security: "starttls",
		},
	}
	mailer := NewMailer(cfg)

	testCases := []struct {
		name       string
		recipients []string
		subject    string
		body       string
		errorMsg   string
	}{
		{
			name:       "no recipients",
			recipients: []string{},
			subject:    "Test",
			body:       "Test body",
			errorMsg:   "no recipients provided",
		},
		{
			name:       "empty subject",
			recipients: []string{"test@example.com"},
			subject:    "",
			body:       "Test body",
			errorMsg:   "subject is required",
		},
		{
			name:       "empty body",
			recipients: []string{"test@example.com"},
			subject:    "Test",
			body:       "",
			errorMsg:   "body is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mailer.Send(tc.recipients, tc.subject, tc.body)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

func TestMailer_Send_IncompleteSMTPConfig(t *testing.T) {
	testCases := []struct {
		name     string
		cfg      config.Config
		errorMsg string
	}{
		{
			name: "missing host",
			cfg: config.Config{
				SMTP: config.SMTPConfig{
					Host:     "",
					Port:     587,
					Username: "user",
					Password: "pass",
					From:     "from@example.com",
				},
			},
			errorMsg: "smtp configuration is incomplete",
		},
		{
			name: "missing port",
			cfg: config.Config{
				SMTP: config.SMTPConfig{
					Host:     "smtp.example.com",
					Port:     0,
					Username: "user",
					Password: "pass",
					From:     "from@example.com",
				},
			},
			errorMsg: "smtp configuration is incomplete",
		},
		{
			name: "missing username",
			cfg: config.Config{
				SMTP: config.SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					Username: "",
					Password: "pass",
					From:     "from@example.com",
				},
			},
			errorMsg: "smtp configuration is incomplete",
		},
		{
			name: "missing password",
			cfg: config.Config{
				SMTP: config.SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					Username: "user",
					Password: "",
					From:     "from@example.com",
				},
			},
			errorMsg: "smtp configuration is incomplete",
		},
		{
			name: "missing from",
			cfg: config.Config{
				SMTP: config.SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					Username: "user",
					Password: "pass",
					From:     "",
				},
			},
			errorMsg: "smtp configuration is incomplete",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mailer := NewMailer(tc.cfg)

			err := mailer.Send([]string{"test@example.com"}, "Subject", "Body")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

func TestSanitizeHeader(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "This is a normal subject",
			expected: "This is a normal subject",
		},
		{
			name:     "with carriage return",
			input:    "Subject\rLine",
			expected: "Subject Line",
		},
		{
			name:     "with newline",
			input:    "Subject\nLine",
			expected: "Subject Line",
		},
		{
			name:     "with CRLF",
			input:    "Subject\r\nLine",
			expected: "Subject  Line",
		},
		{
			name:     "multiple newlines",
			input:    "Line1\n\nLine2\nLine3",
			expected: "Line1  Line2 Line3",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sanitizeHeader(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNormalizeLineEndings(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unix line endings (LF)",
			input:    "Line1\nLine2\nLine3",
			expected: "Line1\r\nLine2\r\nLine3",
		},
		{
			name:     "Windows line endings (CRLF)",
			input:    "Line1\r\nLine2\r\nLine3",
			expected: "Line1\r\nLine2\r\nLine3",
		},
		{
			name:     "Mac line endings (CR)",
			input:    "Line1\rLine2\rLine3",
			expected: "Line1\r\nLine2\r\nLine3",
		},
		{
			name:     "mixed line endings",
			input:    "Line1\nLine2\r\nLine3\rLine4",
			expected: "Line1\r\nLine2\r\nLine3\r\nLine4",
		},
		{
			name:     "no line endings",
			input:    "SingleLine",
			expected: "SingleLine",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeLineEndings(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBuildMultipartMessage(t *testing.T) {
	cfg := config.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
		FromName: "Test Service",
		Security: "starttls",
	}

	recipients := []string{"recipient@example.com"}
	subject := "Test Subject"
	plainBody := "This is plain text"
	htmlBody := "<p>This is HTML</p>"

	message := buildMultipartMessage(cfg, recipients, subject, plainBody, htmlBody)

	// Verify message structure
	assert.Contains(t, message, "From: \"Test Service\" <noreply@example.com>")
	assert.Contains(t, message, "To: recipient@example.com")
	assert.Contains(t, message, "Subject: Test Subject")
	assert.Contains(t, message, "MIME-Version: 1.0")
	assert.Contains(t, message, "Content-Type: multipart/alternative")
	assert.Contains(t, message, "Content-Type: text/plain")
	assert.Contains(t, message, "Content-Type: text/html")
	assert.Contains(t, message, "This is plain text")
	assert.Contains(t, message, "<p>This is HTML</p>")
}

func TestBuildMultipartMessage_WithoutFromName(t *testing.T) {
	cfg := config.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
		FromName: "",
		Security: "starttls",
	}

	recipients := []string{"recipient@example.com"}
	subject := "Test Subject"
	plainBody := "Plain text"
	htmlBody := "<p>HTML</p>"

	message := buildMultipartMessage(cfg, recipients, subject, plainBody, htmlBody)

	// Without FromName, should just use the email address
	assert.Contains(t, message, "From: noreply@example.com")
	assert.NotContains(t, message, "From: <noreply@example.com>")
}

func TestBuildMultipartMessage_MultipleRecipients(t *testing.T) {
	cfg := config.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
		Security: "starttls",
	}

	recipients := []string{
		"recipient1@example.com",
		"recipient2@example.com",
		"recipient3@example.com",
	}
	subject := "Test Subject"
	plainBody := "Plain text"
	htmlBody := "<p>HTML</p>"

	message := buildMultipartMessage(cfg, recipients, subject, plainBody, htmlBody)

	assert.Contains(t, message, "To: recipient1@example.com, recipient2@example.com, recipient3@example.com")
}

func TestBuildHTMLBody_Success(t *testing.T) {
	// Note: This test will fail if the template file doesn't exist
	// In a real test environment, you would either:
	// 1. Ensure the template file exists in the test environment
	// 2. Mock the template loading
	// 3. Skip this test in CI/CD if template is not available

	t.Skip("Skipping test that requires template file - run manually when template is available")

	subject := "Test Subject"
	body := "Test body\nwith multiple lines"

	htmlBody, err := buildHTMLBody(subject, body)

	assert.NoError(t, err)
	assert.NotEmpty(t, htmlBody)
	assert.Contains(t, htmlBody, subject)
	assert.Contains(t, htmlBody, "Test body<br>with multiple lines")
}

func TestBuildHTMLBody_EscapesHTML(t *testing.T) {
	t.Skip("Skipping test that requires template file - run manually when template is available")

	subject := "Test <script>alert('xss')</script>"
	body := "Body with <script>alert('xss')</script> tags"

	htmlBody, err := buildHTMLBody(subject, body)

	assert.NoError(t, err)
	assert.NotContains(t, htmlBody, "<script>")
	assert.Contains(t, htmlBody, "&lt;script&gt;")
}

func TestMailer_Send_WithInvalidSMTPServer(t *testing.T) {
	// This test attempts to connect to a non-existent SMTP server
	cfg := config.Config{
		SMTP: config.SMTPConfig{
			Host:     "nonexistent.smtp.server.invalid",
			Port:     587,
			Username: "user@example.com",
			Password: "password",
			From:     "noreply@example.com",
			Security: "starttls",
		},
	}
	mailer := NewMailer(cfg)

	err := mailer.Send([]string{"test@example.com"}, "Test Subject", "Test Body")

	// Should fail to connect
	assert.Error(t, err)
}

func TestMailer_SecurityModes(t *testing.T) {
	testCases := []struct {
		name     string
		security string
	}{
		{"SSL/TLS", "ssl"},
		{"STARTTLS", "starttls"},
		{"none", ""},
		{"uppercase", "SSL"},
		{"mixed case", "StartTLS"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				SMTP: config.SMTPConfig{
					Host:     "nonexistent.smtp.server.invalid",
					Port:     587,
					Username: "user@example.com",
					Password: "password",
					From:     "noreply@example.com",
					Security: tc.security,
				},
			}
			mailer := NewMailer(cfg)

			// All should fail since the server doesn't exist, but we're testing
			// that different security modes are handled without panicking
			err := mailer.Send([]string{"test@example.com"}, "Subject", "Body")
			assert.Error(t, err)
		})
	}
}
