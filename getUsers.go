package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type License struct {
	SkuId string `json:"skuId"`
}

type User struct {
	DisplayName       string    `json:"displayName"`
	Mail              string    `json:"mail"`
	UserType          string    `json:"userType"`
	Id                string    `json:"id"`
	UserPrincipalName string    `json:"userPrincipalName"`
	AssignedLicenses  []License `json:"assignedLicenses"`
}

type usersResponse struct {
	Value []User `json:"value"`
}

func getUsers(token string) ([]User, error) {
	req, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/v1.0/users?$select=displayName,mail,assignedLicenses,id,userType,userPrincipalName", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status: " + resp.Status)
	}

	var ur usersResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		return nil, err
	}
	return ur.Value, nil
}
