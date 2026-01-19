[![CI](https://github.com/infrasonar/infrasonar-cli/workflows/CI/badge.svg)](https://github.com/infrasonar/infrasonar-cli/actions)
[![Release Version](https://img.shields.io/github/release/infrasonar/infrasonar-cli)](https://github.com/infrasonar/infrasonar-cli/releases)


# InfraSonar Client

The InfraSonar client is a command-line application which can be used to manage assets for a container. The tool has two main features. One is to read all assets from a container to YAML or JSON output. Zones labels and collectors are included. The other feature of this tool is to apply a YAML or JSON file to InfraSonar. These two features combined allow you to easily add new assets as well as managing existing assets for a container.

**1. Download the latest version:**

- [Linux (amd64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.3/infrasonar-linux-amd64-1.0.3.tar.gz)
- [Linux (arm64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.3/infrasonar-linux-arm64-1.0.3.tar.gz)
- [Darwin (amd64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.3/infrasonar-darwin-amd64-1.0.3.tar.gz)
- [Darwin (arm64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.3/infrasonar-darwin-arm64-1.0.3.tar.gz)
- [Windows (amd64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.3/infrasonar-windows-amd64-1.0.3.zip)
- [Windows (arm64)](https://github.com/infrasonar/infrasonar-cli/releases/download/v1.0.3/infrasonar-windows-arm64-1.0.3.zip)

> If your platform is not listed above, refer to the [build from source](#build-from-source) section for instructions.

**2. Extract the contents of the archive using a tool like `tar`. Here's an example for Linux (amd64):**
```bash
tar -xzvf infrasonar-linux-amd64-1.0.3.tar.gz
```

**3. Install:**

The following command will install infrasonar in path and enables bash completion if supported by the OS.

```bash
sudo ./infrasonar install
```

**Mac OS specific:**

If you get the error message:
> Apple could not verify “infrasonar” is free of malware that may harm your Mac or compromise your privacy.

Ensure you downloaded the binary from a trusted source (https://github.com/infrasonar/infrasonar-cli/)

Manually remove the "quarantine" flag that macOS attaches to the file. Open your Terminal and run:

```bash
xattr -d com.apple.quarantine ./infrasonar
```

**4. Create a new configuration:**

```bash
infrasonar config new
```

Finally, give your configuration a name and provide a token

```
Name: foo
Token: ***********
```

### Build from source
Clone this repository and make sure [Go](https://golang.google.cn) is installed.

```bash
CGO_ENABLED=0 go build -o infrasonar
```
