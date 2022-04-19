package metadata

import (
	"errors"
	"time"
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/processors"
)

type MetadataProcessor struct {
	OpenstackTags          []string        `toml:"openstack_tags"`
	Timeout          config.Duration `toml:"timeout"`
	Log              telegraf.Logger `toml:"-"`
}

type OpenstackMetadata struct {
	Project string
}

var m OpenstackMetadata

type Response struct {
	UUID string `json:"uuid"`
	Meta struct {
			Project       string `json:"project"`
			Owner         string `json:"owner"`
			ServiceName   string `json:"service_name"`
			Group         string `json:"group"`
			Groups        string `json:"groups"`
			ExpireAt      string `json:"expire_at"`
			WfExID        string `json:"wf_ex_id"`
			OsType        string `json:"os_type"`
			OsDistro      string `json:"os_distro"`
			OsVersion     string `json:"os_version"`
			SafeExpire    string `json:"safe_expire"`
			Userdata      string `json:"userdata"`
			EmailAlertsCc string `json:"email_alerts_cc"`
			Fqdn          string `json:"fqdn"`
	} `json:"meta"`
	Keys []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Data string `json:"data"`
	} `json:"keys"`
	Hostname         string        `json:"hostname"`
	Name             string        `json:"name"`
	LaunchIndex      int           `json:"launch_index"`
	AvailabilityZone string        `json:"availability_zone"`
	RandomSeed       string        `json:"random_seed"`
	ProjectID        string        `json:"project_id"`
	Devices          []interface{} `json:"devices"`
}

var meta_resp Response

const metadata_url = "http://169.254.169.254/openstack/latest/meta_data.json"

const sampleConfig = `
  ## Available tags to attach to metrics:
  ## * uuid
  ## * project
  ## * owner
  ## * service_name
  ## * group
  ## * fqdn
  ## * hostname
  ## * name
  ## * availability_zone
  ## * project_id
  openstack_tags = [ "project", "availability_zone" ]
`

const (
	DefaultTimeout             = 10 * time.Second
)

func (r *MetadataProcessor) SampleConfig() string {
	return sampleConfig
}

func (r *MetadataProcessor) Description() string {
	return "Attach Openstack metadata to metrics"
}

func (r *MetadataProcessor) Apply(in ...telegraf.Metric) []telegraf.Metric {
	// add tags
	for _, metric := range in {
		r.Log.Debug("length is ",len(r.OpenstackTags))
		r.Log.Debug(r.OpenstackTags)
		for _, tag := range r.OpenstackTags {
			r.Log.Debug("checking tag=",tag)
			if v := getTagFromMetadataResponse(meta_resp, tag); v != "" {
				r.Log.Debug("adding tag=",tag," value=",v)
				metric.AddTag(tag, v)
			}
		}
	}
	return in
}

func (r *MetadataProcessor) Init() error {
	r.Log.Debug("Initializing Openstack Metadata Processor")
	if len(r.OpenstackTags) == 0 {
		return errors.New("no tags specified in configuration")
	}
	meta_resp = getMetadata()
	r.Log.Debug(PrettyPrint(meta_resp))
	return nil
}

func init() {
	processors.Add("metadata", func() telegraf.Processor {
		return &MetadataProcessor{}
	})
}

func getMetadata() Response {
    resp, err := http.Get(metadata_url)
    if err != nil {
        fmt.Println("No response from request")
    }
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {   // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func getTagFromMetadataResponse(o Response, tag string) string {
	switch tag {
	case "uuid":
		return o.UUID
	case "project":
		return o.Meta.Project
	case "owner":
		return o.Meta.Owner
	case "service_name":
		return o.Meta.ServiceName
	case "group":
		return o.Meta.Group
	case "fqdn":
		return o.Meta.Fqdn
	case "hostname":
		return o.Hostname
	case "name":
		return o.Name
	case "availability_zone":
		return o.AvailabilityZone
	case "project_id":
		return o.ProjectID
	default:
		return ""
	}
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}