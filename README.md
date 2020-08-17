# Markdown Lint Checker (mlc)

A tool to check for broken links in markdowns. Because mlc can check links in parallel, it is very fast. Check for links that are described as markdown links or images.

It does not check for links in code like the following.

```
xxx
```

## Usage

```bash
$ mlc ./README.md
```

It supports two ways of loading markdowns.

1. To check the markdown by loading it directly from raw.githubusercontent.com

![mlc_github.gif](./gif/mlc_github.gif)

```
$ mlc https://raw.githubusercontent.com/d-tsuji/flower/master/README.md
FILE: https://raw.githubusercontent.com/d-tsuji/flower/master/README.md
Checking... 13 / 13 [--------------------] 100.00%
[✓] https://img.shields.io/badge/license-MIT-blue.svg
[✓] https://en.wikipedia.org/wiki/Directed_acyclic_graph
[✓] https://github.com/d-tsuji/flower/workflows/build/badge.svg
[✓] https://godoc.org/github.com/d-tsuji/flower
[✓] #post-registertask_id
[✓] /doc/images/system_overview.png
[✓] https://goreportcard.com/badge/github.com/d-tsuji/flower
[✓] https://goreportcard.com/report/github.com/d-tsuji/flower
[✓] https://github.com/d-tsuji/flower/blob/master/LICENSE
[✓] /doc/images/task_structure.png
[✓] https://github.com/jwilder/dockerize
[✓] https://github.com/d-tsuji/flower/actions
[✓] https://godoc.org/github.com/d-tsuji/flower?status.svg
```

2. How to specify a local markdown file

```
$ mlc testdata/README.md
FILE: testdata/README.md
Checking... 13 / 13 [--------------------] 100.00%
[✖] /doc/images/system_overview.png
[✖] /doc/images/task_structure.png
[✓] https://img.shields.io/badge/license-MIT-blue.svg
[✓] #post-registertask_id
[✓] https://github.com/d-tsuji/flower/actions
[✓] https://github.com/d-tsuji/flower/blob/master/LICENSE
[✓] https://github.com/d-tsuji/flower/workflows/build/badge.svg
[✓] https://github.com/jwilder/dockerize
[✓] https://godoc.org/github.com/d-tsuji/flower?status.svg
[✓] https://godoc.org/github.com/d-tsuji/flower
[✓] https://en.wikipedia.org/wiki/Directed_acyclic_graph
[✓] https://goreportcard.com/report/github.com/d-tsuji/flower
[✓] https://goreportcard.com/badge/github.com/d-tsuji/flower

13 links checked.

ERROR: 2 dead links found!
[✖] /doc/images/system_overview.png -> Status: 404
[✖] /doc/images/task_structure.png -> Status: 404
```

## Install

### Binary

If you need the Binary file, download the zip file of the version you want from the [Releases](https://github.com/d-tsuji/markdown-link-check/releases) page.
Unzip the zip file and place the Binary file where the path will take you.

### macOS

```
$ brew tap d-tsuji/mlc
$ brew install mlc
```

### CentOS

```
$ sudo rpm -ivh https://github.com/d-tsuji/markdown-link-check/releases/download/v0.0.4/mlc_0.0.4_Tux-64-bit.rpm
```

### Debian, Ubuntu

```
$ wget https://github.com/d-tsuji/markdown-link-check/releases/download/v0.0.4/mlc_0.0.4_Tux-64-bit.deb
$ sudo dpkg -i mlc_0.0.4_Tux-64-bit.deb
```

### go get

```
$ go get -u github.com/d-tsuji/markdownlink/cmd/mlc
```
