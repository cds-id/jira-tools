package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

type BranchConfig struct {
    ProjectPath        string `json:"project_path"`
    ProductionBranch   string `json:"production_branch,omitempty"`
    DevelopmentBranch  string `json:"development_branch"`
    IsMonorepo         bool   `json:"is_monorepo"`
}

// Add new functions to handle project-specific configs
func GetProjectConfigPath(projectPath string) (string, error) {
    if projectPath == "" {
        return "", fmt.Errorf("project path cannot be empty")
    }
    // Verify the project path exists
    if _, err := os.Stat(projectPath); err != nil {
        return "", fmt.Errorf("invalid project path: %v", err)
    }
    return filepath.Join(projectPath, ".jt-config.json"), nil
}

func LoadProjectBranchConfig(projectPath string) (*BranchConfig, error) {
    configPath, err := GetProjectConfigPath(projectPath)
    if err != nil {
        return nil, err
    }

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

func SaveProjectBranchConfig(config *BranchConfig) error {
    if config.ProjectPath == "" {
        return fmt.Errorf("project path cannot be empty")
    }

    configPath, err := GetProjectConfigPath(config.ProjectPath)
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0600)
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
