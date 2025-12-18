{{- define "replay-platform.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "replay-platform.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "replay-platform.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "replay-platform.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "replay-platform.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "replay-platform.image" -}}
{{- printf "%s/%s:%s" .Values.global.imageRegistry .repository .tag }}
{{- end }}
