# kdump

kdump is a utility designed to facilitate the backup of Kubernetes API resources by dumping them into YAML files. This
tool is particularly useful for those who need a simple method to synchronize their Kubernetes resources with a version
control system like git, serving as a makeshift etcd to git synchronization mechanism.

## Overview

kdump leverages `kubectl` to enumerate all available API resources within the accessible Kubernetes contexts. It then
proceeds to retrieve each resource using `kubectl get <resource> -o yaml` and saves the output into corresponding YAML
files.

One important aspect of kdump is its ability to handle secrets. It can dump secrets as well, but only if the user
provides a specific encryption key for AES GCM encryption. This ensures that sensitive information is handled securely.

## Requirements

- `kubectl` must be installed and configured to interact with your Kubernetes clusters.
- `kubectl neat` is also required for cleaning up the manifests.

## Installation

kdump does not require a specific installation process. Ensure that you have `kubectl` and `kubectl neat` installed on
your system, and you can run kdump directly.

You can use `go install` to install kdump:

```shell
go install github.com/GiGurra/kdump@v1.27.7
```

## Usage

To use kdump, execute the command with the desired options:

```
kdump [global options] [arguments...]
```

### Options

- `--output-dir value, -o value`: Specify the output directory where the YAML files will be created.
- `--delete-previous-dir`: Set this flag to delete the previous output directory before creating a new one. The default
  is `false`.
- `--secrets-encryption-key value`: Provide a symmetric encryption key in hexadecimal format (64 lowercase characters)
  to enable secrets encryption with AES GCM.
- `--help, -h`: Display the help information.
- `--version, -v`: Print the version of kdump.

## Code Quality

The author of kdump describes the code quality as "hacky hack," indicating that while the tool is functional, it may not
adhere to best practices or be elegantly written. Additionally, the author notes that they are new to the Go programming
language and that the tool does not include tests.

kdump is a practical tool for developers and system administrators who need a straightforward solution for backing up
Kubernetes resources. It is not intended to be a polished or enterprise-grade application but rather a quick and
effective way to achieve a specific task.