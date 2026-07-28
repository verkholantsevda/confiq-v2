package warp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiBase   = "https://api.cloudflareclient.com/v0i1909051800"
	userAgent = "okhttp/3.12.1"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type registerRequest struct {
	InstallID string `json:"install_id"`
	TOS       string `json:"tos"`
	Key       string `json:"key"`
	FCMToken  string `json:"fcm_token"`
	Type      string `json:"type"`
	Locale    string `json:"locale"`
}

type registerResponse struct {
	Result struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	} `json:"result"`
}

func (c *Client) RegisterDevice(publicKey string) (*Device, error) {
	reqBody := registerRequest{
		InstallID: "",
		TOS:       time.Now().UTC().Format(time.RFC3339),
		Key:       publicKey,
		FCMToken:  "",
		Type:      "ios",
		Locale:    "en_US",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		apiBase+"/reg",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloudflare returned %d", resp.StatusCode)
	}

	var result registerResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Result.ID == "" || result.Result.Token == "" {
		return nil, fmt.Errorf("invalid cloudflare response")
	}

	return &Device{
		ID:    result.Result.ID,
		Token: result.Result.Token,
	}, nil
}

type enableWarpRequest struct {
	WarpEnabled bool `json:"warp_enabled"`
}

type enableWarpResponse struct {
	Result struct {
		Config struct {
			Peers []struct {
				PublicKey string `json:"public_key"`
			} `json:"peers"`

			Interface struct {
				Addresses struct {
					V4 string `json:"v4"`
					V6 string `json:"v6"`
				} `json:"addresses"`
			} `json:"interface"`
		} `json:"config"`
	} `json:"result"`
}

func (c *Client) EnableWarp(deviceID, token string) (*Interface, error) {
	reqBody := enableWarpRequest{
		WarpEnabled: true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPatch,
		apiBase+"/reg/"+deviceID,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cloudflare returned status %d: %s", resp.StatusCode, string(body))
	}

	var result enableWarpResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Result.Config.Peers) == 0 {
		return nil, fmt.Errorf("peer public key not found")
	}

	return &Interface{
		PeerPublicKey: result.Result.Config.Peers[0].PublicKey,
		ClientIPv4:    result.Result.Config.Interface.Addresses.V4,
		ClientIPv6:    result.Result.Config.Interface.Addresses.V6,
	}, nil
}

func (c *Client) DeleteDevice(deviceID, token string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		apiBase+"/reg/"+deviceID,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Cloudflare может вернуть как 200, так и 204
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {

		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"cloudflare returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return nil
}
