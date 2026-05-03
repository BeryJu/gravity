import { DhcpAPILease, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { first, firstElement } from "../../utils";

@customElement("gravity-dhcp-lease-form")
export class DHCPLeaseForm extends ModelForm<DhcpAPILease, string> {
    @property()
    scope: string | undefined;

    @property()
    nextAddress?: string;

    async loadInstance(pk: string): Promise<DhcpAPILease> {
        const leases = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetLeases({
            scope: this.scope,
            identifier: pk,
        });
        const lease = firstElement(leases.leases);
        if (!lease) throw new Error("No lease");
        return lease;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated lease.";
        } else {
            return "Successfully created lease.";
        }
    }

    needsRecreate(data: DhcpAPILease): boolean {
        if (!this.instance) {
            return false;
        }
        if (data.identifier !== this.instance.identifier) return true;
        return false;
    }

    send = async (data: DhcpAPILease): Promise<void> => {
        if (!data.addressLeaseTime) {
            data.addressLeaseTime = "0";
        }
        if (!this.instance) {
            data.expiry = -1;
        } else {
            data.expiry = this.instance.expiry;
        }
        if (this.instance && this.needsRecreate(data)) {
            await new RolesDhcpApi(DEFAULT_CONFIG).dhcpDeleteLeases({
                scope: this.scope || "",
                identifier: this.instance.identifier,
            });
        }
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutLeases({
            scope: this.scope || "",
            identifier: data.identifier,
            dhcpAPILeasesPutInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html`
            <ak-form-element-horizontal label="Identifier" required name="identifier">
                <input
                    type="text"
                    value=${ifDefined(this.instance?.identifier)}
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    The identifier used by the device this lease is for. In most cases this will be
                    the MAC address of the network device.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Address" required name="address">
                <input
                    type="text"
                    value=${first(this.instance?.address, this.nextAddress, "")}
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Hostname" required name="hostname">
                <input
                    type="text"
                    value=${ifDefined(this.instance?.hostname)}
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Description" name="description">
                <input
                    type="text"
                    value=${ifDefined(this.instance?.description)}
                    class="pf-c-form-control"
                />
            </ak-form-element-horizontal>
        `;
    }
}
