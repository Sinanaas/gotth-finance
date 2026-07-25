package initializers

import "testing"

func fullConfig() Config {
	return Config{
		DBHost:                 "localhost",
		DBUserName:             "postgres",
		DBUserPassword:         "secret",
		DBName:                 "gotth",
		DBPort:                 "5432",
		ServerPort:             "8081",
		SessionSecretKey:       "s3cr3t",
		AccessTokenPrivateKey:  "apriv",
		AccessTokenPublicKey:   "apub",
		RefreshTokenPrivateKey: "rpriv",
		RefreshTokenPublicKey:  "rpub",
	}
}

func TestValidate_AllPresent(t *testing.T) {
	if err := fullConfig().Validate(); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestValidate_MissingSessionSecret(t *testing.T) {
	c := fullConfig()
	c.SessionSecretKey = ""
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for missing session secret")
	}
}

func TestValidate_MissingSigningKey(t *testing.T) {
	c := fullConfig()
	c.AccessTokenPrivateKey = "   "
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for blank signing key")
	}
}

func TestValidate_OptionalFlagsNotRequired(t *testing.T) {
	c := fullConfig() // AllowRegistration / CookieSecure left false
	if err := c.Validate(); err != nil {
		t.Fatalf("optional flags should not be required: %v", err)
	}
}
