import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFNav from "@patternfly/patternfly-v6/components/Nav/nav.css";
import PFPage from "@patternfly/patternfly-v6/components/Page/page.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";

@customElement("ak-sidebar")
export class Sidebar extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFPage,
            PFNav,
            AKElement.GlobalStyle,
            css`
                :host {
                    z-index: 100;
                }
                .pf-v6-c-nav__link.pf-m-current::after,
                .pf-v6-c-nav__link.pf-m-current:hover::after,
                .pf-v6-c-nav__item.pf-m-current:not(.pf-m-expanded) .pf-v6-c-nav__link::after {
                    --pf-v6-c-nav__link--m-current--after--BorderColor: var(--ak-accent);
                }

                .pf-v6-c-nav__section + .pf-v6-c-nav__section {
                    --pf-v6-c-nav__section--section--MarginTop: var(--pf-global--spacer--sm);
                }
                .pf-v6-c-nav__list .sidebar-brand {
                    max-height: 82px;
                    margin-bottom: -0.5rem;
                }
            `,
        ];
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-page__sidebar-body">
            <nav class="pf-v6-c-nav" aria-label="Global">
                <ul class="pf-v6-c-nav__list">
                    <slot></slot>
                </ul>
            </nav>
        </div>`;
    }
}
