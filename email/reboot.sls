
run_patching_{{ .YamlFileName }}:
  salt.runner:
    - name: sumapatch.reboot 
    - reboot_list: /srv/pillar/sumapatch/{{ .YamlFileName }}
    - kwargs:
      delay: 2

