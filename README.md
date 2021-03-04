# sørvør

> lightning fast, zero config build tool for modern Javascript and Typescript.

## Introduction

`sørvør` is a lightweight replacement for heavier, more complex build tools such as webpack or parcel. It is an opinionated take on [`esbuild`](https://esbuild.github.io/) with sane default and plugins to enhance your development workflow.

## Major Features

- **Instant Startup**: `sørvør` is authored in [golang](https://golang.org/), which offers the best startup time for command line applications. Often times, `sørvør` will finish bundling your project by the time a `node.js` bundler starts loading.
- **Easy Installation**: `sørvør` is distributed as a single binary for all major platforms. It's just one command to install `sørvør` using installation method of your choice.
- **Optimize for Production**: Build for production with built-in optimizations and Asset Pipeline.

## Installation

The easiest method to install sørvør is using [go binaries](https://gobinaries.com/):

```bash
curl -sf https://gobinaries.com/osdevisnot/sorvor | sh
```

<details>
  <summary>See other installation methods</summary>

Alternatively, if you have [go](https://golang.org/) installed, use `go get` to install sørvør:

```bash
go get github.com/osdevisnot/sorvor
```

sørvør can also be installed using NPM or yarn package manager:

```bash
npm install sorvor
# or
yarn add sorvor
```

</details>

Once installed, verify your installation using `version` flag

```bash
sorvor --version
```

## Quickstart

You can always refer [example projects](examples) for fully integrated setup using `sørvør`. To get started, let's set up a simple Preact application using `sørvør`. First, create a minimal scaffold using `degit`:

```bash
npx degit osdevisnot/sorvor-minimal preact-hello
```

The minimal scaffold comes with a README.md file with short description of available commands. Let's start the live reloading server using the command from README.md:

```bash
sorvor --serve
```

This should bundle your project and start a live reloading server at [http://localhost:1234](http://localhost:1234).

Now, let's add a simple Preact Component which renders Hello World.

```js
import { h, render } from "https://unpkg.com/preact@10.5.12/dist/preact.module.js";

const Counter = () => <div>Hello World</div>;

render(<Counter />, document.body);
```

You should notice an error on terminal which should look like this:

```log
2021/02/23 09:55:58 Warn: Unexpected "<"
```

This error indicates sørvør was unable to parse JSX syntax we just entered. Let's restart sørvør to be able to parse JSX Syntax:

```bash
sorvor --serve --loader:.js=jsx --jsx-factory=h
```

On restart, your browser should display Hello World rendered using preact.

## Node.js Usage

You can use `sørvør` to bundle browser based applications, but it is equally suitable for node.js applications as well. Let's try to build a simple express server using `sørvør` as build tool.

First of, create a project directory and install `express` as a dependency:

```bash
mkdir hello-server
cd hello-seerver
npm init --yes
npm install express
```

Now, let's create a `server.js` file in src directory, and paste below code in it:

```js
import express from "express";

const app = express();
const port = 3000;

app.get("/", (req, res) => {
  res.send("Hello World !!");
});

app.listen(port, () => {
  console.log(`App Ready on http://localhost:${port}`);
});
```

Now, let's start the build using

```bash
sorvor src/server.js --serve --platform=node
```

The platform node automatically adds all the dependencies from `package.json` as external. The `serve` when combined with node platform, builds the entry file and starts running it once the build is complete. Try visiting [http://localhost:3000](http://localhost:3000) to see "Hello World" rendered in the browser.

## Asset Pipeline

The asset pipeline for `sørvør` is partially inspired by HUGO pipes. To use the asset pipeline, use below annotations in your HTML entrypoint.

### Live Reload

`sørvør` supports live reloading connected browsers based on [SSE Web API](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events).

To use this feature, include the livereload annotation in your HTML entrypoint like this:

```html
<html>
  ...
  <body>
    {{ livereload }}
  </body>
  ...
</html>
```

### Bundle with `esbuild`

To run an entrypoint through `esbuild`, use `esbuild` annotation in your HTML entrypoint like this:

```html
<script type="module" src="{{ esbuild "index.js" }}"></script>
```

## Plugins

`sørvør` enhances esbuild with few quality of life plugins. These plugins are enabled by default and require no configuration for usage.

### esm plugin

This plugin allows you to import HTTP URLs into JavaScript code. The code will be automatically downloaded at build time. For production builds, the resolved file will be bundled together with your code.

With this plugin, you can import your ESM dependencies from CDN urls **without** installing it via NPM:

```js
import { zip } from "https://unpkg.com/lodash-es@4.17.15/lodash.js";
console.log(zip([1, 2], ["a", "b"]));
```

### env plugin

The env plugin imports the current environment variables at build time. You can use the environment variables like this:

```js
import { PATH } from "env";
console.log(`PATH is ${PATH}`);
```

## Configuration

For most part, `sørvør` tries to use sensible defaults, but you can configure the behaviour using command line arguments below:

| cli argument | description                | default value |
| ------------ | -------------------------- | ------------- |
| `--host=...` | host for sørvør            | `localhost`   |
| `--port=...` | port for sørvør            | `1234`        |
| `--serve`    | enable development mode    | `false`       |
| `--secure`   | use https with `localhost` | `false`       |

> `sørvør` forwards all the other command line arguments to `esbuild`.

> `--secure` automatically creates a self-signed certificate for provided host.

> to disable chrome warnings, open chrome://flags/#allow-insecure-localhost and change the setting to "Enabled".

Please refer documentation for [simple esbuild options](https://esbuild.github.io/api/#simple-options) or [advance options](https://esbuild.github.io/api/#advanced-options) to further customize the bundling process.

`sørvør` configures below values for esbuild as defaults which you can override using command line arguments:

| cli argument   | description                         | default value |
| -------------- | ----------------------------------- | ------------- |
| `--outdir=...` | target directory for esbuild output | `dist`        |

`sørvør` configures below values for esbuild which can not be changed:

| cli argument | description                          | default value                                 |
| ------------ | ------------------------------------ | --------------------------------------------- |
| `--bundle`   | enables bundling output files        | `true`                                        |
| `--write`    | enables writing built output to disk | `true`                                        |
| `--define`   | value for `process.env.NODE_ENV`     | `production` and `development` with `--serve` |

## Inspirations

This project is inspired by [servør](https://www.npmjs.com/package/servor) from [Luke Jackson](https://twitter.com/lukejacksonn), which provides similar zero dependency development experience but lacks integration with a bundler/build tools.

## License

**sørvør** is licensed under the [MIT License](http://opensource.org/licenses/MIT).

Documentation is licensed under [Creative Commons License](http://creativecommons.org/licenses/by/4.0/).

Created with ❤️ by [@osdevisnot](https://github.com/osdevisnot) and [all contributors](https://github.com/osdevisnot/sorvor/graphs/contributors).
