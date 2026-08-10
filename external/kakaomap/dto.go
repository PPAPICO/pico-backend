package kakaomap

type addressResponse struct {
	Documents []struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"documents"`
}
