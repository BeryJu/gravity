import { HorizontalFormElement } from "src/elements/forms/HorizontalFormElement";

import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property, query } from "lit/decorators.js";

import PFAlert from "@patternfly/patternfly/components/Alert/alert.css";
import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFCard from "@patternfly/patternfly/components/Card/card.css";
import PFForm from "@patternfly/patternfly/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly/components/FormControl/form-control.css";
import PFInputGroup from "@patternfly/patternfly/components/InputGroup/input-group.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../Base";
import { Form, KeyUnknown, formFiles, serializeForm } from "../forms/Form";
import { WizardPage } from "./WizardPage";

@customElement("ak-wizard-form")
export class WizardForm extends Form<KeyUnknown> {
    viewportCheck = false;

    @property({ attribute: false })
    nextDataCallback!: (data: KeyUnknown) => Promise<boolean>;

    submit(): Promise<boolean> | undefined {
        const data = this.serializeForm();
        if (!data) {
            return;
        }
        const finalData = Object.assign({}, data);
        return this.nextDataCallback(finalData);
    }

    getFormFiles(): { [key: string]: File } {
        const elements = this.querySelectorAll<HorizontalFormElement>("ak-form-element-horizontal");
        if (!elements) {
            return {};
        }
        return formFiles(elements);
    }

    serializeForm(): KeyUnknown | undefined {
        const elements = this.querySelectorAll<HorizontalFormElement>("ak-form-element-horizontal");
        if (!elements) {
            return {} as KeyUnknown;
        }
        return serializeForm(elements) as KeyUnknown;
    }
}

export class WizardFormPage extends WizardPage {
    @query("ak-wizard-form")
    form?: WizardForm;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFCard,
            PFButton,
            PFForm,
            PFAlert,
            PFInputGroup,
            PFFormControl,
            AKElement.GlobalStyle,
        ];
    }

    inputCallback(): void {
        const form = this.shadowRoot?.querySelector<HTMLFormElement>("form");
        if (!form) {
            return;
        }
        const state = form.checkValidity();
        this.host.isValid = state;
    }

    nextCallback = async (): Promise<boolean> => {
        const form = this.shadowRoot?.querySelector<WizardForm>("ak-wizard-form");
        if (!form) {
            console.warn("authentik/wizard: could not find form element");
            return false;
        }
        const response = await form.submit();
        if (response === undefined) {
            return false;
        }
        return response;
    };

    nextDataCallback: (data: KeyUnknown) => Promise<boolean> = async (): Promise<boolean> => {
        return false;
    };

    renderForm(): TemplateResult {
        return html``;
    }

    activeCallback = async () => {
        this.inputCallback();
    };

    render(): TemplateResult {
        return html`
            <ak-wizard-form
                .nextDataCallback=${this.nextDataCallback}
                @input=${() => this.inputCallback()}
            >
                <form
                    class="pf-c-form pf-m-horizontal"
                    @submit=${(ev: Event) => {
                        ev.preventDefault();
                    }}
                >
                    ${this.renderForm()}
                </form>
            </ak-wizard-form>
        `;
    }
}
