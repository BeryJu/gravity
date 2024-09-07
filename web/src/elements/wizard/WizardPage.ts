import { CSSResult, PropertyDeclaration, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";
import { Wizard } from "./Wizard";

@customElement("ak-wizard-page")
export class WizardPage extends AKElement {
    static get styles(): CSSResult[] {
        return [PFBase, AKElement.GlobalStyle];
    }

    @property()
    sidebarLabel: () => string = () => {
        return "UNNAMED";
    };

    get host(): Wizard {
        return this.parentElement as Wizard;
    }

    activeCallback: () => Promise<void> = async () => {
        this.host.isValid = false;
    };

    nextCallback: () => Promise<boolean> = async () => {
        return true;
    };

    requestUpdate(
        name?: PropertyKey,
        oldValue?: unknown,
        options?: PropertyDeclaration<unknown, unknown>,
    ): void {
        this.querySelectorAll("*").forEach((el) => {
            if ("requestUpdate" in el) {
                (el as AKElement).requestUpdate();
            }
        });
        return super.requestUpdate(name, oldValue, options);
    }

    render(): TemplateResult {
        return html`<slot></slot>`;
    }
}
