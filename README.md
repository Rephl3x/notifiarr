<img src="https://raw.githubusercontent.com/wiki/Go-Lift-TV/discordnotifier-client/images/golift-discordnotifier-client.png">

This is the unified client for [DiscordNotifier.com](https://discordnotifier.com).
The client enables content requests from Media Bot in your Discord Server.

In the future it will do even more!

## Installation

### Docker

This project builds automatically in [Docker Cloud](https://hub.docker.com/r/golift/discordnotifier-client)
and creates [ready-to-use multi-architecture images](https://hub.docker.com/r/golift/discordnotifier-client/tags).
The `latest` tag is always a tagged release in GitHub. The `main` tag corresponds
to the `main` branch in GitHub and may be broken.

#### Docker Config File

-   Copy the [example config file](https://github.com/Go-Lift-TV/discordnotifier-client/blob/main/examples/dnclient.conf.example) from this repo.
-   Then grab the image from docker hub and run it using an overlay for the config file.

```shell
docker pull golift/discordnotifier-client
docker run -d \
  -v /your/config/dnclient.conf:/etc/discordnotifier-client/dnclient.conf \
  golift/discordnotifier-client
docker logs <container id from docker run>
```

#### Docker Env Variables

-   Instead of, or in addition to a config file, you may configure the docker
    container with environment variables.
-   Any variable not provided takes the default.
-   Must provide an API key from discordnotifier.com.
-   Must provide URL and API key for Sonarr or Radarr or Readarr or any combination.
-   You may provide multiple sonarr, radarr or readarr instances using
    `DN_SONARR_1_URL`, `DN_SONARR_2_URL`, etc.
-   DiscordNotifier.com may or may not support multiple instances.


|Config Name|Variable Name|Default / Note|
|---|---|---|
api_key|`DN_API_KEY`|**Required** / API Key from DiscordNotifier.com|
quiet|`DN_QUIET`|`false` / Turns off output. Set a log_file if this is true|
bind_addr|`DN_BIND_ADDR`|`0.0.0.0:5454` / The IP and port to listen on|
urlbase|`DN_URLBASE`|default: `/` Change the web root with this setting|
upstreams|`DN_UPSTREAMS_0`|List of upstream networks that can set X-Forwarded-For|
log_file|`DN_LOG_FILE`|None by default. Optionally provide a file path to save app logs|
http_log|`DN_HTTP_LOG`|None by default. Provide a file path to save HTTP request logs|
log_files|`DN_LOG_FILES`|`10` / Log files to keep after rotating. `0` disables rotation|
log_file_mb|`DN_LOG_FILE_MB`|`10` / Max size of log files in megabytes|
timeout|`DN_TIMEOUT`|`60s` / Global API Timeouts (all apps default)|

##### Sonarr

|Config Name|Variable Name|Default / Note|
|---|---|---|
sonarr.url|`DN_SONARR_0_URL`|No Default. Something like: `http://localhost:8989`|
sonarr.api_key|`DN_SONARR_0_API_KEY`|No Default. Provide URL and API key if you use Sonarr|

##### Radarr

|Config Name|Variable Name|Default / Note|
|---|---|---|
radarr.url|`DN_RADARR_0_URL`|No Default. Something like: `http://localhost:7878`|
radarr.api_key|`DN_RADARR_0_API_KEY`|No Default. Provide URL and API key if you use Radarr|

##### Readarr

|Config Name|Variable Name|Default / Note|
|---|---|---|
readarr.url|`DN_READARR_0_URL`|No Default. Something like: `http://localhost:8787`|
readarr.api_key|`DN_READARR_0_API_KEY`|No Default. Provide URL and API key if you use Readarr|

##### Example Usage

```shell
docker pull golift/discordnotifier-client
docker run -d \
  -e "DN_API_KEY=abcdef-12345-bcfead-43312-bbbaaa-123" \
  -e "DN_SONARR_0_URL=http://localhost:8989" \
  -e "DN_SONARR_0_API_KEY=kjsdkasjdaksdj" \
  golift/discordnotifier-client
docker logs <container id from docker run>
```

### Linux and FreeBSD Install

-   Download a package from the [Releases](https://github.com/Go-Lift-TV/discordnotifier-client/releases) page.
-   Install it, edit config, start it.

Example of the above in shell form:

```shell
wget -qO- https://raw.githubusercontent.com/Go-Lift-TV/discordnotifier-client/main/scripts/install.sh | sudo bash

nano /etc/discordnotifier-client/dnclient.conf         # linux
vi /usr/local/etc/discordnotifier-client/dnclient.conf # freebsd

sudo systemctl restart discordnotifier-client    # linux
service discordnotifier-client start             # freebsd
```

On Linux, discordnotifier-client runs as `user:group` `discordnotifier-client:discordnotifier-client`.

On FreeBSD the app runs as `nobody`. That's not very good and will probably change in the future.

### macOS Install

#### Homebrew

-   Edit config file at `/usr/local/etc/discordnotifier-client/dnclient.conf`
-   Start it.
-   Like this:

```shell
brew install golift/mugs/discordnotifier-client
vi /usr/local/etc/discordnotifier-client/dnclient.conf
brew services start discordnotifier-client
```

#### macOS App

-   You can use the unsigned app on the Releases page.
-   When you open it for the first time it will create a config file and log file:
    -   `~/.dnclient/dnclient.conf`
    -   `~/.dnclient/discordnotifier-client.log`
-   Edit the config file and reload or restart the app.

### Windows Install

-   Extract a `.exe.zip` file from [the Releases page](https://github.com/Go-Lift-TV/discordnotifier-client/releases).
-   Run the `discordnotifier-client.amd64.exe` binary. This starts the app in the system tray.
-   When you open it for the first time it will create a config file and log file:
    -   `C:\ProgramData\discordnotifier-client\dnclient.conf`
    -   `<your home folder>\.dnclient\discordnotifier-client.log`
-   Edit the new config file suit your environment then reload or restart the app.

## Reverse Proxy

You'll need to expose this application to the Internet, so DiscordNotifier.com
can make connections to it. While you can certainly poke a hole your firewall
and send the traffic directly to this app, it is recommended that you put it
behind a reverse proxy. It's pretty easy.

You'll want to tune the `upstreams` and `urlbase` settings for your environment.
If your reverse proxy IP is `192.168.3.45` then set `upstreams = ["192.168.3.45/32"]`.
The `urlbase` can be left at `/`, but change it if you serve this app from a
subfolder. We'll assume you want to serve the app from `/discordnotifier/` and
it's running on `192.168.3.33` - here's a sample nginx config to do that:

```
location /discordnotifier/ {
  proxy_set_header X-Forwarded-For $remote_addr;
  proxy_pass http://192.168.3.33:5454$request_uri;
}
```

Make sure the `location` path matches the `urlbase` and ends with a `/`.
That's all there is to it.

## Troubleshooting

[Find help on Discord](https://golift.io/discord).

Log files:

-   Linux: `/var/log/messages` or `/var/log/syslog` (w/ default syslog)
-   FreeBSD: `/var/log/syslog` (w/ default syslog)
-   macOS: `/usr/local/var/log/discordnotifier-client.log`

If transfers are in a Warning or Error state they will not be extracted.
Try the Force Recheck option if you use Deluge.

Still having problems?
[Let me know!](https://github.com/Go-Lift-TV/discordnotifier-client/issues/new)

## Contributing

Yes, please.

## License

[MIT](https://github.com/Go-Lift-TV/discordnotifier-client/blob/main/LICENSE) - Copyright (c) 2020-2021 Go Lift