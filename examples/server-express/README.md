# server-express

This example demonstrates a simple express server with sørvør as build system.

## prerequisites

`sørvør` should be installed, preferably using `npm` or `yarn`:

```bash
npm install sorvor
# or
yarn add sorvor
```

## available commands

`yarn start` - builds the `src/server.js` with esbuild and launches the resulting output file to start express server.

> On every change in source code, sørvør rebuilds the `src/server.js` and restarts the express server

`yarn build` - build the project with minification and sourcemaps support.
