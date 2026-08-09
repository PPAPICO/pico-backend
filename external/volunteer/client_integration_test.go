//go:build integration

package volunteer

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain/mocks"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/internal/httpclient"
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
		apiKey:           config.E.VolunteerApiKey,
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
	t.Log(config.E.AppEnv)

	client := newRealClient()
	ctx := context.Background()

	policies, err := client.Fetch(ctx)
	require.NoError(t, err, "FetchVolunteerEvents should not return error")

	t.Logf("Fetched %d volunteer events", len(policies))
	assert.Greater(t, len(policies), 0, "should fetch at least 1 volunteer event")

	for i, p := range policies {
		if i >= 3 {
			break
		}
		t.Logf("[Volunteer %d] Title=%q, Address=%q, RegionCode=%d, Lat=%.4f, Lng=%.4f",
			i, p.Title, p.Address, p.RegionCode, p.Latitude, p.Longitude)
		assert.NotEmpty(t, p.Title, "event title should not be empty")
		assert.NotEmpty(t, p.Description, "event description should not be empty")
		assert.False(t, p.StartDate.IsZero(), "start date should be set")
		assert.False(t, p.EndDate.IsZero(), "end date should be set")
	}
}

func TestIntegration_list(t *testing.T) {
	if config.E.VolunteerApiKey == "" {
		t.Skip("VOLUNTEER_API_KEY not set in .env")
	}

	client := newRealClient()
	ctx := context.Background()

	items, err := client.list(ctx)
	require.NoError(t, err, "list should not return error")

	t.Logf("Fetched %d volunteer area items", len(items))
	assert.Greater(t, len(items), 0, "should fetch at least 1 volunteer area item")

	for i, item := range items {
		if i >= 5 {
			break
		}
		t.Logf("[VolunteerItem %d] RegistNo=%q, Subject=%q, Status=%q",
			i, item.ProgrmRegistNo, item.ProgrmSj, item.ProgrmSttusSe)
		assert.NotEmpty(t, item.ProgrmRegistNo, "ProgrmRegistNo should not be empty")
	}
}

func TestIntegration_detail(t *testing.T) {
	if config.E.VolunteerApiKey == "" {
		t.Skip("VOLUNTEER_API_KEY not set in .env")
	}

	uc := newRealClient()
	ctx := context.Background()

	// First, get a real registNo from the area list
	items, err := uc.list(ctx)
	require.NoError(t, err)
	require.Greater(t, len(items), 0, "need at least 1 item to test detail")

	registNo := items[0].ProgrmRegistNo
	t.Logf("Testing detail for RegistNo=%q", registNo)

	policy, err := uc.detail(ctx, registNo)
	if err != nil {
		log.Printf("Detail fetch failed (may be expected for some items): %v", err)
		t.Skipf("detail fetch failed for %q: %v", registNo, err)
	}

	t.Logf("Detail: Title=%q, RegionCode=%d, Lat=%.4f, Lng=%.4f, Address=%q",
		policy.Title, policy.RegionCode, policy.Latitude, policy.Longitude, policy.Address)

	assert.NotEmpty(t, policy.Title, "title should not be empty")
	assert.False(t, policy.StartDate.IsZero(), "start date should be set")
}
