
run_patching_{{ .YamlFileName }}:
  salt.runner:
    - name: sumapatch.reboot 
    - reboot_list: /srv/pillar/sumapatch/{{ .YamlFileName }}
    - kwargs:
      delay: 3
      jobchecker_timeout: 20
      jobchecker_emails:
      {{- range .JobcheckerEmails}}
        - {{.}}
      {{- end}}
      t7user: {{ .T7user }}

