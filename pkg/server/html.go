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
#wrapper {
	max-width: 750px;
	margin: 60px auto;
	line-height: 32px;
	font-size: 18px;
	font-family: "Roboto", "Helvetica Neue", helvetica, arial, sans-serif;
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
	margin-bottom: 60px;
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
</style>
</head>
<body>
<div id="wrapper">
{{ if .Frontmatter.Title }}<h1>{{ .Frontmatter.Title }}</h1>{{ end }}
{{ .Content }}
</div>
</body>
</html>`
