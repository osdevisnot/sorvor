# server-fastify

This example demonstrates a simple fastify server with sørvør as build system.

## prerequisites

sørvør should be installed, preferably using Use [go binaries](https://gobinaries.com/):

```bash
curl -sf https://gobinaries.com/osdevisnot/sorvor | sh
```

## available commands

`yarn start` - builds the `src/server.js` with esbuild and launches the resulting output file to start fastify server.

> On every change in source code, sørvør rebuilds the `src/server.js` and restarts the fastify server

`yarn build` - build the project with minification and sourcemaps support.
