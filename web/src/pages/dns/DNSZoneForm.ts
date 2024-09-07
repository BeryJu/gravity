import { DnsAPIZone, RolesDnsApi } from "gravity-api";
import YAML from "yaml";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/CodeMirror";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { first } from "../../utils";

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
        const zone = first(zones.zones);
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
        return html` ${this.instance
                ? html``
                : html` <ak-form-element-horizontal label="Name" ?required=${true} name="name">
                      <input type="text" required />
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
            <ak-form-element-horizontal label="Default TTL" ?required=${true} name="defaultTTL">
                <input type="number" value="${this.instance?.defaultTTL || 0}" required />
                <p class="pf-c-form__helper-text">
                    Default TTL for records which don't specify a non-zero value.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"Handler Configs"} name="handlerConfigs">
                <ak-codemirror
                    mode="yaml"
                    value="${YAML.stringify(
                        this.instance?.handlerConfigs || DEFAULT_HANDLER_CONFIG,
                    )}"
                >
                </ak-codemirror>
            </ak-form-element-horizontal>`;
    }
}
