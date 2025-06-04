# NED - Node Exporter Deployment

**NED** is a lightweight CLI tool designed to simplify the deployment, update, and management of Prometheus Node Exporter instances across your infrastructure.

> ⚠️ **Caution:** Credentials are embedded in the binary for user convenience during manual exporter installation. This approach is used only when automation is not possible due to policy restrictions. The tool is intended for use in secure environments. Use at your own risk.
---

## Features

- **Easy Deployment:** Quickly roll out Node Exporter across multiple host types.
- **Automated Updates:** Seamlessly update existing deployments.
- **Consul Integration:** Fetch DC patterns and manage infrastructure with ease.
- **User-Friendly CLI:** Intuitive interface for streamlined operations.

---

## Getting Started

### 1. Clone the Repository
```bash
git clone <git_link>
cd monitoring-grafana-ned/
```

### 2. Build the Binary
Copy `env.sample` to `.env`, then build the binary:
```bash
cp env.sample .env
./scripts/build.sh
```

This will generate an executable named `ned`.

### 3. Run NED
```bash
./build/ned
```

### 4. Transfer the binary to the file server
To make the binary available for download, transfer it to your secured file server (for example, to `/data/filehost/html` if your file host is at `https://example.com/ui/filehost`):
```bash
./scripts/transfer.sh
```

### 5. Run on Target Server
To run NED on your server, use the following command:
```bash
curl -u user:pass -o ned https://example.com/ui/filehost/ned && chmod +x ned && ./ned && rm -f ned
```

The script will:
- Fetch DC patterns from Consul.
- Present an interactive menu.
- Prompt for confirmation before performing actions.

---

**License:** [MIT](./LICENSE)  
**Author:** [furiatona](https://github.com/furiatona)