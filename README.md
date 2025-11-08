# Goli Action Helper

![Goli](img/GOLI.jpg)

## Get started
This program was implemented in order to have a kind of deployment adapted to people with a server that is not in an **Azure**, **AWS** etc. environment.

It is a very simple REST-API implemented in the **Go** language, which allows to perform actions in a server from a simple request.

## What do you need to have before using this program (prerequisites)?
You need to have a server and have ubuntu running on it.

## Configure your server to use Goli

### Create Domain for Goli


### Create a GitHub personal access token
![Follow this link](https://github.com/settings/tokens/new)

- Select the __read:packages__ scope to download container images and read their metadata.
- Select the __write:packages__ scope to download and upload container images and read and write their metadata.

### Run the setup script
```
chmod +x ./setup.sh && ./setup.sh
```

After this your "Goli Action Helper" is ready for use. Enjoy ;-)


## Make a Test

```
curl -X POST  \
       --header "Authorization: DeepL-Auth-Key dummy_key" \
       --header "Content-Type: application/json" \
       --data '{ 
                "name": "hello_container", 
                "image": "dummy_image:latest", 
                "network": "host", 
                "port_ex": "9000", 
                "port_in": "80",
                "v_map": false,
                "volume_ex": "/var/log", 
                "volume_in": "/var/log",
                "opts": "" }' \
        http://127.0.0.1:<host_port>/api/v1/docker/ps

```

## Available Routers

```
        POST: /api/v1/docker/container/start 
        POST: /api/v1/docker/container/stop
	POST: /api/v1/docker/container/rm
	POST: /api/v1/docker/container/run
	POST: /api/v1/docker/container/pause
	POST: /api/v1/docker/container/unpause
	POST: /api/v1/docker/container/inspect
	POST: /api/v1/docker/container/logs

	POST: /api/v1/docker/image/rm
	POST: /api/v1/docker/image/pull

	POST: /api/v1/docker/ps
	POST: /api/v1/docker/images
```

## Tipps for your Pipeline



