package consul

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/spaghettifunk/norman/internal/common/entities"
	cingestion "github.com/spaghettifunk/norman/internal/common/ingestion"

	"github.com/spaghettifunk/norman/internal/service"
)

var (
	// TTL time to life for service in consul
	TTL int = 15
	// Where the ServiceKVPath resides
	ServicesKVpath string = "services"
	// Where the JobsKVPath resides
	JobsKVPath string = "ingestion-jobs"
	// Where the TablesKVPath resides
	TablesKVPath string = "tables"
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

// Store the consul table configuration
func (c *Consul) PutTableConfiguration(table *entities.Table) error {
	js, err := json.Marshal(table)
	if err != nil {
		return err
	}
	key := formattedTableKey(table.Name)
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: js,
	}
	_, err = c.KV.Put(&pair, nil)
	return err
}

// Retrieve the consul table configuration
func (c *Consul) GetTableConfiguration(name string) (*entities.Table, error) {
	key := formattedTableKey(name)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return nil, err
	}
	cfg := &entities.Table{}
	if err := json.Unmarshal(v.Value, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Store the consul ingestion job configuration
func (c *Consul) PutIngestionJobConfiguration(config *cingestion.IngestionJobConfiguration) error {
	cfg, err := json.Marshal(config)
	if err != nil {
		return err
	}

	key := formattedJobConfigKey(config.ID.String())
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: cfg,
	}
	_, err = c.KV.Put(&pair, nil)
	return err
}

// Retrieve the consul ingestion job configuration
func (c *Consul) GetIngestionJobConfiguration(id string) (*cingestion.IngestionJobConfiguration, error) {
	key := formattedJobConfigKey(id)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return nil, err
	}
	var cfg *cingestion.IngestionJobConfiguration
	if err := json.Unmarshal(v.Value, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Store the consul ingestion job status
func (c *Consul) PutIngestionJobStatus(id string, status cingestion.JobStatus) error {
	key := formattedJobStatusKey(id)
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: []byte(status),
	}
	_, err := c.KV.Put(&pair, nil)
	return err
}

// Retrieve the consul ingestion job status
func (c *Consul) GetIngestionJobStatus(id string) (cingestion.JobStatus, error) {
	key := formattedJobStatusKey(id)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return "", err
	}
	var status *cingestion.JobStatus
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
	name := fmt.Sprintf("%v-%v", s.GetName(), s.GetHost())
	return strings.Replace(name, ".", "-", -1)
}

// formattedKey returns correctly formatted key of the service
func formattedKey(s service.Service) string {
	return fmt.Sprintf("%v/%v/%v/definition", ServicesKVpath, s.GetName(), s.GetID())
}

// formattedTableKey returns correctly formatted key of the table
// TODO: how do we get the tenant name?
func formattedTableKey(name string) string {
	// Format: tables/{tenantId}/tableName/definition
	return fmt.Sprintf("%v/%v/%v/definition", TablesKVPath, "default", name)
}

// formattedJobStatusKey returns correctly formatted key of the job
// TODO: how do we get the tenant name?
func formattedJobStatusKey(id string) string {
	// Format: Jobs/{tenantId}/jobId/status
	return fmt.Sprintf("%v/%v/%v/status", JobsKVPath, "default", id)
}

// formattedJobConfigKey returns correctly formatted key of the job
// TODO: how do we get the tenant name?
func formattedJobConfigKey(id string) string {
	// Format: Jobs/{tenantId}/jobId/configurations
	return fmt.Sprintf("%v/%v/%v/configurations", JobsKVPath, "default", id)
}

// formattedID returns correctly formatted id of the service
func formattedID(s service.Service) string {
	return fmt.Sprintf("%v-%v-%v", formattedName(s), s.GetID(), s.GetPort())
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
