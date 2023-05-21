package consul

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/spaghettifunk/norman/internal/common/ingestion"
	"github.com/spaghettifunk/norman/internal/common/schema"
	"github.com/spaghettifunk/norman/internal/service"
)

var (
	// TTL time to life for service in consul
	TTL int = 15
	// Where the ServiceKVPath resides
	ServicesKVpath string = "services"
	// Where the JobsKVPath resides
	JobsKVPath string = "ingestion-jobs"
	// Where the JobsKVPath resides
	SchemasKVPath string = "schemas"
	// Config falls back to client default config
	ConsulConfig *api.Config = api.DefaultConfig()
)

// Consul structure
type Consul struct {
	// Agent to register service
	Agent *api.Agent
	// KV to save service definition
	KV            *api.KV
	heartBeatKill chan bool
}

// New Consul
func New() *Consul {
	c := new(Consul)
	c.heartBeatKill = make(chan bool)
	return c
}

// Start Register the service in the consul pool of services
func (c *Consul) Start(s service.Service) error {
	p, err := strconv.ParseInt(s.GetPort(), 10, 0)
	if err != nil {
		return err
	}
	AgentService := api.AgentServiceRegistration{
		ID:   formattedID(s),
		Name: formattedName(s),
		Port: int(p),
		Check: &api.AgentServiceCheck{
			TTL: fmt.Sprintf("%vs", TTL),
		},
	}
	// Register the service
	err = c.Agent.ServiceRegister(&AgentService)
	if err != nil {
		return err
	}
	// Initial run for TTL
	c.Agent.PassTTL(fmt.Sprintf("service:%v", formattedID(s)), "TTL heartbeat")

	// Begin TTL refresh
	go c.Heartbeat(s)
	return nil
}

// Kill the heartbeat and remove the service
func (c *Consul) Stop(s service.Service) error {
	c.heartBeatKill <- true
	return c.Agent.ServiceDeregister(formattedID(s))
}

// Init Consul with Default Settings
func (c *Consul) Init() error {
	client, err := api.NewClient(ConsulConfig)
	if err != nil {
		return err
	}
	if c.Agent == nil {
		c.Agent = client.Agent()
	}

	if c.KV == nil {
		c.KV = client.KV()
	}
	return nil
}

// Send service definition to consul
func (c *Consul) Declare(s service.Service) error {
	js, err := json.Marshal(s)
	if err != nil {
		return err
	}
	key := formattedKey(s)
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: js,
	}
	_, err = c.KV.Put(&pair, nil)
	return err
}

// Store the consul schema configuration
func (c *Consul) PutSchemaConfiguration(schema *schema.Schema) error {
	js, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	key := formattedSchemaKey(schema.Name)
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: js,
	}
	_, err = c.KV.Put(&pair, nil)
	return err
}

// Retrieve the consul schema configuration
func (c *Consul) GetSchemaConfiguration(name string) (*schema.Schema, error) {
	key := formattedSchemaKey(name)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return nil, err
	}
	cfg := &schema.Schema{}
	if err := json.Unmarshal(v.Value, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Store the consul ingestion job status
func (c *Consul) PutIngestionJobStatus(id string, status ingestion.JobStatus) error {
	key := formattedJobKey(id)
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: []byte(status),
	}
	_, err := c.KV.Put(&pair, nil)
	return err
}

// Retrieve the consul ingestion job status
func (c *Consul) GetIngestionJobStatus(id string) (ingestion.JobStatus, error) {
	key := formattedJobKey(id)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return "", err
	}
	var status *ingestion.JobStatus
	if err := json.Unmarshal(v.Value, &status); err != nil {
		return "", err
	}
	return *status, nil
}

// Retrieve the consul service definition
func (c *Consul) GetService(s service.Service) error {
	key := formattedKey(s)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return err
	}
	return json.Unmarshal(v.Value, s)
}

// formattedName returns correctly formatted name of the service
func formattedName(s service.Service) string {
	name := fmt.Sprintf("%v-%v", s.GetName(), s.GetID())
	return strings.Replace(name, ".", "-", -1)
}

// formattedKey returns correctly formatted key of the service
func formattedKey(s service.Service) string {
	return fmt.Sprintf("%v/%v/%v/definition", ServicesKVpath, s.GetName(), s.GetID())
}

// formattedSchemaKey returns correctly formatted key of the schema
// TODO: how do we get the tenant name?
func formattedSchemaKey(name string) string {
	// Format: schemas/{tenantId}/schemaName/definition
	return fmt.Sprintf("%v/%v/%v/definition", SchemasKVPath, "default", name)
}

// formattedJobKey returns correctly formatted key of the job
// TODO: how do we get the tenant name?
func formattedJobKey(id string) string {
	// Format: Jobs/{tenantId}/jobId/definition
	return fmt.Sprintf("%v/%v/%v/definition", JobsKVPath, "default", id)
}

// formattedID returns correctly formatted id of the service
func formattedID(s service.Service) string {
	return fmt.Sprintf("%v-%v-%v", formattedName(s), s.GetHost(), s.GetPort())
}

// Heartbeat begins heart beat of health check.
func (c *Consul) Heartbeat(s service.Service) {
	t := time.NewTicker(time.Duration(TTL-1) * time.Second)
	for range t.C {
		select {
		case <-c.heartBeatKill:
			c.Stop(s)
			return
		default:
		}
		c.Agent.PassTTL(fmt.Sprintf("service:%v", formattedID(s)), "TTL heartbeat")
	}
}
