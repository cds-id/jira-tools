# Jira Tools - Git Commit Automation

A CLI tool that integrates Jira with Git commits, automatically fetching issue summaries and creating semantic commit messages with Git Flow support.

## Features

- Fetch Jira issue details directly from the command line
- Automatically create semantic commit messages with Jira issue information
- Lookup detailed Jira issue information
- Supports custom commit types (feat, fix, chore, etc.)
- Git Flow branching strategy support
- Wizard-based configuration setup

## Prerequisites

- Go 1.16 or higher
- Git installed and configured
- Jira account with API access
- Jira API token ([How to generate](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/))

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/cds-id/jira-tools.git
cd jira-tools
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the binary:
```bash
go build -o jt ./cmd/jt
```

4. (Optional) Move to PATH for system-wide access:
```bash
sudo mv jt /usr/local/bin/
```

### Configuration

Run the setup wizard:
```bash
jt setup
```

The wizard will guide you through:
1. Jira Configuration
   - Domain
   - Email
   - API Token

2. Git Branch Configuration
   - Single Branch (Development only)
     - Select your development branch
   - Git Flow Setup (Production/Development)
     - Select production branch (main/master)
     - Select development branch (develop)

## Usage

### Look up Jira Issue Details

```bash
jt lookup PROJ-123
```

This will display:
- Issue Key
- Summary
- Status
- Assignee
- Description

### Branch Management

Create a new branch based on Jira issue:
```bash
jt branch PROJ-123 feature    # Creates feature/PROJ-123-issue-summary
jt branch PROJ-123 bugfix     # Creates bugfix/PROJ-123-issue-summary
jt branch PROJ-123 hotfix     # Creates hotfix/PROJ-123-issue-summary
```

### Create a Commit

```bash
jt commit PROJ-123 feat       # feat(PROJ-123): Issue summary
```

### Push Changes

Push current branch to remote:
```bash
jt push
```

### Branch Types

- `feature` - New feature branch (from development)
- `bugfix` - Bug fix branch (from development)
- `hotfix` - Hot fix branch (from production)
- `release` - Release branch (from production)

### Commit Types

- `feat`: New feature
- `fix`: Bug fix
- `chore`: Maintenance tasks
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or modifying tests

## Example Output

```bash
# Looking up issue
$ jt lookup PROJ-123
Issue Details:
Key: PROJ-123
Summary: Implement user authentication
Status: In Progress
Assignee: John Doe
Description:
Add user authentication using OAuth2...

# Creating a branch
$ jt branch PROJ-123 feature
Created branch feature/PROJ-123-implement-user-authentication from develop

# Creating a commit
$ jt commit PROJ-123 feat
Changes committed with message:
feat(PROJ-123): Implement user authentication

# Pushing changes
$ jt push
Successfully pushed branch feature/PROJ-123-implement-user-authentication to remote
```

## Development

### Project Structure
```
jira-tools/
├── cmd/
│   └── jt/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── git/
│   │   └── git.go
│   └── jira/
│       └── jira.go
├── .env
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

### Building from Source

1. Clone the repository
```bash
git clone https://github.com/cds-id/jira-tools.git
```

2. Navigate to the project directory
```bash
cd jira-tools
```

3. Install dependencies
```bash
go mod tidy
```

4. Build the binary
```bash
go build -o jt ./cmd/jt
```

### Running Tests

```bash
go test -v ./...
```

## GitHub Actions

The project includes a GitHub Actions workflow for:
- Building the application for multiple platforms (Linux, Windows, macOS)
- Creating releases when tags are pushed
- Automated builds and tests

To create a release:
1. Tag your version
```bash
git tag v1.0.0
git push origin v1.0.0
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Security

- Never commit your `.env` file
- Keep your Jira API token secure
- Consider using GitHub Secrets for CI/CD
- Credentials are stored securely in user's home directory

## Troubleshooting

### Common Issues

1. **Authentication Error**
   - Run `jt setup` to reconfigure credentials
   - Check if your API token is valid
   - Verify credentials in ~/.jira-tools/.env

2. **Command Not Found**
   - Ensure the binary is in your PATH
   - Verify the binary has execute permissions
   - Run `which jt` to locate the binary

3. **Branch Creation Issues**
   - Ensure you're in a git repository
   - Check if the base branch exists
   - Verify your git flow configuration

For more issues, please check the GitHub Issues section.
