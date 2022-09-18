import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFExpandableSection from "@patternfly/patternfly/components/ExpandableSection/expandable-section.css";

import { AKElement } from "./Base";

@customElement("ak-expand")
export class Expand extends AKElement {
    @property({ type: Boolean })
    expanded = false;

    @property()
    textOpen = `Show less`;

    @property()
    textClosed = `Show more`;

    static get styles(): CSSResult[] {
        return [PFExpandableSection];
    }

    render(): TemplateResult {
        return html`<div class="pf-c-expandable-section ${this.expanded ? "pf-m-expanded" : ""}">
            <button
                type="button"
                class="pf-c-expandable-section__toggle"
                aria-expanded="${this.expanded}"
                @click=${() => {
                    this.expanded = !this.expanded;
                }}
            >
                <span class="pf-c-expandable-section__toggle-icon">
                    <i class="fas fa-angle-right" aria-hidden="true"></i>
                </span>
                <span class="pf-c-expandable-section__toggle-text"
                    >${this.expanded ? `${this.textOpen}` : `${this.textClosed}`}</span
                >
            </button>
            <slot ?hidden=${!this.expanded} class="pf-c-expandable-section__content"></slot>
        </div>`;
    }
}
