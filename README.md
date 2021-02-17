# sørvør

> extremely fast, zero config server for modern web applications.

## :sparkles: Features

- **Flexible EntryPoints** - use HTML, CSS or JS as entry point for your application.
- **SPA Routing** - redirects path requests to HTML entry point for frontend routing.
- **Asset Pipeline** - great asset processing with simple primitives.
- **Live Reloading** - live reload browsers on code change.
- **Bundle Libraries** - bundle libraries for distribution using JS/TS/JSX/TSX entry points.
- **Secure Server** - Supports https with trusted self signed certificates.

### :muscle: Powered By

- [esbuild](https://esbuild.github.io/) - an extremely fast JavaScript bundler.
- [golang](https://golang.org/) - an expressive, concise, clean, and efficient programming language.

## :zap: Installation

Use [go binaries](https://gobinaries.com/) to install sørvør:

```bash
curl -sf https://gobinaries.com/osdevisnot/sorvor | sh
```

Alternatively, if you have [go](https://golang.org/) installed, use `go get` to install sørvør:

```bash
go get github.com/osdevisnot/sorvor
```

## :plate_with_cutlery: Usage

You can use `sørvør` as a local development server or as a build tool for your applications and/or NPM libraries.

### Live Reloading Server

To serve an application using a live reloading server, use HTML as entrypoint and `--dev` command line argument.

For example:

```bash
sorvor public/index.html --dev
```

### Build Applications for Production

Alternatively, you can build your application with `sørvør` using HTML as entrypoint. For Example:

```bash
sorvor public/undex.html
```

### Bundle NPM Library

You can also bundle your library for distribution on NPM using a JS entrypoint. For example:

```bash
sorvor src/index.js
```

## Example Projects

Check out the [example projects](examples) for fully integrated setup.

## :sunglasses: Asset Pipeline

`sørvør` provides great asset pipeline with simple primitives.

#### Build JS or CSS with esbuild

To run entry points from `public/index.html` through esbuild, use `esbuild` function in the index file

Example:

```html
<script type="module" src="{{ esbuild "index.js" }}"></script>
```

This will bundle `index.js` file and serve the build output on local development server.

#### Enable Livereload

To enable livereload functionality, use `livereload` function in the index file

Example:

```
{{ livereload }}
```

## :anger: Configuration

For most part, `sørvør` tries to use sensible defaults, but you can configure the behaviour using command line arguments below:

| cli argument | description             | default value |
| ------------ | ----------------------- | ------------- |
| `--host=...` | host for sørvør         | `localhost`   |
| `--port=...` | port for sørvør         | `1234`        |
| `--dev`      | enable development mode | `false`       |
| `--secure`   | use https in dev mode   | `false`       |

> `sørvør` forwards all the other command line arguments to `esbuild`.

Note: `--secure` automatically creates a self signed certificate for provided host.

> to disable chrome warnings, open chrome://flags/#allow-insecure-localhost and change the setting to "Enabled".

Please refer documentation for [simple esbuild options](https://esbuild.github.io/api/#simple-options) or [advance options](https://esbuild.github.io/api/#advanced-options) to further customize the bundling process.

`sørvør` configures below values for esbuild as defaults:

| cli argument   | description                          | default value |
| -------------- | ------------------------------------ | ------------- |
| `--bundle`     | enables bundling output files        | `true`        |
| `--write`      | enables writing built output to disk | `true`        |
| `--outdir=...` | target directory for esbuild output  | `dist`        |

## :hatching_chick: Motivations/Inspirations

`sørvør` started with desire to simplify frontend tooling, with strong focus on speed of execution. It uses `esbuild` for bundling modern javascript and typescript syntax to a lower target. The idea here is to implement features that `esbuild` deems as out of scope, but are necessary for a decent development environment.

This project is inspired by [servør](https://www.npmjs.com/package/servor) from [Luke Jackson](https://twitter.com/lukejacksonn), which provides similar zero dependency development experience but lacks integration with bundler/build tools. I chose golang to implement this project to solidify my learning of the language and to achieve a zero dependency model.

## :microscope: Roadmap

This project currently lacks some extended features available in [servør](https://www.npmjs.com/package/servor), some of which will be implemented in the future.

I also want to avoid implementing features that are already on the roadmap for `esbuild`. The idea is to use esbuild as is without duplicating efforts.

## :clinking_glasses: License

**sørvør** is licensed under the [MIT License](http://opensource.org/licenses/MIT).

Documentation is licensed under [Creative Commons License](http://creativecommons.org/licenses/by/4.0/).

Created with ❤️ by [@osdevisnot](https://github.com/osdevisnot) and [all contributors](https://github.com/osdevisnot/sorvor/graphs/contributors).
