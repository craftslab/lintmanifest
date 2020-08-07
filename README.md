# lintmanifest

[![Build Status](https://travis-ci.com/craftslab/lintmanifest.svg?branch=master)](https://travis-ci.com/craftslab/lintmanifest)
[![Coverage Status](https://coveralls.io/repos/github/craftslab/lintmanifest/badge.svg?branch=master)](https://coveralls.io/github/craftslab/lintmanifest?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/craftslab/lintmanifest)](https://goreportcard.com/report/github.com/craftslab/lintmanifest)
[![License](https://img.shields.io/github/license/craftslab/lintmanifest.svg?color=brightgreen)](https://github.com/craftslab/lintmanifest/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/craftslab/lintmanifest.svg?color=brightgreen)](https://github.com/craftslab/lintmanifest/tags)



## Introduction

*Lint Manifest* is a tool to lint Android Manifest using *[Gitiles](https://gerrit.googlesource.com/gitiles)*.



## Features

- Support to lint revision in branch of repository.

- Support to lint revision as HEAD in branch of repository.



## Prerequisites

- Gitiles 0.3+



## Usage

```bash
usage: lintmanifest --gitiles-url=GITILES-URL --manifest-file=MANIFEST-FILE [<flags>]

Lint Manifest

Flags:
  --help                         Show context-sensitive help (also try
                                 --help-long and --help-man).
  --version                      Show application version.
  --gitiles-pass=GITILES-PASS    Gitiles password
  --gitiles-url=GITILES-URL      Gitiles location
  --gitiles-user=GITILES-USER    Gitiles username
  --lint-mode="sync"             Lint mode (async|sync)
  --lint-out="out.txt"           Lint output (.json|.txt|.xlsx)
  --manifest-file=MANIFEST-FILE  Manifest file
```



## Run

```bash
./lintmanifest --gitiles-url=localhost:80 --manifest-file=manifest.xml
```



## Output

### JSON Format

```json
{
  "lintmanifest": [
    {
      "branch": "NAME",
      "commit": "HASH",
      "details": "DETAILS",
      "repo": "NAME",
      "type": "ERROR"
    },
    {
      "branch": "NAME",
      "commit": "HASH",
      "details": "DETAILS",
      "repo": "NAME",
      "type": "WARN"
    }
  ]
}
```



### TXT Format

```txt
TYPE,REPO,BRANCH,COMMIT,DETAILS
```



### XLSX Format

```
TYPE,REPO,BRANCH,COMMIT,DETAILS
```



## License

Project License can be found [here](LICENSE).
