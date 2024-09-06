import { CSSResult, TemplateResult, css, html, nothing } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFSpinner from "@patternfly/patternfly-v6/components/Spinner/spinner.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { ERROR_CLASS, PROGRESS_CLASS, SUCCESS_CLASS } from "../../common/constants";
import { AKElement } from "../Base";
import { PFSize } from "../Spinner";

@customElement("ak-spinner-button")
export class SpinnerButton extends AKElement {
    @property({ type: Boolean })
    isRunning = false;

    @property()
    callAction: (() => Promise<unknown>) | undefined;

    @property({ type: Boolean })
    disabled = false;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFButton,
            PFSpinner,
            AKElement.GlobalStyle,
            css`
                button {
                    height: 100%;
                    /* Have to use !important here, as buttons with pf-m-progress have transition already */
                    transition: all var(--pf-c-button--m-progress--TransitionDuration) ease 0s !important;
                }
            `,
        ];
    }

    constructor() {
        super();
    }

    setLoading(): void {
        this.isRunning = true;
        this.classList.add(PROGRESS_CLASS);
        this.requestUpdate();
    }

    setDone(statusClass: string): void {
        this.isRunning = false;
        setTimeout(() => {
            this.classList.remove(PROGRESS_CLASS);
        }, 2500);
        this.classList.add(statusClass);
        this.requestUpdate();
        setTimeout(() => {
            this.classList.remove(statusClass);
            this.requestUpdate();
        }, 1000);
    }

    render(): TemplateResult {
        return html`<button
            class="pf-v6-c-button pf-m-progress ${this.classList.toString()}"
            ?disabled=${this.disabled}
            @click=${() => {
                if (this.isRunning === true) {
                    return;
                }
                this.setLoading();
                if (this.callAction) {
                    this.callAction()
                        .then(() => {
                            this.setDone(SUCCESS_CLASS);
                        })
                        .catch(() => {
                            this.setDone(ERROR_CLASS);
                        });
                }
            }}
        >
            ${this.isRunning
                ? html`<span class="pf-v6-c-button__progress">
                      <ak-spinner size=${PFSize.Medium}></ak-spinner>
                  </span>`
                : nothing}
            <slot class="pf-v6-c-button__text"></slot>
        </button>`;
    }
}
