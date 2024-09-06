import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFEmptyState from "@patternfly/patternfly-v6/components/EmptyState/empty-state.css";
import PFTitle from "@patternfly/patternfly-v6/components/Title/title.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "./Base";
import { PFSize } from "./Spinner";
import "./Spinner";

@customElement("ak-empty-state")
export class EmptyState extends AKElement {
    @property({ type: String })
    icon = "";

    @property({ type: Boolean })
    loading = false;

    @property({ type: Boolean })
    fullHeight = false;

    @property()
    header = "";

    static get styles(): CSSResult[] {
        return [PFBase, PFEmptyState, PFTitle, AKElement.GlobalStyle];
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-empty-state ${this.fullHeight && "pf-m-full-height"}">
            <div class="pf-v6-c-empty-state__content">
                ${this.loading
                    ? html`<div class="pf-v6-c-empty-state__icon">
                          <ak-spinner size=${PFSize.XLarge}></ak-spinner>
                      </div>`
                    : html`<i
                          class="pf-icon fa ${this.icon ||
                          "fa-question-circle"} pf-v6-c-empty-state__icon"
                          aria-hidden="true"
                      ></i>`}
                <h1 class="pf-v6-c-title pf-m-lg">${this.header}</h1>
                <div class="pf-v6-c-empty-state__body">
                    <slot name="body"></slot>
                </div>
                <div class="pf-v6-c-empty-state__primary">
                    <slot name="primary"></slot>
                </div>
            </div>
        </div>`;
    }
}
