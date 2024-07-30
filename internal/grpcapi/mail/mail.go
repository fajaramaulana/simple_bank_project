package mail

func SendMail(email string, code string) error {
	// Simulate email sending logic
	// Return an error to test retry logic
	return nil // or return fmt.Errorf("simulated error")
}
