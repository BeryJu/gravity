import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { javascript } from "@codemirror/lang-javascript";
import { LanguageSupport, StreamLanguage } from "@codemirror/language";
import * as yamlMode from "@codemirror/legacy-modes/mode/yaml";
import { Compartment, EditorState, Extension } from "@codemirror/state";
import { EditorView, drawSelection, keymap, lineNumbers } from "@codemirror/view";
import { vsCodeDark } from "@fsegurai/codemirror-theme-vscode-dark";
import { vsCodeLight } from "@fsegurai/codemirror-theme-vscode-light";
import YAML from "yaml";

import { customElement, property } from "lit/decorators.js";

import { AKElement } from "./Base";

@customElement("ak-codemirror")
export class CodeMirrorTextarea extends AKElement {
    @property({ type: Boolean })
    readOnly = false;

    @property()
    mode = "yaml";

    @property({ reflect: true })
    name: string | undefined;

    editor?: EditorView;

    _value?: string;

    theme = new Compartment();
    syntaxHighlight = new Compartment();

    themeLight = vsCodeLight;
    themeDark = vsCodeDark;

    @property()
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    set value(v: any) {
        if (v === null || v === undefined) return;
        // Value might be an object if within an iron-form, as that calls the getter of value
        // in the beginning and the calls this setter on reset
        let textValue = v;
        if (!(typeof v === "string" || v instanceof String)) {
            switch (this.mode.toLowerCase()) {
                case "yaml":
                    textValue = YAML.stringify(v);
                    break;
                default:
                    textValue = v.toString();
                    break;
            }
        }
        if (this.editor) {
            this.editor.dispatch({
                changes: { from: 0, to: this.editor.state.doc.length, insert: textValue },
            });
        } else {
            this._value = textValue;
        }
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    get value(): any {
        try {
            switch (this.mode.toLowerCase()) {
                case "yaml":
                    return YAML.parse(this.getInnerValue());
                default:
                    return this.getInnerValue();
            }
        } catch (e) {
            console.warn(e);
            return this.getInnerValue();
        }
    }

    private getInnerValue(): string {
        if (!this.editor) {
            return "";
        }
        return this.editor.state.doc.toString();
    }

    getLanguageExtension(): LanguageSupport | undefined {
        switch (this.mode.toLowerCase()) {
            case "javascript":
                return javascript();
            case "yaml":
                return new LanguageSupport(StreamLanguage.define(yamlMode.yaml));
        }
        return undefined;
    }

    firstUpdated(): void {
        const matcher = window.matchMedia("(prefers-color-scheme: light)");
        const handler = (ev?: MediaQueryListEvent) => {
            let theme;
            if (ev?.matches || matcher.matches) {
                theme = this.themeLight;
            } else {
                theme = this.themeDark;
            }
            this.editor?.dispatch({
                effects: [this.theme.reconfigure(theme)],
            });
        };
        const extensions = [
            history(),
            keymap.of([...defaultKeymap, ...historyKeymap]),
            this.getLanguageExtension(),
            lineNumbers(),
            drawSelection(),
            EditorView.lineWrapping,
            EditorState.readOnly.of(this.readOnly),
            EditorState.tabSize.of(2),
            this.theme.of(this.themeLight),
        ];
        this.editor = new EditorView({
            extensions: extensions.filter((p) => p) as Extension[],
            root: this.shadowRoot || document,
            doc: this._value,
        });
        this.shadowRoot?.appendChild(this.editor.dom);
        matcher.addEventListener("change", handler);
        handler();
    }
}
