# minimal-css

This example demonstrates usage of sørvør with minimal boilerplate for CSS.

The `public/index.html` file includes an annotated reference to `public/index.css`.

The `esbuild` annotation in `public/index.html` instructs `sørvør` to build `public/index.css` file with esbuild.

## prerequisites

`sørvør` should be installed, preferably using `npm` or `yarn`:

```bash
npm install sorvor
# or
yarn add sorvor
```

## available commands

`sorvor` - build the project with minification and sourcemaps

`sorvor --serve` - starts a live reload dev server at [http://localhost:1234](http://localhost:1234), rebuilds the `index.css` file on change and live reloads all connected browsers.
