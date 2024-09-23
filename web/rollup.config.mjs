import commonjs from "@rollup/plugin-commonjs";
import json from "@rollup/plugin-json";
import { nodeResolve } from "@rollup/plugin-node-resolve";
import copy from "rollup-plugin-copy";
import cssimport from "rollup-plugin-cssimport";
import esbuild from "rollup-plugin-esbuild";

export const extensions = [".js", ".jsx", ".ts", ".tsx"];

export const resources = [
    {
        src: "node_modules/rapidoc/dist/rapidoc-min.js",
        dest: "dist/",
    },
    {
        src: "node_modules/@patternfly/patternfly/patternfly.min.css",
        dest: "dist/",
    },
    {
        src: "node_modules/@patternfly/patternfly/patternfly-base.css",
        dest: "dist/",
    },
    {
        src: "node_modules/@patternfly/patternfly/assets/*",
        dest: "dist/assets/",
    },
    { src: "src/elements/styles/gravity.css", dest: "dist/" },
    { src: "src/assets/*", dest: "dist/assets" },
    { src: "./icons/*", dest: "dist/assets/icons" },
];

// eslint-disable-next-line no-undef
export const isProdBuild = process.env.NODE_ENV === "production";

export default {
    input: "./src/main.ts",
    output: [
        {
            format: "es",
            dir: "./dist/",
            sourcemap: true,
        },
    ],

    plugins: [
        cssimport(),
        json(),
        nodeResolve({ extensions, browser: true, preferBuiltins: false }),
        commonjs(),
        esbuild({
            minify: isProdBuild,
        }),
        copy({
            targets: [...resources],
            copyOnce: false,
        }),
    ],
    watch: {
        clearScreen: false,
    },
    preserveEntrySignatures: "strict",
    cache: true,
    context: "window",
    onwarn: function (warning, warn) {
        if (warning.code === "UNRESOLVED_IMPORT") {
            throw Object.assign(new Error(), warning);
        }
        warn(warning);
    },
};
