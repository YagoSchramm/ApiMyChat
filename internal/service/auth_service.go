package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/service/model"
)

type SupabaseAuthService struct {
	Url string
	Key string
}

func NewSupabaseAuthService(url, key string) *SupabaseAuthService {
	return &SupabaseAuthService{
		Url: url,
		Key: key,
	}
}
func (s *SupabaseAuthService) Login(email, password string) (*model.LoginResponse, error) {

	body := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		s.Url+"/auth/v1/token?grant_type=password",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("apikey", s.Key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("invalid credentials")
	}

	var response model.LoginResponse
	json.NewDecoder(res.Body).Decode(&response)

	return &response, nil
}
