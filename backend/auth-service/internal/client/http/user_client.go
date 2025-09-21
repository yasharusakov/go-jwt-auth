package httpClient

import (
	"auth-service/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error)
}

type UserClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPUserClient(baseURL string) UserService {
	return &UserClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (u *UserClient) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	url := u.BaseURL + "/user/get-by-email/" + email

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user by email: status code %d", resp.StatusCode)
	}

	var user model.User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	log.Println(&user)

	return &user, nil
}

func (u *UserClient) GetUserByID(ctx context.Context, id string) (*model.UserWithoutPassword, error) {
	url := u.BaseURL + "/user/get-by-id/" + id

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user by ID: status code %d", resp.StatusCode)
	}
	var user model.UserWithoutPassword
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserClient) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	url := u.BaseURL + "/user/check-by-email/" + email
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to check user existence by email: status code %d", resp.StatusCode)
	}

	var result struct {
		Exists bool `json:"exists"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Exists, nil
}

func (u *UserClient) RegisterUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	url := u.BaseURL + "/user/register"

	payload := map[string]interface{}{
		"email":    email,
		"password": hashedPassword,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to register user: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}
