//go:build integration

package kakaomap

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/janghanul090801/pico-backend/config"
	"github.com/janghanul090801/pico-backend/internal/httpclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if _, err := os.Stat("../../.env"); err != nil {
		return
	}

	viper.SetConfigFile("../../.env")

	if err := viper.ReadInConfig(); err != nil {
		return
	}

	if err := viper.Unmarshal(&config.E); err != nil {
		return
	}

	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		config.E.AppEnv = appEnv
	}
}

func newRealClient() *Client {
	return &Client{
		httpClient: httpclient.NewClient(
			&http.Client{
				Timeout: 10 * time.Second,
			},
		),
		apiKey: config.E.KakaoMapApiKey,
	}
}

func TestIntegration_Geocode(t *testing.T) {
	if config.E.KakaoMapApiKey == "" {
		t.Skip("KAKAO_API_KEY not set in .env")
	}

	t.Logf("APP_ENV=%q", config.E.AppEnv)

	client := newRealClient()
	ctx := context.Background()

	latitude, longitude, err := client.Geocode(
		ctx,
		"서울특별시 종로구 세종대로 175",
	)

	require.NoError(t, err, "Geocode should not return error")

	t.Logf(
		"Geocode result: latitude=%.6f, longitude=%.6f",
		latitude,
		longitude,
	)

	assert.NotZero(t, latitude, "latitude should not be zero")
	assert.NotZero(t, longitude, "longitude should not be zero")

	// 대한민국 좌표 범위에 있는지도 확인
	assert.Greater(t, latitude, float64(33))
	assert.Less(t, latitude, float64(39))

	assert.Greater(t, longitude, float64(124))
	assert.Less(t, longitude, float64(132))
}

func TestIntegration_Geocode_InvalidAddress(t *testing.T) {
	if config.E.KakaoMapApiKey == "" {
		t.Skip("KAKAO_API_KEY not set in .env")
	}

	client := newRealClient()
	ctx := context.Background()

	_, _, err := client.Geocode(
		ctx,
		"존재하지 않는 주소 123456789",
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "address not found")
}
