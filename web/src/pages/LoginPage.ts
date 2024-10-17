import { AuthAPIConfigOutput, AuthAPILoginInput, RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFBackgroundImage from "@patternfly/patternfly/components/BackgroundImage/background-image.css";
import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFDrawer from "@patternfly/patternfly/components/Drawer/drawer.css";
import PFList from "@patternfly/patternfly/components/List/list.css";
import PFLogin from "@patternfly/patternfly/components/Login/login.css";
import PFTitle from "@patternfly/patternfly/components/Title/title.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { DEFAULT_CONFIG } from "../api/Config";
import { AKElement } from "../elements/Base";
import { Form } from "../elements/forms/Form";
import "../elements/forms/HorizontalFormElement";

@customElement("gravity-login-form")
export class LoginForm extends Form<AuthAPILoginInput> {
    send = async (data: AuthAPILoginInput): Promise<void> => {
        const a = await new RolesApiApi(DEFAULT_CONFIG).apiLoginUser({
            authAPILoginInput: data,
        });
        if (a.successful) {
            window.location.hash = "#/";
            window.location.reload();
        }
    };

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal label=${"Username"} name="username">
                <input type="text" class="pf-c-form-control" autocomplete="username" />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"Password"} name="password">
                <input type="password" class="pf-c-form-control" autocomplete="current-password" />
            </ak-form-element-horizontal>
            <button
                class="pf-c-button pf-m-primary pf-m-block"
                @click=${(e: MouseEvent) => {
                    e.preventDefault();
                    this.submit(e);
                }}
            >
                Log in
            </button>`;
    }
}

@customElement("gravity-login")
export class LoginPage extends AKElement {
    @state()
    authConfig: AuthAPIConfigOutput | undefined;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFLogin,
            PFDrawer,
            PFButton,
            PFTitle,
            PFList,
            PFBackgroundImage,
            AKElement.GlobalStyle,
            css`
                :host {
                    background-image: url("/ui/static/assets/images/background.jpg");
                    background-position: center center;
                    background-size: cover;
                }
                .pf-c-login__header {
                    font-size: 3rem;
                    color: var(--ak-accent);
                }
            `,
        ];
    }

    firstUpdated(): void {
        new RolesApiApi(DEFAULT_CONFIG).apiAuthConfig().then((config) => {
            this.authConfig = config;
        });
    }

    render(): TemplateResult {
        return html`
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="pf-c-background-image__filter"
                width="0"
                height="0"
            >
                <filter id="image_overlay">
                    <feColorMatrix
                        in="SourceGraphic"
                        type="matrix"
                        values="1.3 0 0 0 0 0 1.3 0 0 0 0 0 1.3 0 0 0 0 0 1 0"
                    />
                    <feComponentTransfer color-interpolation-filters="sRGB" result="duotone">
                        <feFuncR
                            type="table"
                            tableValues="0.086274509803922 0.43921568627451"
                        ></feFuncR>
                        <feFuncG
                            type="table"
                            tableValues="0.086274509803922 0.43921568627451"
                        ></feFuncG>
                        <feFuncB
                            type="table"
                            tableValues="0.086274509803922 0.43921568627451"
                        ></feFuncB>
                        <feFuncA type="table" tableValues="0 1"></feFuncA>
                    </feComponentTransfer>
                </filter>
            </svg>
            <div class="pf-c-login">
                <div class="pf-c-login__container">
                    <header class="pf-c-login__header">
                        <img src="static/assets/images/logo-color.png" alt="gravity logo" />
                    </header>
                    <main class="pf-c-login__main">
                        <header class="pf-c-login__main-header">
                            <h1 class="pf-c-title pf-m-3xl">Log in to Gravity.</h1>
                        </header>
                        <div class="pf-c-login__main-body">
                            <gravity-login-form></gravity-login-form>
                        </div>
                        ${this.authConfig?.oidc
                            ? html`
                                  <div class="pf-c-login__main-body">
                                      <a
                                          class="pf-c-button pf-m-primary pf-m-block"
                                          href="/auth/oidc"
                                          >Login with SSO</a
                                      >
                                  </div>
                              `
                            : html``}
                    </main>
                    <footer class="pf-c-login__footer">
                        <p></p>
                        <ul class="pf-c-list pf-m-inline">
                            <li>
                                <a
                                    href="https://gravity.beryju.io?utm_source=gravity&amp;utm_medium=login"
                                    >${"Powered by Gravity"}</a
                                >
                            </li>
                        </ul>
                    </footer>
                </div>
            </div>
        `;
    }
}
