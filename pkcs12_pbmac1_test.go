package pkcs12

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func loadTestData(t *testing.T, filename string) []byte {
	base64data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("failed to load test data: %v", err)
	}
	rawData, err := base64.StdEncoding.DecodeString(string(base64data))
	if err != nil {
		t.Fatalf("failed to decode test data %q: %v", filename, err)
	}
	return rawData
}

// RFC 9579 Appendix A.1
func TestDecodePKCS12DataOk(t *testing.T) {
	pfxData := loadTestData(t, "rfc9579-a1.txt")
	password := "1234"

	encodedPassword, err := bmpStringZeroTerminated(password)
	if err != nil {
		t.Fatalf("Failed to encode password to BMP string: %v", err)
	}

	safeContents, updatedPassword, err := GetSafeContents(pfxData, encodedPassword, 1, 10)
	if err != nil {
		t.Fatalf("Failed to load PKCS#12 data with GetSafeContents: %v", err)
	}

	if len(safeContents) == 0 {
		t.Error("Expected non-empty safe contents")
	}

	t.Logf("Successfully loaded PKCS#12 data with %d safe bag(s)", len(safeContents))
	t.Logf("Updated password length: %d", len(updatedPassword))

	for i, bag := range safeContents {
		t.Logf("  Bag %d: ID=%s", i, bag.Id.String())
	}
}

// RFC 9579 Appendix A.2
func TestDecodePKCS12DataSha256Sha512(t *testing.T) {
	pfxData := loadTestData(t, "rfc9579-a2.txt")
	password := "1234"

	encodedPassword, err := bmpStringZeroTerminated(password)
	if err != nil {
		t.Fatalf("Failed to encode password to BMP string: %v", err)
	}

	safeContents, updatedPassword, err := GetSafeContents(pfxData, encodedPassword, 1, 10)
	if err != nil {
		t.Fatalf("Failed to load PKCS#12 data with GetSafeContents: %v", err)
	}

	if len(safeContents) == 0 {
		t.Error("Expected non-empty safe contents")
	}

	t.Logf("Successfully loaded PKCS#12 SHA256/SHA512 data with %d safe bag(s)", len(safeContents))
	t.Logf("Updated password length: %d", len(updatedPassword))

	for i, bag := range safeContents {
		t.Logf("  Bag %d: ID=%s", i, bag.Id.String())
	}
}

// RFC 9579 Appendix A.3
func TestDecodePKCS12DataSha512Sha512(t *testing.T) {
	pfxData := loadTestData(t, "rfc9579-a3.txt")
	password := "1234"

	encodedPassword, err := bmpStringZeroTerminated(password)
	if err != nil {
		t.Fatalf("Failed to encode password to BMP string: %v", err)
	}

	safeContents, updatedPassword, err := GetSafeContents(pfxData, encodedPassword, 1, 10)
	if err != nil {
		t.Fatalf("Failed to load PKCS#12 data with GetSafeContents: %v", err)
	}

	if len(safeContents) == 0 {
		t.Error("Expected non-empty safe contents")
	}

	t.Logf("Successfully loaded PKCS#12 SHA512/SHA512 data with %d safe bag(s)", len(safeContents))
	t.Logf("Updated password length: %d", len(updatedPassword))

	for i, bag := range safeContents {
		t.Logf("  Bag %d: ID=%s", i, bag.Id.String())
	}
}

// RFC 9579 Appendix A.4
// Test with bad iteration count
func TestDecodePKCS12DataBadIterationCount(t *testing.T) {
	pfxData := loadTestData(t, "rfc9579-a4.txt")
	password := "1234"

	encodedPassword, err := bmpStringZeroTerminated(password)
	if err != nil {
		t.Fatalf("Failed to encode password to BMP string: %v", err)
	}

	_, _, err = GetSafeContents(pfxData, encodedPassword, 1, 10)
	if err == nil {
		t.Fatal("Expected GetSafeContents to fail with bad iteration count, but it succeeded")
	}

	if err != ErrIncorrectPassword {
		t.Fatalf("Got error %v but expected %v", err, ErrIncorrectPassword)
	}

	t.Logf("Successfully detected bad iteration count: %v", err)
}

// RFC 9579 Appendix A.5
// Test with incorrect salt
func TestDecodePKCS12DataIncorrectSalt(t *testing.T) {
	pfxData := loadTestData(t, "rfc9579-a5.txt")
	password := "1234"

	encodedPassword, err := bmpStringZeroTerminated(password)
	if err != nil {
		t.Fatalf("Failed to encode password to BMP string: %v", err)
	}

	_, _, err = GetSafeContents(pfxData, encodedPassword, 1, 10)
	if err == nil {
		t.Fatal("Expected GetSafeContents to fail with incorrect salt, but it succeeded")
	}

	if err != ErrIncorrectPassword {
		t.Fatalf("Got error %v but expected %v", err, ErrIncorrectPassword)
	}

	t.Logf("Successfully detected incorrect salt: %v", err)
}

// RFC 9579 Appendix A.6
// Test with missing key length
func TestDecodePKCS12DataMissingKeyLength(t *testing.T) {
	pfxData := loadTestData(t, "rfc9579-a6.txt")
	password := "1234"

	encodedPassword, err := bmpStringZeroTerminated(password)
	if err != nil {
		t.Fatalf("Failed to encode password to BMP string: %v", err)
	}

	_, _, err = GetSafeContents(pfxData, encodedPassword, 1, 10)
	if err == nil {
		t.Fatal("Expected GetSafeContents to fail with missing key length, but it succeeded")
	}

	if expected := "pkcs12: PBMAC1 requires explicit KeyLength parameter in PBKDF2 parameters"; err.Error() != expected {
		t.Fatalf("Got error %v but expected %v", err.Error(), expected)
	}

	t.Logf("Successfully detected missing key length: %v", err)
}
