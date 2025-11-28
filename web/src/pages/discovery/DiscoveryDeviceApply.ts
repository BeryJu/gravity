import {
    DiscoveryAPIDevice,
    DiscoveryAPIDevicesApplyInput,
    DiscoveryAPIDevicesApplyInputToEnum,
    RolesDhcpApi,
    RolesDiscoveryApi,
    RolesDnsApi,
} from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { until } from "lit/directives/until.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { DeleteBulkForm } from "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/DeleteBulkForm";
import { Form } from "../../elements/forms/Form";
import "../../elements/forms/HorizontalFormElement";

@customElement("gravity-discover-apply-form")
export class DiscoveryDeviceApplyForm extends Form<DiscoveryAPIDevicesApplyInput> {
    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal label="To" required name="to">
                <select class="pf-c-form-control">
                    <option value="dhcp">
                        DHCP (will also create a DNS record if the DHCP Scope is DNS integrated)
                    </option>
                    <option value="dns">DNS</option>
                </select>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="DHCP Scope" required name="dhcpScope">
                <select class="pf-c-form-control">
                    <option value="">---</option>
                    ${until(
                        new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes().then((scopes) => {
                            return scopes.scopes?.map((scope) => {
                                return html`<option value=${scope.scope}>
                                    ${scope.scope} (${scope.subnetCidr})
                                </option>`;
                            });
                        }),
                    )}
                </select>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="DNS Zone" required name="dnsZone">
                <select class="pf-c-form-control">
                    <option value="">---</option>
                    ${until(
                        new RolesDnsApi(DEFAULT_CONFIG).dnsGetZones().then((zones) => {
                            return zones.zones?.map((zone) => {
                                return html`<option value=${zone.name}>${zone.name}</option>`;
                            });
                        }),
                    )}
                </select>
            </ak-form-element-horizontal>`;
    }
}

@customElement("gravity-discovery-apply")
export class DiscoveryDeviceApply extends DeleteBulkForm {
    metadata = (item: DiscoveryAPIDevice) => {
        return [
            { key: "IP", value: item.ip },
            { key: "Hostname", value: item.hostname || "-" },
            { key: "MAC", value: item.mac || "-" },
        ];
    };

    preDelete = (): Promise<DiscoveryAPIDevicesApplyInput> => {
        const form = this.shadowRoot?.querySelector<DiscoveryDeviceApplyForm>(
            "gravity-discover-apply-form",
        );
        if (!form) {
            return Promise.reject("No form found");
        }
        const data = form.serializeForm() as DiscoveryAPIDevicesApplyInput;
        if (data.to === DiscoveryAPIDevicesApplyInputToEnum.Dhcp && data.dhcpScope === "") {
            throw Error("DHCP Scope needs to be set to import to DHCP.");
        }
        if (data.to === DiscoveryAPIDevicesApplyInputToEnum.Dns && data.dnsZone === "") {
            throw Error("DNS Zone needs to be set to import to DNS.");
        }
        return Promise.resolve(data);
    };

    delete = (
        item: DiscoveryAPIDevice,
        extraData: DiscoveryAPIDevicesApplyInput,
    ): Promise<void> => {
        if (item.mac === "" && extraData.to === DiscoveryAPIDevicesApplyInputToEnum.Dhcp) {
            return Promise.reject();
        }
        if (item.hostname === "" && extraData.to === DiscoveryAPIDevicesApplyInputToEnum.Dns) {
            return Promise.reject();
        }
        return new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryApplyDevice({
            identifier: item.identifier,
            discoveryAPIDevicesApplyInput: extraData,
        });
    };

    renderModalInner(): TemplateResult {
        return html`<section class="pf-c-modal-box__header pf-c-page__main-section pf-m-light">
                <div class="pf-c-content">
                    <h1 class="pf-c-title pf-m-2xl">${`Apply ${this.objectLabel}`}</h1>
                </div>
            </section>
            <section class="pf-c-modal-box__body pf-c-page__main-section pf-m-light">
                <form class="pf-c-form pf-m-horizontal">
                    <p class="pf-c-title">
                        ${`Are you sure you want to apply ${this.objects.length} ${this.objectLabel}?`}
                    </p>
                    <slot name="notice"></slot>
                </form>
            </section>
            <section class="pf-c-modal-box__body pf-c-page__main-section pf-m-light">
                <gravity-discover-apply-form></gravity-discover-apply-form>
            </section>
            <section class="pf-c-modal-box__body pf-c-page__main-section pf-m-light">
                <!-- @ts-ignore -->
                <ak-delete-objects-table .objects=${this.objects} .metadata=${this.metadata}>
                </ak-delete-objects-table>
            </section>
            <footer class="pf-c-modal-box__footer">
                <ak-spinner-button
                    .callAction=${() => {
                        return this.confirm();
                    }}
                    class="pf-m-primary"
                >
                    ${"Apply"} </ak-spinner-button
                >&nbsp;
                <ak-spinner-button
                    .callAction=${async () => {
                        this.open = false;
                    }}
                    class="pf-m-secondary"
                >
                    ${"Cancel"}
                </ak-spinner-button>
            </footer>`;
    }
}
