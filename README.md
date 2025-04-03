## yarrienet-tools

Repository containing tools used to maintain [yarrie.net](http://yarrie.net).

### Structure

Most of the tools are specialized and built to complete a highly specific task, however most of the code outside of `tools`, `scripts`, and `main.go` is generic and potentially helpful for any project.

### Tools

#### `insertdates.go`

The tool uses the microblog's RSS feed to extract the date of each post and then insert said date semantically into the HTML tree. Created in order to migrate from using a static site generator backed by a database to a semantically storing data within the page's source.

```sh
insertdates <microblog file> <rss file> [output]
```

