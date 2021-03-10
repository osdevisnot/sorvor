# sorvor-minimal

This example demonstrates usage of sørvør with minimal boilerplate.

The `public/index.html` file includes an annotated reference to `public/index.js`.

The `esbuild` annotation in `public/index.html` instructs `sørvør` to build `public/index.js` file with `esbuild`.

## prerequisites

`sørvør` should be installed, preferably using `npm` or `yarn`:

```bash
npm install sorvor
# or
yarn add sorvor
```

## available commands

`sorvor --minify --sourcemap` - build the project with minification and sourcemaps

`sorvor --serve` - starts a live reload dev server at [http://localhost:1234](http://localhost:1234), and rebuilds the index.js file on change.
