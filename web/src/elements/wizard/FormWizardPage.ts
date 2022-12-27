import { customElement } from "lit/decorators.js";

import { Form } from "../forms/Form";
import { WizardPage } from "./WizardPage";

/**
 * This Wizard page is used for proxy forms with the older-style
 * wizards
 */
@customElement("ak-wizard-page-form")
export class FormWizardPage extends WizardPage {
    activePageCallback: (context: FormWizardPage) => Promise<void> = async () => {
        return Promise.resolve();
    };

    activeCallback = async () => {
        this.host.isValid = true;
        this.activePageCallback(this);
    };

    nextCallback = async () => {
        const form = this.querySelector<Form<unknown>>("*");
        if (!form) {
            return Promise.reject("No form found");
        }
        const formPromise = form.submit(new Event("submit"));
        if (!formPromise) {
            return Promise.reject("Form didn't return a promise for submitting");
        }
        return formPromise
            .then((data) => {
                this.host.state[this.slot] = data;
                this.host.canBack = false;
                return true;
            })
            .catch(() => {
                return false;
            });
    };
}
