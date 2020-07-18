package git

import (
	"testing"

	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"github.com/xfxdev/xlog"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var (
	signature = object.Signature{
		Name:  "Golang Tests",
		Email: "gotest@noreply.com",
	}
)

// createRepository creates an in-memory repository that can be used for testing
func createRepository(t *testing.T) *git.Repository {
	// Create the repo
	r, err := git.Init(memory.NewStorage(), memfs.New())
	if err != nil {
		t.Errorf("Unable to create repository for testing: %s", err.Error())
	}

	return r
}

// createBranch creates and checks out a new branch
func createBranch(t *testing.T, r *git.Repository, branch string) {
	err := r.CreateBranch(&config.Branch{
		Name: branch,
	})
	if err != nil {
		t.Errorf("Unable to create branch %s: %s", branch, err.Error())
	}
	xlog.Debugf("Created branch %s", branch)

	checkout(t, r, branch)
}

// checkout checksout the given branch name
func checkout(t *testing.T, r *git.Repository, branch string) {
	w, err := r.Worktree()
	if err != nil {
		t.Errorf("Unable to retrieve worktree: %s", err.Error())
	}

	w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
}

// createCommit creates a commit on the current branch
func createCommit(t *testing.T, r *git.Repository, message string) {
	w, err := r.Worktree()
	if err != nil {
		t.Errorf("Unable to retrieve worktree: %s", err.Error())
	}

	// Create a file with a random name
	s, err := randutil.AlphaString(32)
	if err != nil {
		t.Errorf("Unable to generate random string: %s", err.Error())
	}
	f, err := w.Filesystem.Create(s)
	if err != nil {
		t.Errorf("Unable to create file in worktree: %s", err.Error())
	}
	xlog.Debugf("Created file %s", f.Name())

	// Stage the file
	_, err = w.Add(f.Name())
	if err != nil {
		t.Errorf("Unable to stage file %s: %s", f.Name(), err.Error())
	}

	_, err = w.Commit(message, &git.CommitOptions{
		Author:    &signature,
		Committer: &signature,
	})
	if err != nil {
		t.Errorf("Unable to commit file %s: %s", f.Name(), err.Error())
	}
}

func createTag(t *testing.T, r *git.Repository, tag string) {
	head, err := r.Head()
	if err != nil {
		t.Error("Unable to get HEAD")
	}

	_, err = r.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Tagger:  &signature,
		Message: tag,
	})
	if err != nil {
		t.Errorf("Unable to create tag %s: %s", tag, err.Error())
	}
	xlog.Debugf("Created tag %s", tag)
}

func TestGetCommits(t *testing.T) {
	const (
		b1 = "b1"

		c1 = "first commit"
		c2 = "second commit"
		c3 = "third commit"

		v1 = "1.0.0"

		head  = "HEAD"
		head1 = "HEAD~1"
		head2 = "HEAD~2"
	)

	cases := map[string]struct {
		setup    func(*testing.T) *git.Repository
		from     string
		to       string
		expected []string
	}{
		"simple": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createCommit(t, r, c2)
				createCommit(t, r, c3)
				return r
			},
			from: head1,
			to:   head,
			expected: []string{
				c3,
			},
		},
		"multiple commits": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createCommit(t, r, c2)
				createCommit(t, r, c3)
				return r
			},
			from: head2,
			to:   head,
			expected: []string{
				c2,
				c3,
			},
		},
		"truncated branch": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createCommit(t, r, c2)
				createCommit(t, r, c3)
				return r
			},
			from: head2,
			to:   head1,
			expected: []string{
				c2,
			},
		},
		"from and to are equal": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				return r
			},
			from:     head,
			to:       head,
			expected: []string{},
		},
		"from tag": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createTag(t, r, v1)
				createCommit(t, r, c2)
				createCommit(t, r, c3)
				return r
			},
			from: v1,
			to:   head,
			expected: []string{
				c2,
				c3,
			},
		},
	}

	for name, data := range cases {
		r := data.setup(t)
		commits, err := getCommits(r, data.from, data.to)
		if assert.NoErrorf(t, err, name) {
			assert.Lenf(t, commits, len(data.expected), name)

			// Check all the elements of data.expected are present, but don't worry about the order
			for _, c := range data.expected {
				assert.Containsf(t, commits, c, name)
			}
		}
	}
}

func TestGetCommitsError(t *testing.T) {
	const (
		b1 = "b1"

		c1 = "first commit"
		c2 = "second commit"
		c3 = "third commit"

		head  = "HEAD"
		head1 = "HEAD~1"
		head2 = "HEAD~2"
	)

	cases := map[string]struct {
		setup func(*testing.T) *git.Repository
		from  string
		to    string
	}{
		"bad from": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				return r
			},
			from: "bad",
			to:   head,
		},
		"bad to": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				return r
			},
			from: head1,
			to:   "bad",
		},
		"to before from": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createCommit(t, r, c2)
				createCommit(t, r, c3)
				return r
			},
			from: head,
			to:   head2,
		},
	}

	for name, data := range cases {
		r := data.setup(t)
		_, err := getCommits(r, data.from, data.to)
		assert.Errorf(t, err, name)
	}
}

func TestGetLatestVersion(t *testing.T) {
	const (
		b1 = "b1"
		b2 = "b2"

		c1 = "first commit"
		c2 = "second commit"

		v0 = "0.1.0"
		v1 = "1.0.0"
		v2 = "2.0.0"
	)

	cases := map[string]struct {
		setup    func(*testing.T) *git.Repository
		expected string
	}{
		"simple": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createTag(t, r, v1)
				return r
			},
			expected: v1,
		},
		"multiple tags": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createTag(t, r, v1)
				createTag(t, r, v2)
				return r
			},
			expected: v2,
		},
		"multiple tags on multiple branches": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createTag(t, r, v1)
				createBranch(t, r, b2)
				createCommit(t, r, c2)
				createTag(t, r, v2)
				return r
			},
			expected: v2,
		},
		"contains invalid semver": {
			setup: func(t *testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createTag(t, r, v1)
				createTag(t, r, "Invalid Semver")
				return r
			},
			expected: v1,
		},
	}

	for name, data := range cases {
		r := data.setup(t)
		v, err := getLatestVersion(r)
		if assert.NoError(t, err) {
			assert.Equalf(t, data.expected, v, name)
		}
	}
}

func TestGetLatestVersionError(t *testing.T) {
	const (
		b1 = "b1"

		c1 = "first commit"
	)

	cases := map[string]struct {
		setup func(*testing.T) *git.Repository
	}{
		"no valid semvers": {
			setup: func(*testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				createTag(t, r, "Invalid Semver")
				return r
			},
		},
		"no tags": {
			setup: func(*testing.T) *git.Repository {
				r := createRepository(t)
				createBranch(t, r, b1)
				createCommit(t, r, c1)
				return r
			},
		},
	}

	for name, data := range cases {
		r := data.setup(t)
		_, err := getLatestVersion(r)
		assert.Errorf(t, err, name)
	}
}
