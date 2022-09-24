import { DnsRecord, RolesDnsApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-dns-record-form")
export class DNSRecordForm extends ModelForm<DnsRecord, string> {
    @property()
    zone?: string;

    loadInstance(pk: string): Promise<DnsRecord> {
        return new RolesDnsApi(DEFAULT_CONFIG)
            .dnsGetRecords({
                zone: this.zone,
            })
            .then((records) => {
                const record = records.records?.find((z) => z.hostname + z.uid === pk);
                if (!record) throw new Error("No record");
                return record;
            });
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return `Successfully updated record.`;
        } else {
            return `Successfully created record.`;
        }
    }

    send = (data: DnsRecord): Promise<void> => {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsPutRecords({
            zone: this.zone || "",
            hostname: data.hostname,
            uid: data.uid,
            dnsRecordsInputType2: data,
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Hostname" ?required=${true} name="hostname">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.hostname)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Type" ?required=${true} name="type">
                <select class="pf-c-form-control">
                    <option value="A" ?selected=${this.instance?.type === "A"}>A</option>
                    <option value="AAAA" ?selected=${this.instance?.type === "AAAA"}>AAAA</option>
                    <option value="PTR" ?selected=${this.instance?.type === "PTR"}>PTR</option>
                    <option value="NS" ?selected=${this.instance?.type === "NS"}>NS</option>
                    <option value="MX" ?selected=${this.instance?.type === "MX"}>MX</option>
                </select>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Data" ?required=${true} name="data">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.data)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
        </form>`;
    }
}
