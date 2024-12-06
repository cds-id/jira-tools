package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"jira-tools/internal/config"
	"jira-tools/internal/git"
	"jira-tools/internal/jira"
	"github.com/joho/godotenv"
)

func loadConfig() error {
	configDir, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	envPath := filepath.Join(configDir, ".env")
	return godotenv.Load(envPath)
}

func promptUser(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

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
	projectRoot, err := git.GetProjectRoot()
	if err != nil {
		return fmt.Errorf("not a git repository: %v", err)
	}

	// Handle Jira credentials globally
	configDir, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %v", err)
	}

	// Only do Jira setup if credentials don't exist
	envPath := filepath.Join(configDir, ".env")
	var err2 error
	if _, err2 = os.Stat(envPath); os.IsNotExist(err2) {
		// Jira Configuration
		fmt.Println("\n=== Jira Configuration ===")
		domain := promptUser("Jira Domain (e.g., company.atlassian.net): ")
		email := promptUser("Jira Email: ")
		apiToken := promptUser("Jira API Token: ")

		fmt.Println("\nValidating Jira credentials...")
		if err2 = jira.ValidateCredentials(domain, email, apiToken); err2 != nil {
			return fmt.Errorf("credential validation failed: %v", err2)
		}

		// Save Jira credentials globally
		envContent := fmt.Sprintf("JIRA_DOMAIN=%s\nJIRA_EMAIL=%s\nJIRA_API_TOKEN=%s\n",
			domain, email, apiToken)
		if err2 = os.WriteFile(envPath, []byte(envContent), 0600); err2 != nil {
			return fmt.Errorf("failed to save credentials: %v", err2)
		}
	}

	// Project-specific git configuration
	branchConfig := &config.BranchConfig{
		ProjectPath: projectRoot,
	}

	// Rest of the setup remains the same, but save to project-specific location
	if err := config.SaveProjectBranchConfig(branchConfig); err != nil {
		return fmt.Errorf("failed to save project configuration: %v", err)
	}

	fmt.Printf("\nProject configuration saved in: %s\n", filepath.Join(projectRoot, ".jt-config.json"))

	fmt.Println("Welcome to Jira Tools (jt) Setup!")
	fmt.Println("=================================")

	// 1. Check git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository. Please run this command in a git repository")
	}

	// 2. Jira Configuration
	fmt.Println("\n=== Jira Configuration ===")
	domain := promptUser("Jira Domain (e.g., company.atlassian.net): ")
	email := promptUser("Jira Email: ")
	apiToken := promptUser("Jira API Token: ")

	fmt.Println("\nValidating Jira credentials...")
	if err := jira.ValidateCredentials(domain, email, apiToken); err != nil {
		return fmt.Errorf("credential validation failed: %v", err)
	}
	fmt.Println("✓ Credentials validated successfully")

	// 3. Git Branch Configuration
	fmt.Println("\n=== Git Branch Configuration ===")

	// Get available branches
	branches, err := git.GetAvailableBranches()
	if err != nil {
		return fmt.Errorf("failed to get branches: %v", err)
	}

	if len(branches) == 0 {
		return fmt.Errorf("no branches found in repository")
	}

	fmt.Println("\nAvailable branches:")
	for i, branch := range branches {
		fmt.Printf("%d. %s\n", i+1, branch)
	}

	// 4. Repository Type Configuration
	branchConfig = &config.BranchConfig{}

	fmt.Println("\nRepository Setup Options:")
	fmt.Println("1. Single Branch (development only)")
	fmt.Println("2. Git Flow (production/development branches)")

	repoType := promptUser("Select option (1/2): ")
	branchConfig.IsMonorepo = repoType == "2"

	if branchConfig.IsMonorepo {
		// Git Flow setup
		fmt.Println("\n=== Git Flow Configuration ===")

		// Production branch
		for {
			prodInput := promptUser("Production branch (enter number or name) [main/master]: ")
			if prodInput == "" {
				// Try main or master as default
				if contains(branches, "main") {
					branchConfig.ProductionBranch = "main"
					break
				} else if contains(branches, "master") {
					branchConfig.ProductionBranch = "master"
					break
				}
			} else if branch := getBranchFromInput(prodInput, branches); branch != "" {
				branchConfig.ProductionBranch = branch
				break
			}
			fmt.Println("Invalid branch. Please try again.")
		}

		// Development branch
		for {
			devInput := promptUser("Development branch (enter number or name) [develop]: ")
			if devInput == "" && contains(branches, "develop") {
				branchConfig.DevelopmentBranch = "develop"
				break
			} else if branch := getBranchFromInput(devInput, branches); branch != "" {
				if branch != branchConfig.ProductionBranch {
					branchConfig.DevelopmentBranch = branch
					break
				}
				fmt.Println("Development branch must be different from production branch.")
			} else {
				fmt.Println("Invalid branch. Please try again.")
			}
		}

		// Ask about creating develop branch if it doesn't exist
		if !contains(branches, branchConfig.DevelopmentBranch) {
			createDev := promptUser(fmt.Sprintf("Development branch '%s' doesn't exist. Create it? (Y/n): ", branchConfig.DevelopmentBranch))
			if createDev == "" || strings.ToLower(createDev) == "y" {
				if err := git.CreateBranch(branchConfig.ProductionBranch, "create development branch", git.BranchType("develop")); err != nil {
					return fmt.Errorf("failed to create development branch: %v", err)
				}
				fmt.Printf("Created development branch '%s' from '%s'\n", branchConfig.DevelopmentBranch, branchConfig.ProductionBranch)
			}
		}
	} else {
		// Single branch setup
		fmt.Println("\n=== Single Branch Configuration ===")
		for {
			devInput := promptUser("Development branch (enter number or name): ")
			if branch := getBranchFromInput(devInput, branches); branch != "" {
				branchConfig.DevelopmentBranch = branch
				break
			}
			fmt.Println("Invalid branch. Please try again.")
		}
	}

	// 5. Save configurations
	configDir, err = config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %v", err)
	}

	// Save Jira credentials
	envPath = filepath.Join(configDir, ".env")
	envContent := fmt.Sprintf("JIRA_DOMAIN=%s\nJIRA_EMAIL=%s\nJIRA_API_TOKEN=%s\n",
		domain, email, apiToken)

	if err = os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		return fmt.Errorf("failed to save credentials: %v", err)
	}

	// Save branch configuration
	if err = config.SaveProjectBranchConfig(branchConfig); err != nil {
		return fmt.Errorf("failed to save branch configuration: %v", err)
	}

	// 6. Print configuration summary
	fmt.Println("\nConfiguration Summary")
	fmt.Println("=====================")
	fmt.Printf("Jira Domain: %s\n", domain)
	fmt.Printf("Jira Email: %s\n", email)
	fmt.Println("\nGit Configuration:")
	if branchConfig.IsMonorepo {
		fmt.Printf("Repository Type: Git Flow\n")
		fmt.Printf("Production Branch: %s\n", branchConfig.ProductionBranch)
		fmt.Printf("Development Branch: %s\n", branchConfig.DevelopmentBranch)
	} else {
		fmt.Printf("Repository Type: Single Branch\n")
		fmt.Printf("Development Branch: %s\n", branchConfig.DevelopmentBranch)
	}

	// 7. Print next steps
	fmt.Printf("\nConfiguration saved in: %s\n", configDir)

	fmt.Println("\nNext Steps")
	fmt.Println("==========")
	if branchConfig.IsMonorepo {
		fmt.Println("For features:")
		fmt.Println("1. Create a feature branch:")
		fmt.Printf("   jt branch PROJ-123 feature\n")
		fmt.Println("\nFor bug fixes:")
		fmt.Println("1. Create a bugfix branch:")
		fmt.Printf("   jt branch PROJ-123 bugfix\n")
		fmt.Println("\nFor hotfixes:")
		fmt.Println("1. Create a hotfix branch:")
		fmt.Printf("   jt branch PROJ-123 hotfix\n")
	} else {
		fmt.Println("1. Create a feature branch:")
		fmt.Printf("   jt branch PROJ-123 feature\n")
	}

	fmt.Println("\nCommon commands:")
	fmt.Println("2. Look up issue details:")
	fmt.Println("   jt lookup PROJ-123")
	fmt.Println("3. Create a commit:")
	fmt.Println("   jt commit PROJ-123 feat")
	fmt.Println("4. Push changes:")
	fmt.Println("   jt push")

	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	if err := appendToGitignore(gitignorePath, ".jt-config.json"); err != nil {
		fmt.Printf("Warning: Could not add .jt-config.json to .gitignore: %v\n", err)
	}

	return nil
}

func appendToGitignore(gitignorePath, entry string) error {
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check if entry already exists
	if strings.Contains(string(content), entry) {
		return nil
	}

	// Append entry
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}

	_, err = f.WriteString(entry + "\n")
	return err
}

// Helper functions

func getBranchFromInput(input string, branches []string) string {
	// Try as index
	if index, err := strconv.Atoi(input); err == nil {
		if index > 0 && index <= len(branches) {
			return branches[index-1]
		}
		return ""
	}

	// Try as branch name
	input = strings.TrimSpace(input)
	for _, branch := range branches {
		if branch == input {
			return branch
		}
	}
	return ""
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
