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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/apex/log"
)

// Load maps the configuration provided from a file to a Configuration object
func Load(file string) (*Config, error) {
	// log context
	logger := log.WithFields(log.Fields{
		"caller": "config",
	})

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Create Config object pointer and unmashal JSON into it
	config := &Config{
		Consul:  &ConsulConfig{},
		Webhook: &WebhookServerConfig{},
		Log:     &LogConfig{},
	}
	ext := filepath.Ext(file)
	switch ext {
	case ".json":
		err = json.Unmarshal(content, config)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(content, config)
	default:
		err = fmt.Errorf("invalid config file extension: %s", ext)
	}
	if err != nil {
		return nil, err
	}

	logger.Info("Setting configuration with sane defaults")
	config.setDefaultConfig()
	config.setDefaultLogConfig()

	err = config.checkConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Check for the validity of the configuration file
func (c *Config) checkConfig() error {
	for _, repo := range c.Repos {
		// Check on name
		if repo.Name == "" {
			return fmt.Errorf("Repository array object missing \"name\" value")
		}

		// Check on Url
		if repo.URL == "" {
			return fmt.Errorf("%s does no have a repository URL", repo.Name)
		}

		// Check on hooks
		for _, hook := range repo.Hooks {
			if hook.Type != "polling" && hook.Type != "webhook" {
				return fmt.Errorf("Invalid hook type: %s", hook.Type)
			}

			if hook.Type == "polling" && hook.Interval <= 0 {
				return fmt.Errorf("Invalid interval: %s. Hook interval must be greater than zero", hook.Interval)
			}
		}

		// Check on mount_point
		if repo.MountPoint != "" {
			if strings.HasPrefix(repo.MountPoint, "/") {
				return fmt.Errorf("Invalid mount point format for the %s repository - found \"/\" in the beginning of the path", repo.Name)
			}
			if !strings.HasSuffix(repo.MountPoint, "/") {
				return fmt.Errorf("Invalid mount point format for the %s repository - missing trailing \"/\"", repo.Name)
			}
		}

		// Check on source_root
		if repo.SourceRoot != "" {
			if !strings.HasPrefix(repo.SourceRoot, "/") {
				return fmt.Errorf("Invalid source_root format for the %s repository - missing \"/\" in the beginning of the path", repo.Name)
			}
			if !strings.HasSuffix(repo.SourceRoot, "/") {
				return fmt.Errorf("Invalid source_root format for the %s repository - missing trailing \"/\"", repo.Name)
			}
		}
	}

	return nil
}

// Return a configuration with sane defaults
func (c *Config) setDefaultConfig() {
	// Set the default cache store to be the OS' temp dir
	if c.LocalStore == "" {
		c.LocalStore = os.TempDir()
	}

	// Set the default webhook port
	if c.Webhook.Port == 0 {
		c.Webhook.Port = 9000
	}

	// For each repo, set default branch and hook
	for _, repo := range c.Repos {
		branch := []string{"main"}
		// If there are no branches, set it to master
		if len(repo.Branches) == 0 {
			repo.Branches = branch
		}

		// If there are no hooks, set a 60s polling hook
		if len(repo.Hooks) == 0 {
			hook := &Hook{
				Type:     "polling",
				Interval: 60 * time.Second,
			}

			repo.Hooks = append(repo.Hooks, hook)
		}

		// expand tilde home directory for key path
		if repo.Credentials.PrivateKey.Key != "" {
			if strings.HasPrefix(repo.Credentials.PrivateKey.Key, "~/") {
				dirname, _ := os.UserHomeDir()
				repo.Credentials.PrivateKey.Key = filepath.Join(dirname, repo.Credentials.PrivateKey.Key[2:])
			}
		}
	}
}

func (c *Config) setDefaultLogConfig() {
	if c.Log.Format == "" {
		c.Log.Format = "text"
	}
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
}
