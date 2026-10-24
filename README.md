# Task API — AI PR Review POC

A small Go REST API used as a target repository for testing an AI-powered pull-request review framework.

## Requirements
- Go 1.22 or newer
- Git
- A GitHub repository where you can push branches and open pull requests

## Run locally
```bash
go test ./...
go run .
```

In another terminal:
```bash
curl http://localhost:8080/tasks
curl -X POST http://localhost:8080/tasks -H 'Content-Type: application/json' -d '{"title":"Test AI PR review"}'
```

## Push the clean application
Create an empty repository on GitHub, then run from this directory:
```bash
git init
git add .
git commit -m "Initial Task API application"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO.git
git push -u origin main
```

## Create a sample review pull request
After pushing `main`, run:
```bash
chmod +x scripts/create-review-demo.sh
./scripts/create-review-demo.sh
git push -u origin feature/ai-review-test
```
Then open a PR from `feature/ai-review-test` into `main`.

The demo intentionally logs an Authorization header and uses a nested loop for duplicate detection. Do not deploy that branch to production or use real credentials against its demo endpoint.

## Connect the review framework
In the separate `ai-pr-review` framework directory:
```bash
export GITHUB_REPOSITORY="YOUR_USERNAME/YOUR_REPO"
export PR_NUMBER="1"
```
Use the actual PR number. Ensure `OPENAI_API_KEY`, `OPENAI_MODEL`, and `GITHUB_TOKEN` are set in the same terminal session, then run the review CLI from the framework project.

## Endpoints
- `GET /tasks` — returns the in-memory task list.
- `POST /tasks` — creates a task with a non-empty `title`.

This is a demonstration app; data is stored in memory and resets on restart.
