package security

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	password := "Admin123!"

	encodedHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if encodedHash == password {
		t.Fatal("password tidak boleh disimpan sebagai teks biasa")
	}

	if !VerifyPassword(encodedHash, password) {
		t.Fatal("password yang benar seharusnya diterima")
	}

	if VerifyPassword(encodedHash, "PasswordSalah") {
		t.Fatal("password yang salah seharusnya ditolak")
	}
}

func TestHashPasswordRejectsShortPassword(t *testing.T) {
	if _, err := HashPassword("pendek"); err == nil {
		t.Fatal("password pendek seharusnya ditolak")
	}
}
