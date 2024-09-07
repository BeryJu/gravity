import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFEmptyState from "@patternfly/patternfly-v6/components/EmptyState/empty-state.css";
import PFTitle from "@patternfly/patternfly-v6/components/Title/title.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";

@customElement("ak-router-404")
export class Router404 extends AKElement {
    @property()
    url = "";

    static get styles(): CSSResult[] {
        return [PFBase, PFEmptyState, PFTitle];
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-empty-state pf-m-full-height">
            <div class="pf-v6-c-empty-state__content">
                <i class="fas fa-question-circle pf-v6-c-empty-state__icon" aria-hidden="true"></i>
                <h1 class="pf-v6-c-title pf-m-lg">${"Not found"}</h1>
                <div class="pf-v6-c-empty-state__body">
                    ${`The URL "${this.url}" was not found.`}
                </div>
                <a href="#/" class="pf-v6-c-button pf-m-primary" type="button">${"Return home"}</a>
            </div>
        </div>`;
    }
}
