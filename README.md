# EnvMe

Manage local development environments with ease. Expose your local services to the world.

## Usage

### Create a new service

```shell
Usage:
    envme create service <service-name> <image-name> [flags]

Flags:
    -h, --help          help for service
    -e, --env           environment variables
        --env-file      environment variables file
    -p, --expose        port to expose (format: <port>:<hostname>)
    -i, --interactive   interactive mode
```

### Create a new development environment

```shell
Usage:
    envme create development <environment-name> <directory> [flags]

Flags:
    -h, --help          help for development
    -e, --env           environment variables
        --env-file      environment variables file
    -p, --expose        port to expose (format: <port>:<hostname>)
    -i, --interactive   interactive mode
```

### Expose a service

```shell
Usage:
    envme expose <service-name> <port> <hostname> [flags]

Flags:
    -h, --help          help for expose
    -i, --interactive   interactive mode
```

### List services

```shell
Usage:
    envme list services [flags]

Flags:
    -h, --help            help for services
```
