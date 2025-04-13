# yarrienet-tools

Repository containing tools used to maintain [yarrie.net](http://yarrie.net).

## Semantic Publishing

These tools are built to achieve semantic publishing where writing occurs directly in the webpage's HTML source and then all auto-generated elements are built against that. In practice a single webpage containing each post can be used to produce the associated RSS feed without the backing of an abstract file tree or database.

HTML does not need to be abstracted to produce working RSS unlike other static site generators.

The tool relies heavily on a set schema utilizing semantic elements and as such tools are included to insert the HTML source of schema friendly elements directly into the webpage to be modified.

### Microblog

For example, each post on the [yarrie.net microblog](http://yarrie.net/microblog) uses the following schema:

```html
<div id="posts">
    <div class="post" id="exampleid">
        <div class="date">
            <a href="#exampleid" class="post-link"><time datetime="2025-04-10T17:38:10+01:00"><p>april 10, 2025</p></time></a>
        </div>
        <p></p>
    </div>
    <!--- ... -->
</div>
```

Using the ID from the post, date from the `<time datetime>` and remaining content after the date, the tool is able to determine and produce all necessary information for a valid RSS entry.

## Usage

1. Ensure you have the Go toolchain.
2. `go mod tidy`
3. `go run main.go` to print command help.

### CLI

Usage: `yarrienet [command] [subcommand] <options>`

## Configuration

An optional configuration file can be placed in `~/.config/yarrienet.conf` to replace arguments used repeatedly across commands.

```
# path of the microblog html file
microblog_html_file "~/Documents/yarrie.net/microblog/index.html"
# path of the microblog rss file
microblog_rss_file "~/Documents/yarrie.net/microblog/rss.xml"
```

## Structure

The project is split into two parts:

1. `main.go` is the main command line tool which performs common tasks by command.
2. `tools/` and `scripts/` contain specific tools which are only intended to be ran once, e.g. migrations. See *Tools & Scripts* section below.

Most of the code in either of two is specialized however most of the module code is generic and potentially reusable in any project.

### Tools & Scripts

The following lists the one-time tools and scripts not included in the main tool.

#### `tools/insertdates.go`

The tool uses the microblog's RSS feed to extract the date of each post and then insert said date semantically into the HTML tree. Created in order to migrate from using a static site generator backed by a database to a semantically storing data within the page's source. Requires the microblog's HTML and RSS files.

```sh
insertdates <microblog file> <rss file> [output]
```

