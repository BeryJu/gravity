import { AuthUserLoginInput, RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFBackgroundImage from "@patternfly/patternfly/components/BackgroundImage/background-image.css";
import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFDrawer from "@patternfly/patternfly/components/Drawer/drawer.css";
import PFList from "@patternfly/patternfly/components/List/list.css";
import PFLogin from "@patternfly/patternfly/components/Login/login.css";
import PFTitle from "@patternfly/patternfly/components/Title/title.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { DEFAULT_CONFIG } from "./api/Config";
import { AKElement } from "./elements/Base";
import { Form } from "./elements/forms/Form";
import "./elements/forms/HorizontalFormElement";

@customElement("gravity-login-form")
export class LoginForm extends Form<AuthUserLoginInput> {
    send = (data: AuthUserLoginInput): Promise<void> => {
        return new RolesApiApi(DEFAULT_CONFIG)
            .apiUsersLogin({
                authUserLoginInput: data,
            })
            .then((a) => {
                if (a.successful) {
                    window.location.href = "#/";
                    window.location.reload();
                }
            });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label=${`Username`} name="username">
                <input type="text" class="pf-c-form-control" />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${`Password`} name="password">
                <input type="password" class="pf-c-form-control" />
            </ak-form-element-horizontal>
            <button
                class="pf-c-button pf-m-primary pf-m-block"
                click=${(e: MouseEvent) => {
                    e.preventDefault();
                    this.submit(e);
                }}
            >
                Log in
            </button>
        </form>`;
    }
}

@customElement("gravity-login")
export class LoginPage extends AKElement {
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
                .pf-c-background-image::before {
                    --ak-flow-background: url("./static/assets/images/pfbg_1200.jpg");
                }
            `,
        ];
    }

    render(): TemplateResult {
        return html`
            <div class="pf-c-background-image">
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
            </div>
            <div class="pf-c-login">
                <div class="pf-c-login__container">
                    <header class="pf-c-login__header">Gravity</header>
                    <main class="pf-c-login__main">
                        <header class="pf-c-login__main-header">
                            <h1 class="pf-c-title pf-m-3xl">Log in to your account</h1>
                        </header>
                        <div class="pf-c-login__main-body">
                            <gravity-login-form></gravity-login-form>
                        </div>
                        <footer class="pf-c-login__main-footer">
                            <ul class="pf-c-login__main-footer-links"></ul>
                            <div class="pf-c-login__main-footer-band"></div>
                        </footer>
                    </main>
                    <footer class="pf-c-login__footer"></footer>
                </div>
            </div>
        `;
    }
}
