Deploying Apocalyptica
---
> **IN PROGRESS! DO NOT USE!!!** Still testing out some things.

> :memo: current version of the document is for Docker savvyheads! We will include many different deployment profiles soon.

## Docker (Quickstart & Custom)

###  Requirements
If you want to deploy your own instance of Apocalyptica, you will need:

* A Discord Account and [a registered app](register.md). **Please register app before starting a deploy!** You need `DISCORD_TOKEN` and `CLIENT_ID` in your `.env` file to deploy successfully.
* Docker (or other container runtime engine, though it has only been tested with Docker)
* Instructions are written assuming OSX/Linux access, though it should be easy to adapt to Windows.
* Web server needs access to port 8080.

### Docker Quickstart

> :information_source: The Docker quickstart is when you want to get going quickly and only want/require the default games and playbooks installed.

To build with defaults:

* fill out `env_discord.sample` with `DISCORD_TOKEN` and `CLIENT_ID` values, then copy that file to `.env` under the root file of the repo.
* from the root of the repo, run `bin/docker_quickstart`, which is really a wrapper script for:

`docker run --name apocalyptica --rm --env-file -p8080:8080 ../.env -d gamefiend/apocalyptica`

### Docker Custom

> :information_source: The Docker custom option is when you want to add your own custom games and playbooks to Apocalyptica.  

To build with custom information:

* fill out `env_discord.sample` with `DISCORD_TOKEN` and `CLIENT_ID` values, then copy that file to `.env` under the root file of the repo.
* load your custom games and playbooks into `data/games`
* from the root of the repo, run `bin/docker_custom`, which is really a wrapper script for:

```
docker build -t apocalyptica-custom \
&& docker run --name apocalyptica --rm --env-file -p8080:8080 ../.env -d apocalyptica-custom
```

This builds a custom container image that uses your custom info instead of using the generally available docker image.


## Deploying to Kubernetes 

> :memo: COMING SOON (REALLY SOON)


## Using Go Binary

> :memo: COMING SOON
