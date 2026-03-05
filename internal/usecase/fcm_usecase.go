package usecase

import (
	"time"

	"github.com/YagoSchramm/ApiMyChat/internal/repository"
)

type PushNotifier interface {
	SendToTokens(tokens []string, title, body string, data map[string]string) error
}

type FCMUsecase struct {
	repository repository.FcmRepository
	notifier   PushNotifier
}

func NewFCMUsecase(repo repository.FcmRepository, notifier PushNotifier) *FCMUsecase {
	return &FCMUsecase{repository: repo, notifier: notifier}
}

func (u *FCMUsecase) SaveToken(userID, token string) error {
	err := u.repository.Create(userID, token, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (u *FCMUsecase) GetTokensByUid(userID string) ([]string, error) {
	tokens, err := u.repository.GetTokensByUid(userID)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (u *FCMUsecase) DeleteToken(token string) error {
	return u.repository.Delete(token)
}

func (u *FCMUsecase) NotifyUsers(userIDs []string, title, body string, data map[string]string) error {
	if u.notifier == nil {
		return nil
	}

	tokenSet := make(map[string]struct{})
	for _, userID := range userIDs {
		if userID == "" {
			continue
		}

		tokens, err := u.repository.GetTokensByUid(userID)
		if err != nil {
			return err
		}

		for _, token := range tokens {
			if token == "" {
				continue
			}
			tokenSet[token] = struct{}{}
		}
	}

	if len(tokenSet) == 0 {
		return nil
	}

	tokens := make([]string, 0, len(tokenSet))
	for token := range tokenSet {
		tokens = append(tokens, token)
	}

	return u.notifier.SendToTokens(tokens, title, body, data)
}
