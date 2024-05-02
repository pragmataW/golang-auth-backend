package service

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateSecureRandomCode() (int, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	delta := new(big.Int).Sub(max, min)
	delta = delta.Add(delta, big.NewInt(1))

	n, err := rand.Int(rand.Reader, delta)
	if err != nil {
		return 0, err
	}

	n.Add(n, min)

	return int(n.Int64()), nil
}

func PostRequest(url string, body interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d - Body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func GenerateJWT(username string, issuer string, key string) (string, error) {
	claims := jwt.MapClaims{
		"Username": username,
		"iss":      issuer,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenizedStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenizedStr, nil
}
