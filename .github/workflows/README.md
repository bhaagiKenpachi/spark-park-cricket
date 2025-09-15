# GitHub Workflows

[![Backend CI/CD](https://github.com/luffybhaagi/spark-park-cricket/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/luffybhaagi/spark-park-cricket/actions/workflows/backend-ci.yml)
[![Frontend CI/CD](https://github.com/luffybhaagi/spark-park-cricket/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/luffybhaagi/spark-park-cricket/actions/workflows/frontend-ci.yml)
[![Self-Hosted Runner](https://img.shields.io/badge/Runner-Self--Hosted-orange.svg)](https://github.com/bhaagiKenpachi/spark-park-cricket/settings/actions/runners)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![Node.js Version](https://img.shields.io/badge/Node.js-20-green.svg)](https://nodejs.org/)

This directory contains GitHub Actions workflows for the Spark Park Cricket project.

## Workflows

### Backend CI/CD (`backend-ci.yml`)

Comprehensive CI/CD pipeline for the Go backend service using self-hosted runners.

#### Triggers
- **Push**: `main`, `develop`, `feat/*`, `fix/*` branches (only when backend files change)
- **Pull Request**: `main`, `develop` branches (only when backend files change)
- **Manual Trigger**: On-demand execution with customizable options

#### Jobs

1. **Code Quality** - Runs first, includes:
   - Code formatting check (`gofmt`)
   - Static analysis (`go vet`)
   - Linting (`golangci-lint`)
   - Security analysis (`staticcheck`)

2. **Unit Tests** - Runs in parallel with code quality:
   - Unit tests (no database required)
   - Race condition detection
   - Code coverage reporting
   - Coverage upload to Codecov

3. **Integration Tests** - Runs after unit tests pass:
   - Integration tests (requires Supabase database)
   - E2E tests (requires Supabase database)
   - Database migration setup
   - Test database schema setup

4. **Build Verification** - Runs after unit tests pass:
   - Application build verification
   - Binary testing
   - Build artifact upload

5. **Security Scan** - Runs after code quality:
   - Gosec security scanner
   - SARIF report upload

6. **Deploy** - Runs on manual trigger only:
   - Builds for selected environment
   - Deploys to chosen environment (development, staging, production, testing)
   - Environment-specific build optimizations
   - Uploads deployment artifacts

7. **Notify on Failure** - Runs if any job fails:
   - Failure notification (placeholder)

### Frontend CI/CD (`frontend-ci.yml`)

Comprehensive CI/CD pipeline for the Next.js frontend application using self-hosted runners.

#### Triggers
- **Push**: `main`, `develop`, `feat/*`, `fix/*` branches (only when frontend files change)
- **Pull Request**: `main`, `develop` branches (only when frontend files change)
- **Manual Trigger**: On-demand execution with customizable options

#### Jobs

1. **Code Quality** - Runs first, includes:
   - Code formatting check (Prettier)
   - Linting (ESLint)
   - Type checking (TypeScript)

2. **Unit Tests** - Runs in parallel with code quality:
   - Jest unit tests
   - Code coverage reporting
   - Coverage upload to Codecov

3. **E2E Tests** - Runs after unit tests pass:
   - Cypress end-to-end tests
   - Application build and start
   - E2E test results upload

4. **Build Verification** - Runs after unit tests pass:
   - Next.js build verification
   - Build artifact upload

5. **Security Scan** - Runs after code quality:
   - npm audit for dependency vulnerabilities
   - Snyk security scanning (optional)

6. **Deploy to Staging** - Runs on `develop` branch push:
   - Builds for staging environment
   - Deploys to staging (placeholder)

7. **Deploy to Production** - Runs on `main` branch push:
   - Builds for production environment
   - Deploys to production (placeholder)

8. **Notify on Failure** - Runs if any job fails:
   - Failure notification (placeholder)

#### Required Secrets

The following secrets can be configured in your GitHub repository (optional):

```
SNYK_TOKEN=your_snyk_token_for_security_scanning
```

#### Environment Variables

The workflow uses these environment variables:
- `NODE_VERSION`: Node.js version to use (currently 20)
- `WORKING_DIRECTORY`: Frontend directory path (`./web`)

## Backend Required Secrets

The following secrets must be configured in your GitHub repository for the backend workflow:

```
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_API_KEY=your_supabase_anon_key
```

#### Environment Variables

The workflow uses these environment variables:
- `GO_VERSION`: Go version to use (currently 1.23)
- `WORKING_DIRECTORY`: Backend directory path (`./backend`)

## Manual Trigger Usage

Both workflows support manual triggering with customizable options:

### Backend Manual Trigger Options
- **Environment**: Choose deployment target (`development`, `staging`, `production`, `testing`, `none`)
- **Skip Integration Tests**: Skip database-dependent tests
- **Skip Security Scan**: Skip security scanning

### Frontend Manual Trigger Options
- **Environment**: Choose deployment target (`staging`, `production`, `none`)
- **Skip E2E Tests**: Skip end-to-end tests
- **Skip Security Scan**: Skip security scanning

### How to Trigger Manually
1. Go to your GitHub repository
2. Navigate to Actions tab
3. Select the workflow you want to run
4. Click "Run workflow" button
5. Choose your options and click "Run workflow"

## Setup Instructions

### 1. Configure Repository Secrets

1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Add the following repository secrets:

**Backend Secrets (Required):**
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_API_KEY`: Your Supabase anonymous key

**Frontend Secrets (Optional):**
   - `SNYK_TOKEN`: Your Snyk token for security scanning

**Self-Hosted Runner Secrets (Required):**
   - `GITHUB_PAT`: Personal Access Token for self-hosted runner authentication

### 2. Self-Hosted Runner Setup

The workflows are configured to use self-hosted runners. To set up a self-hosted runner:

1. **Add the PAT as a repository secret**:
   - Go to Settings → Secrets and variables → Actions
   - Add `GITHUB_PAT` with your Personal Access Token: `your_github_personal_access_token_here`

2. **Set up the self-hosted runner**:
   - Go to Settings → Actions → Runners
   - Click "New self-hosted runner"
   - Follow the setup instructions for your operating system
   - Configure the runner with appropriate labels (e.g., `self-hosted`)

3. **Runner Requirements**:
   - **Operating System**: Linux, macOS, or Windows
   - **Go**: Version 1.23 or later
   - **Node.js**: Version 20 or later
   - **Memory**: At least 4GB RAM
   - **Storage**: At least 10GB free space
   - **Network**: Internet access for downloading dependencies

### 3. Test Database Setup

For integration tests to work properly, ensure your Supabase database has:

1. A `testing_db` schema created
2. The test schema setup script run:
   ```sql
   -- Run the setup script from backend/scripts/setup_test_db.sql
   ```

### 4. Branch Protection Rules (Recommended)

Set up branch protection rules for `main` and `develop` branches:

1. Go to Settings → Branches
2. Add rule for `main`:
   - Require status checks: All CI jobs must pass
   - Require branches to be up to date
   - Restrict pushes to matching branches
3. Add rule for `develop`:
   - Require status checks: Code quality, unit tests, build
   - Allow force pushes (for development flexibility)

## Self-Hosted Runner Benefits

### Advantages
- **Cost Control**: No GitHub Actions minutes usage
- **Custom Environment**: Full control over runner configuration
- **Faster Builds**: Local resources and caching
- **Network Access**: Direct access to internal services
- **Custom Tools**: Install any required software

### Considerations
- **Maintenance**: You're responsible for runner updates and security
- **Availability**: Runner must be online for workflows to execute
- **Security**: Ensure runner is properly secured and isolated
- **Resource Management**: Monitor CPU, memory, and storage usage

## Workflow Features

### Caching
- Go modules are cached for faster builds
- Go build cache is preserved between runs

### Parallel Execution
- Code quality and unit tests run in parallel
- Integration tests run after unit tests pass
- Build verification runs in parallel with integration tests

### Conditional Deployment
- Staging deployment only on `develop` branch
- Production deployment only on `main` branch
- Integration tests only run when Supabase secrets are available

### Error Handling
- Each job has proper error handling
- Failure notifications (extensible)
- Graceful degradation when secrets are missing

## Customization

### Adding New Test Types

To add new test types:

1. Add the test files to the appropriate directory
2. Update the workflow to include the new test job
3. Add proper dependencies and environment setup

### Deployment Configuration

The deployment steps are currently placeholders. To implement actual deployment:

1. Replace the placeholder deployment steps with your deployment logic
2. Add necessary secrets for deployment (e.g., server credentials, Docker registry)
3. Configure environment-specific settings

### Notification Setup

To implement failure notifications:

1. Add notification service secrets (e.g., Slack webhook, email credentials)
2. Replace the placeholder notification step with actual notification logic
3. Configure notification channels and recipients

## Troubleshooting

### Common Issues

1. **Integration tests failing**: Check Supabase secrets and database connectivity
2. **Build failures**: Verify Go version compatibility and dependencies
3. **Security scan failures**: Review and fix security issues reported by Gosec
4. **Deployment failures**: Check deployment credentials and environment setup

### Debug Mode

To enable debug logging, add this to any step:
```yaml
env:
  ACTIONS_STEP_DEBUG: true
```

### Local Testing

To test the workflow locally:

1. Use `act` tool to run GitHub Actions locally
2. Set up local environment variables
3. Test individual workflow steps

## Performance Optimization

### Build Time Optimization
- Go modules caching reduces download time
- Parallel job execution reduces total pipeline time
- Conditional job execution based on file changes

### Resource Optimization
- Jobs run on `ubuntu-latest` (2-core, 7GB RAM)
- Timeout settings prevent hanging jobs
- Artifact retention limits storage usage

## Security Considerations

- Secrets are properly masked in logs
- SARIF reports are uploaded for security analysis
- Dependencies are verified before use
- Security scanning is integrated into the pipeline

## Monitoring and Metrics

- Code coverage reports are generated and uploaded
- Build artifacts are preserved for debugging
- Workflow status is visible in pull requests
- Failure notifications help with quick issue resolution
