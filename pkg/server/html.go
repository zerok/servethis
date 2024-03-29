package server

import "html/template"

type htmlTemplateContent struct {
	Content     template.HTML
	Frontmatter frontmatter
}

var htmlTemplate = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
.file-listing {
	list-style: none;
	padding: 0;
	margin: 0;
}
.file-listing li {
	padding: 5px 0;
	border-bottom: 1px solid #DDD;
}
.file-listing__name {
	display: block;
}
.file-listing__modtime {
	position: absolute;
	bottom: 5px;
	right: 0px;
	font-size: 80%;
	opacity: 0.6;
}
.file-listing__item {
	position: relative;
}
.file-listing__item--dir {
	opacity: 0.5;
}
#wrapper {
	max-width: 750px;
	margin: 60px auto;
	line-height: 32px;
	font-size: 18px;
	font-family: "Liberation Mono", "Courier New", Courier, monospace;
	color: #333;
}
pre {
	font-family: "Courier New", Courier, monospace;
	font-size: 14px;
	margin: 30px 0;
	border: 1px solid #DDD;
	padding: 1em;
	line-height: 1.3em;
	overflow: hidden;
}
pre code {
	display: block;
	background: #444;
	color: #EEE;
	padding: 1em;
	overflow: hidden;
}
blockquote {
	font-style: italic;
	padding-left: 2em;
	font-size: 95%;
	margin: 3em 0;
	color: #4d4d4d;
	border-left: 8px solid #e9f3fd;
}
h1 {
	text-align: center;
}
header {
	margin-bottom: 60px;
}
header h1 {
	margin-bottom: 0;
}
header p {
	margin: 0;
	text-align: center;
	font-size: 14px;
	color: #999;
}

header p span {
	color: #1077D7;
	text-decoration: underline;
}
h1, h2, h3, h4, h5, h6 {
	color: #1177d7;
}
a {
	box-shadow: inset 0 -4px 0 #8ac2f6;
	text-decoration: none;
	color: #333;
}
a:hover {
	background: #8ac2f6;
	color: #FFF;
}
p {
	margin-bottom: 1em;
}
figure img {
	max-width: 98%;
	padding: 1%;
}

figure {
	border: 1px solid #DDD;
	text-align: center;
}

figcaption {
	background: #DDD;
	color: #333;
	font-size: 14px;
	padding: 2px 5px;
	text-align: center;
}
@media (prefers-color-scheme: dark) {
  body {
		background: #333;
		color: #FEFEFE;
	}
	a, #wrapper {
		color: #FEFEFE;
	}
	h1, h2, h3, h4, h5, h6,
	header p span {
		color: #449ef2;
	}
	blockquote {
		color: #DFDFDF;
	}
}
</style>
</head>
<body>
<div id="wrapper">
<header>
{{ if .Frontmatter.Title }}<h1>{{ .Frontmatter.Title }}</h1>{{ end }}
{{ if .Frontmatter.Tags }}
<p>Tags: {{ range .Frontmatter.Tags }}<span>{{ . }}</span> {{ end }}</p>
{{ end }}
</header>
{{ .Content }}
</div>
</body>
</html>`

var listingTemplate = `<ul class="file-listing">
{{ range . }}
<li class="file-listing__item {{ if .IsDir }}file-listing__item--dir{{ end }}">
	<span class="file-listing__name"><a href="{{ .Name }}{{ if .IsDir }}/{{ end }}">{{ .Name }}</a></span>
	<span class="file-listing__modtime">{{ datetime .ModTime }}</span>
</li>
{{ end }}
</ul>
`
