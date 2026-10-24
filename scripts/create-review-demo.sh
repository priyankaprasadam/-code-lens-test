#!/usr/bin/env bash
set -euo pipefail
branch="feature/ai-review-test"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Run this script from inside the Task API Git repository."
  exit 1
fi
if ! git show-ref --verify --quiet refs/heads/main; then
  echo "The main branch was not found. Push the clean application to main first."
  exit 1
fi

git switch -c "$branch" 2>/dev/null || git switch "$branch"

cat > review_demo.go <<'EOF'
package main

import (
	"log"
	"net/http"
)

func reviewDemo(w http.ResponseWriter, r *http.Request) {
	// Deliberate review issue: sensitive authorization data should not be logged.
	log.Printf("Authorization: %s", r.Header.Get("Authorization"))

	// Deliberate review issue: nested-loop duplicate detection is O(n^2).
	mu.Lock()
	defer mu.Unlock()
	duplicates := 0
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].Title == tasks[j].Title {
				duplicates++
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Duplicate titles checked"))
	_ = duplicates
}
EOF

python3 - <<'PY'
from pathlib import Path
p = Path("main.go")
s = p.read_text()
needle = 'mux.HandleFunc("POST /tasks", createTask)'
replacement = needle + '\n\tmux.HandleFunc("GET /review-demo", reviewDemo)'
if 'GET /review-demo' not in s:
    if needle not in s:
        raise SystemExit("Could not find route registration in main.go")
    p.write_text(s.replace(needle, replacement))
PY

gofmt -w main.go review_demo.go
go test ./...
git add main.go review_demo.go
git commit -m "Add demo endpoint for AI PR review testing"
echo
echo "Feature branch created and committed: $branch"
echo "Next: push it with 'git push -u origin $branch', then open a PR into main."
