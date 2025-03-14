import { DnsAPIZonesPutInput, RolesDnsApi } from "gravity-api";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dns-wizard-initial")
export class ZoneInitialWizardPage extends WizardFormPage {
    sidebarLabel = () => "Zone details";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        let name = data.name as string;
        if (!name.endsWith(".")) {
            name += ".";
        }
        const zone: DnsAPIZonesPutInput = {
            authoritative: data.authoritative as boolean,
            handlerConfigs: [],
            defaultTTL: 86400,
            hook: "",
        };
        this.host.state["handlerConfigs"] = [];
        this.host.addActionBefore("Create zone", "create-zone", async (): Promise<boolean> => {
            zone.handlerConfigs = this.host.state["handlerConfigs"] as {
                [key: string]: string;
            }[];
            this.host.state["zone"] = await new RolesDnsApi(DEFAULT_CONFIG).dnsPutZones({
                zone: name,
                dnsAPIZonesPutInput: zone,
            });
            this.host.state["name"] = name;
            return true;
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label=${"Name"} required name="name">
                <input type="text" value="" class="pf-c-form-control" required />
                <p class="pf-c-form__helper-text">
                    The zone name specifies which DNS namespace this zone is responsible for. This
                    might be a domain name (beryju.io), a subdomain (gravity.beryju.io), or
                    everything (.).
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal name="authoritative">
                <div class="pf-c-check">
                    <input type="checkbox" class="pf-c-check__input" />
                    <label class="pf-c-check__label"> ${"Authoritative"} </label>
                </div>
                <p class="pf-c-form__helper-text">
                    Determines whether Gravity holds the source of truth for the domain specified.
                </p>
            </ak-form-element-horizontal>`;
    }
}
