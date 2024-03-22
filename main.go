package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetAccessToken(clientId string, clientSecret string) (*GetAccessTokenResponse, error) {
	req, err := http.NewRequest("GET", "https://api.domo.com/oauth/token?grant_type=client_credentials&scope=data", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(clientId, clientSecret)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body *GetAccessTokenResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}
	return body, nil
}

func ListDataSets(accessToken string) []byte {
	req, _ := http.NewRequest("GET", "https://api.domo.com/v1/datasets?offet=149", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func ListDataSetPolicies(accessToken string, datasetId string) []Policy {
	req, _ := http.NewRequest("GET", "https://api.domo.com/v1/datasets/"+datasetId+"/policies", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	policies := make([]Policy, 0)
	err = json.NewDecoder(res.Body).Decode(&policies)
	if err != nil {
		log.Fatal(err)
	}
	return policies
}

func CreateDataSetPolicy(accessToken string, datasetId string, policy CreatePolicyRequest) Policy {
	body, _ := json.Marshal(policy)
	req, _ := http.NewRequest("POST", "https://api.domo.com/v1/datasets/"+datasetId+"/policies", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Create Policy Response: %s", resBody)
	var p Policy
	err = json.Unmarshal(resBody, &p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func DeletePolicy(accessToken string, datasetId string, policyId int) bool {
	req, _ := http.NewRequest("DELETE", "https://api.domo.com/v1/datasets/"+datasetId+"/policies/"+strconv.Itoa(policyId), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	return res.StatusCode == 204
}

func main() {
	id, secret := "...", "..."
	t, _ := GetAccessToken(id, secret)
	log.Printf("Access Token: %s", t.AccessToken)
	os.WriteFile("datasets.json", ListDataSets(t.AccessToken), 0644)
	log.Printf("Datasets loaded")

	pdpTestDatasetId := "da00c975-3a56-47bc-8bf8-046751dd52ea"
	ds1Id := "358a54a2-1d4e-47b5-a31d-6c92a4d1a5bc"
	policies := ListDataSetPolicies(t.AccessToken, ds1Id)

	for _, p := range policies {
		fmt.Printf("Creating a policy: %v", p)
		created := CreateDataSetPolicy(t.AccessToken, pdpTestDatasetId, CreatePolicyRequest{
			Type:         "user",
			Name:         p.Name,
			Filters:      p.Filters,
			Users:        p.Users,
			VirtualUsers: p.VirtualUsers,
			Groups:       p.Groups,
		})
		fmt.Printf("; ID: %v\n", created)
	}

}
