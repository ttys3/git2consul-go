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
	"gopkg.in/yaml.v2"
	"io"
	"time"
)

// Config is used to represent the passed in configuration
type Config struct {
	LocalStore string               `json:"local_store" yaml:"local_store"`
	Webhook    *WebhookServerConfig `json:"webhook" yaml:"webhook"`
	Repos      []*Repo              `json:"repos" yaml:"repos"`
	Consul     *ConsulConfig        `json:"consul" yaml:"consul"`
}

func (c Config) DumpSampleConfig(w io.Writer) error {
	c.LocalStore = "/var/lib/git2consul"

	c.Webhook = &WebhookServerConfig{
		Address: "",
		Port:    8484,
	}

	c.Repos = []*Repo{
		&Repo{
			Name:     "consul-kv-config",
			URL:      "ssh://git@git.nomad.lan:2222/ttys3/consul-kv-config.git",
			Branches: []string{"main"},
			Hooks: []*Hook{
				&Hook{
					Type:     "webhook",
					Interval: 30 * time.Second,
					URL:      "",
				},
			},
			SourceRoot:     "/",
			MountPoint:     "",
			ExpandKeys:     false,
			SkipBranchName: false,
			SkipRepoName:   false,
			Credentials: Credentials{
				Username: "",
				Password: "",
				PrivateKey: PrivateKey{
					Key:      "~/.ssh/id_ed25519",
					Username: "git",
					Password: "",
				},
			},
		},
	}

	c.Consul = &ConsulConfig{
		Address:   "127.0.0.1:8500",
		Token:     "",
		SSLEnable: false,
		SSLVerify: false,
	}
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

// Credentials is the representation of git authentication
type Credentials struct {
	Username   string     `json:"username,omitempty" yaml:"username,omitempty"`
	Password   string     `json:"password,omitempty" yaml:"password,omitempty"`
	PrivateKey PrivateKey `json:"private_key,omitempty" yaml:"private_key,omitempty"`
}

// PrivateKey is the representation of private key used for the authentication
type PrivateKey struct {
	Key      string `json:"pk_key" yaml:"key"`
	Username string `json:"pk_username,omitempty" yaml:"username,omitempty"`
	Password string `json:"pk_password,omitempty" yaml:"password,omitempty"`
}

// Hook is the configuration for hooks
type Hook struct {
	Type string `json:"type" yaml:"type"`

	// Specific to polling
	Interval time.Duration `json:"interval" yaml:"interval"`

	// Specific to webhooks
	URL string `json:"url,omitempty" yaml:"url"`
}

// Repo is the configuration for the repository
type Repo struct {
	Name           string      `json:"name" yaml:"name"`
	URL            string      `json:"url" yaml:"url"`
	Branches       []string    `json:"branches" yaml:"branches"`
	Hooks          []*Hook     `json:"hooks" yaml:"hooks"`
	SourceRoot     string      `json:"source_root" yaml:"source_root"`
	MountPoint     string      `json:"mount_point" yaml:"mount_point"`
	ExpandKeys     bool        `json:"expand_keys,omitempty" yaml:"expand_keys,omitempty"`
	SkipBranchName bool        `json:"skip_branch_name,omitempty" yaml:"skip_branch_name,omitempty"`
	SkipRepoName   bool        `json:"skip_repo_name,omitempty" yaml:"skip_repo_name,omitempty"`
	Credentials    Credentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
}

// WebhookServerConfig is the configuration for the git hoooks server
type WebhookServerConfig struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	Port    int    `json:"port" yaml:"port"`
}

// ConsulConfig is the configuration for the Consul client
type ConsulConfig struct {
	Address   string `json:"address" yaml:"address"` // default is 127.0.0.1:8500
	Token     string `json:"token,omitempty" yaml:"token,omitempty"`
	SSLEnable bool   `json:"ssl" yaml:"ssl_enable"`
	SSLVerify bool   `json:"ssl_verify,omitempty" yaml:"ssl_verify,omitempty"`
	// can also set from env, see github.com/hashicorp/consul/api@v1.12.0/api.go
	// consul api.DefaultConfig() will handle these env vars
	// if mTLS is enabled on consul, below env vars should be configured
	// api.TLSConfig.Address CONSUL_TLS_SERVER_NAME
	// api.TLSConfig.CAFile CONSUL_CACERT
	// api.TLSConfig.CertFile CONSUL_CLIENT_CERT
	// api.TLSConfig.KeyFile CONSUL_CLIENT_KEY
	// api.TLSConfig.InsecureSkipVerify CONSUL_HTTP_SSL_VERIFY
}

func (r *Repo) String() string {
	if r != nil {
		return r.Name
	}
	return ""
}
