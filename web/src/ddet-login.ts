import { css, html, LitElement, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";
import "@spectrum-web-components/theme/theme-light.js";
import "@spectrum-web-components/theme/scale-medium.js";
import "@spectrum-web-components/theme/sp-theme.js";
import "@spectrum-web-components/button/sp-button.js";

import "@spectrum-web-components/textfield/sp-textfield.js";
import "@spectrum-web-components/field-label/sp-field-label.js";
import { isLoggedIn, login } from "src/services/api";

@customElement("gravity-login")
export class Login extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
            }
            form {
                width: 300px;
                margin: 0 auto;
                padding: 10rem 0;
            }
            .group {
                margin-bottom: 10px;
            }
        `;
    }

    login(): void {
        const username = this.shadowRoot?.querySelector<HTMLInputElement>(
            "sp-textfield[id=username]",
        );
        const password = this.shadowRoot?.querySelector<HTMLInputElement>(
            "sp-textfield[id=password]",
        );
        if (!username || !password) {
            return;
        }
        login(username?.value, password?.value);
    }

    render(): TemplateResult {
        return html`
            <sp-theme color="light" scale="medium">
                <form>
                    <div class="group">
                        <sp-field-label for="username">Name</sp-field-label>
                        <sp-textfield id="username" placeholder="Username"></sp-textfield>
                    </div>
                    <div class="group">
                        <sp-field-label for="password">Passowrd</sp-field-label>
                        <sp-textfield
                            id="password"
                            type="password"
                            placeholder="Enter your name"
                        ></sp-textfield>
                    </div>
                    <div class="group">
                        <sp-button
                            size="l"
                            @click=${() => {
                                this.login();
                            }}
                            >Login</sp-button
                        >
                    </div>
                </form>
            </sp-theme>
        `;
    }
}
