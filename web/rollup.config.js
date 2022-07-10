import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import cssimport from "rollup-plugin-cssimport";
import copy from "rollup-plugin-copy";
import minifyHTML from "rollup-plugin-minify-html-literals";
import { terser } from "rollup-plugin-terser";

const resources = [
    { src: "assets/*", dest: "./dist/assets" },
];

module.exports = [
    {
        input: "./src/main.js",
        output: [
            {
                format: "es",
                dir: "./dist/",
                sourcemap: true,
            },
        ],
        plugins: [
            resolve({ browser: true }),
            commonjs(),
            cssimport(),
            process.env.NODE_ENV === "production" && minifyHTML(),
            process.env.NODE_ENV === "production" && terser(),
            copy({
                targets: [...resources],
                copyOnce: false,
            }),
        ],
        watch: {
            clearScreen: false,
        },
    },
];
