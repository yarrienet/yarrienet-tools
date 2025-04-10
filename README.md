# yarrienet-tools

Repository containing tools used to maintain [yarrie.net](http://yarrie.net).

## Structure

The project is split into two parts:

1. `main.go` is the main command line tool which performs common tasks by command.
2. `tools/` and `scripts/` contain specific tools which are only intended to be ran once, e.g. migrations. See *Tools & Scripts* section below.

Most of the code in either of two is specialized however most of the module code is generic and potentially reusable in any project.

## Usage

1. Ensure you have the Go toolchain.
2. `go mod tidy`
3. `go run main.go` to print command help.

### CLI

Usage: `yarrienet [command] [subcommand]`

## Configuration

An optional configuration file can be placed in `~/.config/yarrienet.conf` to replace arguments used repeatedly across commands.

```
# path of the microblog html file
microblog_html_file "~/Documents/yarrie.net/microblog/index.html"
```

## Tools & Scripts

The following lists the one-time tools and scripts not included in the main tool.

### `tools/insertdates.go`

The tool uses the microblog's RSS feed to extract the date of each post and then insert said date semantically into the HTML tree. Created in order to migrate from using a static site generator backed by a database to a semantically storing data within the page's source. Requires the microblog's HTML and RSS files.

```sh
insertdates <microblog file> <rss file> [output]
```

