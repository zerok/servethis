# ServeThis

A simple HTTP server for the current working directory with some extras 🙂 I
often have the use-case that I need a non-"file:" served HTML page for proper
JavaScript handling. I thought it might be nice to also combine that with
something like a Markdown processor and ideally some other content-specific
"renderers".

## Usage

```
$ servethis
9:35PM INF Serving content from /Users/user/some/path
9:35PM INF Starting listener on 127.0.0.1:9980
```

## Content renderers

- Markdown files are rendered using [gomarkdown][].


[gomarkdown]: https://github.com/gomarkdown/markdown