package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type BMIMetricResponse struct {
	BMI float64 `json:"bmi"`
}

type BMIAPIResponse struct {
	WeightCategory string `json:"weightCategory"`
}

func GetBMIAndCategory(weight int, heightCm int) (float64, string, error) {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		return 0, "", fmt.Errorf("RAPIDAPI_KEY not set in env")
	}

	heightMeters := float64(heightCm) / 100
	bmiURL := fmt.Sprintf("https://body-mass-index-bmi-calculator.p.rapidapi.com/metric?weight=%d&height=%.2f", weight, heightMeters)

	req, _ := http.NewRequest("GET", bmiURL, nil)
	req.Header.Add("x-rapidapi-host", "body-mass-index-bmi-calculator.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	var bmiResp BMIMetricResponse
	if err := json.NewDecoder(res.Body).Decode(&bmiResp); err != nil {
		return 0, "", err
	}

	// second request: category
	categoryURL := fmt.Sprintf("https://body-mass-index-bmi-calculator.p.rapidapi.com/weight-category?bmi=%.1f", bmiResp.BMI)
	req2, _ := http.NewRequest("GET", categoryURL, nil)
	req2.Header.Add("x-rapidapi-host", "body-mass-index-bmi-calculator.p.rapidapi.com")
	req2.Header.Add("x-rapidapi-key", apiKey)

	res2, err := client.Do(req2)
	if err != nil {
		return bmiResp.BMI, "", err
	}
	defer res2.Body.Close()

	var catResp BMIAPIResponse
	if err := json.NewDecoder(res2.Body).Decode(&catResp); err != nil {
		return bmiResp.BMI, "", err
	}

	return bmiResp.BMI, catResp.WeightCategory, nil
}
