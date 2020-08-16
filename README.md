# Markdown Lint Checker

A tool to check for broken links in markdowns.

## Usage

```bash
$ mlc ./README.md
```

It supports two ways of loading markdowns.

1. To check the markdown by loading it directly from raw.githubusercontent.com

```
$ mlc https://raw.githubusercontent.com/d-tsuji/flower/master/README.md
[✓] https://goreportcard.com/report/github.com/d-tsuji/flower
[✓] https://goreportcard.com/badge/github.com/d-tsuji/flower
[✓] https://img.shields.io/badge/license-MIT-blue.svg
[✓] https://github.com/d-tsuji/flower/actions
[✓] https://github.com/d-tsuji/flower/workflows/build/badge.svg
[✓] https://godoc.org/github.com/d-tsuji/flower
[✓] https://godoc.org/github.com/d-tsuji/flower?status.svg
[✓] /doc/images/system_overview.png
[✓] https://en.wikipedia.org/wiki/Directed_acyclic_graph
[✓] /doc/images/task_structure.png
[✓] https://github.com/jwilder/dockerize
[✓] #post-registertask_id
[✓] https://github.com/d-tsuji/flower/blob/master/LICENSE
```

2. How to specify a local markdown file

```
$ mlc testdata/README.md
[✓] https://goreportcard.com/report/github.com/d-tsuji/flower
[✓] https://goreportcard.com/badge/github.com/d-tsuji/flower
[✓] https://img.shields.io/badge/license-MIT-blue.svg
[✓] https://github.com/d-tsuji/flower/actions
[✓] https://github.com/d-tsuji/flower/workflows/build/badge.svg
[✓] https://godoc.org/github.com/d-tsuji/flower
[✓] https://godoc.org/github.com/d-tsuji/flower?status.svg
[✓] /doc/images/system_overview.png
[✓] https://en.wikipedia.org/wiki/Directed_acyclic_graph
[✓] /doc/images/task_structure.png
[✓] https://github.com/jwilder/dockerize
[✓] #post-registertask_id
[✓] https://github.com/d-tsuji/flower/blob/master/LICENSE
```
