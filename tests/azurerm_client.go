package test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const AZURERM_URL = "https://management.azure.com"

type AzureRMClient struct {
	GetAzureRMToken func() (string, error)
}

type RequestDefinition struct {
	HttpMethod         string
	Url                string
	Timeout            time.Duration
	SuccessStatusCodes []int
}

type AzureToken struct {
	Token     string `json:"access_token"`
	ExpiresOn int64  `json:"expires_on,string"`
}

type PolicyState struct {
	PolicyAssignmentName string `json:"policyAssignmentName"`
	ComplianceState      string `json:"complianceState"`
}

type PolicyStateQueryResults struct {
	Value []PolicyState `json:"value"`
}

func CreateAzureRMClient(tenantId string, clientId string, clientSecret string) *AzureRMClient {
	return &AzureRMClient{
		GetAzureRMToken: azureAccessToken(tenantId, clientId, clientSecret),
	}
}

func (client *AzureRMClient) request(requestDefinition RequestDefinition) ([]byte, error) {
	token, err := client.GetAzureRMToken()
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(requestDefinition.HttpMethod, requestDefinition.Url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	httpClient := http.Client{
		Timeout: time.Duration(requestDefinition.Timeout * time.Second),
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("Send request. Reason: %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read response body. Reason: %v", err)
	}

	succeeded := false
	for _, statusCode := range requestDefinition.SuccessStatusCodes {
		if statusCode == resp.StatusCode {
			succeeded = true
			break
		}
	}
	if !succeeded {
		err = fmt.Errorf("Unexpected HTTP status code %v. Reason: %v", resp.StatusCode, string(respBody))
	}

	return respBody, err
}

func (client *AzureRMClient) TriggerPolicyEvaluation(resourceGroupId string) error {
	resourceUrl := fmt.Sprintf(
		"%s/%s/providers/Microsoft.PolicyInsights/policyStates/latest/triggerEvaluation?api-version=2019-10-01",
		AZURE_MANAGEMENT_URL,
		resourceGroupId,
	)
	requestDefinition := RequestDefinition{
		HttpMethod:         http.MethodPost,
		Url:                resourceUrl,
		SuccessStatusCodes: []int{200, 202},
		Timeout:            10,
	}
	_, err := client.request(requestDefinition)
	return err
}

func (client *AzureRMClient) getComplianceState(resourceId string, policyAssignmentName string) (string, error) {
	resourceUrl := fmt.Sprintf(
		"%s/%s/providers/Microsoft.PolicyInsights/policyStates/latest/queryResults?api-version=2019-10-01",
		AZURE_MANAGEMENT_URL,
		resourceId,
	)
	requestDefinition := RequestDefinition{
		HttpMethod:         http.MethodPost,
		Url:                resourceUrl,
		SuccessStatusCodes: []int{200},
		Timeout:            20,
	}
	respBody, err := client.request(requestDefinition)
	if err != nil {
		return "", fmt.Errorf("Request error. Reason: %v %v", err, string(respBody))
	}

	policyStateQueryResults := PolicyStateQueryResults{}
	err = json.Unmarshal(respBody, &policyStateQueryResults)
	if err != nil {
		return "", fmt.Errorf("Unmarshal response body. Reason: %v %v", err, string(respBody))
	}

	complianceState := ""
	for _, policyState := range policyStateQueryResults.Value {
		if policyState.PolicyAssignmentName == policyAssignmentName {
			complianceState = policyState.ComplianceState
			break
		}
	}

	return complianceState, nil
}

func (client *AzureRMClient) GetComplianceState(resourceId string, policyAssignmentName string) (string, error) {
	for retries := 0; retries < 30; retries++ {
		complianceState, err := client.getComplianceState(resourceId, policyAssignmentName)
		if complianceState != "" || err != nil {
			return complianceState, err
		}
		retries += 1
		time.Sleep(30 * time.Second)
	}
	return "", fmt.Errorf("GetComplianceState timeout")
}

func azureAccessToken(tenantId string, clientId string, clientSecret string) func() (string, error) {
	loginUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenantId)
	loc, _ := time.LoadLocation("UTC")
	azureToken := AzureToken{
		Token:     "",
		ExpiresOn: 0,
	}
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	return func() (string, error) {
		expiresOn := time.Unix(azureToken.ExpiresOn, 0)
		if expiresOn.After(time.Now().In(loc).Add(1 * time.Minute)) {
			return azureToken.Token, nil
		}
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		data.Set("client_id", clientId)
		data.Set("client_secret", clientSecret)
		data.Set("resource", AZURERM_URL)
		req, err := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(data.Encode()))
		if err != nil {
			return "", fmt.Errorf("Create request. Reason: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret))))
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("Send request. Reason: %v", err)
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("Read response body. Reason: %v", err)
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return "", fmt.Errorf("Unexpected HTTP status code %v. Reason: %v", resp.StatusCode, string(respBody))
		}

		err = json.Unmarshal(respBody, &azureToken)
		if err != nil {
			return "", fmt.Errorf("Unmarshal response body. Reason: %v %v", err, string(respBody))
		}
		return azureToken.Token, nil
	}
}
