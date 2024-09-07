import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFForm from "@patternfly/patternfly-v6/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly-v6/components/FormControl/form-control.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";

@customElement("ak-form-group")
export class FormGroup extends AKElement {
    @property({ type: Boolean })
    expanded = false;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFForm,
            PFButton,
            PFFormControl,
            AKElement.GlobalStyle,
            css`
                slot[name="body"][hidden] {
                    display: none !important;
                }
            `,
        ];
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-form__field-group ${this.expanded ? "pf-m-expanded" : ""}">
            <div class="pf-v6-c-form__field-group-toggle">
                <div class="pf-v6-c-form__field-group-toggle-button">
                    <button
                        class="pf-v6-c-button pf-m-plain"
                        type="button"
                        aria-expanded="${this.expanded}"
                        aria-label="Details"
                        @click=${() => {
                            this.expanded = !this.expanded;
                        }}
                    >
                        <span class="pf-v6-c-button__icon">
                            <span class="pf-v6-c-form__field-group-toggle-icon">
                                <i class="fas fa-angle-right" aria-hidden="true"></i>
                            </span>
                        </span>
                    </button>
                </div>
            </div>
            <div class="pf-v6-c-form__field-group-header">
                <div class="pf-v6-c-form__field-group-header-main">
                    <div class="pf-v6-c-form__field-group-header-title">
                        <div class="pf-v6-c-form__field-group-header-title-text">
                            <slot name="header"></slot>
                        </div>
                    </div>
                    <div class="pf-v6-c-form__field-group-header-description">
                        <slot name="description"></slot>
                    </div>
                </div>
            </div>
            <slot
                ?hidden=${!this.expanded}
                class="pf-v6-c-form__field-group-body"
                name="body"
            ></slot>
        </div>`;
    }
}
