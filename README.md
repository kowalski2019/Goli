## Get started

This program was implemented in order to have a kind of deployment adapted to people with a server that is not in an **Azure**, **AWS** etc. environment.

It is a very simple REST-API implemented in the **Go** language, which allows to perform actions in a server from a simple request.

## What do you need to have before using this program (prerequisites)?
You need to have a server and have **Go** installed in it. *[How to install Golang](https://go.dev/doc/install)*.

## How to use it ? 
This program works most efficiently as a UNIX/Linux service (will be detailed a bit below). The most important thing is that you have a folder __/csmkactionhelper/config/__ in your server and in this folder a file called **config.toml**, because without this file the application will not work.

The content of the configuration file should look like this:

```
[constants]
auth_key = "dummy_key"
```

### Build the Programm
**Attention!** The program listens by default on port 8125, you can change this before compiling it.

```
go get // recover dependencies if some are missing

go build main.go // build the programm and create an output named "main"
```

Once the program is compiled, an executable named **main** will be visible in the folder.

### Create a UNIX/Linux Service for the program

Create a file named **action-helper.service** with the following content.
This is what a UNIX/Linux service file might look like:

```
[Unit]
Description=Github action-runner Helper for Deployment
After=network.target

[Service]
Type=simple
WorkingDirectory=<path to the directory where the main file is>
ExecStart=<full path to the main file compiled beofre>
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```
This is what the values of variables **WorkingDirectory** and **ExecStart** might look like:

```
WorkingDirectory=/home/dummy/Deployment/action_helper
ExecStart=/home/dummy/Deployment/action_helper/main
```


### Final steps

```
sudo cp action-helper.service /etc/system/systemd/

sudo systemctl enable --now action-helper.service

```
After these steps your "action helper" is ready for use. Enjoy ;-)

### Make a Test

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


