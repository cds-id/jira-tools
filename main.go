package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "jira-tools/internal/config"
    "jira-tools/internal/git"
    "jira-tools/internal/jira"
    "github.com/joho/godotenv"
)

func main() {
    if err := loadConfig(); err != nil {
        fmt.Printf("Error loading configuration: %v\n", err)
        os.Exit(1)
    }

    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    switch os.Args[1] {
    case "setup":
        if err := runSetup(); err != nil {
            fmt.Printf("Setup failed: %v\n", err)
            os.Exit(1)
        }

    case "lookup":
        if len(os.Args) < 3 {
            fmt.Println("Usage: jt lookup <card-number>")
            os.Exit(1)
        }
        if err := handleLookup(os.Args[2]); err != nil {
            fmt.Printf("Error looking up issue: %v\n", err)
            os.Exit(1)
        }

    case "branch":
        if len(os.Args) < 4 {
            fmt.Println("Usage: jt branch <card-number> <type>")
            fmt.Println("Types: feature, bugfix, hotfix, release")
            os.Exit(1)
        }
        if err := handleBranch(os.Args[2], git.BranchType(os.Args[3])); err != nil {
            fmt.Printf("Error creating branch: %v\n", err)
            os.Exit(1)
        }

    case "commit":
        if len(os.Args) < 3 {
            fmt.Println("Usage: jt commit <card-number> [type]")
            os.Exit(1)
        }
        commitType := "chore"
        if len(os.Args) >= 4 {
            commitType = os.Args[3]
        }
        if err := handleCommit(os.Args[2], commitType); err != nil {
            fmt.Printf("Error committing changes: %v\n", err)
            os.Exit(1)
        }

    case "push":
        if err := git.PushBranch(); err != nil {
            fmt.Printf("Error pushing branch: %v\n", err)
            os.Exit(1)
        }

    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        printUsage()
        os.Exit(1)
    }
}

func handleLookup(issueKey string) error {
    issue, err := jira.FetchIssue(issueKey)
    if err != nil {
        return err
    }

    printIssueDetails(issue)
    return nil
}

func handleBranch(issueKey string, branchType git.BranchType) error {
    issue, err := jira.FetchIssue(issueKey)
    if err != nil {
        return err
    }

    return git.CreateBranch(issueKey, issue.Fields.Summary, branchType)
}

func handleCommit(issueKey, commitType string) error {
    issue, err := jira.FetchIssue(issueKey)
    if err != nil {
        return err
    }

    return git.CommitChanges(issueKey, commitType, issue.Fields.Summary)
}

func printUsage() {
    fmt.Println("Usage:")
    fmt.Println("  jt setup                        - Run setup wizard")
    fmt.Println("  jt lookup <card-number>         - Look up Jira issue details")
    fmt.Println("  jt branch <card-number> <type>  - Create branch from Jira issue")
    fmt.Println("  jt commit <card-number> [type]  - Commit changes with Jira issue summary")
    fmt.Println("  jt push                         - Push current branch to remote")
    fmt.Println("\nBranch types:")
    fmt.Println("  feature  - New feature branch (from development)")
    fmt.Println("  bugfix   - Bug fix branch (from development)")
    fmt.Println("  hotfix   - Hot fix branch (from production)")
    fmt.Println("  release  - Release branch (from production)")
    fmt.Println("\nCommit types:")
    fmt.Println("  feat     - New feature")
    fmt.Println("  fix      - Bug fix")
    fmt.Println("  chore    - Maintenance")
    fmt.Println("  docs     - Documentation")
    fmt.Println("  style    - Code style")
    fmt.Println("  refactor - Code refactoring")
    fmt.Println("  test     - Testing")
}

func printIssueDetails(issue *jira.JiraIssue) {
    fmt.Println("Issue Details:")
    fmt.Printf("Key: %s\n", issue.Key)
    fmt.Printf("Summary: %s\n", issue.Fields.Summary)
    fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
    fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.DisplayName)
    fmt.Printf("Description:\n%s\n", issue.Fields.Description)
}

func runSetup() error {
    // Implementation from previous setup wizard
    // Move the setup wizard implementation here
    return nil
}
