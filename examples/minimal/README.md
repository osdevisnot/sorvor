# minimal

This example demonstrates usage of sørvør with minimal boilerplate.

The public/index.html file includes an annotated reference to public/index.js.

The `esbuild` annotation in public/index.html instructs sorvor to build public/index.js file with esbuild.

## prerequisites

sørvør should be installed, preferably using Use [go binaries](https://gobinaries.com/):

```bash
curl -sf https://gobinaries.com/osdevisnot/sorvor | sh
```

## available commands

`sorvor` - build the project with minification and sourcemaps

`sorvor --serve` - starts a live reload dev server at [http://localhost:1234](http://localhost:1234), rebuilds the index.js file on change.

> Note: To enable live reload, sørvør needs a livereload annotation as shown in other examples.
