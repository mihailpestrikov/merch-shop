package e2e

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080/api"

// TestAuth_BuyMerch_GetInfo_Success проверяет сценарий покупки товара и проверки инвентаря.
func TestAuth_BuyMerch_GetInfo_Success(t *testing.T) {
	client := &http.Client{}

	token, err := registerAndLogin(client, "testuser2", "password123")
	assert.NoError(t, err)

	itemName := "t-shirt"
	err = buyItem(client, token, itemName)
	assert.NoError(t, err)

	inventory, err := getUserInventory(client, token)
	assert.NoError(t, err)
	assert.True(t, containsItem(inventory, itemName), "Purchased item not found in inventory")
}

// TestAuth_SendCoins_Success проверяет сценарий передачи монет самому себе.
func TestAuth_SendToYourself_Fail(t *testing.T) {
	client := &http.Client{}

	token, err := registerAndLogin(client, "testuser2", "password123")
	assert.NoError(t, err)

	amount := 100
	err = sendCoins(client, token, "testuser2", amount)
	assert.Error(t, err)
}

// TestAuth_SendCoins_Success проверяет сценарий передачи монет между пользователями.
func TestAuth_SendCoins_Success(t *testing.T) {
	client := &http.Client{}

	senderToken, err := registerAndLogin(client, "sender", "password123")
	assert.NoError(t, err)

	receiverToken, err := registerAndLogin(client, "receiver", "password123")
	assert.NoError(t, err)

	amount := 100
	err = sendCoins(client, senderToken, "receiver", amount)
	assert.NoError(t, err)

	senderBalance, err := getUserBalance(client, senderToken)
	assert.NoError(t, err)
	assert.Equal(t, 900, senderBalance, "Incorrect sender balance after transaction")

	receiverBalance, err := getUserBalance(client, receiverToken)
	assert.NoError(t, err)
	assert.Equal(t, 1100, receiverBalance, "Incorrect receiver balance after transaction")

	err = sendCoins(client, receiverToken, "sender", amount)
	assert.NoError(t, err)
}

func registerAndLogin(client *http.Client, username, password string) (string, error) {
	reqBody := map[string]interface{}{
		"username": username,
		"password": password,
	}

	resp, err := sendPostRequest(client, baseURL+"/auth", reqBody, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	token, exists := data["token"]
	if !exists {
		return "", errors.New("token not found in response")
	}

	return token, nil
}

func buyItem(client *http.Client, token, itemName string) error {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/merch/"+itemName, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to buy item")
	}

	return nil
}

func getUserInventory(client *http.Client, token string) ([]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/info", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	inventory, exists := info["inventory"].([]interface{})
	if !exists {
		return nil, errors.New("inventory not found in response")
	}

	return inventory, nil
}

func containsItem(inventory []interface{}, itemName string) bool {
	for _, item := range inventory {
		itemMap := item.(map[string]interface{})
		if itemMap["type"] == itemName {
			return true
		}
	}
	return false
}

func sendCoins(client *http.Client, senderToken, toUser string, amount int) error {
	reqBody := map[string]interface{}{
		"toUser": toUser,
		"amount": amount,
	}

	resp, err := sendPostRequest(client, baseURL+"/sendCoin", reqBody, &senderToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send coins")
	}

	return nil
}

func getUserBalance(client *http.Client, token string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/info", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return 0, err
	}

	balance, ok := info["coins"].(float64)
	if !ok {
		return 0, errors.New("balance not found in response")
	}

	return int(balance), nil
}

func sendPostRequest(client *http.Client, url string, body map[string]interface{}, token *string) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}
	return client.Do(req)
}
