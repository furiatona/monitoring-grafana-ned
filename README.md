# NED - Node Exporter Deployment

**NED** is a lightweight CLI tool designed to simplify the deployment, update, and management of Prometheus Node Exporter instances across your infrastructure.

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
cp env.sampe .env
./scripts/build.sh
```

This will generate an executable named `ned`.

### 3. Run NED

```bash
./build/ned
```

- Fetches DC patterns from Consul.
- Presents an interactive menu.
- Prompts for confirmation before actions.

---

**License:** [MIT](./LICENSE)  
**Author:** [furiatona](https://github.com/furiatona)