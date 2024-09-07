import { LitElement, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";



import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFContent from "@patternfly/patternfly-v6/components/Content/content.css";
import PFMasthead from "@patternfly/patternfly-v6/components/Masthead/masthead.css";
import PFMenuToggle from "@patternfly/patternfly-v6/components/MenuToggle/menu-toggle.css";
import PFToggleGroup from "@patternfly/patternfly-v6/components/ToggleGroup/toggle-group.css";
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

enum Theme {
    Automatic,
    Light,
    Dark
}

@customElement("ak-header")
export class Header extends AKElement {
    @state()
    _title: Title = {
        title: "Loading",
    };

    @state()
    _theme: Theme = Theme.Automatic;

    get theme(): Theme {
        return this._theme;
    }
    set theme(v: Theme) {
        this._theme = v;
        if (v === Theme.Dark) {
            document.querySelector("html")?.classList.add("pf-v6-theme-dark");
        } else {
            document.querySelector("html")?.classList.remove("pf-v6-theme-dark");
        }
    }

    static get styles() {
        return [
            PFBase,
            PFButton,
            PFContent,
            PFMenuToggle,
            PFToolbar,
            PFToggleGroup,
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
        if (this.theme !== Theme.Automatic) {
            return;
        }
        const matcher = window.matchMedia("(prefers-color-scheme: light)");
        const handler = (ev?: MediaQueryListEvent) => {
            this.theme = ev?.matches ? Theme.Light : Theme.Dark;
        };
        handler();
        matcher.addEventListener("change", handler);
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
                                <div
                                    class="pf-v6-c-toggle-group"
                                    role="group"
                                    aria-label="Dark theme toggle group"
                                >
                                    <div class="pf-v6-c-toggle-group__item">
                                        <button
                                            type="button"
                                            class="pf-v6-c-toggle-group__button ${this.theme === Theme.Light ? "pf-m-selected" : ""}"
                                            aria-label="light theme toggle"
                                            @click=${() => {
                                                this.theme = Theme.Light;
                                            }}
                                        >
                                            <span class="pf-v6-c-toggle-group__icon"
                                                ><svg
                                                    class="pf-v6-svg"
                                                    viewBox="0 0 512 512"
                                                    fill="currentColor"
                                                    aria-hidden="true"
                                                    role="img"
                                                    width="1em"
                                                    height="1em"
                                                >
                                                    <path
                                                        d="M256 160c-52.9 0-96 43.1-96 96s43.1 96 96 96 96-43.1 96-96-43.1-96-96-96zm246.4 80.5l-94.7-47.3 33.5-100.4c4.5-13.6-8.4-26.5-21.9-21.9l-100.4 33.5-47.4-94.8c-6.4-12.8-24.6-12.8-31 0l-47.3 94.7L92.7 70.8c-13.6-4.5-26.5 8.4-21.9 21.9l33.5 100.4-94.7 47.4c-12.8 6.4-12.8 24.6 0 31l94.7 47.3-33.5 100.5c-4.5 13.6 8.4 26.5 21.9 21.9l100.4-33.5 47.3 94.7c6.4 12.8 24.6 12.8 31 0l47.3-94.7 100.4 33.5c13.6 4.5 26.5-8.4 21.9-21.9l-33.5-100.4 94.7-47.3c13-6.5 13-24.7.2-31.1zm-155.9 106c-49.9 49.9-131.1 49.9-181 0-49.9-49.9-49.9-131.1 0-181 49.9-49.9 131.1-49.9 181 0 49.9 49.9 49.9 131.1 0 181z"
                                                    ></path></svg
                                            ></span>
                                        </button>
                                    </div>
                                    <div class="pf-v6-c-toggle-group__item">
                                        <button
                                            type="button"
                                            class="pf-v6-c-toggle-group__button ${this.theme === Theme.Dark ? "pf-m-selected" : ""}"
                                            aria-label="dark theme toggle"
                                            @click=${() => {
                                                this.theme = Theme.Dark;
                                            }}
                                        >
                                            <span class="pf-v6-c-toggle-group__icon"
                                                ><svg
                                                    class="pf-v6-svg"
                                                    viewBox="0 0 512 512"
                                                    fill="currentColor"
                                                    aria-hidden="true"
                                                    role="img"
                                                    width="1em"
                                                    height="1em"
                                                >
                                                    <path
                                                        d="M283.211 512c78.962 0 151.079-35.925 198.857-94.792 7.068-8.708-.639-21.43-11.562-19.35-124.203 23.654-238.262-71.576-238.262-196.954 0-72.222 38.662-138.635 101.498-174.394 9.686-5.512 7.25-20.197-3.756-22.23A258.156 258.156 0 0 0 283.211 0c-141.309 0-256 114.511-256 256 0 141.309 114.511 256 256 256z"
                                                    ></path></svg
                                            ></span>
                                        </button>
                                    </div>
                                </div>
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
