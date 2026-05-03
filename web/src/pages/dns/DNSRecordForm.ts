import { DnsAPIRecord, RolesDnsApi } from "gravity-api";

import { TemplateResult, html, nothing } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-dns-record-form")
export class DNSRecordForm extends ModelForm<DnsAPIRecord, string> {
    @property()
    zone: string | undefined;

    @property()
    recordType = "A";

    async loadInstance(pk: string): Promise<DnsAPIRecord> {
        const records = await new RolesDnsApi(DEFAULT_CONFIG).dnsGetRecords({
            zone: this.zone,
        });
        const record = records.records?.find((z) => z.hostname + z.uid === pk);
        if (!record) throw new Error("No record");
        this.recordType = record.type;
        return record;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated record.";
        } else {
            return "Successfully created record.";
        }
    }

    needsRecreate(data: DnsAPIRecord): boolean {
        if (!this.instance) {
            return false;
        }
        if (data.hostname !== this.instance.hostname) return true;
        if (data.uid !== this.instance.uid) return true;
        if (data.type !== this.instance.type) return true;
        return false;
    }

    send = async (data: DnsAPIRecord): Promise<void> => {
        if (this.instance && this.needsRecreate(data)) {
            await new RolesDnsApi(DEFAULT_CONFIG).dnsDeleteRecords({
                zone: this.zone || "",
                ...this.instance,
            });
        }
        return new RolesDnsApi(DEFAULT_CONFIG).dnsPutRecords({
            zone: this.zone || "",
            ...data,
            dnsAPIRecordsPutInput: data,
        });
    };

    getLabel(): string {
        switch (this.recordType) {
            case "CNAME":
                return "CNAME Target";
            case "SRV":
                return "SRV Target";
            case "MX":
                return "Mail server";
            case "AAAA":
            case "A":
                return "IP Address";
            default:
                return "Data";
        }
    }

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal label="Hostname" required name="hostname">
                <div class="pf-c-input-group">
                    <input
                        type="text"
                        value=${ifDefined(this.instance?.hostname)}
                        class="pf-c-form-control"
                        required
                    />
                    ${this.zone !== "."
                        ? html`<span class="pf-c-input-group__text">.${this.zone}</span>`
                        : nothing}
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="UID" required name="uid">
                <input
                    type="text"
                    value=${this.instance?.uid || ""}
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    Unique identifier to configure multiple records for the same hostname.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Type" required name="type">
                <select
                    class="pf-c-form-control"
                    @change=${(ev: Event) => {
                        const current = (ev.target as HTMLInputElement).value;
                        this.recordType = current;
                    }}
                >
                    <option value="A" ?selected=${this.recordType === "A"}>A</option>
                    <option value="AAAA" ?selected=${this.recordType === "AAAA"}>AAAA</option>
                    <option value="CNAME" ?selected=${this.recordType === "CNAME"}>CNAME</option>
                    <option value="PTR" ?selected=${this.recordType === "PTR"}>PTR</option>
                    <option value="NS" ?selected=${this.recordType === "NS"}>NS</option>
                    <option value="MX" ?selected=${this.recordType === "MX"}>MX</option>
                    <option value="SRV" ?selected=${this.recordType === "SRV"}>SRV</option>
                    <option value="TXT" ?selected=${this.recordType === "TXT"}>TXT</option>
                </select>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${this.getLabel()} required name="data">
                <input
                    type="text"
                    value=${ifDefined(this.instance?.data)}
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            ${this.recordType === "MX"
                ? html`
                      <ak-form-element-horizontal
                          label="MX Preference"
                          required
                          name="mxPreference"
                      >
                          <input
                              type="number"
                              value=${ifDefined(this.instance?.mxPreference)}
                              class="pf-c-form-control"
                              required
                          />
                      </ak-form-element-horizontal>
                  `
                : html``}
            ${this.recordType === "SRV"
                ? html`
                      <ak-form-element-horizontal label="SRV Port" required name="srvPort">
                          <input
                              type="number"
                              value=${ifDefined(this.instance?.srvPort)}
                              class="pf-c-form-control"
                              required
                          />
                      </ak-form-element-horizontal>
                      <ak-form-element-horizontal label="SRV Priority" required name="srvPriority">
                          <input
                              type="number"
                              value=${ifDefined(this.instance?.srvPriority)}
                              class="pf-c-form-control"
                              required
                          />
                      </ak-form-element-horizontal>
                      <ak-form-element-horizontal label="SRV Weight" required name="srvWeight">
                          <input
                              type="number"
                              value=${ifDefined(this.instance?.srvWeight)}
                              class="pf-c-form-control"
                              required
                          />
                      </ak-form-element-horizontal>
                  `
                : html``}`;
    }
}
