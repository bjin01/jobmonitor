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
    healthcheck_interval: 60
    healthcheck_email:
      - my-addr@domain.com
```

## Features:
* api endpoint - monitor SUSE Manager scheduled jobs, upon completion email notification will be sent.
* api endpoint - one can make HTTPS POST to the api to delete a system from SUSE Manager.
* health check - the program periodically makes SUSE Manager HTTP GET request to make health check.

## Delete system from SUMA via jobchecker api
```
curl --location 'http://suma1.bo2go.home:12345/delete_system' \
--data '{
  "minion_name": "pxesap02.bo2go.home",
  "authentication_token": "e2J8anZ4G4n4IM="
}'
```


