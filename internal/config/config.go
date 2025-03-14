package config

import (
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
)

const RepositoryToken = "config"

// Config represents the configurable options of the system
type Config struct {
	item.Item

	Name                    string `json:"name"`
	Domain                  string `json:"domain"`
	BindAddress             string `json:"bind_addr"`
	HTTPPort                string `json:"http_port"`
	HTTPSPort               string `json:"https_port"`
	AdminEmail              string `json:"admin_email"`
	ClientSecret            string `json:"client_secret"`
	Etag                    string `json:"etag"`
	DisableCORS             bool   `json:"cors_disabled"`
	DisableGZIP             bool   `json:"gzip_disabled"`
	DisableHTTPCache        bool   `json:"cache_disabled"`
	CacheMaxAge             int64  `json:"cache_max_age"`
	BackupBasicAuthUser     string `json:"backup_basic_auth_user"`
	BackupBasicAuthPassword string `json:"backup_basic_auth_password"`
}

func (c *Config) GetTitle() string { return c.Name }

func (*Config) GetRepositoryToken() string {
	return RepositoryToken
}

func (*Config) EntityName() string {
	return "Config"
}

func (*Config) Time() int64 {
	return 0
}

// MarshalEditor writes a buffer of templates to edit a Post and partially implements editor.Editable
func (c *Config) MarshalEditor(publicPath string) ([]byte, error) {
	return editor.Form(c,
		editor.Field{
			View: editor.Input("Name", c, map[string]string{
				"label":       "Site Name",
				"placeholder": "Add a name to this site (internal use only)",
			}, nil),
		},
		editor.Field{
			View: editor.Input("Domain", c, map[string]string{
				"label":       "Domain Name (required for SSL certificate)",
				"placeholder": "e.g. www.example.com or example.com",
			}, nil),
		},
		editor.Field{
			View: editor.Input("BindAddress", c, map[string]string{
				"type": "hidden",
			}, nil),
		},
		editor.Field{
			View: editor.Input("HTTPPort", c, map[string]string{
				"type": "hidden",
			}, nil),
		},
		editor.Field{
			View: editor.Input("HTTPSPort", c, map[string]string{
				"type": "hidden",
			}, nil),
		},
		editor.Field{
			View: editor.Input("AdminEmail", c, map[string]string{
				"label": "Administrator Email (notified of internal system information)",
			}, nil),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"label":    "Client Secret (used to validate requests, DO NOT SHARE)",
				"disabled": "true",
			}, nil),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"type": "hidden",
			}, nil),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"label":    "Etag Header (used to cache resources)",
				"disabled": "true",
			}, nil),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"type": "hidden",
			}, nil),
		},
		editor.Field{
			View: editor.Checkbox("DisableCORS", c, map[string]string{
				"label": "Disable CORS",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableGZIP", c, map[string]string{
				"label": "Disable GZIP",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableHTTPCache", c, map[string]string{
				"label": "Disable HTTP Cache",
			}),
		},
		editor.Field{
			View: editor.Input("CacheMaxAge", c, map[string]string{
				"label": "Max-Age value for HTTP caching (in seconds, 0 = 2592000)",
				"type":  "text",
			}, nil),
		},
		editor.Field{
			View: editor.Input("BackupBasicAuthUser", c, map[string]string{
				"label":       "HTTP Basic Auth User",
				"placeholder": "Enter a user name for Basic Auth access",
				"type":        "text",
			}, nil),
		},
		editor.Field{
			View: editor.Input("BackupBasicAuthPassword", c, map[string]string{
				"label":       "HTTP Basic Auth Password",
				"placeholder": "Enter a password for Basic Auth access",
				"type":        "password",
			}, nil),
		},
	)
}

var Builder = func() interface{} {
	return new(Config)
}
