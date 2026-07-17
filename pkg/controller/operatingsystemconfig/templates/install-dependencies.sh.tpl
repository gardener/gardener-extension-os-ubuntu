UBUNTU_VERSION=""
BUILD_SERIAL=""
if [ -f /etc/os-release ]; then
  UBUNTU_VERSION=$(grep '^VERSION_ID=' /etc/os-release | cut -d= -f2 | tr -d '"')
fi
if [ -f /etc/cloud/build.info ]; then
  BUILD_SERIAL=$(grep '^serial:' /etc/cloud/build.info | awk '{print $2}')
fi

until apt-get update -qq; do sleep 1; done

install_package() {
  local name=$1
  local version=$2
  local hold=$3

  if [ -n "$version" ]; then
    until apt-get install --no-upgrade -qqy "${name}=${version}"; do sleep 1; done
  else
    until apt-get install --no-upgrade -qqy "${name}"; do sleep 1; done
  fi

  if [ "$hold" = "true" ]; then
    apt-mark hold "$name"
  fi
}

{{ range . -}}
{{ $name := .Name -}}
{{ if .Blocks -}}
{{ range $idx, $block := .Blocks -}}
{{ if eq $idx 0 -}}
if [[ {{ $block.Condition }} ]]; then
{{- else }}
elif [[ {{ $block.Condition }} ]]; then
{{- end }}
{{- if $block.Disabled }}
  echo "Skipping disabled package {{ $name }} (UBUNTU_VERSION=$UBUNTU_VERSION, BUILD_SERIAL=$BUILD_SERIAL)"
{{- else }}
  install_package "{{ $name }}" "{{ $block.Command.Version }}" {{ $block.Command.Hold }}
{{- end }}
{{- end }}
{{- if .Fallback }}
else
  install_package "{{ $name }}" "{{ .Fallback.Version }}" {{ .Fallback.Hold }}
fi
{{- else }}
else
  echo "No matching pinned dependency for {{ $name }} (UBUNTU_VERSION=$UBUNTU_VERSION, BUILD_SERIAL=$BUILD_SERIAL)" >&2
  exit 1
fi
{{- end }}
{{ else -}}
install_package "{{ $name }}" "{{ if .Fallback }}{{ .Fallback.Version }}{{ end }}" {{ if .Fallback }}{{ .Fallback.Hold }}{{ else }}false{{ end }}
{{ end -}}
{{ end -}}
