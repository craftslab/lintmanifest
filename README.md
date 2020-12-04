# lintmanifest

[![Build Status](https://travis-ci.com/craftslab/lintmanifest.svg?branch=master)](https://travis-ci.com/craftslab/lintmanifest)
[![Coverage Status](https://coveralls.io/repos/github/craftslab/lintmanifest/badge.svg?branch=master)](https://coveralls.io/github/craftslab/lintmanifest?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/craftslab/lintmanifest)](https://goreportcard.com/report/github.com/craftslab/lintmanifest)
[![License](https://img.shields.io/github/license/craftslab/lintmanifest.svg?color=brightgreen)](https://github.com/craftslab/lintmanifest/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/craftslab/lintmanifest.svg?color=brightgreen)](https://github.com/craftslab/lintmanifest/tags)



## Introduction

*lintmanifest* is a tool to lint Android Manifest using *[Gerrit](https://gerrit.googlesource.com/gerrit)* and *[Gitiles](https://gerrit.googlesource.com/gitiles)*.



## Features

- Support to lint revision in branch of repository.
- Support to lint revision as HEAD in branch of repository.



## Prerequisites

- Gerrit >= 3.2.0
- Gitiles >= 0.3



## Build

```bash
git clone https://github.com/craftslab/lintmanifest.git

cd lintmanifest
make build
```



## Run

```bash
./lintmanifest --config-file=config.json --lint-mode="sync" --lint-out=out.json --manifest-file=Manifest.xml
```



## Usage

```bash
usage: lintmanifest --config-file=CONFIG-FILE --manifest-file=MANIFEST-FILE [<flags>]

Lint Manifest

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --version                      Show application version.
  --config-file=CONFIG-FILE      Config file, format: .json
  --lint-mode="sync"             Lint mode (async|sync)
  --lint-out="out.json"          Lint output (.json|.txt|.xlsx)
  --manifest-file=MANIFEST-FILE  Manifest file
```



## Settings

*lintmanifest* parameters can be set in the directory [config](https://github.com/craftslab/lintmanifest/blob/master/config).

An example of configuration in [config.json](https://github.com/craftslab/lintmanifest/blob/master/config/config.json):

```
{
  "gerrit": {
    "pass": "pass",
    "query": {
      "option": ["CURRENT_REVISION"]
    },
    "url": "http://127.0.0.1:80",
    "user": "user"
  },
  "gitiles": {
    "pass": "pass",
    "url": "http://127.0.0.1:80/a/plugins/gitiles",
    "user": "user"
  }
}
```



## Output

### JSON

```json
{
  "lintmanifest": [
    {
      "commitHash": "HASH",
      "commitLocalDate": "DATE",
      "commitRemoteDate": "DATE",
      "commitUrl": "URL",
      "headHash": "HASH",
      "headLocalDate": "DATE",
      "headRemoteDate": "DATE",
      "headUrl": "URL",
      "repoBranch": "NAME",
      "repoName": "NAME",
      "reportDetails": "DETAILS",
      "reportType": "ERROR"
    }
  ]
}
```



### Text

```txt
REPORTTYPE, REPONAME, REPOBRANCH, HEADREMOTEDATE, HEADLOCALDATE, HEADHASH, HEADURL, COMMITREMOTEDATE, COMMITLOCALDATE, COMMITHASH, COMMITURL, REPORTDETAILS
```



### Excel

```xlsx
REPORTTYPE, REPONAME, REPOBRANCH, HEADREMOTEDATE, HEADLOCALDATE, HEADHASH, HEADURL, COMMITREMOTEDATE, COMMITLOCALDATE, COMMITHASH, COMMITURL, REPORTDETAILS
```



## License

Project License can be found [here](LICENSE).
