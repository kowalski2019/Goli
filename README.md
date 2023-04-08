# Goli Action Helper
## Get started
This program was implemented in order to have a kind of deployment adapted to people with a server that is not in an **Azure**, **AWS** etc. environment.

It is a very simple REST-API implemented in the **Go** language, which allows to perform actions in a server from a simple request.

## What do you need to have before using this program (prerequisites)?
You need to have a server and have ubuntu running on it.

## Configure your server to use Goli
You have juste to run the script **setup.sh**

```
./setup.sh
```

After this your "Goli Action Helper" is ready for use. Enjoy ;-)


## Make a Test

```
curl -d 'auth_key=dummy_key' \
        -d 'name=hello_container' \
        -d 'image=dommy_image:latest' \
        -d 'network=host' \
        -d 'port_ex=9000' \
        -d 'port_in=80' \
        -d 'volume_ex=/var/log' \
        -d 'volume_in=/var/log' \
        -d 'v_map=yes' \

        -X POST \
        -H 'Content-Type: application/xwww-form-urlencoded' http://127.0.0.1:<host_port>/api/v1/docker/ps

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

