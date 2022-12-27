import { DnsAPIZone, RolesDnsApi } from "gravity-api";

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
        const zone: DnsAPIZone = {
            authoritative: false,
            handlerConfigs: [],
            defaultTTL: 0,
            name: name,
        };
        this.host.state["handlerConfigs"] = [];
        this.host.addActionBefore("Create zone", async (): Promise<boolean> => {
            zone.handlerConfigs = this.host.state["handlerConfigs"] as {
                [key: string]: string;
            }[];
            this.host.state["zone"] = await new RolesDnsApi(DEFAULT_CONFIG).dnsPutZones({
                zone: zone.name,
                dnsAPIZonesPutInput: zone,
            });
            return true;
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`
            <form class="pf-c-form pf-m-horizontal">
                <ak-form-element-horizontal label=${"Name"} ?required=${true} name="name">
                    <input type="text" value="" class="pf-c-form-control" required />
                    <p class="pf-c-form__helper-text">
                        The zone name specifies which DNS namespace this zone is responsible for.
                        This might be a domain name (beryju.io), a subdomain (gravity.beryju.io), or
                        everything (.).
                    </p>
                </ak-form-element-horizontal>
                <ak-form-element-horizontal name="authoritative">
                    <div class="pf-c-check">
                        <input type="checkbox" class="pf-c-check__input" />
                        <label class="pf-c-check__label"> ${"Authoritative"} </label>
                    </div>
                    <p class="pf-c-form__helper-text">
                        Determines whether Gravity holds the source of truth for the domain
                        specified.
                    </p>
                </ak-form-element-horizontal>
            </form>
        `;
    }
}
