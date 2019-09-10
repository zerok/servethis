package server

import "html/template"

type htmlTemplateContent struct {
	Content template.HTML
}

var htmlTemplate = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
#wrapper {
	max-width: 750px;
	margin: auto;
	line-height: 32px;
	font-size: 18px;
	font-family: "Roboto", "Helvetica Neue", helvetica, arial, sans-serif;
	color: #333;
}
h1 {
	text-align: center;
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
{{ .Content }}
</div>
</body>
</html>`
