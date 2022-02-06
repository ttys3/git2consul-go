/*
Copyright 2019 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kv

import (
	"fmt"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// HandleUpdate handles the update of a particular repository.
func (h *KVHandler) HandleUpdate(repo repository.Repo) error {
	w, err := repo.Worktree()
	config := repo.GetConfig()
	repo.Lock()
	defer repo.Unlock()

	if err != nil {
		return err
	}

	for _, branch := range config.Branches {
		ref := fmt.Sprintf("refs/heads/%s", branch)
		err := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(ref),
			Force:  true,
		})
		if err != nil {
			return fmt.Errorf("checkout %s failed: %w", ref, err)
		}
		err = h.UpdateToHead(repo)
		if err != nil {
			return fmt.Errorf("updateToHead %s failed: %w", repo.Name(), err)
		}
	}
	return nil
}

// UpdateToHead handles update to current HEAD comparing diffs against the KV.
func (h *KVHandler) UpdateToHead(repo repository.Repo) error {
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("get repo head failed, err=%w", err)
	}
	refName := head.Name().Short()
	if err != nil {
		return fmt.Errorf("get repo short ref name failed, err=%w", err)
	}

	h.logger.Infof("KV GET ref: %s/%s", repo.Name(), refName)
	kvRef, err := h.getKVRef(repo, refName)
	if err != nil {
		return fmt.Errorf("getKVRef failed, refName=%v err=%w", refName, err)
	}

	// Local ref
	headRefHash := head.Hash().String()
	// log.Debugf("(consul) kvRef: %s | localRef: %s", kvRef, localRef)

	if kvRef == "" {
		log.Infof("init KV PUT branch: %s/%s", repo.Name(), refName)
		err := h.putBranch(repo, plumbing.ReferenceName(head.Name().Short()))
		if err != nil {
			return err
		}

		err = h.putKVRef(repo, refName)
		if err != nil {
			return err
		}
		h.logger.Infof("init KV PUT ref: %s/%s", repo.Name(), refName)
	} else if kvRef != headRefHash {
		// Check if the ref belongs to that repo
		err := repo.CheckRef(refName)
		if err != nil {
			return err
		}

		// Handle modified and deleted files
		deltas, err := repo.DiffStatus(kvRef)
		if err != nil {
			return err
		}
		err = h.handleDeltas(repo, deltas)
		if err != nil {
			h.logger.Errorf("handleDeltas error: %v, repo=%v", err, repo)
			// TODO should we return err here?
		}

		err = h.putKVRef(repo, refName)
		if err != nil {
			return err
		}
		h.logger.Infof("KV PUT ref change: %s/%s", repo.Name(), refName)
	} else {
		h.logger.Infof("KV ref is update to date: %s/%s", repo.Name(), refName)
	}

	return nil
}
