import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { CSSResult, TemplateResult, html } from "lit";
import { property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFRadio from "@patternfly/patternfly/components/Radio/radio.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../../elements/Base";
import "../../../elements/wizard/Wizard";
import "./ClusterJoinCompose";
import "./ClusterJoinInitial";

@customElement("gravity-cluster-join-wizard")
export class ClusterJoinWizard extends AKElement {
    static get styles(): CSSResult[] {
        return [PFBase, PFButton, AKElement.GlobalStyle, PFRadio];
    }

    @property({ type: Boolean })
    accessor open = false;

    @property()
    accessor createText = "Join";

    @property({ type: Boolean })
    accessor showButton = true;

    render(): TemplateResult {
        return html`
            <ak-wizard
                .open=${this.open}
                .steps=${["gravity-cluster-join-initial", "gravity-cluster-join-compose"]}
                header=${"Join a node"}
                description=${"Join a gravity node to the cluster."}
            >
                <button slot="trigger" class="pf-c-button pf-m-primary">${this.createText}</button>
            </ak-wizard>
        `;
    }
}
