[![Go Report Card](https://goreportcard.com/badge/github.com/bitnami/ini-file)](https://goreportcard.com/report/github.com/bitnami/ini-file)
[![CI](https://github.com/bitnami/ini-file/actions/workflows/main.yml/badge.svg)](https://github.com/bitnami/ini-file/actions/workflows/main.yml)

# ini-file

This tool allows manipulating INI files.

# Basic usage

```console
$ ini-file --help

Usage:
  ini-file [OPTIONS] <del | get | set>

Help Options:
  -h, --help  Show this help message

Available commands:
  del  INI FILE Delete
  get  INI FILE Get
  set  INI File Set
```

# Examples

## Set values in INI file

```console
$ ini-file set --section "My book" --key "title" --value "A wonderful book" ./my.ini
$ ini-file set -s "My book" -k "author" -v "Bitnami" ./my.ini
$ ini-file set -s "My book" -k "rate" -v "very good" ./my.ini
$ cat ./my.ini
[My book]
title=A wonderful book
author=Bitnami
rate=very good
```

## Set boolean value in INI file

```console
$ ini-file set -s "My book" -k "already_read" --boolean ./my.ini
$ cat ./my.ini
[My book]
...
already_read
```

## Get values from INI file

```console
$ cat > ./my.ini <<"EOF"
[My book]
title=A wonderful book
author=Bitnami
rate=very good
already_read
EOF
$ ini-file get --section "My book" --key title ./my.ini
A wonderful book
$ ini-file get --section "My book" --key missing_key ./my.ini

$ ini-file get --section "My book" --key author ./my.ini
Bitnami
$ ini-file get --section "My book" --key already_read ./my.ini
true
```

## Deletes values from INI file

```console
$ cat > ./my.ini <<"EOF"
[My book]
title=A wonderful book
author=Bitnami
rate=very good
already_read
EOF
$ ini-file del --section "My book" --key title ./my.ini
$ ini-file del --section "My book" --key missing_key ./my.ini
$ ini-file del --section "My book" --key author ./my.ini
$ cat ./my.ini
[My book]
rate=very good
already_read
```

## Working with identical keys in the file

```console
$ cat > ./my.ini <<"EOF"
[My book]
title=A wonderful book
author[]=Bitnami
author[]=Contributors
EOF

# get retrieves the first mention of the key
$ ini-file get --section "My book" --key "author[]" ./my.ini
Bitnami

# If the original file contains the key more than once, set adds the new value at the end
$ ini-file set --section "My book" --key "author[]" --value "Other" ./my.ini
$ cat ./my.ini
title=A wonderful book
author[]=Bitnami
author[]=Contributors
author[]=Other

# del removes all keys with the given name
$ ini-file del --section "My book" --key "author[]" ./my.ini
$ cat ./my.ini
title=A wonderful book
```

## License

Copyright &copy; 2025 Broadcom. The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
