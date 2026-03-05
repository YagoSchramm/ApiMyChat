package service

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	googleTokenURL = "https://oauth2.googleapis.com/token"
	fcmOAuthScope  = "https://www.googleapis.com/auth/firebase.messaging"
)

type FCMService struct {
	projectID   string
	clientEmail string
	privateKey  *rsa.PrivateKey
	tokenURL    string
	client      *http.Client

	tokenMu     sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

func NewFCMService() *FCMService {
	svc := &FCMService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenURL: googleTokenURL,
	}

	svc.loadV1CredentialsFromEnv()
	return svc
}

type serviceAccountCredentials struct {
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type fcmV1Request struct {
	Message fcmV1Message `json:"message"`
}

type fcmV1Message struct {
	Token        string            `json:"token"`
	Notification fcmNotification   `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *FCMService) SendToTokens(tokens []string, title, body string, data map[string]string) error {
	if len(tokens) == 0 {
		return nil
	}

	if !s.hasV1Config() {
		return fmt.Errorf("fcm v1 not configured: set FCM_PROJECT_ID, FCM_CLIENT_EMAIL and FCM_PRIVATE_KEY (or FCM_SERVICE_ACCOUNT_JSON)")
	}

	return s.sendV1(tokens, title, body, data)
}

func (s *FCMService) hasV1Config() bool {
	return s.projectID != "" && s.clientEmail != "" && s.privateKey != nil
}

func (s *FCMService) sendV1(tokens []string, title, body string, data map[string]string) error {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return err
	}

	var firstErr error
	for _, token := range tokens {
		if token == "" {
			continue
		}

		if err := s.sendToV1Token(accessToken, token, title, body, data); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func (s *FCMService) sendToV1Token(accessToken, token, title, body string, data map[string]string) error {
	endpoint := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", s.projectID)
	payload := fcmV1Request{
		Message: fcmV1Message{
			Token: token,
			Notification: fcmNotification{
				Title: title,
				Body:  body,
			},
			Data: data,
		},
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(rawPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("fcm v1 request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	return nil
}

func (s *FCMService) loadV1CredentialsFromEnv() {
	rawProjectID := strings.TrimSpace(os.Getenv("FCM_PROJECT_ID"))
	rawClientEmail := strings.TrimSpace(os.Getenv("FCM_CLIENT_EMAIL"))
	rawPrivateKey := strings.TrimSpace(os.Getenv("FCM_PRIVATE_KEY"))
	rawTokenURI := strings.TrimSpace(os.Getenv("FCM_TOKEN_URI"))
	rawJSON := strings.TrimSpace(os.Getenv("FCM_SERVICE_ACCOUNT_JSON"))

	s.projectID = rawProjectID
	if rawTokenURI != "" {
		s.tokenURL = rawTokenURI
	}

	if rawJSON != "" {
		var creds serviceAccountCredentials
		if err := json.Unmarshal([]byte(rawJSON), &creds); err == nil {
			if s.projectID == "" {
				s.projectID = strings.TrimSpace(creds.ProjectID)
			}
			if rawClientEmail == "" {
				rawClientEmail = strings.TrimSpace(creds.ClientEmail)
			}
			if rawPrivateKey == "" {
				rawPrivateKey = strings.TrimSpace(creds.PrivateKey)
			}
			if s.tokenURL == googleTokenURL && strings.TrimSpace(creds.TokenURI) != "" {
				s.tokenURL = strings.TrimSpace(creds.TokenURI)
			}
		}
	}

	s.clientEmail = rawClientEmail
	privateKey, err := parseRSAPrivateKey(rawPrivateKey)
	if err == nil {
		s.privateKey = privateKey
	}
}

func parseRSAPrivateKey(rawKey string) (*rsa.PrivateKey, error) {
	if strings.TrimSpace(rawKey) == "" {
		return nil, fmt.Errorf("empty private key")
	}

	decoded := strings.ReplaceAll(rawKey, "\\n", "\n")
	block, _ := pem.Decode([]byte(decoded))
	if block == nil {
		return nil, fmt.Errorf("invalid private key pem")
	}

	pkcs8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		rsaKey, ok := pkcs8.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not rsa")
		}
		return rsaKey, nil
	}

	pkcs1, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return pkcs1, nil
	}

	return nil, fmt.Errorf("failed parsing rsa private key")
}

func (s *FCMService) getAccessToken() (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	if s.accessToken != "" && time.Now().Before(s.tokenExpiry.Add(-1*time.Minute)) {
		return s.accessToken, nil
	}

	assertion, err := s.createSignedJWT()
	if err != nil {
		return "", err
	}

	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", assertion)

	req, err := http.NewRequest(http.MethodPost, s.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("oauth token request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	var tokenResp oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("oauth token response missing access_token")
	}

	s.accessToken = tokenResp.AccessToken
	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	s.tokenExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return s.accessToken, nil
}

func (s *FCMService) createSignedJWT() (string, error) {
	now := time.Now().Unix()
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	}
	claims := map[string]interface{}{
		"iss":   s.clientEmail,
		"scope": fcmOAuthScope,
		"aud":   s.tokenURL,
		"iat":   now,
		"exp":   now + 3600,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims

	hash := crypto.SHA256.New()
	_, _ = hash.Write([]byte(signingInput))
	digest := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}

	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)
	return signingInput + "." + encodedSignature, nil
}
