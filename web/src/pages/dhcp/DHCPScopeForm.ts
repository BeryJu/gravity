import { DhcpAPIScope, RolesDhcpApi } from "gravity-api";
import YAML from "yaml";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/CodeMirror";
import "../../elements/forms/FormGroup";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { KV, first, firstElement } from "../../utils";

@customElement("gravity-dhcp-scope-form")
export class DHCPScopeForm extends ModelForm<DhcpAPIScope, string> {
    async loadInstance(pk: string): Promise<DhcpAPIScope> {
        const scopes = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes({
            name: pk,
        });
        const zone = firstElement(scopes.scopes);
        if (!zone) throw new Error("No scope");
        return zone;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated scope.";
        } else {
            return "Successfully created scope.";
        }
    }

    send = (data: DhcpAPIScope): Promise<void> => {
        if (data.ipam) {
            data.ipam.type = "internal";
        }
        if (!data.options) {
            data.options = [];
        }
        const routerOpts = data.options.filter((op) => op.tagName === "router");
        if (routerOpts.length < 1) {
            data.options.push({
                tagName: "router",
                value: (data as unknown as KV)["router"],
            });
        }
        routerOpts
            .filter((op) => op.tagName === "router")
            .forEach((op) => {
                op.value = (data as unknown as KV)["router"];
            });
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutScopes({
            scope: this.instance?.scope || data.scope,
            dhcpAPIScopesPutInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html` ${this.instance
                ? html``
                : html`<ak-form-element-horizontal label="Name" required name="scope">
                      <input type="text" class="pf-c-form-control" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal label="Subnet CIDR" required name="subnetCidr">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.subnetCidr)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    CIDR for which this scope is authoritative for.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Router" name="router">
                <input
                    type="text"
                    value="${ifDefined(
                        this.instance?.options
                            ?.filter((op) => op.tagName === "router")
                            .map((op) => op.value || "")
                            .join(""),
                    )}"
                    class="pf-c-form-control"
                />
                <p class="pf-c-form__helper-text">Router for the subnet.</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal name="_default">
                <div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        ?checked=${this.instance?._default}
                    />
                    <label class="pf-c-check__label"> ${"Default"} </label>
                </div>
                <p class="pf-c-form__helper-text">
                    If checked, this scope will be used for clients when their network can't be
                    determined.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="TTL" required name="ttl">
                <input
                    type="number"
                    value="${this.instance?.ttl || 86400}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">Default TTL of leases, in seconds.</p>
            </ak-form-element-horizontal>
            <ak-form-group expanded>
                <span slot="header">IPAM</span>
                <div slot="body" class="pf-c-form">
                    <ak-form-element-horizontal
                        label="IP Range Start"
                        required
                        name="ipam.range_start"
                    >
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.ipam?.range_start)}"
                            class="pf-c-form-control"
                            required
                        />
                        <p class="pf-c-form__helper-text">Start of the IP range, inclusive.</p>
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="IP Range End" required name="ipam.range_end">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.ipam?.range_end)}"
                            class="pf-c-form-control"
                            required
                        />
                        <p class="pf-c-form__helper-text">End of the IP range, exclusive.</p>
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal name="ipam.should_ping">
                        <div class="pf-c-check">
                            <input
                                type="checkbox"
                                class="pf-c-check__input"
                                ?checked=${first(
                                    this.instance?.ipam?.should_ping as unknown as boolean,
                                    false,
                                )}
                            />
                            <label class="pf-c-check__label"
                                >Ping IP Address before assigning it.</label
                            >
                        </div>
                    </ak-form-element-horizontal>
                </div>
            </ak-form-group>
            <ak-form-group>
                <span slot="header">DNS settings</span>
                <div slot="body" class="pf-c-form">
                    <ak-form-element-horizontal label="DNS Zone" name="dns.zone">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.dns?.zone)}"
                            class="pf-c-form-control"
                        />
                        <p class="pf-c-form__helper-text">
                            Optional, set to a DNS zone configured in Gravity to create DNS records.
                            If the configured zone does not exist in Gravity, it is only used as
                            domain for the leases.
                        </p>
                    </ak-form-element-horizontal>
                </div>
            </ak-form-group>
            <ak-form-group>
                <span slot="header">Advanced settings</span>
                <div slot="body" class="pf-c-form">
                    <ak-form-element-horizontal label=${"DHCP Options"} name="options">
                        <ak-codemirror
                            mode="yaml"
                            value="${YAML.stringify(this.instance?.options)}"
                        >
                        </ak-codemirror>
                        <p class="pf-c-form__helper-text">
                            Add additional DHCP options
                            <a href="https://gravity.beryju.io/docs/dhcp/options/" target="_blank"
                                >Documentation</a
                            >
                        </p>
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label=${"Hook"} name="hook">
                        <ak-codemirror mode="javascript" value="${this.instance?.hook || ""}">
                        </ak-codemirror>
                        <p class="pf-c-form__helper-text">
                            Dynamically alter the DHCP request/response after it is received and
                            before it is sent.
                            <a href="https://gravity.beryju.io/docs/dhcp/hooks/" target="_blank"
                                >Documentation</a
                            >
                        </p>
                    </ak-form-element-horizontal>
                </div>
            </ak-form-group>`;
    }
}
