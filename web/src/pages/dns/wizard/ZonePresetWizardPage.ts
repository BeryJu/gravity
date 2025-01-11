import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { CSSResult, TemplateResult, html } from "lit";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFForm from "@patternfly/patternfly/components/Form/form.css";
import PFRadio from "@patternfly/patternfly/components/Radio/radio.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../../elements/Base";
import { Wizard } from "../../../elements/wizard/Wizard";
import { WizardPage } from "../../../elements/wizard/WizardPage";

export interface TypeCreate {
    name: string;
    description: string;
    components: string[];
    callback: (host: Wizard) => void;
}

@customElement("gravity-dns-wizard-type")
export class ZonePresetWizardPage extends WizardPage {
    applicationTypes: TypeCreate[] = [
        {
            components: [],
            name: "Gravity-hosted",
            description: "Records are stored in Gravity, optionally integrated with DHCP.",
            callback: (host: Wizard) => {
                host.state["handlerConfigs"] = [
                    {
                        type: "memory",
                    },
                    {
                        type: "etcd",
                    },
                ];
            },
        },
        {
            components: ["gravity-dns-wizard-forward", "gravity-dns-wizard-cache"],
            name: "Forwarder (direct)",
            description: "Forward requests to an upstream resolver, and optionally cache them.",
            callback: (host: Wizard) => {
                host.state["handlerConfigs"] = [
                    {
                        type: "forward_ip",
                    },
                ];
            },
        },
        {
            components: ["gravity-dns-wizard-forward", "gravity-dns-wizard-cache"],
            name: "Forwarder (Blocky)",
            description: "Forward requests through Blocky, providing Ad and privacy blocking.",
            callback: (host: Wizard) => {
                host.state["handlerConfigs"] = [
                    {
                        type: "forward_blocky",
                    },
                    {
                        type: "forward_ip",
                    },
                ];
            },
        },
        {
            components: ["gravity-dns-wizard-import"],
            name: "Import",
            description: "Import DNS records from various different formats.",
            callback: (host: Wizard) => {
                host.state["handlerConfigs"] = [
                    {
                        type: "etcd",
                    },
                ];
            },
        },
    ];

    sidebarLabel = () => "Zone configuration preset";

    static get styles(): CSSResult[] {
        return [PFBase, PFButton, PFForm, PFRadio, AKElement.GlobalStyle];
    }

    activeCallback = async () => {
        this.host.isValid = true;
        this.applicationTypes[0].callback(this.host);
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
                                "gravity-dns-wizard-initial",
                                "gravity-dns-wizard-type",
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
