import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFBackdrop from "@patternfly/patternfly-v6/components/Backdrop/backdrop.css";
import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFCard from "@patternfly/patternfly-v6/components/Card/card.css";
import PFContent from "@patternfly/patternfly-v6/components/Content/content.css";
import PFForm from "@patternfly/patternfly-v6/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly-v6/components/FormControl/form-control.css";
import PFModalBox from "@patternfly/patternfly-v6/components/ModalBox/modal-box.css";
import PFPage from "@patternfly/patternfly-v6/components/Page/page.css";
import PFTitle from "@patternfly/patternfly-v6/components/Title/title.css";
import PFBullseye from "@patternfly/patternfly-v6/layouts/Bullseye/bullseye.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";
import { PFSize } from "../Spinner";

export const MODAL_BUTTON_STYLES = css`
    :host {
        text-align: left;
        font-size: var(--pf-global--FontSize--md);
    }
    .pf-v6-c-modal-box.pf-m-lg {
        overflow-y: auto;
    }
    .pf-v6-c-modal-box > .pf-v6-c-button + * {
        margin-right: 0;
    }
    /* fix multiple selects height */
    select[multiple] {
        height: 15em;
    }
`;

@customElement("ak-modal-button")
export class ModalButton extends AKElement {
    @property()
    size: PFSize = PFSize.Large;

    @property({ type: Boolean })
    open = false;

    @property({ type: Boolean })
    locked = false;

    handlerBound = false;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFButton,
            PFModalBox,
            PFForm,
            PFTitle,
            PFFormControl,
            PFBullseye,
            PFBackdrop,
            PFPage,
            PFCard,
            PFContent,
            AKElement.GlobalStyle,
            MODAL_BUTTON_STYLES,
            css`
                .locked {
                    overflow-y: hidden !important;
                }
            `,
        ];
    }

    firstUpdated(): void {
        if (this.handlerBound) return;
        window.addEventListener("keyup", this.keyUpHandler);
        this.handlerBound = true;
    }

    keyUpHandler = (e: KeyboardEvent): void => {
        if (e.code === "Escape") {
            this.resetForms();
            this.open = false;
        }
    };

    disconnectedCallback(): void {
        super.disconnectedCallback();
        window.removeEventListener("keyup", this.keyUpHandler);
    }

    resetForms(): void {
        this.querySelectorAll<HTMLFormElement>("[slot=form]").forEach((form) => {
            if ("resetForm" in form) {
                form?.resetForm();
            }
        });
    }

    onClick(): void {
        this.open = true;
        this.querySelectorAll("*").forEach((child) => {
            if ("requestUpdate" in child) {
                (child as AKElement).requestUpdate();
            }
        });
    }

    renderModalInner(): TemplateResult {
        return html`<slot name="modal"></slot>`;
    }

    renderModal(): TemplateResult {
        return html`<div class="pf-v6-c-backdrop">
            <div class="pf-v6-l-bullseye">
                <div
                    class="pf-v6-c-modal-box ${this.size} ${this.locked ? "locked" : ""}"
                    role="dialog"
                    aria-modal="true"
                >
                    <div class="pf-v6-c-modal-box__close">
                        <button
                            class="pf-v6-c-button pf-m-plain"
                            type="button"
                            aria-label="Close"
                            @click=${() => {
                                this.resetForms();
                                this.open = false;
                            }}
                        >
                            <span class="pf-v6-c-button__icon">
                                <i class="fas fa-times" aria-hidden="true"></i>
                            </span>
                        </button>
                    </div>
                    ${this.renderModalInner()}
                </div>
            </div>
        </div>`;
    }

    render(): TemplateResult {
        return html` <slot name="trigger" @click=${() => this.onClick()}></slot>
            ${this.open ? this.renderModal() : ""}`;
    }
}
