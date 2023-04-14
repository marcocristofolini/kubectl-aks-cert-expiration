# kubectl-aks-cert-expiration

A kubectl krew plugin that checks the expiration of Kubernetes API server TLS certificates for AKS clusters across all subscriptions in an Azure account.

## Prerequisites

- Go 1.16 or higher
- An Azure account with appropriate permissions to list subscriptions and AKS clusters
- kubectl and krew installed

## Installation

### Using Krew

1. Clone the repository:

```bash
git clone https://github.com/marcocristofolini/kubectl-aks-cert-expiration.git
```


1. Install the plugin using Krew:

```bash

kubectl krew install --manifest=./kubectl-aks-cert-expiration/plugin.yaml
```


### Manual Installation
1. Clone the repository:

```bash

git clone https://github.com/marcocristofolini/kubectl-aks-cert-expiration.git
```


1. Change to the project directory and build the plugin:

```bash

cd kubectl-aks-cert-expiration
go build -o kubectl-aks_cert_expiration main.go
```


1. Add the compiled binary to your `PATH`:

```bash

export PATH=$PATH:/path/to/your/kubectl-aks-cert-expiration
```


## Usage

To use the plugin, set the following environment variables with your Azure account's credentials:
- `AZURE_CLIENT_ID`: Your Azure Client ID
- `AZURE_CLIENT_SECRET`: Your Azure Client Secret
- `AZURE_TENANT_ID`: Your Azure Tenant ID

Run the plugin:

```bash

kubectl aks-cert-expiration --client-id=<your-client-id> --client-secret=<your-client-secret> --tenant-id=<your-tenant-id>
```



The plugin will authenticate with Azure, iterate through each subscription, and check the certificate expiration for all AKS clusters by connecting to the Kubernetes API server directly and inspecting the TLS certificate.
## Contributing

If you'd like to contribute to the project, please submit a pull request or open an issue on the [GitHub repository](https://github.com/marcocristofolini/kubectl-aks-cert-expiration) .
## License

MIT License
