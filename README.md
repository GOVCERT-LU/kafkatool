[![codebeat badge](https://codebeat.co/badges/5fda5ffb-7433-4821-95dc-1b1b6b260ce4)](https://codebeat.co/projects/github-com-govcert-lu-kafkatool-master)


# Kafka GOVCERT.LU Tools

The kafkatool leverages the Sarama Go library for Apache Kafka to propose an easy to use command line tool to facilitate the daily management of a Kafka cluster.

## Compile the kafka-tool binary

```bash
make
```

The binary is available as **build/kafkatool**.

The bash completion file is available as **./kafkatool_completion.sh**.

## Create a Debian package

### Add a new entry to the changelog

```bash
# with a specific version
dch --urgency medium --distribution stable --newversion "2.4.6"

# only bump the package
dch --urgency medium --distribution stable
```

### Create the source tarball

```bash
make tar
```

### Create the Debian package

```bash
dpkg-buildpackage -b --no-sign
```
