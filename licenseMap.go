package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SubscribedSku struct {
	SkuId         string `json:"skuId"`
	SkuPartNumber string `json:"skuPartNumber"`
}

type SubscribedSkusResponse struct {
	Value []SubscribedSku `json:"value"`
}

// getLicenseMap retrieves a map of SKU ID (uppercase) to SKU part number.
func getLicenseMap(token string) (map[string]string, error) {
	const endpoint = "https://graph.microsoft.com/v1.0/subscribedSkus?$select=skuId,skuPartNumber"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var skus SubscribedSkusResponse
	if err := json.Unmarshal(body, &skus); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	licenseMap := make(map[string]string, len(skus.Value))
	for _, s := range skus.Value {
		licenseMap[strings.ToUpper(s.SkuId)] = s.SkuPartNumber
	}
	return licenseMap, nil
}

var prettyNames = map[string]string{
	"EXCHANGEONLINE_PLAN1": "Exchange Online (Plan 1)",
	"EXCHANGEONLINE_PLAN2": "Exchange Online (Plan 2)",
	"ENTERPRISEPACK":       "Microsoft 365 E3",
	"SPE_E3":               "Microsoft 365 E3",
	"ENTERPRISEPREMIUM":    "Microsoft 365 E5",
	"BUSINESS_PREMIUM":     "Microsoft 365 Business Premium",
	"SPB":                  "Microsoft 365 Business Premium",
	"BUSINESS_STANDARD":    "Microsoft 365 Business Standard",
	"VISIOCLIENT":          "Visio Plan 2",
	"POWER_BI_PRO":         "Power BI Pro",
	"FLOW_FREE":            "Microsoft Power Automate (kostenlos)",
}
