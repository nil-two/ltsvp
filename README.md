ltsvp
=====

[![Build Status](https://travis-ci.org/kusabashira/ltsvp.svg?branch=master)](https://travis-ci.org/kusabashira/ltsvp)

Print selected parts of LTSV from each FILE to standard output.

```
$ cat items.ltsv
name:Apple	price:60	quantity:20
name:Grapes	price:140	quantity:8
name:Pineapple	price:400	quantity:2
name:Orange	price:50	quantity:14

$ cat items.ltsv | ltsvp -k price,quantity
60	20
140	8
400	2
50	14
```

Usage
-----

```
$ ltsvp OPTION... [FILE]...

Options:
  -k, --keys=LIST
                 select only these keys (required)
  -D, --output-delimiter=STRING
                 use STRING to separate parts (default: \t)
  -r, --remain-ltsv
                 print selected parts as LTSV
  -h, --help
                 display this help text and exit
  -v, --version
                 output version information and exit

LIST is made up of keys separated by commas.
  host           # Select host
  host,time,ua   # Select host, time, and ua
```

Installation
------------

### compiled binary

See [releases](https://github.com/kusabashira/ltsvp/releases)

### go get

```
go get github.com/kusabashira/ltsvp
```

Options
-------

### -h, --help

Display a help message.

### -v, --version

Output the version of ltsvp.

### -k, --keys=LIST

Select only specified keys.

Keys separated by a `,`.

```sh
# select only host
ltsvp --keys=host

# select only host, time, and ua
ltsvp --keys=name,time,ua

# select only "foo,bar" and "baz"
ltsvp --keys="foo\,bar,baz"
```

#### syntax of keys list

Here is the syntax of keys in extended BNF.

```
keys = key , { "," , key } ;
key  = { [ "\" ] , ? unicode character ? - "," | "\," } ;
```

### -D, --output-delimiter=STRING

Change the output delimiter to `STRING`.
`STRING` is unicode characters.

```sh
# Outputs with a slash delimited
ltsvp --keys=time,host --output-delimiter=/

# Outputs with a "::" delimited
ltsvp --keys=time,host --output-delimiter=::
```

### -r, --remain-ltsv

Output selected parts in LTSV format.

```sh
# Output as LTSV
ltsvp --keys=time,host --remain-ltsv
```

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
