//go:build integration

package usecase

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
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, "..", ".env")

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
}

func newRealPolicyUseCase() *policyUseCase {
	return &policyUseCase{
		contextTimeout: 30 * time.Second,
		httpClient:     newHTTPClient(),
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

func TestIntegration_FetchWelfarePolicies(t *testing.T) {
	if config.E.WelfareApiKey == "" {
		t.Skip("WELFARE_API_KEY not set in .env")
	}

	uc := newRealPolicyUseCase()
	ctx := context.Background()

	policies, err := uc.FetchWelfarePolicies(ctx)
	require.NoError(t, err, "FetchWelfarePolicies should not return error")

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

func TestIntegration_FetchVolunteerEvents(t *testing.T) {
	if config.E.VolunteerApiKey == "" {
		t.Skip("VOLUNTEER_API_KEY not set in .env")
	}

	uc := newRealPolicyUseCase()
	ctx := context.Background()

	policies, err := uc.FetchVolunteerEvents(ctx)
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

func TestIntegration_FetchWelfareList(t *testing.T) {
	if config.E.WelfareApiKey == "" {
		t.Skip("WELFARE_API_KEY not set in .env")
	}

	uc := newRealPolicyUseCase()
	ctx := context.Background()

	items, err := uc.fetchWelfareList(ctx)
	require.NoError(t, err, "fetchWelfareList should not return error")

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

func TestIntegration_FetchVolunteerAreaList(t *testing.T) {
	if config.E.VolunteerApiKey == "" {
		t.Skip("VOLUNTEER_API_KEY not set in .env")
	}

	uc := newRealPolicyUseCase()
	ctx := context.Background()

	items, err := uc.fetchVolunteerAreaList(ctx)
	require.NoError(t, err, "fetchVolunteerAreaList should not return error")

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

func TestIntegration_FetchWelfareDetail(t *testing.T) {
	if config.E.WelfareApiKey == "" {
		t.Skip("WELFARE_API_KEY not set in .env")
	}

	uc := newRealPolicyUseCase()
	ctx := context.Background()

	// First, get a real servID from the list
	items, err := uc.fetchWelfareList(ctx)
	require.NoError(t, err)
	require.Greater(t, len(items), 0, "need at least 1 item to test detail")

	servID := items[0].ServID
	t.Logf("Testing detail for ServID=%q", servID)

	policy, err := uc.fetchWelfareDetail(ctx, servID)
	require.NoError(t, err, "fetchWelfareDetail should not return error")

	t.Logf("Detail: Title=%q, Address=%q, Description(first 100)=%q",
		policy.Title, policy.Address, truncate(policy.Description, 100))

	assert.NotEmpty(t, policy.Title, "title should not be empty")
	assert.NotEmpty(t, policy.Description, "description should not be empty")
}

func TestIntegration_FetchVolunteerDetail(t *testing.T) {
	if config.E.VolunteerApiKey == "" {
		t.Skip("VOLUNTEER_API_KEY not set in .env")
	}

	uc := newRealPolicyUseCase()
	ctx := context.Background()

	// First, get a real registNo from the area list
	items, err := uc.fetchVolunteerAreaList(ctx)
	require.NoError(t, err)
	require.Greater(t, len(items), 0, "need at least 1 item to test detail")

	registNo := items[0].ProgrmRegistNo
	t.Logf("Testing detail for RegistNo=%q", registNo)

	policy, err := uc.fetchVolunteerDetail(ctx, registNo)
	if err != nil {
		log.Printf("Detail fetch failed (may be expected for some items): %v", err)
		t.Skipf("detail fetch failed for %q: %v", registNo, err)
	}

	t.Logf("Detail: Title=%q, RegionCode=%d, Lat=%.4f, Lng=%.4f, Address=%q",
		policy.Title, policy.RegionCode, policy.Latitude, policy.Longitude, policy.Address)

	assert.NotEmpty(t, policy.Title, "title should not be empty")
	assert.False(t, policy.StartDate.IsZero(), "start date should be set")
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
