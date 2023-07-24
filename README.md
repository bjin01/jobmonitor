# SUSE Manager - jobmonitor

This is a jobmonitor written in go.

The program needs an SUSE Manager configuration file with login and others.
SUMA Config file:
```
/etc/salt/master.d/spacewalk.conf 
suma_api:
  suma1.bo2go.home:
    username: 'admin'
    password: tAOdyzcvQ==
    email_to:
      - my-addr@domain.com
    healthcheck_interval: 120
    healthcheck_email:
      - my-addr@domain.com
```
email_to: is used for /delete_system API to send notifications to.

healthcheck_email: is used to notify admins if SUMA is health check failed.

healthcheck_interval: 120 - is provided in seconds. 2 Minutes is a good interval for health check without overwhelming SUMA API.


## Features:
* api endpoint - monitor SUSE Manager scheduled jobs, upon completion email notification will be sent.
* api endpoint - one can make HTTPS POST to the api to delete a system from SUSE Manager.
* health check - the program periodically makes SUSE Manager HTTP GET request to make health check.
* product migration - is under development.

## systemd service for jobchecker
Feel free to use the systemd service file provided in this repo. [jobchecker.service](./etc/systemd/system/jobchecker.service)

Inside the service file, you need to 
* change the path to the binary and the path to the config file.
* change the Enrivonment variable SUMAKEY to your own key. This key is used to authenticate the api calls.
* change the path of templates to your own path. Examples: [templates](./templates)
* change the interval to your own interval.

```
[Service]
Environment="SUMAKEY=R2bfp223Qa="
Type=simple
Restart=always
ExecStart=/usr/local/bin/jobmonitor -config /etc/salt/master.d/spacewalk.conf -interval 60 -templates /srv/jobmonitor
```

```
cp jobchecker.service /etc/systemd/system/
systemctl daemon-reload
```

## Delete system from SUMA via jobchecker api
```
curl http://suma1.bo2go.home:12345/delete_system \
--data '{ \
  "minion_name": "pxesap02.bo2go.home", \
  "authentication_token": "e2J8anZ4G4n4IM=" \
}'
```
