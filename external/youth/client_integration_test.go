//go:build integration

package youth

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
	"github.com/janghanul090801/pico-backend/internal/httpclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, "../..", ".env")

	if _, err := os.Stat(envPath); err != nil {
		log.Printf(
			"Integration test: .env not found at %s, skipping env load",
			envPath,
		)
		return
	}

	viper.SetConfigFile(envPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf(
			"Integration test: failed to read .env: %v",
			err,
		)
		return
	}

	if err := viper.Unmarshal(&config.E); err != nil {
		log.Printf(
			"Integration test: failed to unmarshal env: %v",
			err,
		)
	}

	// 환경변수로 APP_ENV를 직접 지정한 경우 우선 적용
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		config.E.AppEnv = appEnv
	}
}

func newRealClient() *Client {
	return NewClient(
		newHTTPClient(),
		config.E.YouthApiKey,
	)
}

func newHTTPClient() *httpclient.Client {
	return httpclient.NewClient(
		&http.Client{
			Timeout: 10 * time.Second,
		},
	)
}

func TestIntegration_Fetch(t *testing.T) {
	if config.E.YouthApiKey == "" {
		t.Skip("YOUTH_API_KEY not set in .env")
	}

	client := newRealClient()
	ctx := context.Background()

	policies, err := client.Fetch(ctx)

	require.NoError(t, err, "Fetch should not return error")
	assert.NotNil(t, policies)
	assert.Greater(t, len(policies), 0, "should fetch at least 1 youth policy")

	t.Logf("Fetched %d youth policies", len(policies))

	for i, policy := range policies {
		if i >= 5 {
			break
		}

		t.Logf(
			"[Youth %d] Title=%q, Description=%q, RegionCode=%d, Address=%q, Lat=%.4f, Lng=%.4f",
			i,
			policy.Title,
			policy.Description,
			policy.RegionCode,
			policy.Address,
			policy.Latitude,
			policy.Longitude,
		)

		assert.NotEmpty(
			t,
			policy.Title,
			"policy title should not be empty",
		)
	}
}
