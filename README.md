# Openstack Metadata Processor Plugin

Openstack Metadata processor plugin appends metadata gathered from Openstack
to metrics.

## Configuration

```toml
[[processors.metadata]]
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
```

## Example

Append `project` and `availability_zone` to metrics tags:

```toml
[[processors.metadata]]
  tags = [ "project", "availability_zone"]
```

```diff
- cpu,hostname=localhost time_idle=42
+ cpu,hostname=localhost,project=webshop,availability_zone=primary time_idle=42
```
