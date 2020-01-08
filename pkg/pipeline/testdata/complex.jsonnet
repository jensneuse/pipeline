{
  steps: [
    {
      kind: 'JSON',
      config: {
        template: |||
          {
              "url": "{{ .url }}",
              "method": "POST",
              "body" : {
                  "policies": {
                      {{ range $i, $e := .body.policies }}
                      {{ if $i }},{{ end }}
                      {{ with $e }}
                      "{{ .id }}": {
                          "id": {{ .id }},
                          "name": "{{ .name }}"
                      }
                      {{ end }}
                      {{ end }}
                  }
              }
          }
        |||,
      },
    },
    {
      kind: 'HTTP',
      config: {
        default_timeout: 10000000000,
        default_method: 'GET',
      },
    },
    {
      kind: 'JSON',
      config: {
        template: |||
          {
              "policies": [
                  {{ $first := true }}
                  {{ range $key, $value := .body.policies }}
                  {{ if $first }}{{ $first = false }}{{ else }},{{ end }}
                  {
                      "id": {{ $key }},
                      "name": "{{ $value.name }}"
                  }
                  {{ end }}
              ]
          }
        |||,
      },
    },
    {
      kind: 'JSON',
      config: {
        template: '{{ toPrettyJson . }}',
      },
    },
  ],
}
