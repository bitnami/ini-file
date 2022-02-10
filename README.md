[![CI](https://github.com/bitnami/ini-file/actions/workflows/main.yml/badge.svg)](https://github.com/bitnami/ini-file/actions/workflows/main.yml)

# ini-file

This tool allows manipulating INI files.

# Basic usage

~~~bash
$> ini-file --help

Usage:
  ini-file [OPTIONS] <del | get | set>

Help Options:
  -h, --help  Show this help message

Available commands:
  del  INI FILE Delete
  get  INI FILE Get
  set  INI File Set
~~~

# Examples

## Set values in INI file

~~~bash
$> ini-file set --section "My book" --key "title" --value "A wonderful book" ./my.ini
$> ini-file set -s "My book" -k "author" -v "Bitnami" ./my.ini
$> ini-file set -s "My book" -k "rate" -v "very good" ./my.ini
$> cat ./my.ini
[My book]
title=A wonderful book
author=Bitnami
rate=very good

~~~

## Set boolean value in INI file

~~~bash
$> ini-file set -s "My book" -k "already_read" --boolean ./my.ini
$> cat ./my.ini
[My book]
...
already_read
~~~

## Get values from INI file

~~~bash
$> cat > ./my.ini <<"EOF"
[My book]
title=A wonderful book
author=Bitnami
rate=very good
already_read
EOF
$> ini-file get --section "My book" --key title ./my.ini
A wonderful book
$> ini-file get --section "My book" --key missing_key ./my.ini

$> ini-file get --section "My book" --key author ./my.ini
Bitnami
$> ini-file get --section "My book" --key already_read ./my.ini
true
~~~

## Deletes values from INI file

~~~bash
$> cat > ./my.ini <<"EOF"
[My book]
title=A wonderful book
author=Bitnami
rate=very good
already_read
EOF
$> ini-file del --section "My book" --key title ./my.ini
$> ini-file del --section "My book" --key missing_key ./my.ini
$> ini-file del --section "My book" --key author ./my.ini
$> cat ./my.ini
[My book]
rate=very good
already_read
~~~

## Working with identical keys in the file

~~~bash
$> cat > ./my.ini <<"EOF"
[My book]
title=A wonderful book
author[]=Bitnami
author[]=Contributors
EOF

# get retrieves the first mention of the key
$> ini-file get --section "My book" --key "author[]" ./my.ini
Bitnami

# If the original file contains the key more than once, set adds the new value at the end
$> ini-file set --section "My book" --key "author[]" --value "Other" ./my.ini
$> cat ./my.ini
title=A wonderful book
author[]=Bitnami
author[]=Contributors
author[]=Other

# del removes all keys with the given name
$> ini-file del --section "My book" --key "author[]" ./my.ini
$> cat ./my.ini
title=A wonderful book

~~~


