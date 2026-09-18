import { cpSync, readFileSync } from "node:fs";
import * as esbuild from "esbuild";

export const isProdBuild = process.env.NODE_ENV === "production";
const isWatch = process.argv.includes("--watch");

// Mirrors rollup-plugin-cssimport: turn a CSS import into a constructable
// CSSStyleSheet, which is how lit's `static styles` consumes patternfly/gravity css.
const cssStyleSheetPlugin = {
    name: "css-stylesheet",
    setup(build) {
        build.onLoad({ filter: /\.css$/ }, (args) => {
            const text = readFileSync(args.path, "utf8");
            return {
                contents: `const sheet = new CSSStyleSheet();\nsheet.replaceSync(${JSON.stringify(text)});\nexport default sheet;\n`,
                loader: "js",
            };
        });
    },
};

const resources = [
    ["node_modules/rapidoc/dist/rapidoc-min.js", "dist/rapidoc-min.js"],
    ["node_modules/@patternfly/patternfly/patternfly.min.css", "dist/patternfly.min.css"],
    ["node_modules/@patternfly/patternfly/patternfly-base.css", "dist/patternfly-base.css"],
    ["src/elements/styles/gravity.css", "dist/gravity.css"],
    ["node_modules/@patternfly/patternfly/assets", "dist/assets"],
    ["src/assets", "dist/assets"],
];

function copyResources() {
    for (const [src, dest] of resources) {
        cpSync(src, dest, { recursive: true });
    }
}

const buildOptions = {
    entryPoints: ["./src/main.ts"],
    outdir: "./dist/",
    bundle: true,
    splitting: true,
    format: "esm",
    platform: "browser",
    sourcemap: true,
    minify: isProdBuild,
    plugins: [cssStyleSheetPlugin],
    loader: { ".json": "json" },
    logLevel: "info",
};

if (isWatch) {
    const ctx = await esbuild.context(buildOptions);
    copyResources();
    await ctx.watch();
} else {
    await esbuild.build(buildOptions);
    copyResources();
}
