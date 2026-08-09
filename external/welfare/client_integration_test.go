//go:build integration

package welfare

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/janghanul090801/pico-backend/config"
	"github.com/janghanul090801/pico-backend/domain/mocks"
	"github.com/janghanul090801/pico-backend/internal/httpclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, "../..", ".env")

	if _, err := os.Stat(envPath); err != nil {
		log.Printf("Integration test: .env not found at %s, skipping env load", envPath)
		return
	}

	viper.SetConfigFile(envPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Integration test: failed to read .env: %v", err)
		return
	}
	if err := viper.Unmarshal(&config.E); err != nil {
		log.Printf("Integration test: failed to unmarshal env: %v", err)
	}

	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		config.E.AppEnv = appEnv
	}
}

func newRealClient() *Client {
	policyRepository := new(mocks.PolicyRepository{})
	policyRepository.On(
		"FindByTitle",
		mock.Anything,
		mock.AnythingOfType("string"),
	).Return(nil, nil)

	return &Client{
		httpClient:       newHTTPClient(),
		policyRepository: policyRepository,
		apiKey:           config.E.WelfareApiKey,
	}
}

func newHTTPClient() *httpclient.Client {
	return httpclient.NewClient(
		&http.Client{
			Timeout: time.Second * 10,
		},
	)
}

func TestIntegration_Fetch(t *testing.T) {
	client := newRealClient()
	ctx := context.Background()

	policies, err := client.Fetch(ctx)
	require.NoError(t, err, "Fetch should not return error")

	t.Logf("Fetched %d welfare policies", len(policies))
	assert.Greater(t, len(policies), 0, "should fetch at least 1 welfare policy")

	for i, p := range policies {
		if i >= 3 {
			break
		}
		t.Logf("[Welfare %d] Title=%q, Address=%q, RegionCode=%d", i, p.Title, p.Address, p.RegionCode)
		assert.NotEmpty(t, p.Title, "policy title should not be empty")
		assert.NotEmpty(t, p.Description, "policy description should not be empty")
		assert.False(t, p.StartDate.IsZero(), "start date should be set")
		assert.False(t, p.EndDate.IsZero(), "end date should be set")
	}
}

func TestIntegration_list(t *testing.T) {
	client := newRealClient()
	ctx := context.Background()

	items, err := client.list(ctx)
	require.NoError(t, err, "list should not return error")

	t.Logf("Fetched %d welfare items from list API", len(items))
	assert.Greater(t, len(items), 0, "should fetch at least 1 welfare item")

	for i, item := range items {
		if i >= 5 {
			break
		}
		t.Logf("[WelfareItem %d] ServID=%q, ServNm=%q", i, item.ServID, item.ServNm)
		assert.NotEmpty(t, item.ServID, "ServID should not be empty")
	}
}

func TestIntegration_detail(t *testing.T) {
	if config.E.WelfareApiKey == "" {
		t.Skip("WELFARE_API_KEY not set in .env")
	}

	client := newRealClient()
	ctx := context.Background()

	// First, get a real servID from the list
	items, err := client.list(ctx)
	require.NoError(t, err)
	require.Greater(t, len(items), 0, "need at least 1 item to test detail")

	servID := items[0].ServID
	t.Logf("Testing detail for ServID=%q", servID)

	policy, err := client.detail(ctx, servID)
	require.NoError(t, err, "detail should not return error")

	t.Logf("Detail: Title=%q, Address=%q, Description(first 100)=%q",
		policy.Title, policy.Address, truncate(policy.Description, 100))

	assert.NotEmpty(t, policy.Title, "title should not be empty")
	assert.NotEmpty(t, policy.Description, "description should not be empty")
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
