package main

import "text/template"

func parseStrings(name string, base *template.Template, templs ...string) *template.Template {
	var t *template.Template
	if base == nil {
		t = template.New(name)
	} else {
		t = template.Must(base.Clone())
	}
	for _, templ := range templs {
		t = template.Must(t.Parse(templ))
	}
	return t
}

// Main template
var templ = `
<!DOCTYPE html>
<html>
<head>
{{ block "title" . }}<title>Siegfried</title>{{ end }}
<link rel="icon" type="image/png" href="/img/richard.png">
{{ block "incCSS" . }}{{ end }}
<link rel="stylesheet" href="/css/maia.css">
<meta name="viewport" content="width=device-width, initial-scale=1">
{{ block "incJS" . }}{{ end }}
</head>
<body>
<div class="topcorner">
	<a href="/siegfried">Back to siegfried</a> 
</div>
<div>
{{ block "content" . }}{{ end }}
</div>
</body>
</html>`

// Results templates
var rTitleTempl = `{{ define "title" }}<title>{{ if len .Title | eq 0  }}Siegfried results chart{{ else }}{{ .Title }}{{ end }}</title>{{ end }}`

// to refresh these: go to https://datatables.net/download/index, choose Datatables styling, jquery3, Datatables, Buttons-> HTML5 export, CDN -> Minify + Concatentate
var rChartCSSTempl = `{{ define "incCSS" -}} 
<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/v/dt/jq-3.6.0/dt-1.12.0/b-2.2.3/b-html5-2.2.3/datatables.min.css"/>
<style>
.chart {
	height: 320px;
}
</style>
{{- end }} `

var rChartJSTempl = `{{ define "incJS" -}} 
<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
<script type="text/javascript" src="https://cdn.datatables.net/v/dt/jq-3.6.0/dt-1.12.0/b-2.2.3/b-html5-2.2.3/datatables.min.js"></script>
<script type="text/javascript">var RESULTS = {{ .JSON }};</script>
<script src="/js/results.js"></script>
{{- end }} `

var rShareTempl = `{{ define "share" -}} 
<div class="l-box">
	<h1>{{ if len .Title | lt 0 }}{{ .Title }}{{ else }}Untitled{{ end }}</h1>
	<p><i>{{ if len .Name | lt 0 }}{{ .Name }}{{ end }}</i></p>
	<p>{{ if len .Desc | lt 0 }}{{ .Desc }}{{ end }}</p>
	<p>
		<a class= "signature" href="https://twitter.com/intent/tweet?text=
			{{- urlquery "I'm charting my formats! " .Title " (https://www.itforarchivists.com/siegfried/results/" .UUID ")" -}}">
		Post to twitter
		</a>
	</p>
</div>
{{- end }} `

var rContent = `{{ define "content" -}}
<div class="chart-container"> 
{{ block "share" . -}}
	<div class="chart-box">
	<h1>Share results</h1>
		<form id="share-form">
		  <fieldset>
			<input id="sharename" type="text" name="name" maxlength="128" size="20" placeholder="Name (or organisation)">
			<input id="sharetitle" type="text" name="title" maxlength="128" size="20" placeholder="Title">
			<textarea id="sharedesc" name="description" maxlength="256" rows="3" cols="20" placeholder="Description"></textarea>
			<label for="redact">
				<input id="redact" type="checkbox" name="redact" value="true" checked> Redact filenames <a href="#" id="redactNow">(redact now)</a>
			</label>
			<input type="submit" value="Publish" class="publish-button">
			<button style="display: none" class="publish-button">Publish</button>
		  </fieldset>
		</form>
	</div>
{{- end }}
<div class="chart-box">
<h1>Identifiers</h1>
{{- range $idx, $el := .Identifiers -}}
		<p><a href="#" onclick="load({{ $idx }}); return false;"><strong>{{ $el.Name }}</strong></a><br />{{ $el.Details }}</p>
{{- end -}}
</div>
<div class="chart-box">
<h1>Details</h1>
<p>
{{- range $idx, $el := .Metadata }}
	{{- if index $el 1 | len | ne 0 -}}
	  {{- if gt $idx 0 -}}<br />{{ end -}}
	  <strong>{{ index $el 0 }}</strong>: {{ index $el 1 }}
	{{- end -}}
{{ end -}}
</p>
  <p>
    <a href="#" id="errNo"><span></span> errors</a><br />
	<a href="#" id="warnNo"><span></span> warnings</a><br />
	<a href="#" id="unkNo"><span></span> unknowns</a><br />
	<a href="#" id="multiNo"><span></span> multiple IDs</a><br />
	<a href="#" id="dupesNo"><span></span> duplicates</a>
  </p>
  <p>
    <a href="#" id="reset">Reset (show all)</a>
  </p>
</div>
<div class="chart-box" id="charts">
  <div id="fmtchart" class="chart"></div>
  <div id="mimechart" class="chart"></div>
  <div class="centre" role="group">
    <button onclick="reveal('fmtchart'); return false;">Format IDs</button>
    <button onclick="reveal('mimechart'); return false;">MIME-types</button>
  </div>
</div>
</div>
<div><table id="table" class="display" width="100%"></table></div>
{{- end }} `

// Log templates

var lTitleTempl = `{{ define "title" }}<title>{{ if len .Title | eq 0  }}Siegfried logs{{ else }}{{ .Title }}{{ end }}</title>{{ end }}`

var lCSSTempl = `{{ define "incCSS" -}} 
<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/v/dt/jq-3.6.0/dt-1.12.0/b-2.2.3/b-html5-2.2.3/datatables.min.css"/>
{{- end -}}
`
var lJSTempl = `{{ define "incJS" -}} 
<script type="text/javascript" src="https://cdn.datatables.net/v/dt/jq-3.6.0/dt-1.12.0/b-2.2.3/b-html5-2.2.3/datatables.min.js"></script>
{{- end -}}
`

var lContent = `{{ define "content" -}} 
<div class="content">
<h1>{{ .Title }}</h1>
<h2>{{ .Time }}</h2>
<h3>Environment</h3>
<p>These benchmarks were <a href="https://github.com/richardlehane/runner">run</a> on a <a href="{{.Machine.Link}}">{{ .Machine.Label}}</a> machine that was <a href="https://github.com/richardlehane/provisioner">automatically provisioned</a>.</p>
<p>Specs for the <a href="{{.Machine.Link}}">{{ .Machine.Label}}</a>: {{.Machine.Description}}.</p>
<p>You can inspect the commands that were run to generate these benchmarks <a href="/siegfried/jobs/{{ .Prefix }}">here</a>.</p>
{{ if len .Versions | lt 0 -}}
<table>
	<caption>List of tools benchmarked</caption>
	<thead>
		<tr>
			<th>Tool</th>
			<th>Version</th>
		</tr>
	</thead>
	<tbody>
		{{- range .Versions -}}
		<tr>
			<td>{{ .Label }}</td>
			<td>{{ .Version }}</td>
		</tr>
		{{- end -}}
	</tbody>
</table>
{{- end -}}
{{- range $idx, $el := .Benchmarks -}}
<div>
<h2>{{ .Title }}</h2>
<p>{{ .Description }}</p>
<table>
	<caption>Results</caption>
	<thead>
		<tr>
			<th>Tool</th>
			<th>Description</th>
			<th>Duration</th>
		</tr>
	</thead>
	<tbody>
		{{- range .Tools -}}
		<tr>
			<td>{{ .Label }}</td>
			<td>{{ .Description }}</td>
			<td>{{ .Duration }}</td>
		</tr>
		{{- end -}}
	</tbody>
</table>
<p>{{ .CompareDesc }}</p>
{{ if len .Compare | lt 0 -}}
<table id="cmp{{ $idx }}" style="width:100%">
	<caption>Differences between tools</caption>
	<thead>
		<tr>
		<td>file</td>
		{{- range .CompareHdrs -}}
		<td>{{ . }}</td>
		{{- end -}}
		</tr>
	</thead>
	<tbody>
		{{- range $row := .Compare -}}
		<tr>
		{{- range $row -}}
			<td>{{ . }}</td>
		{{- end -}}
		</tr>
		{{- end -}}
	</tbody>
</table>
<script>
$(document).ready(function() {
    $('#cmp{{ $idx }}').DataTable();
} );
</script>
{{ end }}
<p><a href="/siegfried/logs/{{ .Src }}">Raw output</a></p>
</div>			
{{- end -}}
{{ if len .Profile | lt 0 -}}
<div>
    <h2>Profile</h2>
    <img width="1200" src="data:image/png;base64, {{ .Profile }}" alt="profiler information for siegfried development branch" />
</div> 
{{- end -}}
<div>
<h2>History</h2>
{{ $prefix := .Prefix }}
{{- range .History -}}
<p><a href="/siegfried/{{ if eq $prefix "develop" }}develop{{ else }}benchmarks{{ end }}/{{ index . 0 }}">{{ index . 1 }}</a></p>
{{- end -}}
</div>
</div>
{{- end -}}
`
