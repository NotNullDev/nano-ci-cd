package auth

import "testing"

func BenchmarkTokenCreation(b *testing.B) {

	for i := 0; i < b.N; i++ {
		CreateToken()
	}
}

func TestTokenCreation(t *testing.T) {
	t.Run("should not return error", func(t *testing.T) {
		_, err := CreateToken()
		if err != nil {
			t.Errorf("Error creating token: %v", err)
		}
		t.Log("Token created successfully")
	})
}

func TestShouldNotValidateRandomStringAsToken(t *testing.T) {
	randomString := "randomString"

	if ValidateToken(randomString) == nil {
		t.Errorf("Token %v is valid", randomString)
	}
}

func TestShouldValidateValidJWT(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.R89MVOJwuDNW_UnFL_cpAV04P004BETdSLhMPNxyNsI"
	secret := "hahaah"

	if validateTokenInternal(jwt, []byte(secret)) != nil {
		t.Errorf("Token %v is not valid", jwt)
	}
}

func TestShouldNotValidateValidJWTWithInvalidSecret(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.R89MVOJwuDNW_UnFL_cpAV04P004BETdSLhMPNxyNsI"
	secret := "hahaah1"

	if validateTokenInternal(jwt, []byte(secret)) == nil {
		t.Errorf("Token %v is not valid", jwt)
	}
}

func TestShouldValidateToken(t *testing.T) {
	token, err := CreateToken()

	if err != nil {
		t.Errorf("Error creating token: %v", err)
	}

	if ValidateToken(token) != nil {
		t.Errorf("Token %v is not valid", token)
	}
}
