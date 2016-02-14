## v0.4.2 - 2016-02-01

- Fix typo of short option for --output-delimiter.

## v0.4.1 - 2016-01-11

- Disallow empty LIST.

## v0.4.0 - 2016-01-11

- Release compiled binary for Windows, OSX, and Linux.
- Allow mixed flag like "-rk host,time".
- Revive short option for --help.
- Revive short option for --version.
- Interpret equal sign between short option and value as a part of value.
  - later:   "-k=host,time" #=> LIST = ["=host", "time"]
  - earlier: "-k=host,time" #=> LIST = ["host", "time"]

## v0.3.0 - 2015-11-18

- Rename -d, --delimiter to -D, --output-delimiter.
- Remove short option for --help.
- Remove short option for --version.
- Change the format of version from "v0.3.1" to "0.3.1".

## v0.2.0 - 2015-10-21

- Require specifying LIST with -k, --keys.
- Allow changing output-delimiter with -d, --delimiter.
- Allow outputting as LTSV with -r, --remain-ltsv.
- Allow specifying input files before options.

## v0.1.0 - 2015-09-29

- Initial release.
