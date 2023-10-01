Automated bulk SLES Service Pack Migration engine
=================================================

This is a bulk migration engine for SUSE Linux Enterprise Server (SLES) Service Pack (SP) migration written in go. 

The program needs an SUSE Manager / Uyuni v.4.3.6 or higher.

## Pre-requisites:
* SUSE Manager / Uyuni v.4.3.6 or higher
* Salt minions running on SLES 12 SP4 or higher

## Concept:
The one binary program **jobchecker** is listening on port 12345 for HTTP requests. The program uses xmlrpc api to talk to SUSE Manager. On the other hand side jobchecker uses salt tornado rest api to talk to salt-master in order to run salt state, runner and execution modules.

The name jobchecker was founded at the time when I started it for just checking update and reboot jobs periodically. Over time I added more features to it. Now it is a bulk migration engine.

The program is using SUSE Manager / Uyuni Salt API to execute salt states on minions. The program is using the following salt states:
* [spmigration_pre](./salt/spmigration_pre/init.sls) - this state is used to execute pre migration tasks.
* [spmigration_post](./salt/spmigration_post/init.sls) - this state is used to execute post migration tasks.


## Features:
* api endpoint - monitor SUSE Manager scheduled jobs, upon completion email notification will be sent.
* api endpoint - one can make HTTPS POST to the api to delete a system from SUSE Manager.
* health check - the program periodically makes SUSE Manager HTTP GET request to make health check.
* api endpoint - product migration - Upgrade systems within given groups in SUSE Manager to a defined service pack.
* email notifications - about tracking file, job results.
* api endpoint - salt states, grains execution for pre and post tasks

## systemd service for jobchecker
Feel free to use the systemd service file provided in this repo. [jobchecker.service](./etc/systemd/system/jobchecker.service)

The suma-jobchecker runs non-stop. Upone received HTTP requests it will processes the requests in sub-routines concurrently.
It no any job is in the queue the jobcheck is simply check the health of SUSE Manager, every 10 seconds.


Inside the service file, you need to
* change the path to the binary and the path to the config file.
* change the Enrivonment variable SUMAKEY to your own key. This key is used to authenticate the api calls.
* change the path of templates to your own path. Examples: [templates](./templates)
* change the interval to your own interval.



