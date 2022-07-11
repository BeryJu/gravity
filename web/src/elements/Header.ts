import SPTypographyIndexVars from "@spectrum-css/typography/dist/index-vars.css";
import SPTypographyVars from "@spectrum-css/typography/dist/vars.css";

import { CSSResult, LitElement, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

@customElement("gravity-header")
export class Table extends LitElement {
    static get styles(): CSSResult[] {
        return [SPTypographyIndexVars, SPTypographyVars];
    }

    render(): TemplateResult {
        return html`
            <h1 class="spectrum-Heading spectrum-Heading--sizeXXL">
                <slot></slot>
            </h1>
        `;
    }
}
