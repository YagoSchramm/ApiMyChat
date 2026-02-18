package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	UID         string    `json:"id"`
	Email       string    `json:"email" binding:"required,email"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" binding:"required"`
	Password    string    `json:"password" binding:"required"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type userAlias struct {
		UID         string      `json:"id"`
		Email       string      `json:"email"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		CreatedAt   interface{} `json:"createdAt"`
		Password    string      `json:"password"`
	}

	var aux userAlias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedCreatedAt, err := parseCreatedAt(aux.CreatedAt)
	if err != nil {
		return err
	}

	u.UID = aux.UID
	u.Email = aux.Email
	u.Name = aux.Name
	u.Description = aux.Description
	u.CreatedAt = parsedCreatedAt
	u.Password = aux.Password
	return nil
}

func parseCreatedAt(raw interface{}) (time.Time, error) {
	switch v := raw.(type) {
	case string:
		layouts := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05.000 -0700",
			"2006-01-02 15:04:05 -0700",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid createdAt format: %q", v)
	case float64:
		// Accept Unix timestamp in seconds or milliseconds.
		ts := int64(v)
		if ts > 1_000_000_000_000 {
			return time.UnixMilli(ts), nil
		}
		return time.Unix(ts, 0), nil
	case nil:
		return time.Time{}, fmt.Errorf("createdAt is required")
	default:
		return time.Time{}, fmt.Errorf("invalid createdAt type: %T", raw)
	}
}
