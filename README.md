# Introduction
This project is an intro to ScyllaDB in the context of Big Data Architecture

## Project Structure
The project consists of the follow directories/files:
- cassandra (docker-compose for cassandra deployment)
- scylla (docker-compose for cassandra deployment)
- kubernetes
- demo-scripts (basic CQL scripts)
- loader (golang CQL client implementation to load data to cassandra/scylladb)
- terraform (terraform modules to deploy on AWS EC2 the demo)
- bootstrap.sh (setup kernel params for fs.aio-max-nr)

## Usage

The project has some lfs files that are required to be downloaded before you can run the `loader` golang project.
- ./demo-scripts 

Please ensure that you have lfs installed before cloning the project.
```sh
git lfs --version
```

If you have already cloned the project, delete the following files and check them out to get a copy of the full file from lfs:
```sh
# delete the files
rm loader/airline.csv.sample
rm scylla-introduction-slides.pdf

# download from lfs
git checkout loader/airline.csv.sample
git checkout scylla-introduction-slides.pdf
```

Starts Cassandra and ScyllaDB locally with Docker:
```sh
# deploy scylla locally
docker compose -f scylla/docker-compose.yml up

# (optional) deploy cassandra locally
docker compose -f cassandra/docker-compose.yml up

# view the containers
docker compose ps

# clean up scylla
docker compose -f scylla/docker-compose.yml down

# (optional) clean up cassandra
docker compose -f cassandra/docker-compose.yml down
```

