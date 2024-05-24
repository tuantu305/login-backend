package main

type Sanitizer interface {
	// Sanitize input to prevent SQL Injection
	Sanitize(input string) string
	// Verify input to prevent SQL Injection or weak password
	Verify(input string) bool
}

type mockSanitizer struct {
}

func (ms *mockSanitizer) Verify(input string) bool {
	return true
}

func (ms *mockSanitizer) Sanitize(input string) string {
	return input
}

func NewMockSanitizer() Sanitizer {
	return &mockSanitizer{}
}
