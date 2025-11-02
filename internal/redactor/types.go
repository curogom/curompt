package redactor

// Redactor masks sensitive information in text
type Redactor interface {
	Redact(text string) string
}
