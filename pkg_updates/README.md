SLES Service Pack Migration and Package Update Engine - jobchecker
=================================================

This is a service pack migration and package update engine for [SUSE Linux Enterprise Server (SLES)](https://www.suse.com/products/server/) Service Pack (SP) migration written in Go. 

Systems from the given groups could be updated and migrated to the next service pack.

Reason choosing Go over Python is the concurrency capabilities and stabilities I experienced. Additionally Go is a compiled language and the binary is statically linked. This makes it easy to deploy the binary to any Linux systems.

## Pre-requisites:
* [SUSE Manager / Uyuni](https://www.suse.com/products/suse-manager/) v.4.3.6 or higher
* [Salt-master](https://docs.saltproject.io/en/latest/contents.html) running on SLES 15 or higher
* SMTP (e.g. postfix) where jobchecker is running must be configured to allow sending emails to the recipients
* Port 12345 for local or remote access should be open on the firewall
* Jobchecker must run as local root user

## Concept:
The one binary **jobchecker** is listening at port ```12345``` for HTTP requests. The program uses xmlrpc api to call SUSE Manager. On the other hand side jobchecker uses salt tornado rest api to execute salt state, runner and execution modules.

The name jobchecker was originally founded at the time when I started it for checking update and reboot jobs periodically. Over time I added more features to it. Now it is a bulk migration engine but can also be used for mass update and patching workflows.

## The workflow at a glance:
* service pack migration for given groups in SUSE Manager.
* email notification with system workflow information will be sent every 10 minutes.
* salt-run presence check to identify really online salt minions.
* salt grains check if btrfs root disk has enough free space. (special customer requirement)
* exclude systems by predefined no_patch grains key. (special customer requirement)
* run service pack pre-migration salt states to prepare systems for updates.
* assigne predefined software channels to systems to apply latest updates prior to service pack migration.
* update systems with latest updates of newly assigned software channels.
* reboot systems after update.
* run package refresh job.
* run service pack migration dry-run.
* run service pack migration, using ident and base channel with optional channels.
* schedule system reboot.
* run salt state to set certain values e.g. patch level, service pack level, etc. (customer requirement, for cmdb etc.)
* run service pack post-migration salt states.

Disqualified systems will get a "note" message recorded under the system -> Notes in SUSE Manager so that the reason is also visible in SUMA UI.

During the workflow cycles admins can dump data from the DB file or via web browser e.g. http://localhost:12345/pkg_update?filename=/srv/sumapatch/08112023_testgrp_t7udp.db


## systemd service for jobchecker
Jobchecker runs as a systemd service. [jobchecker.service](./etc/systemd/system/jobchecker.service)

The suma-jobchecker runs non-stop. Upone received HTTP requests it will processes the requests in sub-routines concurrently.
If no any job is in the queue then only health check is running every 10 seconds.


Inside the service file, you need to
* change the path to the binary and the path to the config file.
* change the Enrivonment variable SUMAKEY to your own key. This key is used to decrypt and encrypt the password value in the SUSE Manager configuration file.
* change the path of templates to your own path. Examples: [templates](./templates)
* change the interval to your own interval.

## SUMA Configuration and Password encryption

The password is encrypted with the key (SUMAKEY) provided in the systemd service file. The key is used to decrypt the password value.

To encrypt the password, you can use the following command:
https://github.com/bjin01/salt-sap-patching/blob/master/encrypt.py

```
python3.6 encrypt.py <YOUR-PASSWD>
```
Output:
```
Randomly generated key! Keep it safely!: 
taZk-X-MRuUSB-xYAzPys41Hi0X1iFDf0wBWynLTodw=

Save this encrypted password in your configuration file.
gAAAAABlGn1RxFaE9rRVJqVRehxTIJ6sPxPSSFuEvW4GGzmEXpT_b39D6yAQx5Us_FLLsthgUInR0UE0TPl79yf5Dsv-MNM0Bw==
```

With the encrypted password and the key, you can create the configuration file.
Example:

SUSE Manager configuration file:
```
cat /etc/salt/master.d/suma.conf 
suma_api:
  suma1.bo2go.home:
    username: 'admin'
    password: gAAAAABj_xzeu23IpzKM-mYOYO
    email_to:
      - bo.jin@example.com
    healthcheck_interval: 10
    healthcheck_email:
      - bo.jin@example.com
```

## Run spmigration
In order to start the spmigration workflow which in turn is a HTTP POST Request sent to jobchecker a [salt runner module](https://github.com/bjin01/salt-sap-patching/blob/master/srv/salt/_runners/get_spmigration_targets.py) was created. It takes a given configuration file as input parameter.

The runner module needs to be placed runners directory on salt master node. e.g. /srv/salt/_runners/get_spmigration_targets.py

```
salt-run start_spmigration.run config=/srv/salt/spmigration/spmigration_config.yaml
```

If jobchecker is not running on the same host as SUSE Manager then you need to specify the api_server parameter.

```
salt-run start_spmigration.run config=/srv/salt/spmigration/spmigration_config.yaml api_server='192.168.122.23'
```

## Configuration file
The configuration file is a yaml file. Example: [spmigration_config.yaml](./spmigration_config.yaml)

```
# The configuration file for spmigration
groups:
- testgrp
sqlite_db: "/srv/sumapatch/07112023_testgrp_t7udp.db"
qualifying_only: false
log_level: info
timeout: 3
gather_job_timeout: 15
logfile: "/var/log/patching/sumapatching.log"
salt_master_address: 192.168.122.23
salt_api_port: 8000
salt_user: mysalt
salt_password: mytest
salt_diskspace_grains_key: btrfs:for_patching
salt_diskspace_grains_value: ok
salt_no_upgrade_exception_key: no_patch
salt_no_upgrade_exception_value: 'true'
salt_prep_state: orch.prepatch_states
salt_post_state: orch.postpatch_states
jobchecker_timeout: 0
jobchecker_emails:
- bo.jin@example.com
- bo.jin@example.com
t7user: t7udp
authentication_token: R2bfp223Qsk-pX970Jw8tyJUChT4-e2J8anZ4G4n4IM=
tracking_file_directory: "/var/log/sumapatch/"
patch_level: 2023-Q4
workflow:
- assigne_channels: 1
- package_updates: 2
- package_update_reboot: 3
- package_refresh: 4
- waiting: 5
- spmigration_dryrun: 6
- spmigration_run: 7
- spmigration_reboot: 8
- spmigration_package_refresh: 9
- post_migration: 10
assigne_channels:
- assigne_channel:
    current_base_channel: mytest-prd-sle-product-sles_sap15-sp4-pool-x86_64
    new_base_prefix: mytest-test-
- assigne_channel:
    current_base_channel: mysle15sp3-prd-sle-product-sles15-sp3-pool-x86_64
    new_base_prefix: ''
products:
- product:
    ident: '1565,1576,1609,1592,1572,1616,1580,935,1633,1614'
    name: SUSE Linux Enterprise Server for SAP Applications 15 SP5 x86_64
    base_channel_label: sle-product-sles_sap15-sp5-pool-x86_64
    clm_project_label: mysap15sp5
    optionalChildChannels:
      - rke2-sles15sp5
- product:
    ident: '1396,1423,1431,1411,1415,935,1403,1438'
    name: SUSE Linux Enterprise Server 15 SP4 x86_64
    base_channel_label: sle-product-sles15-sp4-pool-x86_64
    clm_project_label:
    optionalChildChannels:
      - rke2-sles15sp4
      - rke2-common-15sp4
```
Now let's go through the configuration file.

### groups
The groups parameter is a list of groups in SUSE Manager. The systems in these groups will be migrated to the next service pack.

### sqlite_db
The sqlite_db parameter is a string value. This parameter is used to specify the sqlite database file path. The sqlite database file will be used to store the status of the systems during the workflow. The sqlite database file will be created if it does not exist. The sqlite database file will not be deleted if the workflow is finished. The jobchecker can be restarted and the process will re-read the db file and continue the workflow.

If you want to start a new workflow then you need to delete the sqlite database file or use a new db file name prior to start the workflow.

### qualifying_only
The qualifying_only parameter is a boolean value. If it is set to true then the workflow only checks if the systems have valid migration targets or not. Admins need to read the log file to see the result. If it is set to false then the workflow will continue to run the whole workflow.

### log_level
The log_level parameter is a string value. This parameter is used to specify the log level. The log level can be debug, info, warning, error, critical.

### timeout
The timeout parameter is an integer value. This parameter is used in the salt-run manage.status timeout= and gather_job_timeout= parameters. The timeout is in seconds. If number of systems is high and salt minion take more time to respond then one might consider to increase this value.

### gather_job_timeout
The gather_job_timeout parameter is an integer value. This parameter is used in the salt-run manage.status gather_job_timeout= parameter. The timeout is in seconds. If number of systems is high and salt minion take more time to respond then one might consider to increase this value.

### logfile
The logfile parameter is a string value. This parameter is used to specify the log file path.

### salt_master_address
The salt_master_address parameter is a string value. This parameter is used to specify the salt master address.

### salt_api_port
The salt_api_port parameter is an integer value. This parameter is used to specify the salt api port.

### salt_user
The salt_user parameter is a string value. This parameter is used to specify the salt user.
The salt user is defined in the salt master configuration file.
e.g.
```
external_auth:
  file:
    ^filename: /etc/salt/master.d/susemanager-users.txt
    ^hashtype: sha512
    admin:
      - .*
      - '@wheel'
      - '@runner'
      - '@jobs'
  sharedsecret:
    mysalt:
       - .*
       - '@wheel'
       - '@runner'
       - '@jobs'
```

### salt_password
The password for the mysalt user is stored in a file called ```/etc/salt/master.d/sharedsecret.conf```
```
sharedsecret: mytest
```

### salt_diskspace_grains_key
The salt_diskspace_grains_key parameter is a string value. This parameter is used to specify the salt grains key name for the disk space check. The name will be used in the salt grains.get function.

### salt_diskspace_grains_value
The salt_diskspace_grains_value parameter is a string value. This parameter is used to specify the salt grains value for the disk space check. The value will be used in the salt grains.get function. If the grains.get result key value match the pre-defined value then the system is qualified for the migration. A note will be recorded under the system -> Notes in SUSE Manager.

### salt_no_upgrade_exception_key
The salt_no_upgrade_exception_key parameter is a string value. This parameter is used to specify the salt grains key name for the no_upgrade exception. The name will be used in the salt grains.get function. If the grains.get result key value match the pre-defined value then the system is not qualified for the migration. A note will be recorded under the system -> Notes in SUSE Manager.    

### salt_no_upgrade_exception_value
The salt_no_upgrade_exception_value parameter is a string value. This parameter is used to specify the salt grains value for the no_upgrade exception. The value will be used in the salt grains.get function. If the grains.get result key value match the pre-defined value then the system is not qualified for the migration. A note will be recorded under the system -> Notes in SUSE Manager.

### salt_prep_state
The salt_prep_state parameter is a string value. This parameter is used to specify the salt state for the pre-migration tasks. The state will be executed via salt state.apply. Within the state one can define To-Do's that are needed to prepare the system. e.g. stop some heavy running processes, etc.

### salt_post_state
The salt_post_state parameter is a string value. This parameter is used to specify the salt state for the post-migration tasks. The state will be executed via salt state.apply. Within the state one can define To-Do's that are needed to clean up the system. e.g. start some heavy running processes, etc.

### jobchecker_timeout
The jobchecker_timeout parameter is an integer value. This parameter is used to specify the timeout for the jobchecker api call. The timeout is in minutes. This timeout will be used in the jobcheck loops for update and spmigration jobs. After timeout the workflow will continue with next steps. For systems which jobs are still pending the status will be recorded.

### jobchecker_emails
The jobchecker_emails parameter is a list of email addresses. This parameter is used to specify the email addresses for the jobchecker api call. The email addresses will be used in the jobchecker api call. At beginning and after the workflow is finished the jobchecker will send an email to the email recipients.

### t7user
This is a customer specific parameter to indicate which admin user is running the workflow. The value will be used in the tracking file.

### authentication_token
The authentication_token parameter is a string value. This parameter is used to specify the authentication token for the jobchecker api call. The authentication token will be used in the jobchecker api call. The authentication token is used to authenticate the caller.

### tracking_file_directory
The tracking_file_directory parameter is a string value. This parameter is used to specify the tracking file directory. The tracking file directory will be used to store the tracking file.
The purpose of this tracking file is to allow admins to follow the status of the systems while workflow is running. This file will be rewritten by jobchecker.

### patch_level
This is a customer specific parameter. The value of this parameter will be used in a built-in salt execution module to set the patch level of the system. 

### workflow
The workflow parameter is a list of workflow steps. The workflow steps will be executed in the given order.
The number indicates the order of each step. If spmigration_run or spmigration_dryrun will not be executed on systems which ident value is empty. The ident value is the product ident that SUSE Manager detects. The ident value can be obtained by a script I wrote.

### assigne_channels
```
assigne_channels:
- assigne_channel:
    current_base_channel: mytest-prd-sle-product-sles_sap15-sp4-pool-x86_64
    new_base_prefix: mytest-test-
```
In this section admins can define the channels for the systems to assigne to.
The systems will be assigned to the channels if the current_base_channel matches the system base channel.

current_base_channel: must have the parent channel label of the original SUSE vendor channel. e.g. sle-product-sles15-sp4-pool-x86_64

The new_base_prefix can have the value e.g. "myclm-test-" myclm is the content lifecycle management project label. test is the environment label. This value will be prepended to the current_base_channel and all child channels. Be careful to not forget the dash at the end of the prefix.

If new_base_prefix is left empty then the original SUSE vendor parent and child channels will be assigned if the system not already has the channels.

Make sure all child channels the systems need are available under the new parent channel.

You can define multiple assigne_channel sections to match systems with different base channels in the given groups.

The channels will be used to apply latest updates/patches before running service pack migration. It is recommended to update the system as much as possible before running service pack migration. If a system has not been patched for more than 3 months then the service pack migration might fail because certain bugs have been fixed in the latest patches.

### products
```
products:
- product:
    ident: '1565,1576,1609,1592,1572,1616,1580,935,1633,1614'
    name: SUSE Linux Enterprise Server for SAP Applications 15 SP5 x86_64
    base_channel_label: sle-product-sles_sap15-sp5-pool-x86_64
    clm_project_label: mysap15sp5
    optionalChildChannels:
      - rke2-sles15sp5
```

In this section admins can define the products (service packs) for the systems to migrate to.
The systems will be migrated to the product if the given ident matches the system product ident that SUSE Manager detects.

The ident value holds IDs that match to the varios products that each system has. The ident value of every system within SUSE Manager can be obtained by a script I wrote.

https://github.com/bjin01/salt-sap-patching/blob/master/srv/salt/_runners/get_spmigration_targets.py
```
    salt-run get_spmigration_targets.list_targets groups="a_group1 testgrp"
```
The output will be ident values and products names that will be needed in order to define the products section in the configuration file.
```
- ident: [1563,1609,1602,1592,1572,1584,1580,935], friendly: [base: SUSE Linux Enterprise Server 15 SP5 x86_64, addon: SUSE Linux Enterprise Live Patching 15 SP5 x86_64, Public Cloud Module 15 SP5 x86_64, Web and Scripting Module 15 SP5 x86_64, Basesystem Module 15 SP5 x86_64, Containers Module 15 SP5 x86_64, Server Applications Module 15 SP5 x86_64, SUSE Manager Client Tools for SLE 15 x86_64]
- ident: [1563,1609,1602,1592,1572,1584,1580,935,1633], friendly: [base: SUSE Linux Enterprise Server 15 SP5 x86_64, addon: SUSE Linux Enterprise Live Patching 15 SP5 x86_64, Public Cloud Module 15 SP5 x86_64, Web and Scripting Module 15 SP5 x86_64, Basesystem Module 15 SP5 x86_64, Containers Module 15 SP5 x86_64, Server Applications Module 15 SP5 x86_64, SUSE Manager Client Tools for SLE 15 x86_64, Python 3 Module 15 SP5 x86_64]
```
The __ident__ value in the configuration file does not need the [] brackets. The ident value is a string value. 

The **name value** is a string value. The name value is the base product name.

The **base_channel_label** value is a string value. __The base channel label is the parent channel label of the original SUSE vendor channel. e.g. sle-product-sles15-sp4-pool-x86_64__

The **clm_project_label** value is a string value. This value is the content lifecycle management project label. This value will be prepended to the environment label and base_channel_label including all child channels.

If the clm_project_label is left empty then the original SUSE vendor parent and child channels will be assigned if the system not already has the channels.

The **optionalChildChannels** value is a list of string values. This value is the optional child channels of the base channel. The child channels will be assigned to the systems. Make sure all child channels the systems need are available under the new parent channel. We don't check if the channels exist under the new parent channel. If the channels do not exist then these channels will not be assigned to the systems.

## Email notification
The email notification is a feature that sends emails to the recipients every 10 minutes within workflow deadline. The email contains the system's workflow steps information as well as any remarks, if system is online or not etc.

The email content is generated by a template. The template is located in the templates directory. The directory is specified in the jobchecker systemd service file.



## Command overview
```
To get the migration target ident:

salt-run get_spmigration_targets.list_targets target_system=minion_name

To get unique migration target ident for all systems in a suma group:

salt-run get_spmigration_targets.list_targets groups="spmigration_test"
-------------------------------------------------------------------------

It is good to delete squid cache before starting spmigration:
salt -N sumaproxy state.apply squid.delete_cache

-------------------------------------------------------------------------

To start spmigration run:

salt-run start_spmigration.run config=spmigration_config.yaml
-------------------------------------------------------------------------

To create single btrfs snapshot for a group of systems (suma)
salt-run btrfs_snapshot.create_single groups="testgrp"

-------------------------------------------------------------------------

Configuration file:
/srv/salt/base/spmigration/spmigration_config.yaml

if this parameter is set true then only target ident validation is run. 
qualifying_only: true

-------------------------------------------------------------------------

There is a bash script that will rollback snapshot to the one before spmigration. 
Nodes must be rebooted after snapshot rollback.
In SUMA the nodes must perform pkg list refresh after reboot.

btrfs_rollback.sh
-------------------------------------------------------------------------

Delete  spmigration_t7*  groups in SUSE Manager. This is a cleanup step.

salt-run get_spmigration_targets.delete_groups

-------------------------------------------------------------------------

Run salt-run module post_patching.report to get the patching report in csv format.

salt-run post_patching.report file=/var/log/sumapatch/all_a_group1_minions.yaml presence_check=True
```
