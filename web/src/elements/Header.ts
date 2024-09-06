import { LitElement, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFContent from "@patternfly/patternfly-v6/components/Content/content.css";
import PFMasthead from "@patternfly/patternfly-v6/components/Masthead/masthead.css";
import PFMenuToggle from "@patternfly/patternfly-v6/components/MenuToggle/menu-toggle.css";
import PFToolbar from "@patternfly/patternfly-v6/components/Toolbar/toolbar.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { EVENT_SIDEBAR_TOGGLE, EVENT_TMP_TITLE } from "../common/constants";
import { AKElement } from "../elements/Base";

export interface Title {
    title: string;
    subtext?: string;
    icon?: string;
}
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type AbstractConstructor<T = object> = abstract new (...args: any[]) => T;

export function WithHeaderTitle<T extends AbstractConstructor<LitElement>>(superclass: T) {
    abstract class WithHeader extends superclass {
        abstract get title(): string;
    }
    return WithHeader;
}

@customElement("ak-header")
export class Header extends AKElement {
    @state()
    _title: Title = {
        title: "Loading",
    };

    static get styles() {
        return [
            PFBase,
            PFButton,
            PFContent,
            PFMenuToggle,
            PFToolbar,
            PFMasthead,
            css`
                .pf-v6-c-toolbar__group.title {
                    height: 100%;
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    padding-left: 24px;
                }
            `,
        ];
    }

    firstUpdated() {
        document.addEventListener(EVENT_TMP_TITLE, ((ev: CustomEvent<Title>) => {
            this._title = ev.detail;
        }) as EventListener);
    }

    render() {
        return html`<header class="pf-v6-c-masthead">
            <div class="pf-v6-c-masthead__main">
                <span class="pf-v6-c-masthead__toggle">
                    <button
                        class="pf-v6-c-button pf-m-plain"
                        type="button"
                        @click=${() => {
                            this.dispatchEvent(
                                new CustomEvent(EVENT_SIDEBAR_TOGGLE, {
                                    bubbles: true,
                                    composed: true,
                                }),
                            );
                        }}
                    >
                        <span class="pf-v6-c-button__icon">
                            <i class="fas fa-bars" aria-hidden="true"></i>
                        </span>
                    </button>
                </span>
                <div class="pf-v6-c-masthead__brand">
                    <a class="pf-v6-c-masthead__logo" href="#/">
                        <img src="static/assets/images/logo-color.png" alt="gravity logo" />
                    </a>
                </div>
            </div>
            <div class="pf-v6-c-masthead__content">
                <div class="pf-v6-c-toolbar pf-m-static">
                    <div class="pf-v6-c-toolbar__content">
                        <div class="pf-v6-c-toolbar__content-section">
                            <div class="pf-v6-c-toolbar__group pf-m-filter-group title">
                                <div class="pf-v6-c-toolbar__item">
                                    <div class="pf-v6-c-content">
                                        <h1>${this._title.title}</h1>
                                    </div>
                                </div>
                            </div>
                            <div class="pf-v6-c-toolbar__item pf-m-align-end">
                                <a class="pf-v6-c-menu-toggle pf-m-plain" href="/auth/logout">
                                    <span class="pf-v6-c-menu-toggle__icon">
                                        <i class="fas fa-sign-out-alt" aria-hidden="true"></i>
                                    </span>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </header>`;
    }
}
