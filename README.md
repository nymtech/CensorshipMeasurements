# Censorship Measurement of the Nym network

## OONI-Probe tests for the Nym network

This repository introduces a new type of Go test for Nym network censorship. Our tests are performed using the [OONI probe](https://ooni.org/).
* First, the test evaluates the connectivity to the to the validator API, which is a necessary step to allow users access to the Nym network, as it allows the Nym client to retrieve vital information such as the list of active relay nodes,
gateways, their topology, and the necessary credentials for network access.
* After successfully fetching the list of available gateways in the initial test,
a second test is conducted to determine the reachability of these gateways.
Since the Nym network topology mandates that packets are routed through
a gateway before entering the mixnet, being able to connect to at least one
gateway is crucial for network access. To assess this, the test attempts to
establish connections with each gateway individually.

## Requirements
To be able able to run the tests, you should have `docker`  and `docker-compose` installed. You can follow [this procedure](https://docs.docker.com/desktop) to install docker desktop which contains `docker` and `docker-compose` and a GUI to manage docker. Once installed and running, you can run the docker tests.

## How to run the tests ?

From the command line, enter the `Nym_Censorship_Measurements` folder by running the command
```
  cd Nym_Censorship_Measurements
```

By using `docker-compose`, it will build the image, install the requirements inside the image, compile the binaries for the test and run the test. The result will be in the directory `./results` in file `report.jsonl`.

```bash
docker-compose up
```

To force to build again each time, you can do

```bash
docker-compose up --build
```

Please share your `report.jsonl` with us.

## Licensing and copyright information
This program is available as open source under the terms of the Apache 2.0 license.

## Notes
The environment variable `OONI_NYMVALIDATORURL` can be used to modify the validator server to reach for the tests. By default, it is `"https://validators.nymtech.net"`.

```bash
OONI_NYMVALIDATORURL="https://validators.nymtech.net" docker-compose up
```
