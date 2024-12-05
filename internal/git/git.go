package git

import (
	"fmt"
    "os/exec"
    "strings"
    "jira-tools/internal/config"
)

type BranchType string

const (
    FeatureBranch BranchType = "feature"
    BugfixBranch  BranchType = "bugfix"
    HotfixBranch  BranchType = "hotfix"
    ReleaseBranch BranchType = "release"
)

func GetAvailableBranches() ([]string, error) {
    cmd := exec.Command("git", "branch", "--format=%(refname:short)")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to get branches: %v", err)
    }

    branches := strings.Split(strings.TrimSpace(string(output)), "\n")
    return branches, nil
}

func CreateBranch(issueKey, summary string, branchType BranchType) error {
    branchConfig, err := config.LoadBranchConfig()
    if err != nil {
        return fmt.Errorf("failed to load branch configuration: %v", err)
    }

    // Determine base branch
    var baseBranch string
    if branchConfig.IsMonorepo {
        switch branchType {
        case FeatureBranch, BugfixBranch:
            baseBranch = branchConfig.DevelopmentBranch
        case HotfixBranch, ReleaseBranch:
            baseBranch = branchConfig.ProductionBranch
        default:
            return fmt.Errorf("invalid branch type: %s", branchType)
        }
    } else {
        baseBranch = branchConfig.DevelopmentBranch
    }

    // Create branch name
    branchName := fmt.Sprintf("%s/%s-%s", branchType, issueKey, formatBranchName(summary))

    // Checkout base branch
    if err := exec.Command("git", "checkout", baseBranch).Run(); err != nil {
        return fmt.Errorf("failed to checkout %s branch: %v", baseBranch, err)
    }

    // Create and checkout new branch
    if err := exec.Command("git", "checkout", "-b", branchName).Run(); err != nil {
        return fmt.Errorf("failed to create branch: %v", err)
    }

    fmt.Printf("Created branch %s from %s\n", branchName, baseBranch)
    return nil
}

func CommitChanges(issueKey, commitType, summary string) error {
    // Stage all changes
    if err := exec.Command("git", "add", ".").Run(); err != nil {
        return fmt.Errorf("failed to stage changes: %v", err)
    }

    // Create commit message
    message := fmt.Sprintf("%s(%s): %s", commitType, issueKey, summary)

    // Commit changes
    if err := exec.Command("git", "commit", "-m", message).Run(); err != nil {
        return fmt.Errorf("failed to commit changes: %v", err)
    }

    fmt.Printf("Changes committed with message: %s\n", message)
    return nil
}

func PushBranch() error {
    // Get current branch
    cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
    branchBytes, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("failed to get current branch: %v", err)
    }
    branch := strings.TrimSpace(string(branchBytes))

    // Push to remote
    if err := exec.Command("git", "push", "-u", "origin", branch).Run(); err != nil {
        return fmt.Errorf("failed to push branch: %v", err)
    }

    fmt.Printf("Successfully pushed branch %s to remote\n", branch)
    return nil
}

func formatBranchName(name string) string {
    // Convert to lowercase
    name = strings.ToLower(name)
    // Replace spaces and special characters with hyphens
    name = strings.Map(func(r rune) rune {
        if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
            return r
        }
        return '-'
    }, name)
    // Remove consecutive hyphens
    for strings.Contains(name, "--") {
        name = strings.ReplaceAll(name, "--", "-")
    }
    // Trim hyphens from ends
    return strings.Trim(name, "-")
}
