{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "aegisgate.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "aegisgate.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "aegisgate.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels - compatible with ML and standard chart
*/}}
{{- define "aegisgate.labels" -}}
helm.sh/chart: {{ include "aegisgate.chart" . }}
{{ include "aegisgate.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: gateway
app.kubernetes.io/part-of: aegisgate
{{- end }}

{{/*
Selector labels
*/}}
{{- define "aegisgate.selectorLabels" -}}
app.kubernetes.io/name: {{ include "aegisgate.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "aegisgate.serviceAccountName" -}}
{{- $serviceAccount := .Values.serviceAccount | default .Values.global.serviceAccount | default dict }}
{{- $create := $serviceAccount.create | default false }}
{{- if $create }}
{{- default (include "aegisgate.fullname" .) $serviceAccount.name }}
{{- else }}
{{- default "default" $serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the configmap name.
*/}}
{{- define "aegisgate.configMapName" -}}
{{ include "aegisgate.fullname" . }}-config
{{- end }}

{{/*
Create the secret name for TLS.
*/}}
{{- define "aegisgate.tlsSecretName" -}}
{{- if .Values.tls.existingSecret }}
{{- .Values.tls.existingSecret }}
{{- else }}
{{ include "aegisgate.fullname" . }}-tls
{{- end }}
{{- end }}

{{/*
Generate a self-signed TLS certificate for testing.
*/}}
{{- define "aegisgate.genSelfSignedCert" -}}
{{- $host := index .Values.ingress.hosts 0 }}
{{- $cert := genSelfSignedCert $host.host nil (list $host.host) 365 }}
tls.crt: {{ $cert.Cert | b64enc }}
tls.key: {{ $cert.Key | b64enc }}
{{- end }}

{{/*
Get the image repository.
*/}}
{{- define "aegisgate.imageRepository" -}}
{{- .Values.global.imageRegistry | default "docker.io" }}/{{ .Values.global.imageRepository | default "aegisgate" }}
{{- end }}

{{/*
Get the image tag.
*/}}
{{- define "aegisgate.imageTag" -}}
{{- .Values.global.tag | default .Chart.AppVersion | default "latest" }}
{{- end }}
