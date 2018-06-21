# Imhio Task (CFGM)

# Dependencies

- [Dep](https://github.com/golang/dep)
- Make

# Env

| Name             | Description                                                            | Default |
| ---------------- | ---------------------------------------------------------------------- | ------- |
| CFGM_HTTP_ADDR   | Service bind address                                                   | :9000   |
| CFGM_DB_CONN_STR | PostgreSQL connection URL: postgres://username:passwd@host:port/dbname |         |

# Build binary

```bash
make dep build
```

# Docker build image

```bash
make docker pk="$(cat ~/.ssh/id_rsa)" c=cfgm
```

this is command build `cfgm` docker image

# Migration

Migration command run on service start. Supported commands are:

- up - runs all available migrations.
- down - reverts last migration.
- reset - reverts all migrations.

Default command is: `up`.
