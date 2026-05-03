import { DhcpAPIScopesPutInput, DnsAPIZonesPutInput, RolesDnsApi } from "gravity-api";
import { IPv4CidrRange } from "ip-num";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";
import { subnetToDnsZone } from "../../../utils";

const defaultDNSZoneSettings: DnsAPIZonesPutInput = {
    authoritative: true,
    defaultTTL: 86400,
    hook: "",
    handlerConfigs: [
        {
            type: "memory",
        },
        {
            type: "etcd",
        },
    ],
};

@customElement("gravity-dhcp-wizard-dns")
export class ScopeDNSWizardPage extends WizardFormPage {
    sidebarLabel = () => "DNS configuration";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        const req = this.host.state["scopeReq"] as DhcpAPIScopesPutInput;
        req.dns = {
            zone: data.zone as string,
        };
        this.host.state["scopeReq"] = req;
        if (data.createZone) {
            this.host.addActionAfter(
                "Create DNS Zone (forward)",
                "create-dns-forward",
                async (): Promise<boolean> => {
                    new RolesDnsApi(DEFAULT_CONFIG).dnsPutZones({
                        zone: data.zone as string,
                        dnsAPIZonesPutInput: defaultDNSZoneSettings,
                    });
                    return true;
                },
            );
        }
        if (data.createReverseZone) {
            this.host.addActionAfter(
                "Create DNS Zone (reverse)",
                "create-dns-reverse",
                async (): Promise<boolean> => {
                    new RolesDnsApi(DEFAULT_CONFIG).dnsPutZones({
                        zone: subnetToDnsZone(IPv4CidrRange.fromCidr(req.subnetCidr)),
                        dnsAPIZonesPutInput: defaultDNSZoneSettings,
                    });
                    return true;
                },
            );
        }
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="DNS Zone" name="zone" required>
                <input type="text" value="" class="pf-c-form-control" required />
                <p class="pf-c-form__helper-text">
                    DNS Zone which records in this scope should be registered in. Reverse records
                    are created automatically if a matching reverse-DNS zone exists.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal name="createZone">
                <div class="pf-c-check">
                    <input type="checkbox" class="pf-c-check__input" checked />
                    <label class="pf-c-check__label">${"Create DNS Zone"}</label>
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal name="createReverseZone">
                <div class="pf-c-check">
                    <input type="checkbox" class="pf-c-check__input" checked />
                    <label class="pf-c-check__label">${"Create reverse DNS Zone"}</label>
                </div>
            </ak-form-element-horizontal>`;
    }
}
