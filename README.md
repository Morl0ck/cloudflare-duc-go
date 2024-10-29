# cloudflare-duc-go

## Description

This is a simple Dynamic DNS Update Client (DUC) for Cloudflare written in Go.

## Usage

Edit the docker-compose.yml file and set the environment variables for the Cloudflare API key, zone name, and record name to keep up to date.

Then run the following command to start the container:

```bash
docker compose up -d
```

## Quick Setup

```bash
docker run -d \
  --name cloudflare-duc \
  --restart unless-stopped \
  -e CF_API_TOKEN=<YOUR_API_TOKEN> \
  -e CF_ZONE_NAME=<YOUR_ZONE_NAME> \
  -e CF_RECORD_NAME=<YOUR_RECORD_NAME> \
  -e UPDATE_INTERVAL=${UPDATE_INTERVAL:-5} \
  morl0ck/cloudflare-duc-go
```
