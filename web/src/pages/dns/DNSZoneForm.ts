import { DnsAPIZone, RolesDnsApi } from "gravity-api";
import YAML from "yaml";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/CodeMirror";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { firstElement } from "../../utils";

export const DEFAULT_HANDLER_CONFIG = [
    {
        type: "memory",
    },
    {
        type: "etcd",
    },
    {
        type: "forward_ip",
        to: "8.8.8.8:53",
    },
];

@customElement("gravity-dns-zone-form")
export class DNSZoneForm extends ModelForm<DnsAPIZone, string> {
    async loadInstance(pk: string): Promise<DnsAPIZone> {
        const zones = await new RolesDnsApi(DEFAULT_CONFIG).dnsGetZones({
            name: pk,
        });
        const zone = firstElement(zones.zones);
        if (!zone) throw new Error("No zone");
        return zone;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated zone.";
        } else {
            return "Successfully created zone.";
        }
    }

    send = (data: DnsAPIZone): Promise<void> => {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsPutZones({
            zone: this.instance?.name || data.name,
            dnsAPIZonesPutInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html`
            ${this.instance
                ? html``
                : html` <ak-form-element-horizontal label="Name" required name="name">
                      <input type="text" class="pf-c-form-control" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal name="authoritative">
                <div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        ?checked=${this.instance?.authoritative}
                    />
                    <label class="pf-c-check__label"> ${"Authoritative"} </label>
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Default TTL" required name="defaultTTL">
                <input
                    type="number"
                    value=${this.instance?.defaultTTL || 0}
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    Default TTL for records which don't specify a non-zero value.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"Handler Configs"} name="handlerConfigs">
                <ak-codemirror
                    mode="yaml"
                    value=${YAML.stringify(this.instance?.handlerConfigs || DEFAULT_HANDLER_CONFIG)}
                >
                </ak-codemirror>
                <p class="pf-c-form__helper-text">
                    Configure where requests to this zone will be routed to and how they should be
                    answered.
                    <a href="https://gravity.beryju.io/docs/dns/zones/#handlers" target="_blank"
                        >Documentation</a
                    >
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"Hook"} name="hook">
                <ak-codemirror mode="javascript" value=${this.instance?.hook || ""}>
                </ak-codemirror>
                <p class="pf-c-form__helper-text">
                    Dynamically alter the DNS request/response after it is received and before it is
                    sent.
                    <a href="https://gravity.beryju.io/docs/dns/hooks/" target="_blank"
                        >Documentation</a
                    >
                </p>
            </ak-form-element-horizontal>
        `;
    }
}
