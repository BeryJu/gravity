import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { CSSResult, TemplateResult, html } from "lit";
import { property } from "lit/decorators.js";

import AKGlobal from "@goauthentik/common/styles/authentik.css";
import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFRadio from "@patternfly/patternfly/components/Radio/radio.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../elements/Base";
import "../../elements/wizard/Wizard";

@customElement("gravity-dns-zone-wizard")
export class DNSZoneWizard extends AKElement {
    static get styles(): CSSResult[] {
        return [PFBase, PFButton, AKGlobal, PFRadio];
    }

    @property({ type: Boolean })
    open = false;

    @property()
    createText = "Create";

    @property({ type: Boolean })
    showButton = true;

    @property({ attribute: false })
    finalHandler: () => Promise<void> = () => {
        return Promise.resolve();
    };

    render(): TemplateResult {
        return html`
            <ak-wizard
                .open=${this.open}
                .steps=${["ak-application-wizard-initial", "ak-application-wizard-type"]}
                header=${"New zone"}
                description=${"Create a new DNS zone."}
                .finalHandler=${() => {
                    return this.finalHandler();
                }}
            >
                ${this.showButton
                    ? html`<button slot="trigger" class="pf-c-button pf-m-primary">
                          ${this.createText}
                      </button>`
                    : html``}
            </ak-wizard>
        `;
    }
}
