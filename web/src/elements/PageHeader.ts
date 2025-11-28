import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFContent from "@patternfly/patternfly/components/Content/content.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { EVENT_SIDEBAR_TOGGLE, TITLE_DEFAULT } from "../common/constants";
import { AKElement } from "./Base";
import "./ak-nav-buttons";

@customElement("ak-page-header")
export class PageHeader extends AKElement {
    @property()
    icon: string | undefined;

    @property({ type: Boolean })
    iconImage = false;

    @property({ type: Boolean })
    hasNotifications = false;

    @property()
    description: string | undefined;

    @property()
    set header(value: string) {
        let title = TITLE_DEFAULT;
        title = `Admin - ${title}`;
        if (value !== "") {
            title = `${value} - ${title}`;
        }
        document.title = title;
        this._header = value;
    }

    get header(): string {
        return this._header;
    }

    _header = "";

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFButton,
            PFPage,
            PFContent,
            AKElement.GlobalStyle,
            css`
                :host {
                    position: sticky;
                    top: 0;
                    z-index: 100;
                }
                .bar {
                    border-bottom: var(--pf-global--BorderWidth--sm);
                    border-bottom-style: solid;
                    border-bottom-color: var(--pf-global--BorderColor--100);
                    display: flex;
                    flex-direction: row;
                    min-height: var(--navbar-height);
                    background-color: var(--pf-c-page--BackgroundColor);
                }
                .pf-c-page__main-section {
                    flex-grow: 1;
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                }
                img.pf-icon {
                    max-height: 24px;
                }
                .pf-c-page__header-tools {
                    flex-shrink: 0;
                }
                .pf-c-page__header-tools-group {
                    height: 100%;
                }
            `,
        ];
    }

    renderIcon(): TemplateResult {
        if (this.icon) {
            if (this.iconImage && !this.icon.startsWith("fa://")) {
                return html`<img class="pf-icon" src=${this.icon} />&nbsp;`;
            }
            const icon = this.icon.replaceAll("fa://", "fa ");
            return html`<i class=${icon}></i>&nbsp;`;
        }
        return html``;
    }

    render(): TemplateResult {
        return html`<div class="bar">
            <button
                class="sidebar-trigger pf-c-button pf-m-plain"
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
            <section class="pf-c-page__main-section">
                <div class="pf-c-content">
                    <h1>
                        ${this.renderIcon()}
                        <slot name="header"> ${this.header} </slot>
                    </h1>
                    ${this.description ? html`<p>${this.description}</p>` : html``}
                </div>
            </section>
            <div class="pf-c-page__header-tools">
                <div class="pf-c-page__header-tools-group">
                    <ak-nav-buttons></ak-nav-buttons>
                </div>
            </div>
        </div>`;
    }
}
