import commonjs from "@rollup/plugin-commonjs";
import resolve from "@rollup/plugin-node-resolve";
import typescript from "@rollup/plugin-typescript";
import copy from "rollup-plugin-copy";
import cssimport from "rollup-plugin-cssimport";
import minifyHTML from "rollup-plugin-minify-html-literals";
import { terser } from "rollup-plugin-terser";

const resources = [
    { src: "src/style.css", dest: "./dist" },
    { src: "assets/*", dest: "./dist/assets" },
];

module.exports = [
    {
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
            resolve({ browser: true }),
            commonjs(),
            typescript(),
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
        preserveEntrySignatures: false,
        cache: true,
        context: "window",
    },
];
