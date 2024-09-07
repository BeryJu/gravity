import { DiscoveryAPISubnet, RolesDiscoveryApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-discovery-subnet-form")
export class DiscoverySubnetForm extends ModelForm<DiscoveryAPISubnet, string> {
    async loadInstance(pk: string): Promise<DiscoveryAPISubnet> {
        const subnets = await new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryGetSubnets();
        const subnet = subnets.subnets?.find((z) => z.name === pk);
        if (!subnet) throw new Error("No subnet");
        return subnet;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated subnet.";
        } else {
            return "Successfully created subnet.";
        }
    }

    send = (data: DiscoveryAPISubnet): Promise<void> => {
        return new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryPutSubnets({
            identifier: this.instance?.name || data.name,
            discoveryAPISubnetsPutInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html` ${this.instance
                ? html``
                : html` <ak-form-element-horizontal label="Name" ?required=${true} name="name">
                      <input type="text" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal label="Discovery CIDR" ?required=${true} name="subnetCidr">
                <input type="text" value="${ifDefined(this.instance?.subnetCidr)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="DNS Resolver" ?required=${true} name="dnsResolver">
                <input type="text" value="${ifDefined(this.instance?.dnsResolver)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Default TTL"
                ?required=${true}
                name="discoveryTTL"
                helperText="Default TTL of discovered devices, in seconds."
            >
                <input
                    type="number"
                    value="${ifDefined(this.instance?.discoveryTTL || 86400)}"
                    required
                />
            </ak-form-element-horizontal>`;
    }
}
