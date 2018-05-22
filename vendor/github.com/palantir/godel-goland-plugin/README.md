godel-goland-plugin
===================
godel-goland-plugin is a gödel plugin that generates [Goland](https://www.jetbrains.com/go/) project files for a Go project.

Plugin Tasks
------------
godel-goland-plugin provides the `goland` task, which generates GoLand project files (`.iml` and `.ipr` files) for the project. It assumes that a global Go SDK named "Go" has been defined. It also creates a file watchers task that applies the `./godelw format` on modified files on save. The `goland` task has a `clean` subcommand that removes all of the generated project files.
