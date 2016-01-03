ace-util
===

Command line utility for the [Ace][ACE] HTML template engine.


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

Creates the corresponding file base.html in your **${PWD}**. You can change the location to store through the *--output* flag.


To fill the variables **{{.Title}}** and **{{.Msgs}}** we are going to add a map:

    $ ace -m example/mappings.map example/base.ace

The *mappings.map* file content is parsed per line and is separated to key/values by the given unicode separator which defaults to the middle dot (U+00B7). You can customize this behaviour using the *--separator* flag to use a **$** symbol instead of the middle dot or to whatever you want.

    $ ace -s $ -m example/mappings.map example/base.ace


The first entry is always the keyname! Two value types are available for the map. **string** as single value or **[]string** as multivalue:

    Title·Actions
    Msgs·Message1·Message2·Message3


You are able to use the inner.ace as well:

    $ ace -i example/inner.ace -m example/mappings.map example/base.ace


##License:##

[&copy; Antonino Catinello][HOME] - [MIT-License][MIT]

[MIT]:https://github.com/catinello/ace-util/blob/master/LICENSE
[HOME]:http://antonino.catinello.eu
[ACE]:https://github.com/yosssi/ace
