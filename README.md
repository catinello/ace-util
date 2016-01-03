ace-util
===

Command line utility for the [Ace][ACE] HTML template engine.


##Usage:##

    Usage:
      ace [-i | --inner=<FILE>] [-m | --map=<FILE>] [-s | --separator=<CHAR>] [-p | --stdout] [-o | --output=<FILE>] [-w | --httpd] <FILE>
      ace [-h | --help]
      ace [-v | --version]
    Options:
      -i --inner=<FILE>     Path to the inner.ace file.
      -m --map=<FILE>       Path to the mappings.map file.
      -s --separator=<CHAR> Separator for key/value map file.
      -p --stdout           Print to stdout.
      -o --output=<FILE>    Write to custom file.
      -w --httpd            Start temporary webserver.
      -h --help             Show this help.
      -v --version          Display version.


##Examples:##

Easy development feature (builtin local-webserver):

    $ ace -i example/inner.ace -m example/mappings.map example/base.ace -w

Open your browser and visit http://127.0.0.1:8080 to see a preview of your generated template.

Simple call of an Ace template:

    $ ace example/base.ace

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

Result:

    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="utf-8">
        <title>Actions</title>
      </head>
      <body>
        <h1>Actions</h1>
        <ul>
          <li>Message1</li>
          <li>Message2</li>
          <li>Message3</li>
        </ul>
        <h2>This is a content named "main" of an inner template.</h2>
        <div>&lt;div&gt;Escaped String&lt;/div&gt;</div>
        <div>
          <div>Non-Escaped String</div>
        </div>
        <h2>This is a content named "sub" of an inner template.</h2>
      </body>
    </html>


##License:##

[&copy; Antonino Catinello][HOME] - [MIT-License][MIT]

[MIT]:https://github.com/catinello/ace-util/blob/master/LICENSE
[HOME]:http://antonino.catinello.eu
[ACE]:https://github.com/yosssi/ace
