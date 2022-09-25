import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFTitle from "@patternfly/patternfly/components/Title/title.css";
import PFGlobal from "@patternfly/patternfly/patternfly-base.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { EVENT_SIDEBAR_TOGGLE } from "../../common/constants";
import { AKElement } from "../Base";

// If the viewport is wider than MIN_WIDTH, the sidebar
// is shown besides the content, and not overlaid.
export const MIN_WIDTH = 1200;

@customElement("ak-sidebar-brand")
export class SidebarBrand extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFGlobal,
            PFPage,
            PFButton,
            PFTitle,
            AKElement.GlobalStyle,
            css`
                :host {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    height: 114px;
                    min-height: 114px;
                }
                .pf-c-brand img {
                    width: 100%;
                    padding: 0 0.5rem;
                    height: 42px;
                }
                button.pf-c-button.sidebar-trigger {
                    background-color: transparent;
                    border-radius: 0px;
                    height: 100%;
                    color: var(--ak-dark-foreground);
                }
                .ak-brand {
                    width: 100%;
                    font-size: 3rem;
                    color: var(--ak-accent);
                    text-align: center;
                }
                .pf-c-page__header-brand-link:hover {
                    text-decoration-color: var(--ak-accent);
                }
            `,
        ];
    }

    constructor() {
        super();
        window.addEventListener("resize", () => {
            this.requestUpdate();
        });
    }

    render(): TemplateResult {
        return html` ${window.innerWidth <= MIN_WIDTH
                ? html`
                      <button
                          class="sidebar-trigger pf-c-button"
                          @click=${() => {
                              this.dispatchEvent(
                                  new CustomEvent(EVENT_SIDEBAR_TOGGLE, {
                                      bubbles: true,
                                      composed: true,
                                  }),
                              );
                          }}
                      >
                          <i class="fas fa-bars"></i>
                      </button>
                  `
                : html``}
            <a href="#/" class="pf-c-page__header-brand-link">
                <div class="pf-c-brand ak-brand">Gravity.</div>
            </a>`;
    }
}
