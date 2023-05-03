
run_patching_{{ .YamlFileName }}:
  salt.runner:
    - name: sumapatch.reboot 
    - reboot_list: /srv/pillar/sumapatch/{{ .YamlFileName }}
    - kwargs:
      delay: 3
      jobchecker_timeout: 20
      jobchecker_emails:
        - bo.jin@jinbo01.com
        - bo.jin@suseconsulting.ch
      t7user: t7udp

