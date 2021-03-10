# server-fastify

This example demonstrates a simple fastify server with sørvør as build system.

## prerequisites

`sørvør` should be installed, preferably using `npm` or `yarn`:

```bash
npm install sorvor
# or
yarn add sorvor
```

## available commands

`yarn start` - builds the `src/server.js` with esbuild and launches the resulting output file to start fastify server.

> On every change in source code, sørvør rebuilds the `src/server.js` and restarts the fastify server

`yarn build` - build the project with minification and sourcemaps support.
