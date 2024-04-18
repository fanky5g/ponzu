package entities

import (
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
)

// Config represents the configurable options of the system
type Config struct {
	item.Item

	Name                    string   `json:"name"`
	Domain                  string   `json:"domain"`
	BindAddress             string   `json:"bind_addr"`
	HTTPPort                string   `json:"http_port"`
	HTTPSPort               string   `json:"https_port"`
	AdminEmail              string   `json:"admin_email"`
	ClientSecret            string   `json:"client_secret"`
	Etag                    string   `json:"etag"`
	DisableCORS             bool     `json:"cors_disabled"`
	DisableGZIP             bool     `json:"gzip_disabled"`
	DisableHTTPCache        bool     `json:"cache_disabled"`
	CacheMaxAge             int64    `json:"cache_max_age"`
	CacheInvalidate         []string `json:"cache"`
	BackupBasicAuthUser     string   `json:"backup_basic_auth_user"`
	BackupBasicAuthPassword string   `json:"backup_basic_auth_password"`
}

func (c *Config) Title() string { return c.Name }

// MarshalEditor writes a buffer of templates to edit a Post and partially implements editor.Editable
func (c *Config) MarshalEditor(paths config.Paths) ([]byte, error) {
	view, err := editor.Form(c,
		paths,
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
				"label": "Disable CORS (so only " + c.Domain + " can fetch your data)",
			}, map[string]string{
				"true": "Disable CORS",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableGZIP", c, map[string]string{
				"label": "Disable GZIP (will increase server speed, but also bandwidth)",
			}, map[string]string{
				"true": "Disable GZIP",
			}),
		},
		editor.Field{
			View: editor.Checkbox("DisableHTTPCache", c, map[string]string{
				"label": "Disable HTTP Cache (overrides 'Cache-Control' header)",
			}, map[string]string{
				"true": "Disable HTTP Cache",
			}),
		},
		editor.Field{
			View: editor.Input("CacheMaxAge", c, map[string]string{
				"label": "Max-Age value for HTTP caching (in seconds, 0 = 2592000)",
				"type":  "text",
			}, nil),
		},
		editor.Field{
			View: editor.Checkbox("CacheInvalidate", c, map[string]string{
				"label": "Invalidate cache on save",
			}, map[string]string{
				"invalidate": "Invalidate Cache",
			}),
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
	if err != nil {
		return nil, err
	}

	openingTag := []byte(`
	<div class="card">
		<div class="card-content">
			<div class="card-title">System Configuration</div>
		</div>
		<form action="` + paths.PublicPath + `/configure" method="post">
	`)

	closingTag := []byte(`</form></div>`)
	script := []byte(`
	<script>
		$(function() {
			// hide default fields & labels unnecessary for the config
			var fields = $('.default-fields');
			fields.css('position', 'relative');
			fields.find('input:not([type=submit])').remove();
			fields.find('label').remove();
			fields.find('button').css({
				position: 'absolute',
				top: '-10px',
				right: '0px'
			});

			var contentOnly = $('.entities-only.__ponzu');
			contentOnly.hide();
			contentOnly.find('input, textarea, select').attr('name', '');

			// adjust layout of td so save button is in same location as usual
			fields.find('td').css('float', 'right');

			// stop some fixed config settings from being modified
			fields.find('input[name=client_secret]').attr('name', '');
		});
	</script>
	`)

	view = append(openingTag, view...)
	view = append(view, closingTag...)
	view = append(view, script...)

	return view, nil
}

var ConfigBuilder = func() interface{} {
	return new(Config)
}
