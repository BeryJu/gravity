import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFSpinner from "@patternfly/patternfly/components/Spinner/spinner.css";

import { AKElement } from "./Base";

export enum PFSize {
    Small = "pf-m-sm",
    Medium = "pf-m-md",
    Large = "pf-m-lg",
    XLarge = "pf-m-xl",
}

@customElement("ak-spinner")
export class Spinner extends AKElement {
    @property()
    size: PFSize = PFSize.Medium;

    static get styles(): CSSResult[] {
        return [PFSpinner];
    }

    render(): TemplateResult {
        return html`<svg
            class="pf-v6-c-spinner ${this.size.toString()}"
            role="progressbar"
            viewBox="0 0 100 100"
            aria-label="Loading..."
        >
            <circle class="pf-v6-c-spinner__path" cx="50" cy="50" r="45" fill="none" />
        </svg> `;
    }
}
