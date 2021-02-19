# minimal-typescript

This example demonstrates usage of sørvør with minimal boilerplate for Typescript projects.

The public/index.html file includes an annotated reference to public/index.ts.

The `esbuild` annotation in public/index.html instructs sorvor to build public/index.ts file with esbuild.

## prerequisites

sørvør should be installed, preferably using Use [go binaries](https://gobinaries.com/):

```bash
curl -sf https://gobinaries.com/osdevisnot/sorvor | sh
```

## available commands

`sorvor` - build the project with minification and sourcemaps

`sorvor --serve` - starts a live reload dev server at [http://localhost:1234](http://localhost:1234), rebuilds the index.ts file on change and live reloads all connected browsers.
