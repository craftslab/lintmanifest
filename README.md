# lintmanifest

[![Build Status](https://travis-ci.com/craftslab/lintmanifest.svg?branch=master)](https://travis-ci.com/craftslab/lintmanifest)
[![Coverage Status](https://coveralls.io/repos/github/craftslab/lintmanifest/badge.svg?branch=master)](https://coveralls.io/github/craftslab/lintmanifest?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/craftslab/lintmanifest)](https://goreportcard.com/report/github.com/craftslab/lintmanifest)
[![License](https://img.shields.io/github/license/craftslab/lintmanifest.svg?color=brightgreen)](https://github.com/craftslab/lintmanifest/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/craftslab/lintmanifest.svg?color=brightgreen)](https://github.com/craftslab/lintmanifest/tags)



## Introduction

*Lint Manifest* is a tool to lint Android Manifest using *[Gitiles](https://gerrit.googlesource.com/gitiles)*.



## Features

- Support to lint revision in Android Manifest.



## Prerequisites

- Gitiles 0.3+



## Usage

```bash
usage: lintmanifest --manifest-file=MANIFEST-FILE --gitiles-url=GITILES-URL --lint-out=LINT-OUT [<flags>]

Lint Manifest

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --version                      Show application version.
  --manifest-file=MANIFEST-FILE  Manifest file
  --gitiles-pass=GITILES-PASS    Gitiles password
  --gitiles-url=GITILES-URL      Gitiles location
  --gitiles-user=GITILES-USER    Gitiles username
  --lint-out=LINT-OUT            Lint output
```



## Run

```bash
./lintmanifest --manifest-file=manifest.xml --gitiles-url=localhost:80 --lint-out=out.txt
```



## License

Project License can be found [here](LICENSE).
