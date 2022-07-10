import { html, LitElement, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";

@customElement("ddet-router")
export class Router extends LitElement {

    constructor() {
        super();
        window.addEventListener("hashchange", () => {
            // TODO: routing
        });
    }

    render(): TemplateResult {
        return html`

        `;
    }

}
