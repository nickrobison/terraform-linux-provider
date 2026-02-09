# terraform-linux-provider
Terraform Provider for Bare-Metal Linux

## Features

This provider allows you to manage Linux system resources via Terraform using DBus interfaces.

### Supported Resources

#### Firewalld
- **Zones**: Create and manage firewalld zones with custom configurations
  - Set zone descriptions and default targets
  - Query existing zones via data sources
- **Rules**: Manage firewall rules including:
  - Rich rules for complex firewall configurations
  - Port rules (TCP/UDP)
  - Service rules

#### ZFS
- **ZPools**: Create and manage ZFS storage pools

## Usage

### Firewall Zone Management

```hcl
# Create a firewall zone
resource "linux_firewall_zone" "dmz" {
  name        = "dmz"
  description = "DMZ zone for public-facing services"
  target      = "default"
}

# Query existing zones
data "linux_firewall_zones" "all" {}
```

### Firewall Rule Management

```hcl
# Add a rich rule
resource "linux_firewall_rule" "allow_admin_ssh" {
  zone      = "public"
  rule_type = "rich"
  rule      = "rule family=ipv4 source address=10.0.0.0/8 service name=ssh accept"
}

# Add a port rule
resource "linux_firewall_rule" "allow_http" {
  zone      = "public"
  rule_type = "port"
  port      = "80"
  protocol  = "tcp"
}

# Add a service rule
resource "linux_firewall_rule" "allow_https" {
  zone      = "public"
  rule_type = "service"
  service   = "https"
}
```

## Requirements

- Firewalld must be installed and running on the target system
- The provider communicates with system services via DBus
- Appropriate permissions are required to manage firewall and ZFS resources
