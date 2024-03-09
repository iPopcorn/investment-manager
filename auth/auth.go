package auth

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/iPopcorn/investment-manager/config"
)

var max = big.NewInt(math.MaxInt64)

type APIKeyClaims struct {
	*jwt.Claims
	URI string `json:"uri"`
}

type APIKey struct {
	Name          string   `json:"name"`
	Principal     string   `json:"principal"`
	PrincipalType string   `json:"principalType"`
	PublicKey     string   `json:"publicKey"`
	PrivateKey    string   `json:"privateKey"`
	CreateTime    string   `json:"createTime"`
	ProjectId     string   `json:"projectId"`
	Nickname      string   `json:"nickname"`
	Scopes        []string `json:"scopes"`
	AllowedIps    []string `json:"allowedIps"`
	KeyType       string   `json:"keyType"`
	Enabled       bool     `json:"enabled"`
}

type BuildJWTOptions struct {
	Service    string
	Uri        string
	PrivateKey string
	Name       string
}

type nonceSource struct{}

func GetApiKey() (*APIKey, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	apiKeyJsonFile, err := os.Open(config.ApiKeyPath)

	if err != nil {
		return nil, err
	}

	defer apiKeyJsonFile.Close()

	data, _ := ioutil.ReadAll(apiKeyJsonFile)
	var apiKey APIKey
	json.Unmarshal(data, &apiKey)

	return &apiKey, nil
}

func BuildJWT(options BuildJWTOptions) (string, error) {
	block, _ := pem.Decode([]byte(options.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("jwt: Could not decode private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: key},
		(&jose.SignerOptions{NonceSource: nonceSource{}}).WithType("JWT").WithHeader("kid", options.Name),
	)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	cl := &APIKeyClaims{
		Claims: &jwt.Claims{
			Subject:   options.Name,
			Issuer:    "coinbase-cloud",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Expiry:    jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			Audience:  jwt.Audience{options.Service},
		},
		URI: options.Uri,
	}
	jwtString, err := jwt.Signed(sig).Claims(cl).Serialize()
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil
}

func (n nonceSource) Nonce() (string, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return r.String(), nil
}
