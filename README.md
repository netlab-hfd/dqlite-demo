# DQLite Demo Application
This repository contains a containerized version of the [go-dqlite](https://github.com/canonical/go-dqlite) demo application.
It showcases DQLite’s clustering and failover behavior through a simple key-value store, packaged with Docker Compose for convenience.

- DQLite is a lightweight, distributed SQLite database with Raft-based replication for fault tolerance.

## Getting started

### Prerequisites
To run this demo, you need the following requirements on your machine:

- Git must be installed
- [Docker](https://www.docker.com/get-started/) or [Podman](https://podman.io/docs/installation) must be installed


### Setup
To run the demo, clone the repository by running 
```bash
git clone https://github.com/netlab-hfd/dqlite-demo.git
```

Then, open a console terminal in the folder you just cloned and start the application

```bash
docker compose up
```

This starts three containers running the sample application with a clustered DQLite database. The sample application 
exposes an HTTP API which can be used to set and read key-value pairs. It is reachable at ports 8001, 8002 and 8003.

### Usage
For example, you can set a value with this curl:

```bash
curl -X PUT -d value localhost:8001/mykey
```

You can then read your set value with this curl command:
```bash
curl localhost:8001/mykey
```

Because the databases run as a cluster, your changes get synced to the other instances automatically. You can test this
by querying the other two containers with curl:

```bash
curl localhost:8003/mykey
```

You can also observe failover behavior by killing a container from the cluster and then observing that the other nodes
are still working and retain your data. For example, you could terminate the container at 8001 and then check that the 
service at 8003 is still working:

```bash
docker container stop dqlite1
curl localhost:8003/mykey
```

You can play around a little bit more by killing and restarting containers and testing how the application behaves.
Look at the output of the compose command for feedback on what DQLite is currently doing.

## Why this demo?
The official Canonical demo is not containerized and assumes a local Linux setup (only packaged for Ubuntu).
This makes it harder for others to experiment with DQLite.

By wrapping the demo in containers, this project lowers the barrier to entry and enables quick experimentation without local setup headaches.