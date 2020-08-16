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
$ sudo rpm -ivh https://github.com/d-tsuji/markdown-link-check/releases/download/v0.0.3/mlc_0.0.3_Tux-64-bit.rpm
```

### Debian, Ubuntu

```
$ wget https://github.com/d-tsuji/markdown-link-check/releases/download/v0.0.3/mlc_0.0.3_Tux-64-bit.deb
$ sudo dpkg -i mlc_0.0.3_Tux-64-bit.deb
```

### go get

```
$ go get -u github.com/d-tsuji/markdownlink/cmd/mlc
```
