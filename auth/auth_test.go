package auth_test

import (
	"fmt"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/iPopcorn/investment-manager/auth"
)

func TestAuth(t *testing.T) {
	t.Run("Generates JWT with correct claims", func(t *testing.T) {
		testOptions := getJWTOptions(t)
		signatureAlgorithms := []jose.SignatureAlgorithm{jose.ES256}
		jwtString, err := auth.BuildJWT(testOptions)

		if err != nil {
			t.Fatalf("Failed to build JWT\n%v", err)
		}

		jwtParsed, err := jwt.ParseSigned(jwtString, signatureAlgorithms)

		if err != nil {
			t.Fatalf("Failed to parse JWT\n%v", err)
		}

		expectedKid := testOptions.Name
		actualKid := jwtParsed.Headers[0].KeyID

		if expectedKid != actualKid {
			t.Fatalf("Expected %s Received %s", expectedKid, actualKid)
		}

		claims := &jwt.Claims{}
		fmt.Printf("jwtParsed: %v\n", jwtParsed)
		err = jwtParsed.UnsafeClaimsWithoutVerification(claims)

		if err != nil {
			t.Fatalf("Could not get claims\n%v", err)
		}

		expectedAudience := testOptions.Service
		expectedSubject := testOptions.Name
		actualAudience := claims.Audience[0]
		actualSubject := claims.Subject

		if expectedAudience != actualAudience {
			t.Fatalf("Expected %s Actual %s", expectedAudience, actualAudience)
		}

		if expectedSubject != actualSubject {
			t.Fatalf("Expected %s Actual %s", expectedSubject, actualSubject)
		}
	})
}

func getJWTOptions(t *testing.T) auth.BuildJWTOptions {
	t.Helper()
	buildOptions := auth.BuildJWTOptions{
		Service: "test-service",
		Uri:     "GET api.coinbase.com/api/v3/brokerage/accounts",
		Name:    "test-key",
		PrivateKey: `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIASRuLbKmWPx6wwIoT5IAiIG7KKGwWWBGNYwT22b2nxXoAoGCCqGSM49
AwEHoUQDQgAEclqeylf+zHBL8MkgtwBGjK73rbdn4+lrU4QUWibnYssM4puGDAQ6
nVWm7slZHLGEWSf84+ZLuofj3Eei4oCZrw==
-----END EC PRIVATE KEY-----`,
	}

	return buildOptions
}
