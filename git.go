package main

import (
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	StateDir = "state"
	GitRepo  = "https://git.kernel.org/pub/scm/infra/transparency-logs/gitolite/git/1.git"
	GitDir   = path.Join(StateDir, path.Base(GitRepo))
)

type TlogRepo struct {
	Git *git.Repository
}

func GetRepo() (*TlogRepo, error) {
	if _, err := os.Stat(GitDir); err == nil {
		repo, err := git.PlainOpen(GitDir)
		if err != nil {
			return nil, err
		}
		return &TlogRepo{repo}, nil
	}
	repo, err := git.PlainClone(GitDir, false, &git.CloneOptions{
		URL:      GitRepo,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}
	return &TlogRepo{repo}, nil
}

func (tr *TlogRepo) Pull() error {
	w, err := tr.Git.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}
	return nil
}

func (tr *TlogRepo) LogSince(since *time.Time, fn func(r io.Reader) error) error {
	cIter, err := tr.Git.Log(&git.LogOptions{Since: since})
	if err != nil {
		return err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		f, err := c.File("m")
		if err != nil {
			return err
		}
		blob, err := f.Contents()
		if err != nil {
			return err
		}
		r := strings.NewReader(blob)
		if err := fn(r); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
