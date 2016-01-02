ace-util
===

Command line utility for the Ace HTML template engine.

##Usage:##

    Usage:
      ace [-i | --inner=<FILE>] [-m | --map=<FILE>] [-s | --separator=<SYMBOL>] [-p | --stdout] [-o | --output=<FILE>] <FILE>
      ace [-h | --help]
      ace [-v | --version]
    Options:
      -i --inner            Path to the inner.ace file.
      -m --map              Path to the mappings.map file.
      -s --separator        Separator for key/value map file.
      -p --stdout           Print to stdout.
      -o --output           Write to custom file.
      -h --help             Show this help.
      -v --version          Display version.

##Examples:##

Simple call of an Ace template (with or without .ace suffix):

    $ ace example/base

Creates the corresponding file base.html in pwd.


To fill in the variables {{.Title}} and {{.Msgs}} where are going add a map:

    $ ace -m example/mappings.map example/base.ace

The mappings.map file content is paresed per line and separates key/values by the given unicode separator which defaults to the middle dot (U+00B7).


You are able to use the inner.ace as well:

    $ ace -i example/inner.ace -m example/mappings.map example/base.ace

##License:##

[&copy; Antonino Catinello][HOME] - [MIT-License][MIT]

[MIT]:https://github.com/catinello/ace-util/blob/master/LICENSE
[HOME]:http://antonino.catinello.eu
