# Openstack Metadata Processor Plugin

Openstack Metadata processor plugin appends metadata gathered from Openstack
to metrics.
When this plugin is loaded for the first time, it takes metadata from http://169.254.169.254/openstack/latest/meta_data.json. 
Later, it adds configured parameters from plugin.conf as tags to all metrics.

## Usage
Let's download repository
```
git clone https://github.com/lawdt/metadata.git
cd metadata
```
And build plugin binary with name metadata
```
go build -o metadata cmd/main.go
```

You should be able to call plugin from telegraf now using execd processor plugin, add this to your telegraf.conf. 
Just replace paths with your real paths:
```
[[processors.execd]]
  command = ["path/to/metadata/binary", "--config", "path/to/plugin.conf"]
```
now, restart telegraf.

## Advanced Plugin Configuration
To change the composition of metrics added as tags just edit plugin.conf and restart telegraf:
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

## Full Example of plugin usage
###### Building example
```
git clone https://github.com/lawdt/metadata.git
cd metadata
go build -o metadata cmd/main.go
```
###### Installation example
```
mkdir -p /var/lib/telegraf/openstack
mv metadata /var/lib/telegraf/openstack/
mv plugin.conf /var/lib/telegraf/openstack/
chown -R telegraf:telegraf /var/lib/telegraf
```
###### Usage example
Edit /etc/telegraf/telegraf.conf
Append `project` and `availability_zone` to metrics tags:
```toml
[[processors.metadata]]
  tags = [ "project", "availability_zone"]
```
now restart telegraf
```
systemctl restart telegraf
```
And now all metrics contain the specified tags. 
For example:
```diff
- cpu,hostname=localhost time_idle=42
+ cpu,hostname=localhost,project=webshop,availability_zone=primary time_idle=42
```

