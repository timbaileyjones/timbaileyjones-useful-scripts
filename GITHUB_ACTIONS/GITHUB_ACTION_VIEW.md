# GitHub Actions Run Comparison Guide

This guide explains how to find and compare GitHub Actions (GHA) runs between different commits, including how to identify which commits actually triggered runs.

## Understanding GitHub Actions Triggers

Not every commit triggers a GitHub Actions run. The CI workflow in this repository (`ci.yml`) only runs on:
- Pushes to `main`, `alpha`, and `beta` branches
- Manual workflow dispatch

**Important**: Release commits that only update `package.json` files and lock files typically do NOT trigger CI runs.

## Finding GitHub Actions Runs

### Method 1: Using GitHub CLI (Recommended)

#### 1. List Recent Runs
```bash
gh run list --limit 20
```

#### 2. Find Runs for Specific Commits
```bash
# Get commit hash
git rev-parse HEAD
git rev-parse HEAD~1

# Find runs for specific commits
gh run list --commit <commit-hash>
```

#### 3. View Run Details
```bash
# View run summary
gh run view <run-id>

# View detailed logs
gh run view <run-id> --log

# View specific job details
gh run view --job=<job-id>
```

### Method 2: Using GitHub Web Interface

#### 1. Direct Actions Page
- Go to: `https://github.com/TanStack/router/actions`
- Look for runs with commit hashes in the titles

#### 2. Commit-Specific Search
- Use GitHub's search: `https://github.com/TanStack/router/actions?query=sha%3A<commit-hash>`

#### 3. Direct Run URLs
- Format: `https://github.com/TanStack/router/actions/runs/<run-id>`

## Finding the Previous Commit with GHA Runs

Since not every commit triggers a run, you need to find the most recent commit that actually had a GitHub Actions run:

### Step-by-Step Process

1. **Get current commit info**:
   ```bash
   git log --oneline -5
   git rev-parse HEAD
   ```

2. **Check if current commit has runs**:
   ```bash
   gh run list --commit <current-commit-hash>
   ```

3. **If no runs found, check previous commits**:
   ```bash
   # Check HEAD~1
   git rev-parse HEAD~1
   gh run list --commit <HEAD~1-hash>
   
   # If still no runs, check HEAD~2, HEAD~3, etc.
   git rev-parse HEAD~2
   gh run list --commit <HEAD~2-hash>
   ```

4. **Continue until you find a commit with runs**

### Example: Finding Previous GHA Run

```bash
# Current commit (HEAD)
$ git rev-parse HEAD
3c6633a7059197882911876286345b275a3bf64c

# Check if current commit has runs
$ gh run list --commit 3c6633a7059197882911876286345b275a3bf64c
completed	success	chore: silence error in test (#4794)	ci	main	push	16543871244	2m32s	2025-07-26T21:05:20Z
completed	success	chore: silence error in test (#4794)	autofix.ci	main	push	16543871241	56s	2025-07-26T21:05:20Z

# Check previous commit (HEAD~1)
$ git rev-parse HEAD~1
77fa17bdc01d66977c2ddc52eed93af7d55295dd

$ gh run list --commit 77fa17bdc01d66977c2ddc52eed93af7d55295dd
# No output - this commit didn't trigger runs (it's a release commit)

# Check HEAD~2
$ git rev-parse HEAD~2
e83199310ad3cdfdaef26ac5567fa5e4b04797af

$ gh run list --commit e83199310ad3cdfdaef26ac5567fa5e4b04797af
completed	success	feat(react-router): add disableGlobalCatchBoundary option (#4765)	autofix.ci	main	push	16543721322	59s	2025-07-26T20:45:24Z
completed	success	feat(react-router): add disableGlobalCatchBoundary option (#4765)	ci	main	push	16543721323	9m59s	2025-07-26T20:45:24Z
```

## Comparing Runs

### Key Information to Compare

1. **Run Status**: Success ✅, Failure ❌, or Cancelled ⏹️
2. **Duration**: How long each run took
3. **Job Results**: Individual job statuses and durations
4. **Log Output**: Test results, build outputs, error messages

### Example Comparison

**Current Commit** (`3c6633a70`):
- Run ID: `16543871244`
- Status: ✅ Success
- Duration: 2m28s
- URL: https://github.com/TanStack/router/actions/runs/16543871244

**Previous Commit** (`e83199310`):
- Run ID: `16543721323`
- Status: ✅ Success
- Duration: 9m41s
- URL: https://github.com/TanStack/router/actions/runs/16543721323

**Key Observation**: The current commit's run was **4x faster** (2m28s vs 9m41s), suggesting performance improvements or test optimizations.

## Common Scenarios

### Scenario 1: Release Commits
Release commits (like `release: v1.130.0`) typically only update version numbers in `package.json` files and don't trigger CI runs. Look for the commit before the release.

### Scenario 2: Multiple Runs Per Commit
Some commits may trigger multiple workflows:
- `ci` - Main CI workflow
- `autofix.ci` - Auto-fix workflow
- `pr` - Pull request workflow

### Scenario 3: Failed Runs
If a commit has failed runs, you might want to compare with the last successful run before that commit.

## Troubleshooting

### No Runs Found
- Check if the commit was pushed to `main`, `alpha`, or `beta` branch
- Verify the commit wasn't a release-only commit
- Look further back in the commit history

### Permission Issues
- Ensure you have access to the repository
- Use `gh auth login` if GitHub CLI authentication fails

### Finding Specific Information
- Use `gh run view <run-id> --log` for detailed logs
- Use `gh run view --job=<job-id>` for specific job details
- Filter runs by workflow: `gh run list --workflow=ci`

## Quick Reference Commands

```bash
# Get current commit hash
git rev-parse HEAD

# Get previous commit hash
git rev-parse HEAD~1

# List recent runs
gh run list --limit 10

# Find runs for specific commit
gh run list --commit <hash>

# View run details
gh run view <run-id>

# View run logs
gh run view <run-id> --log

# View job details
gh run view --job=<job-id>
```

## Workflow Configuration

The CI workflow (`ci.yml`) triggers on:
- `push` to branches: `[main, alpha, beta]`
- `workflow_dispatch` (manual trigger)

Jobs include:
- **Build & Test**: Runs tests and builds packages
- **Publish**: Publishes packages to npm (on success)

This means only commits pushed to the main branches will trigger CI runs, and release commits that only update package versions typically won't trigger runs.