ltsvp
=====

Print selected parts of LTSV from each FILE to standard output.

Usage
-----

```
$ ltsvp OPTION... [FILE]...

Options:
  -k, --keys=LIST          select only these keys
                           LIST is made up of keys separated by commas
                             host           # Select host
                             host,time,ua   # Select host, time, and ua
  -d, --delimiter=STRING   use STRING to separate parts (default: \t)
  -r, --remain-ltsv        print selected parts as LTSV
  -h, --help               display this help text and exit
  -v, --version            display version information and exit
```

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
