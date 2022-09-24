import { DhcpLease, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-dhcp-lease-form")
export class DHCPLeaseForm extends ModelForm<DhcpLease, string> {
    @property()
    scope?: string;

    loadInstance(pk: string): Promise<DhcpLease> {
        return new RolesDhcpApi(DEFAULT_CONFIG)
            .dhcpGetLeases({
                scope: this.scope,
            })
            .then((leases) => {
                const lease = leases.leases?.find((z) => z.identifier === pk);
                if (!lease) throw new Error("No lease");
                return lease;
            });
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated lease.";
        } else {
            return "Successfully created lease.";
        }
    }

    send = (data: DhcpLease): Promise<void> => {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutLeases({
            scope: this.scope || "",
            identifier: data.identifier,
            dhcpLeasesInputType2: data,
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Identifier" ?required=${true} name="identifier">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.identifier)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Address" ?required=${true} name="address">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.address)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Hostname" ?required=${true} name="hostname">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.hostname)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
        </form>`;
    }
}
