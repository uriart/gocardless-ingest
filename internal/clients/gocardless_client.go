package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gocardless-ingest/internal/constants"
	"github.com/gocardless-ingest/internal/models"
)

const (
	baseURL              = constants.GoCardlessAPIBaseURL
	newTokenEndpoint     = "/api/v2/token/new/"
	refreshTokenEndpoint = "/api/v2/token/refresh/"
	transactionsEndpoint = "/api/v2/accounts/%s/transactions"
)

type GoCardlessClient struct {
	secretID       string
	secretKey      string
	accessToken    string
	refreshToken   string
	accessExpires  time.Time
	refreshExpires time.Time
	mu             sync.Mutex
	http           *http.Client
}

func NewGoCardlessClient(secretID, secretKey string) *GoCardlessClient {
	return &GoCardlessClient{
		secretID:  secretID,
		secretKey: secretKey,
		http:      &http.Client{},
	}
}

// tokenResponse models GoCardless token responses
type tokenResponse struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}

// GetTransactions obtiene las transacciones de una cuenta
func (c *GoCardlessClient) GetTransactions(accountID string) ([]models.Transaction, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf(baseURL+transactionsEndpoint, accountID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch transactions: %s", string(body))
	}

	var data struct {
		Transactions struct {
			Booked  []models.Transaction `json:"booked"`
			Pending []models.Transaction `json:"pending"`
		} `json:"transactions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Transactions.Booked, nil

}

// ensureValidToken checks token validity and refreshes or gets a new one
func (c *GoCardlessClient) ensureValidToken() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	if c.accessToken == "" || now.After(c.accessExpires) {
		// intenta refrescar si se puede
		if c.refreshToken != "" && now.Before(c.refreshExpires) {
			if err := c.refreshAccessToken(); err == nil {
				return nil
			}
		}
		// si no se puede refrescar o falla, se obtiene uno nuevo
		return c.fetchNewToken()
	}
	return nil
}

// fetchNewToken obtiene un nuevo token desde /token/new/
func (c *GoCardlessClient) fetchNewToken() error {
	payload := map[string]string{
		"secret_id":  c.secretID,
		"secret_key": c.secretKey,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", baseURL+newTokenEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to obtain token: %s", string(respBody))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	now := time.Now()
	c.accessToken = tokenResp.Access
	c.accessExpires = now.Add(time.Duration(tokenResp.AccessExpires-60) * time.Second) // 1 minuto de margen
	c.refreshToken = tokenResp.Refresh
	c.refreshExpires = now.Add(time.Duration(tokenResp.RefreshExpires-60) * time.Second)

	return nil
}

// refreshAccessToken refresca el access token usando el refresh token
func (c *GoCardlessClient) refreshAccessToken() error {
	payload := map[string]string{
		"refresh": c.refreshToken,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", baseURL+refreshTokenEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to refresh token: %s", string(respBody))
	}

	var tokenResp struct {
		Access        string `json:"access"`
		AccessExpires int    `json:"access_expires"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	c.accessToken = tokenResp.Access
	c.accessExpires = time.Now().Add(time.Duration(tokenResp.AccessExpires-60) * time.Second)

	return nil
}
