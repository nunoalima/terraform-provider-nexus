/*
Use this resource to create a Nexus Repository.

Example Usage

Example Usage - Apt proxy repository

```hcl
resource "nexus_repository_apt_proxy" "apt_proxy" {
  name   = "apt-proxy"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  cleanup {
	  policy_names = [
		  "string"
	  ]
  }

  proxy {
     remote_url = "https://remote.repository.com"
	 content_max_age = 1440
	 metadata_max_age = 1440
  }

  negative_cache {
	  enabled      = true
	  time_to_live = 1440
  }

  http_client {
	  blocked    = false
	  auto_block = true
	  connection {
		  retries                   = 0
		  user_agent_suffix         = "string"
		  timeout                   = 60
		  enable_circular_redirects = false
		  enable_cookies            = false
		  use_trust_store           = false
	  }
	  authentication {
		  type        = "username"
		  username    = "string"
		  password    = "string"
		  ntlm_host   = "string"
		  ntlm_domain = "string"
	  }
  }

  routing_rule = "string"

  apt {
    distribution = "bionic"
	flat         = false
  }

}
```

*/
package nexus

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRepositoryAptProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryAptProxyCreate,
		Read:   resourceRepositoryAptProxyRead,
		Update: resourceRepositoryAptProxyUpdate,
		Delete: resourceRepositoryAptProxyDelete,
		Exists: resourceRepositoryAptProxyExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Default:     true,
				Description: "Whether this repository accepts incoming requests",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"storage": {
				Description: "The storage configuration of the repository",
				DefaultFunc: repositoryStorageDefault,
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_store_name": {
							Description: "Blob store used to store repository contents",
							Required:    true,
							Type:        schema.TypeString,
						},
						"strict_content_type_validation": {
							Default:     true,
							Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
							Optional:    true,
							Type:        schema.TypeBool,
						},
					},
				},
			},
			"cleanup": {
				Description: "Cleanup policies",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Description: "List of policy names",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Set: func(v interface{}) int {
								return schema.HashString(strings.ToLower(v.(string)))
							},
							Type: schema.TypeSet,
						},
					},
				},
			},
			"http_client": {
				Description: "HTTP Client configuration for proxy repositories",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blocked": {
							Description: "Whether to block outbound connections on the repository",
							Required:    true,
							Type:        schema.TypeBool,
						},
						"auto_block": {
							Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
							Required:    true,
							Type:        schema.TypeBool,
						},
						"connection": {
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retries": {
										Description:  "Total retries if the initial connection attempt suffers a timeout",
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(0, 10),
									},
									"user_agent_suffix": {
										Description: "Custom fragment to append to User-Agent header in HTTP requests",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"timeout": {
										Description:  "Seconds to wait for activity before stopping and retrying the connection",
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(1, 3600),
									},
									"enable_circular_redirects": {
										Default:     false,
										Description: "Whether to enable circular redirects ",
										Optional:    true,
										Type:        schema.TypeBool,
									},
									"enable_cookies": {
										Default:     false,
										Description: "Whether to allow cookies to be stored and used",
										Optional:    true,
										Type:        schema.TypeBool,
									},
									"use_trust_store": {
										Default:     false,
										Description: "Whether to use trust store",
										Optional:    true,
										Type:        schema.TypeBool,
									},
								},
							},
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
						"authentication": {
							Description: "Authentication configuration of the HTTP client",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description:  "Authentication type",
										Optional:     true,
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"ntlm", "username"}, false),
									},
									"username": {
										Description: "The username used by the proxy repository",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"password": {
										Description: "The password used by the proxy repository",
										Optional:    true,
										Sensitive:   true,
										Type:        schema.TypeString,
									},
									"ntlm_host": {
										Description: "The ntlm host to connect",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"ntlm_domain": {
										Description: "The ntlm domain to connect",
										Optional:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"negative_cache": {
				Description: "Configuration of the negative cache handling",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Whether to cache responses for content not present in the proxied repository",
							Required:    true,
							Type:        schema.TypeBool,
						},
						"ttl": {
							Default:     1440,
							Description: "How long to cache the fact that a file was not found in the repository (in minutes)",
							Optional:    true,
							Type:        schema.TypeInt,
						},
					},
				},
			},
			"proxy": {
				Description: "Configuration for the proxy repository",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_max_age": {
							Description: "How long (in minutes) to cache artifacts before rechecking the remote repository",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1440,
						},
						"metadata_max_age": {
							Description: "How long (in minutes) to cache metadata before rechecking the remote repository.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1440,
						},
						"remote_url": {
							Description: "Location of the remote repository being proxied",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"routing_rule": {
				Description: "assign an existing routing rule",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"apt": {
				Description: "Apt specific configuration of the repository",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distribution": {
							Description: "The linux distribution",
							Type:        schema.TypeString,
							Required:    true,
						},
						"flat": {
							Description: "Whether this repository is flat",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceRepositoryAptProxyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	repo := getRepositoryAptProxyFromResourceData(d)

	if err := client.RepositoryCreate(repo); err != nil {
		return err
	}

	if err := setRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceRepositoryAptProxyRead(d, m)
}

func resourceRepositoryAptProxyRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	repo, err := nexusClient.RepositoryRead(d.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		d.SetId("")
		return nil
	}

	return setRepositoryToResourceData(repo, d)
}

func resourceRepositoryAptProxyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	repoName := d.Id()
	repo := getRepositoryAptProxyFromResourceData(d)

	if err := client.RepositoryUpdate(repoName, repo); err != nil {
		return err
	}

	if err := setRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryAptProxyDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	return nexusClient.RepositoryDelete(d.Id())
}

func resourceRepositoryAptProxyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	repo, err := nexusClient.RepositoryRead(d.Id())
	return repo != nil, err
}

func getRepositoryAptProxyFromResourceData(d *schema.ResourceData) nexus.Repository {
	repo := nexus.Repository{
		Format: "apt",
		Name:   d.Get("name").(string),
		Online: d.Get("online").(bool),
		Type:   "proxy",
	}

	if _, ok := d.GetOk("apt"); ok {
		aptList := d.Get("apt").([]interface{})
		aptConfig := aptList[0].(map[string]interface{})

		repo.RepositoryApt = &nexus.RepositoryApt{
			Distribution: aptConfig["distribution"].(string),
			Flat:         aptConfig["flat"].(bool),
		}
	}

	if _, ok := d.GetOk("cleanup"); ok {
		cleanupList := d.Get("cleanup").([]interface{})
		cleanupConfig := cleanupList[0].(map[string]interface{})
		repo.RepositoryCleanup = &nexus.RepositoryCleanup{
			PolicyNames: interfaceSliceToStringSlice(cleanupConfig["policy_names"].(*schema.Set).List()),
		}
	}

	if _, ok := d.GetOk("group"); ok {
		groupList := d.Get("group").([]interface{})
		groupMemberNames := make([]string, 0)

		if len(groupList) == 1 && groupList[0] != nil {
			groupConfig := groupList[0].(map[string]interface{})
			groupConfigMemberNames := groupConfig["member_names"].(*schema.Set)

			for _, v := range groupConfigMemberNames.List() {
				groupMemberNames = append(groupMemberNames, v.(string))
			}
		}
		repo.RepositoryGroup = &nexus.RepositoryGroup{
			MemberNames: groupMemberNames,
		}
	}

	if _, ok := d.GetOk("http_client"); ok {
		httpClientList := d.Get("http_client").([]interface{})
		httpClientConfig := httpClientList[0].(map[string]interface{})

		repo.RepositoryHTTPClient = &nexus.RepositoryHTTPClient{
			AutoBlock: httpClientConfig["auto_block"].(bool),
			Blocked:   httpClientConfig["blocked"].(bool),
		}

		if v, ok := httpClientConfig["authentication"]; ok {
			authList := v.([]interface{})
			if len(authList) == 1 && authList[0] != nil {
				authConfig := authList[0].(map[string]interface{})

				repo.RepositoryHTTPClient.Authentication = &nexus.RepositoryHTTPClientAuthentication{
					NTLMDomain: authConfig["ntlm_domain"].(string),
					NTLMHost:   authConfig["ntlm_host"].(string),
					Type:       authConfig["type"].(string),
					Username:   authConfig["username"].(string),
					Password:   authConfig["password"].(string),
				}
			}
		}

		if v, ok := httpClientConfig["connection"]; ok {
			connList := v.([]interface{})
			if len(connList) == 1 && connList[0] != nil {
				connConfig := connList[0].(map[string]interface{})

				repo.RepositoryHTTPClient.Connection = &nexus.RepositoryHTTPClientConnection{
					//EnableCookies:   connConfig["enable_cookis"].(bool),
					Retries: connConfig["retries"].(*int),
					Timeout: connConfig["timeout"].(*int),
					//UserAgentSuffix: connConfig["user_agent_suffix"].(*string),
				}
			}
		}
	}

	if _, ok := d.GetOk("negative_cache"); ok {
		negativeCacheList := d.Get("negative_cache").([]interface{})
		negativeCacheConfig := negativeCacheList[0].(map[string]interface{})

		repo.RepositoryNegativeCache = &nexus.RepositoryNegativeCache{
			Enabled: negativeCacheConfig["enabled"].(bool),
			TTL:     negativeCacheConfig["ttl"].(int),
		}
	}

	if _, ok := d.GetOk("proxy"); ok {
		proxyList := d.Get("proxy").([]interface{})
		proxyConfig := proxyList[0].(map[string]interface{})
		repo.RepositoryProxy = &nexus.RepositoryProxy{
			ContentMaxAge:  proxyConfig["content_max_age"].(int),
			MetadataMaxAge: proxyConfig["metadata_max_age"].(int),
			RemoteURL:      proxyConfig["remote_url"].(string),
		}
	}

	if _, ok := d.GetOk("storage"); ok {
		storageList := d.Get("storage").([]interface{})
		storageConfig := storageList[0].(map[string]interface{})

		repo.RepositoryStorage = &nexus.RepositoryStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		}
		// Only hosted repository has attribute WritePolicy
		if repo.Type == nexus.RepositoryTypeHosted {
			writePolicy := storageConfig["write_policy"].(string)
			repo.RepositoryStorage.WritePolicy = &writePolicy
		}
	}

	return repo
}
