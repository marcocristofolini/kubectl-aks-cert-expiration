package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubscriptionID(t *testing.T) {
	// Mock the Azure API response
	mockResponse := `
	{
		"id": "/subscriptions/{subscription-id}"
	}`

	// Mock the Azure CLI command output
	mockOutput := `
	{
		"id": "/subscriptions/{subscription-id}"
	}
	`

	// Create a fake command executor that returns the mock output
	executor := func(cmd string, args ...string) ([]byte, error) {
		return []byte(mockOutput), nil
	}

	// Create a fake Azure API client that returns the mock response
	client := &fakeClient{
		response: mockResponse,
	}

	// Create a new azureSubscription object
	subscription := &AzureSubscription{
		client:   client,
		executor: executor,
	}

	// Call the GetSubscriptionID method
	subscriptionID, err := subscription.GetSubscriptionID()

	// Check that the method returns the correct subscription ID and no error
	assert.NoError(t, err)
	assert.Equal(t, "{subscription-id}", subscriptionID)
}

// A fake client that always returns the same response
type fakeClient struct {
	response string
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(c.response))),
	}, nil
}
