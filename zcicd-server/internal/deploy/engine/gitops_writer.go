package engine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	gitopsLockPrefix = "gitops:lock:"
	gitopsLockTTL    = 60 * time.Second
)

// GitOpsWriter handles GitOps repository updates with Redis distributed locking.
type GitOpsWriter struct {
	redisClient *redis.Client
}

// NewGitOpsWriter creates a new GitOpsWriter.
func NewGitOpsWriter(redisClient *redis.Client) *GitOpsWriter {
	return &GitOpsWriter{
		redisClient: redisClient,
	}
}

// UpdateValues acquires a distributed lock, updates values in the GitOps repo,
// commits and pushes. Git operations are currently stubbed.
func (w *GitOpsWriter) UpdateValues(
	ctx context.Context,
	repoURL, branch, filePath string,
	values map[string]interface{},
) (string, error) {
	lockKey := gitopsLockPrefix + repoURL

	// Acquire distributed lock
	acquired, err := w.redisClient.SetNX(ctx, lockKey, "locked", gitopsLockTTL).Result()
	if err != nil {
		return "", fmt.Errorf("failed to acquire gitops lock for %s: %w", repoURL, err)
	}
	if !acquired {
		return "", fmt.Errorf("gitops lock already held for %s, retry later", repoURL)
	}

	// Ensure lock is released when done
	defer func() {
		if err := w.redisClient.Del(ctx, lockKey).Err(); err != nil {
			log.Printf("warning: failed to release gitops lock for %s: %v", repoURL, err)
		}
	}()

	// --- Stubbed git operations ---
	// TODO: Wire git credentials and implement real clone/pull/commit/push.
	log.Printf("gitops: would clone/pull repo=%s branch=%s", repoURL, branch)
	log.Printf("gitops: would update file=%s with %d value overrides", filePath, len(values))
	log.Printf("gitops: would commit and push changes")

	stubCommitSHA := "stub-sha-pending-git-integration"
	return stubCommitSHA, nil
}