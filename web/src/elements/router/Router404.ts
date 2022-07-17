import { CSSResult, LitElement, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

@customElement("gravity-router-404")
export class Router404 extends LitElement {
    @property()
    url = "";

    static get styles(): CSSResult[] {
        return [];
    }

    render(): TemplateResult {
        return html`<div class="pf-c-empty-state pf-m-full-height">
            <div class="pf-c-empty-state__content">
                <i class="fas fa-question-circle pf-c-empty-state__icon" aria-hidden="true"></i>
                <h1 class="pf-c-title pf-m-lg">${`Not found`}</h1>
                <div class="pf-c-empty-state__body">${`The URL "${this.url}" was not found.`}</div>
                <a href="#/" class="pf-c-button pf-m-primary" type="button">${`Return home`}</a>
            </div>
        </div>`;
    }
}
