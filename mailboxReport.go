package main

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"time"
)

type usageData struct {
	lastActivity string
	storageUsed  string
}

func parseCSV(r io.Reader) (map[string]usageData, error) {
	cr := csv.NewReader(r)

	// skip header
	if _, err := cr.Read(); err != nil {
		return nil, err
	}

	out := make(map[string]usageData)

	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 9 {
			continue
		}
		key := rec[1]
		out[key] = usageData{
			lastActivity: rec[6],
			storageUsed:  rec[8],
		}
	}
	return out, nil
}

const mailboxReportURL = "https://graph.microsoft.com/v1.0/reports/getMailboxUsageDetail(period='D7')"

func getMailboxReport(token string) (map[string]usageData, error) {
	req, err := http.NewRequest(http.MethodGet, mailboxReportURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected HTTP status: " + resp.Status)
	}
	return parseCSV(resp.Body)
}
