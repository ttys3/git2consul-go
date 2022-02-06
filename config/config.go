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
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is used to represent the passed in configuration
type Config struct {
	LocalStore string               `json:"local_store" yaml:"local_store"`
	Webhook    *WebhookServerConfig `json:"webhook" yaml:"webhook"`
	Repos      []*Repo              `json:"repos" yaml:"repos"`
	Consul     *ConsulConfig        `json:"consul,omitempty" yaml:"consul,omitempty"`
	Log        *LogConfig           `json:"log,omitempty" yaml:"log,omitempty"`
}

func (c Config) String() string {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (c Config) Dump(w io.Writer) error {
	return dumping(w, c)
}

func dumping(w io.Writer, c Config) error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

func (c Config) DumpSampleConfig(w io.Writer) error {
	c.LocalStore = "/var/lib/git2consul"

	c.Webhook = &WebhookServerConfig{
		Address: "",
		Port:    8484,
	}

	c.Repos = []*Repo{
		{
			Name:     "consul-kv-config",
			URL:      "ssh://git@git.nomad.lan:2222/ttys3/consul-kv-config.git",
			Branches: []string{"main"},
			Hooks: []*Hook{
				{
					Type:     "webhook",
					Interval: 30 * time.Second,
					URL:      "",
				},
			},
			SourceRoot:     "/",
			MountPoint:     "",
			ExpandKeys:     false,
			SkipBranchName: true,
			SkipRepoName:   false,
			Credentials: Credentials{
				Username: "",
				Password: "",
				PrivateKey: PrivateKey{
					Key:              "~/.ssh/id_ed25519",
					SkipHostKeyCheck: true,
					Username:         "git",
					Password:         "",
				},
			},
		},
	}

	c.Consul = &ConsulConfig{
		Address:   "127.0.0.1:8500",
		Token:     "",
		SSLEnable: false,
	}
	return dumping(w, c)
}

// Credentials is the representation of git authentication
type Credentials struct {
	Username   string     `json:"username,omitempty" yaml:"username,omitempty"`
	Password   string     `json:"password,omitempty" yaml:"password,omitempty"`
	PrivateKey PrivateKey `json:"private_key,omitempty" yaml:"private_key,omitempty"`
}

// PrivateKey is the representation of private key used for the authentication
type PrivateKey struct {
	Key              string `json:"key" yaml:"key"`
	SkipHostKeyCheck bool   `json:"skip_host_key_check,omitempty" yaml:"skip_host_key_check,omitempty"`
	Username         string `json:"username,omitempty" yaml:"username,omitempty"`
	Password         string `json:"password,omitempty" yaml:"password,omitempty"`
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

func (r *Repo) String() string {
	if r != nil {
		return r.Name
	}
	return ""
}

// WebhookServerConfig is the configuration for the git hoooks server
type WebhookServerConfig struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	Port    int    `json:"port" yaml:"port"`
}

// ConsulConfig is the configuration for the Consul client
type ConsulConfig struct {
	Address   string          `json:"address,omitempty" yaml:"address,omitempty"` // default to 127.0.0.1:8500 according to consul go SDK
	Token     string          `json:"token,omitempty" yaml:"token,omitempty"`
	SSLEnable bool            `json:"ssl_enable" yaml:"ssl_enable"`
	TLSConfig ConsulTLSConfig `json:"tls_config" yaml:"tls_config,omitempty"`
}

// ConsulTLSConfig used for consul mTLS auth
// can also set from env, see github.com/hashicorp/consul/api@v1.12.0/api.go
// consul api.DefaultConfig() will handle these env vars
// if mTLS is enabled on consul, below env vars should be configured
type ConsulTLSConfig struct {
	ServerName         string `json:"server_name,omitempty" yaml:"server_name,omitempty"`                   // api.TLSConfig.Address CONSUL_TLS_SERVER_NAME
	CAFile             string `json:"ca_file,omitempty" yaml:"ca_file,omitempty"`                           // api.TLSConfig.CAFile CONSUL_CACERT
	CertFile           string `json:"cert_file,omitempty" yaml:"cert_file,omitempty"`                       // api.TLSConfig.CertFile CONSUL_CLIENT_CERT
	KeyFile            string `json:"key_file,omitempty" yaml:"key_file,omitempty"`                         // api.TLSConfig.KeyFile CONSUL_CLIENT_KEY
	InsecureSkipVerify bool   `json:"insecure_skip_verify,omitempty" yaml:"insecure_skip_verify,omitempty"` // api.TLSConfig.InsecureSkipVerify CONSUL_HTTP_SSL_VERIFY
}

type LogConfig struct {
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	Level  string `json:"level,omitempty" yaml:"level,omitempty"`
}
