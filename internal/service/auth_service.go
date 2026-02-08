package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/service/model"
)

type SupabaseService struct {
	Url string
	Key string
}

func NewSupabaseAuthService(url, key string) *SupabaseService {
	return &SupabaseService{
		Url: url,
		Key: key,
	}
}
func (s *SupabaseService) Login(email, password string) (*model.LoginResponse, error) {

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
func (s *SupabaseService) CreateUser(email, password string) (*entity.LoginUserResponse, error) {

	body := entity.LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		s.Url+"/auth/v1/admin/users",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", s.Key)
	req.Header.Set("Authorization", "Bearer "+s.Key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, errors.New("erro ao criar usuário no supabase")
	}

	var response entity.LoginUserResponse
	json.NewDecoder(res.Body).Decode(&response)

	return &response, nil
}
func (s *SupabaseService) UserExists(userID string) (bool, error) {

	req, err := http.NewRequest(
		"GET",
		s.Url+"/auth/v1/admin/users/"+userID,
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("apikey", s.Key)
	req.Header.Set("Authorization", "Bearer "+s.Key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return true, nil
	}

	if res.StatusCode == 404 {
		return false, nil
	}

	return false, errors.New("erro ao verificar usuário no supabase")
}
