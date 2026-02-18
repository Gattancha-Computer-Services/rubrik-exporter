# rubrik-exporter
Rubrik metrics exporter for Prometheus

A lightweight exporter that exposes Rubrik backup metrics in Prometheus format. Supports deployment as standalone binary or systemd service.

## Requirements

- **For building from source**: Go 1.25 or later
- **For running**: Linux/macOS/Windows with network access to your Rubrik cluster
- **For Prometheus integration**: Prometheus server to scrape metrics
- **For Grafana**: Grafana with Prometheus data source

## Installation & Usage

### Option 1: Standalone Binary

**Prerequisites:**
- Go 1.25 or later (only needed for building)

**Build the binary:**
```bash
git clone https://github.com/claranet/rubrik-exporter.git
cd rubrik-exporter
go build -o rubrik-exporter .
```

Or with make:
```bash
make build
```

**Run directly:**
```bash
./rubrik-exporter \
  -rubrik.url https://myrubrik.company.org \
  -rubrik.username "prometheus@local" \
  -rubrik.password 'MyPassword'
```

**Run with environment variables (using make):**
```bash
RUBRIK_URL=https://myrubrik.company.org \
RUBRIK_USER=prometheus@local \
RUBRIK_PASSWORD=MyPassword \
make run
```

**Clean up binary:**
```bash
make clean
```

**Prerequisites:**
- Go 1.25 or later (only needed for building)

**Build the binary:**
```bash
git clone https://github.com/claranet/rubrik-exporter.git
cd rubrik-exporter
go build -o rubrik-exporter .
```

Or with make:
```bash
make build
```

**Run directly:**
```bash
./rubrik-exporter \
  -rubrik.url https://myrubrik.company.org \
  -rubrik.username "prometheus@local" \
  -rubrik.password 'MyPassword'
```

**Run with environment variables (using make):**
```bash
RUBRIK_URL=https://myrubrik.company.org \
RUBRIK_USER=prometheus@local \
RUBRIK_PASSWORD=MyPassword \
make run
```

**Clean up binary:**
```bash
make clean
```

### Option 2: Install as Systemd Service (Linux)

**Automated installation:**
```bash
sudo ./install.sh
```

This script will:
1. Build the binary
2. Install to `/usr/bin/rubrik-exporter`
3. Create systemd unit at `/etc/systemd/system/rubrik-exporter.service`
4. Create config template at `/etc/default/rubrik-exporter`
5. Set up `prometheus` user and group

**After installation, configure credentials:**
```bash
sudo nano /etc/default/rubrik-exporter
```

Edit the file and uncomment/set one of the authentication methods:

**Option 1: Username/Password Authentication:**
```bash
RUBRIK_URL=https://myrubrik.company.org
RUBRIK_USER=prometheus@local
RUBRIK_PASSWORD=MyPassword
```

**Option 2: Service Account Authentication (recommended):**
```bash
RUBRIK_URL=https://myrubrik.company.org
RUBRIK_SERVICE_ACCOUNT=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...
```

**Start the service:**
```bash
# Enable on startup
sudo systemctl enable rubrik-exporter

# Start the service
sudo systemctl start rubrik-exporter

# Check status
sudo systemctl status rubrik-exporter

# View logs
sudo journalctl -u rubrik-exporter -f
```

**Manual installation (without install.sh):**
```bash
# Build and install binary
go build -o rubrik-exporter .
sudo install -m 755 rubrik-exporter /usr/bin/rubrik-exporter

# Copy systemd unit file
sudo cp rubrik-exporter.service /etc/systemd/system/

# Copy and edit config
sudo cp rubrik-exporter.env.example /etc/default/rubrik-exporter
sudo nano /etc/default/rubrik-exporter

# Reload and start
sudo systemctl daemon-reload
sudo systemctl enable rubrik-exporter
sudo systemctl start rubrik-exporter
```

## Configuration

All configuration is done via command-line flags or environment variables (when using systemd):

| Flag | Environment | Default | Required | Description |
|------|-------------|---------|----------|-------------|
| `-rubrik.url` | `RUBRIK_URL` | - | ✓ | Rubrik cluster URL (https://rubrik.example.com) |
| `-rubrik.username` | `RUBRIK_USER` | - | * | Rubrik API username (not required if using service account) |
| `-rubrik.password` | `RUBRIK_PASSWORD` | - | * | Rubrik API password (not required if using service account) |
| `-rubrik.service-account-client-id` | `RUBRIK_SERVICE_ACCOUNT_CLIENT_ID` | - | * | Rubrik service account client ID (alternative to username/password) |
| `-rubrik.service-account-client-secret` | `RUBRIK_SERVICE_ACCOUNT_CLIENT_SECRET` | - | * | Rubrik service account client secret (alternative to username/password) |
| `-listen-address` | `LISTEN_ADDRESS` | `:9477` | | HTTP binding address |

**Authentication Options:**

1. **Username/Password Authentication (default):**
```bash
./rubrik-exporter \
  -rubrik.url https://rubrik.example.com \
  -rubrik.username prometheus@local \
  -rubrik.password secure_password
```

2. **Service Account Authentication (recommended - OAuth2 client credentials):**
```bash
./rubrik-exporter \
  -rubrik.url https://rubrik.example.com \
  -rubrik.service-account-client-id "your-client-id" \
  -rubrik.service-account-client-secret "your-client-secret"
```

**Example:**

## Prometheus Integration

Add the exporter to your `prometheus.yaml`:

```yaml
scrape_configs:
  - job_name: 'rubrik-exporter'
    static_configs:
      - targets: ['localhost:9477']
    scrape_interval: 30s
    scrape_timeout: 10s
```

Multiple Rubrik clusters:
```yaml
scrape_configs:
  - job_name: 'rubrik-dc1'
    static_configs:
      - targets: ['rubrik-exporter-dc1:9477']
  
  - job_name: 'rubrik-dc2'
    static_configs:
      - targets: ['rubrik-exporter-dc2:9477']
```

Then reload Prometheus to pick up the new targets.

## Grafana Integration

1. Add Prometheus as a data source in Grafana (if not already configured)
2. Create dashboards querying `rubrik_*` metrics
3. Example queries:
   - `rubrik_count_streams` - Number of backup streams
   - `rubrik_system_storage_size` - Total storage capacity
   - `rubrik_system_storage_used` - Storage used
   - `rubrik_vm_protected` - VM protection status

## Command Line Options

```
Usage of rubrik-exporter:
  -listen-address string
        HTTP address to listen on (default ":9477")
  -rubrik.password string
        Rubrik API password (required)
  -rubrik.url string
        Rubrik cluster URL, e.g., https://rubrik.example.com (required)
  -rubrik.username string
        Rubrik API username (required)
```

## Make Commands

Convenience commands available via make:

```bash
make build           # Build the binary
make run             # Build and run (requires environment variables)
make clean           # Remove the binary
make deps            # Download Go dependencies
make docker-build    # Build Docker image locally
make help            # Show all available commands
```

Example usage:
```bash
# Build only
make build

# Run with environment variables
RUBRIK_URL=https://rubrik.example.com \
RUBRIK_USER=prometheus@local \
RUBRIK_PASSWORD=MyPassword \
make run

# Clean up
make clean
```

Exported Metrics
==================

        # HELP rubrik_archive_location_status Archive Loction Status - 1: Active, 0: Inactive
        # TYPE rubrik_archive_location_status gauge
        rubrik_archive_location_status{bucket="archive",name="NFS:archive",target="<ip-address>"} 1
        # HELP rubrik_archive_storage_archived_fileset ...
        # TYPE rubrik_archive_storage_archived_fileset gauge
        rubrik_archive_storage_archived_fileset{name="NFS:archive",target="<ip-address>",type="fileset"} 0
        rubrik_archive_storage_archived_fileset{name="NFS:archive",target="<ip-address>",type="linux"} 0
        rubrik_archive_storage_archived_fileset{name="NFS:archive",target="<ip-address>",type="share"} 0
        rubrik_archive_storage_archived_fileset{name="NFS:archive",target="<ip-address>",type="windows"} 0
        # HELP rubrik_archive_storage_archived_vm ...
        # TYPE rubrik_archive_storage_archived_vm gauge
        rubrik_archive_storage_archived_vm{name="NFS:archive",target="<ip-address>",type="hyperv"} 0
        rubrik_archive_storage_archived_vm{name="NFS:archive",target="<ip-address>",type="nutanix"} 0
        rubrik_archive_storage_archived_vm{name="NFS:archive",target="<ip-address>",type="vmware"} 0
        # HELP rubrik_archive_storage_data_archived ...
        # TYPE rubrik_archive_storage_data_archived gauge
        rubrik_archive_storage_data_archived{name="NFS:archive",target="<ip-address>"} 0
        # HELP rubrik_archive_storage_data_downloaded ...
        # TYPE rubrik_archive_storage_data_downloaded gauge
        rubrik_archive_storage_data_downloaded{name="NFS:archive",target="<ip-address>"} 0
        # HELP rubrik_count_nodes Count Rubrik Nodes in a Brick
        # TYPE rubrik_count_nodes gauge
        rubrik_count_nodes{brik="<brik-id>"} 4
        # HELP rubrik_count_streams Count Rubrik Backup Streams
        # TYPE rubrik_count_streams gauge
        rubrik_count_streams 0
        # HELP rubrik_node_io_read Node Read IO per second
        # TYPE rubrik_node_io_read gauge
        rubrik_node_io_read{node="<node-id>"} 281
        # HELP rubrik_node_io_write Node Write IO per second
        # TYPE rubrik_node_io_write gauge
        rubrik_node_io_write{node="<node-id>"} 1676
        # HELP rubrik_node_network_received Node Network Byte received
        # TYPE rubrik_node_network_received gauge
        rubrik_node_network_received{node="<node-id>"} 2.3782571e+07
        # HELP rubrik_node_network_transmitted Node Network Byte transmitted
        # TYPE rubrik_node_network_transmitted gauge
        rubrik_node_network_transmitted{node="<node-id>"} 2.0558578e+07
        # HELP rubrik_node_throughput_read Node Read Throughput per second
        # TYPE rubrik_node_throughput_read gauge
        rubrik_node_throughput_read{node="<node-id>"} 3.0397071e+07
        # HELP rubrik_node_throughput_write Node Write Throughput per second
        # TYPE rubrik_node_throughput_write gauge
        rubrik_node_throughput_write{node="<node-id>"} 2.5802489e+07
        # HELP rubrik_system_storage_available ...
        # TYPE rubrik_system_storage_available gauge
        rubrik_system_storage_available 7.0788234461184e+13
        # HELP rubrik_system_storage_live_mount ...
        # TYPE rubrik_system_storage_live_mount gauge
        rubrik_system_storage_live_mount 0
        # HELP rubrik_system_storage_miscellaneous ...
        # TYPE rubrik_system_storage_miscellaneous gauge
        rubrik_system_storage_miscellaneous 4.48470796529e+12
        # HELP rubrik_system_storage_snapshot ...
        # TYPE rubrik_system_storage_snapshot gauge
        rubrik_system_storage_snapshot 4.5390468081302e+13
        # HELP rubrik_system_storage_total ...
        # TYPE rubrik_system_storage_total gauge
        rubrik_system_storage_total 1.20663410507776e+14
        # HELP rubrik_system_storage_used ...
        # TYPE rubrik_system_storage_used gauge
        rubrik_system_storage_used 4.8594657722368e+13
        # HELP rubrik_vm_consumed_exclusive_bytes ...
        # TYPE rubrik_vm_consumed_exclusive_bytes gauge
        rubrik_vm_consumed_exclusive_bytes{vmname="<vm-name>""} 0
        # HELP rubrik_vm_consumed_index_storage_bytes ...
        # TYPE rubrik_vm_consumed_index_storage_bytes gauge
        rubrik_vm_consumed_index_storage_bytes{vmname="<vm-name>""} 0
        # HELP rubrik_vm_consumed_ingested_bytes ...
        # TYPE rubrik_vm_consumed_ingested_bytes gauge
        rubrik_vm_consumed_ingested_bytes{vmname="<vm-name>"} 0
        # HELP rubrik_vm_consumed_logical_bytes ...
        # TYPE rubrik_vm_consumed_logical_bytes gauge
        rubrik_vm_consumed_logical_bytes{vmname="<vm-name>""} 0
        # HELP rubrik_vm_consumed_shared_physical_bytes ...
        # TYPE rubrik_vm_consumed_shared_physical_bytes gauge
        rubrik_vm_consumed_shared_physical_bytes{vmname="<vm-name>"} 0
        # HELP rubrik_vm_protected ...
        # TYPE rubrik_vm_protected gauge
        rubrik_vm_protected{vmname="<vm-name>"} 0|1
