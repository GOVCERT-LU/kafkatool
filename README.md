# Kafka GOVCERT.LU Tools

## Compile the kafka-tool binary

Set the GOPATH environment variable.

```bash
make
```

The binary is available as **./bin/kafkatool**.

The bash completion file is available as **./kafkatool_completion.sh**.

## Create a Debian package

### Add a new entry to the changelog

```bash
dch
```

### Create the source tarball

```bash
make tar
```

### Create the Debian package

```bash
dpkg-buildpackage -us -uc
```
