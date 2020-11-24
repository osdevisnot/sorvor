# sørvør

> fast, zero config server for single page applications.

## :sparkles: Features

- **HTML EntryPoint** - use `src/index.html` as an entry point for an application.
- **SPA Routing** - redirects path requests to `src/index.html` for frontend routing.
- **Asset Pipeline** - strong asset processing with simple primitives.
- **Live Reloading** - reloads the browsers on code change.

### :muscle: Powered By

- [esbuild](https://esbuild.github.io/) - an extremely fast JavaScript bundler.
- [golang](https://golang.org/) - an expressive, concise, clean, and efficient programming language.

## :zap: Installation

Use [gobinaries](https://gobinaries.com/) to install sørvør:

```bash
curl -sf https://gobinaries.com/osdevisnot/sorvor/cmd/sorvor | sh
```

Alternatively, if you have [go](https://golang.org/) installed, use `go get` to install sørvør:

```bash
go get github.com/osdevisnot/sorvor
```

## :plate_with_cutlery: Usage

You can use `sørvør` as a local development server or as a build tool. By default, the `sørvør` command will build your project and exit. To start a local development server, pass `--dev` as a command line argument.

```bash
sorvor --dev
```

## :sunglasses: Asset Pipeline

`sørvør` provides strong asset pipeline with simple premitives.

Currently, only [esbuild](https://esbuild.github.io/) assets are supported.

For Example: configure `index.html` to use `esbuild` bundling for `index.js`

```html
<script type="module" src="{{ esbuild "index.js" }}"></script>
```

## :anger: Configuration

For most part, `sørvør` tries to use sensible defaults, but you can configure the behaviour using command line arguments below:

| cli argument | description                 | default value |
| ------------ | --------------------------- | ------------- |
| `--src=...`  | source directory for sørvør | `src`         |
| `--port=...` | port for sørvør             | `1234`        |
| `--dev`      | enable development mode     | `false`       |

`sørvør` forwards all the other command line arguments to `esbuild`. Please refer documentation for [simple esbuild options](https://esbuild.github.io/api/#simple-options) or [advance options](https://esbuild.github.io/api/#advanced-options) to further customize your builds.

For example, to use `esbuild` with modern `esm` format, use a command like this:

```bash
sorvor --format=esm --dev
```

`sørvør` configures below values for esbuild as defaults:

| cli argument   | description                          | default value                     |
| -------------- | ------------------------------------ | --------------------------------- |
| `--bundle`     | enables bundling output files        | `true`                            |
| `--write`      | enables writing built output to disk | `true`                            |
| `--port=...`   | port to start esbuild in serve mode  | `1234` (if --dev mode is enabled) |
| `--outdir=...` | target directory for esbuild output  | `dist`                            |

## :clinking_glasses: License

**sørvør** is licensed under the [MIT License](http://opensource.org/licenses/MIT).

Documentation is licensed under [Creative Commons License](http://creativecommons.org/licenses/by/4.0/).

Created with ❤️ by [@osdevisnot](https://github.com/osdevisnot) and [all contributors](https://github.com/osdevisnot/sorvor/graphs/contributors).
