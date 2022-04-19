# Openstack Metadata Processor Plugin

Openstack Metadata processor plugin appends metadata gathered from Openstack
to metrics.

## Building
```bash
git clone https://github.com/lawdt/metadata.git
cd metadata
go build -o metadata cmd/main.go
```

## Installation
```bash
mkdir -p /var/lib/telegraf/openstack
mv metadata /var/lib/telegraf/openstack/
mv plugin.conf /var/lib/telegraf/openstack/
chown -R telegraf:telegraf /var/lib/telegraf
```

## Plugin Configuration
edit /var/lib/telegraf/openstack/plugin.conf:

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

## Usage
You should be able to call this from telegraf now using execd:
```
[[processors.execd]]
  command = ["/var/lib/telegraf/metadata", "--config", "/var/lib/telegraf/plugin.conf"]
```