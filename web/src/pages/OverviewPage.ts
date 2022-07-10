import { html, LitElement, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";

@customElement("gravity-overview")
export class OverviewPage extends LitElement {
    render(): TemplateResult {
        return html` Hello overview `;
    }
}
