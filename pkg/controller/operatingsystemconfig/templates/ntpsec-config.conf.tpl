{{- range .Servers }}
server {{ . }} iburst
{{- end }}

driftfile /var/lib/ntpsec/ntp.drift
restrict default kod nomodify noquery limited
restrict 127.0.0.1
restrict ::1

{{ if .Interfaces -}}
interface ignore wildcard
interface listen 127.0.0.1
{{ range .Interfaces -}}
interface listen {{ . }}
{{ end -}}
{{- end -}}
