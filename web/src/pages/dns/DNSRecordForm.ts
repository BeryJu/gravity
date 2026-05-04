import { DnsAPIRecord, RolesDnsApi, TypesDNSRecordType } from "gravity-api";

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
    recordType: TypesDNSRecordType = TypesDNSRecordType.A;

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
            case TypesDNSRecordType.Cname:
                return "CNAME Target";
            case TypesDNSRecordType.Srv:
                return "SRV Target";
            case TypesDNSRecordType.Mx:
                return "Mail server";
            case TypesDNSRecordType.Aaaa:
            case TypesDNSRecordType.A:
                return "IP Address";
            default:
                return "Data";
        }
    }

    renderTypeSpecific() {
        switch (this.recordType) {
            case TypesDNSRecordType.Mx:
                return html`<ak-form-element-horizontal
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
                </ak-form-element-horizontal>`;
            case TypesDNSRecordType.Srv:
                return html`<ak-form-element-horizontal label="SRV Port" required name="srvPort">
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
                    </ak-form-element-horizontal>`;
            case TypesDNSRecordType.Soa:
                return html`<ak-form-element-horizontal
                        label="SOA Expire"
                        required
                        name="soaExpire"
                    >
                        <input
                            type="number"
                            value=${ifDefined(this.instance?.soaExpire)}
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="SOA Mailbox" required name="soaMbox">
                        <input
                            type="text"
                            value=${ifDefined(this.instance?.soaMbox)}
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="SOA Refresh" required name="soaRefresh">
                        <input
                            type="number"
                            value=${ifDefined(this.instance?.soaRefresh)}
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="SOA Retry" required name="soaRetry">
                        <input
                            type="number"
                            value=${ifDefined(this.instance?.soaRetry)}
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="SOA Serial" required name="soaSerial">
                        <input
                            type="number"
                            value=${ifDefined(this.instance?.soaSerial)}
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>`;
        }
        return html``;
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
                        this.recordType = current as TypesDNSRecordType;
                    }}
                >
                    <option
                        value=${TypesDNSRecordType.A}
                        ?selected=${this.recordType === TypesDNSRecordType.A}
                    >
                        A
                    </option>
                    <option
                        value=${TypesDNSRecordType.Aaaa}
                        ?selected=${this.recordType === TypesDNSRecordType.Aaaa}
                    >
                        AAAA
                    </option>
                    <option
                        value=${TypesDNSRecordType.Cname}
                        ?selected=${this.recordType === TypesDNSRecordType.Cname}
                    >
                        CNAME
                    </option>
                    <option
                        value=${TypesDNSRecordType.Ptr}
                        ?selected=${this.recordType === TypesDNSRecordType.Ptr}
                    >
                        PTR
                    </option>
                    <option
                        value=${TypesDNSRecordType.Mx}
                        ?selected=${this.recordType === TypesDNSRecordType.Mx}
                    >
                        MX
                    </option>
                    <option
                        value=${TypesDNSRecordType.Srv}
                        ?selected=${this.recordType === TypesDNSRecordType.Srv}
                    >
                        SRV
                    </option>
                    <option
                        value=${TypesDNSRecordType.Txt}
                        ?selected=${this.recordType === TypesDNSRecordType.Txt}
                    >
                        TXT
                    </option>
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
            ${this.renderTypeSpecific()}`;
    }
}
