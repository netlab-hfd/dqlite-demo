# DQLite Demo Application
This Repository contains the example application from [go-dqlite](https://github.com/canonical/go-dqlite), packaged
as a container image with a compose file for convenience.

## Getting started
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

## Why create this when a demo already exists?
The approved demo application from Canonical is not containerized and expects that everything is installed locally.
DQLite is only working on Linux as of now and only packaged for Ubuntu, making it hard to play around with it. By
wrapping it in containers, the hurdle to get started is lowered significantly and a lot of headaches for the users are 
saved.