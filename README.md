# InfraSonar Client

1. Download the client
**1. Download the latest installer:**

- [Linux (amd64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.1/infrasonar-cli-linux-amd64-1.0.0.tar.gz)
- [Linux (arm64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.1/infrasonar-cli-linux-arm64-1.0.0.tar.gz)
- [Darwin (amd64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.1/infrasonar-cli-darwin-amd64-1.0.0.tar.gz)
- [Darwin (arm64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.1/infrasonar-cli-darwin-arm64-1.0.0.tar.gz)
- [Windows (amd64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.1/infrasonar-cli-windows-amd64-1.0.0.zip)
- [Windows (arm64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.1/infrasonar-cli-windows-arm64-1.0.0.zip)

> If your platform is not listed above, refer to the [build from source](#build-from-source) section for instructions.

**2. Extract the contents of the archive using a tool like `tar`. Here's an example for Linux (amd64):**
```bash
tar -xzvf infrasonar-cli-linux-amd64-1.0.0.tar.gz
```

2. Install
The following command will install infrasonar in path and enables bash completion if supported by the OS.

```bash
sudo ./infrasonar install
```

3. Create a new configuration

```bash
infrasonar config new
```


### Build
```bash
CGO_ENABLED=0 go build -o appliance-installer
```