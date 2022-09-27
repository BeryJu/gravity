import { DnsZone, RolesDnsApi } from "gravity-api";
import YAML from "yaml";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/CodeMirror";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

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
export class DNSZoneForm extends ModelForm<DnsZone, string> {
    loadInstance(pk: string): Promise<DnsZone> {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsGetZones().then((zones) => {
            const zone = zones.zones?.find((z) => z.name === pk);
            if (!zone) throw new Error("No zone");
            return zone;
        });
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated zone.";
        } else {
            return "Successfully created zone.";
        }
    }

    send = (data: DnsZone): Promise<void> => {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsPutZones({
            zone: this.instance?.name || "",
            dnsZoneInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            ${this.instance
                ? html``
                : html` <ak-form-element-horizontal label="Name" ?required=${true} name="name">
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
            <ak-form-element-horizontal label="Default TTL" ?required=${true} name="defaultTTL">
                <input
                    type="number"
                    value="${this.instance?.defaultTTL || 0}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"Handler Configs"} name="handlerConfigs">
                <ak-codemirror
                    mode="yaml"
                    value="${YAML.stringify(
                        this.instance?.handlerConfigs || DEFAULT_HANDLER_CONFIG,
                    )}"
                >
                </ak-codemirror>
            </ak-form-element-horizontal>
        </form>`;
    }
}
