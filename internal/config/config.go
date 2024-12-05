package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type BranchConfig struct {
    ProductionBranch   string `json:"production_branch,omitempty"`
    DevelopmentBranch  string `json:"development_branch"`
    IsMonorepo         bool   `json:"is_monorepo"`
}

type Config struct {
    JiraDomain  string `json:"jira_domain"`
    JiraEmail   string `json:"jira_email"`
    JiraToken   string `json:"jira_token"`
    Branch      BranchConfig `json:"branch_config"`
}

func GetConfigPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    configDir := filepath.Join(homeDir, ".jira-tools")
    if err := os.MkdirAll(configDir, 0700); err != nil {
        return "", err
    }

    return configDir, nil
}

func LoadBranchConfig() (*BranchConfig, error) {
    configDir, err := GetConfigPath()
    if err != nil {
        return nil, err
    }

    configPath := filepath.Join(configDir, "branch-config.json")
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config BranchConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}

func SaveBranchConfig(config *BranchConfig) error {
    configDir, err := GetConfigPath()
    if err != nil {
        return err
    }

    configPath := filepath.Join(configDir, "branch-config.json")
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0600)
}
