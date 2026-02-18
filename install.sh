#!/bin/bash

# Rubrik Exporter Installation Script
# This script installs the rubrik-exporter binary and creates a systemd service

set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

echo "=== Rubrik Exporter Installation ==="

# Build the binary
echo "[1/5] Building rubrik-exporter..."
cd "$(dirname "$0")"
go build -o rubrik-exporter .

if [ ! -f rubrik-exporter ]; then
    echo "Error: Failed to build rubrik-exporter"
    exit 1
fi

# Copy binary to /usr/bin
echo "[2/5] Installing binary to /usr/bin/rubrik-exporter..."
install -m 755 rubrik-exporter /usr/bin/rubrik-exporter

# Create systemd unit file
echo "[3/5] Creating systemd unit file..."
cat > /etc/systemd/system/rubrik-exporter.service << 'EOF'
[Unit]
Description=Rubrik Exporter
Documentation=https://github.com/Gattancha-Computer-Services/rubrik-exporter
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=prometheus
Group=prometheus
ExecStart=/usr/bin/rubrik-exporter \
  -rubrik.url=${RUBRIK_URL} \
  -rubrik.username=${RUBRIK_USER} \
  -rubrik.password=${RUBRIK_PASSWORD} \
  -rubrik.service-account-client-id=${RUBRIK_SERVICE_ACCOUNT_CLIENT_ID} \
  -rubrik.service-account-client-secret=${RUBRIK_SERVICE_ACCOUNT_CLIENT_SECRET} \
  -listen-address=:9477

# Environment file - create /etc/default/rubrik-exporter with your config
EnvironmentFile=-/etc/default/rubrik-exporter

StandardOutput=journal
StandardError=journal

# Auto-restart on failure
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

# Create environment file template
echo "[4/5] Creating configuration template..."
if [ ! -f /etc/default/rubrik-exporter ]; then
    cat > /etc/default/rubrik-exporter << 'EOF'
# Rubrik Exporter Configuration
# Uncomment and update the following variables:

# REQUIRED: Rubrik URL (e.g., https://rubrik.example.com)
# RUBRIK_URL=https://rubrik.example.com

# REQUIRED: Rubrik API Username
# RUBRIK_USER=prometheus@local

# REQUIRED: Rubrik API Password
# RUBRIK_PASSWORD=your-password-here

# OPTIONAL: Listen address (default: :9477)
# LISTEN_ADDRESS=:9477
EOF
    chmod 600 /etc/default/rubrik-exporter
    echo "Configuration template created at /etc/default/rubrik-exporter"
else
    echo "Configuration file already exists at /etc/default/rubrik-exporter"
fi

# Create prometheus user/group if doesn't exist
echo "[5/5] Setting up prometheus user..."
if ! id "prometheus" &>/dev/null; then
    useradd --no-create-home --shell /bin/false prometheus 2>/dev/null || true
fi

# Reload systemd
systemctl daemon-reload

echo ""
echo "=== Installation Complete ==="
echo ""
echo "Next steps:"
echo "1. Edit /etc/default/rubrik-exporter and set your Rubrik credentials"
echo "2. Enable the service: sudo systemctl enable rubrik-exporter"
echo "3. Start the service: sudo systemctl start rubrik-exporter"
echo "4. Check status: sudo systemctl status rubrik-exporter"
echo "5. View logs: sudo journalctl -u rubrik-exporter -f"
echo ""
echo "The exporter will be available at http://localhost:9477/metrics"
