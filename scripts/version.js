#!/usr/bin/env node

const path = require("path");
const fs = require("fs");
const cmd = require("child_process").execSync;
const npm = path.join(__dirname, "..", "npm");

const version = process.argv[2];

if (!version) {
  throw new Error(`no version supplied.`);
}

for (let dir of [
  "sorvor",
  "sorvor-darwin-64",
  "sorvor-darwin-arm64",
  "sorvor-freebsd-64",
  "sorvor-freebsd-arm64",
  "sorvor-linux-32",
  "sorvor-linux-64",
  "sorvor-linux-arm",
  "sorvor-linux-arm64",
  "sorvor-windows-32",
  "sorvor-windows-64",
]) {
  const pkg = require(path.join(npm, dir, "package.json"));
  pkg.version = version;
  if (pkg.optionalDependencies) {
    for (let dep in pkg.optionalDependencies) {
      pkg.optionalDependencies[dep] = version;
    }
  }
  fs.writeFileSync(
    path.join(npm, dir, "package.json"),
    JSON.stringify(pkg, null, "  "),
    "utf8"
  );
  cmd("npm publish", { stdio: "pipe", cwd: path.join(npm, dir) });
}
