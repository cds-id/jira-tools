package jira

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type JiraIssue struct {
    Key    string `json:"key"`
    Fields struct {
        Summary     string `json:"summary"`
        Description string `json:"description"`
        Status      struct {
            Name string `json:"name"`
        } `json:"status"`
        Assignee struct {
            DisplayName string `json:"displayName"`
        } `json:"assignee"`
    } `json:"fields"`
}

func ValidateCredentials(domain, email, apiToken string) error {
    url := fmt.Sprintf("https://%s/rest/api/2/myself", domain)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    auth := email + ":" + apiToken
    encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
    req.Header.Add("Authorization", "Basic "+encodedAuth)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("invalid credentials (HTTP %d)", resp.StatusCode)
    }

    return nil
}

func FetchIssue(issueKey string) (*JiraIssue, error) {
    domain := os.Getenv("JIRA_DOMAIN")
    email := os.Getenv("JIRA_EMAIL")
    apiToken := os.Getenv("JIRA_API_TOKEN")

    url := fmt.Sprintf("https://%s/rest/api/2/issue/%s", domain, issueKey)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    auth := email + ":" + apiToken
    encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
    req.Header.Add("Authorization", "Basic "+encodedAuth)
    req.Header.Add("Accept", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("failed to fetch issue: %s", string(bodyBytes))
    }

    var issue JiraIssue
    if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
        return nil, err
    }

    return &issue, nil
}
