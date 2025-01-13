import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { CSSResult, TemplateResult, html } from "lit";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFForm from "@patternfly/patternfly/components/Form/form.css";
import PFRadio from "@patternfly/patternfly/components/Radio/radio.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../../elements/Base";
import { WizardPage } from "../../../elements/wizard/WizardPage";
import { TypeCreate } from "../../dns/wizard/ZonePresetWizardPage";

@customElement("gravity-dhcp-wizard-type")
export class ScopePresetWizardPage extends WizardPage {
    applicationTypes: TypeCreate[] = [
        {
            components: ["gravity-dhcp-wizard-internal", "gravity-dhcp-wizard-dns"],
            name: "Gravity-hosted",
            description: "Leases are stored in Gravity, optionally integrated with DNS.",
            callback: () => {},
        },
        {
            components: ["gravity-dhcp-wizard-import"],
            name: "Import",
            description: "Import DHCP leases from various different formats.",
            callback: () => {},
        },
    ];

    sidebarLabel = () => "Scope configuration preset";

    static get styles(): CSSResult[] {
        return [PFBase, PFButton, PFForm, PFRadio, AKElement.GlobalStyle];
    }

    activeCallback = async () => {
        this.host.isValid = true;
        this.applicationTypes[0].callback(this.host);
        this.host.steps = ["gravity-dhcp-wizard-initial", "gravity-dhcp-wizard-type"].concat(
            ...this.applicationTypes[0].components,
        );
    };

    render(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            ${this.applicationTypes.map((type, idx) => {
                return html`<div class="pf-c-radio">
                    <input
                        class="pf-c-radio__input"
                        type="radio"
                        name="type"
                        id=${type.name}
                        ?checked=${idx === 0}
                        @change=${() => {
                            type.callback(this.host);
                            this.host.steps = [
                                "gravity-dhcp-wizard-initial",
                                "gravity-dhcp-wizard-type",
                            ].concat(...type.components);
                            this.host.isValid = true;
                        }}
                    />
                    <label class="pf-c-radio__label" for=${type.name}>${type.name}</label>
                    <span class="pf-c-radio__description">${type.description}</span>
                </div>`;
            })}
        </form>`;
    }
}
