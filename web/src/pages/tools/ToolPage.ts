import {
    ApiAPIToolPingOutput,
    ApiAPIToolPortmapOutput,
    ApiAPIToolTracerouteOutput,
    RolesApiApi,
} from "gravity-api";

import { CSSResult, TemplateResult, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFCard from "@patternfly/patternfly-v6/components/Card/card.css";
import PFContent from "@patternfly/patternfly-v6/components/Content/content.css";
import PFDataList from "@patternfly/patternfly-v6/components/DataList/data-list.css";
import PFDescriptionList from "@patternfly/patternfly-v6/components/DescriptionList/description-list.css";
import PFForm from "@patternfly/patternfly-v6/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly-v6/components/FormControl/form-control.css";
import PFInputGroup from "@patternfly/patternfly-v6/components/InputGroup/input-group.css";
import PFPage from "@patternfly/patternfly-v6/components/Page/page.css";
import PFSidebar from "@patternfly/patternfly-v6/components/Sidebar/sidebar.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKElement } from "../../elements/Base";
import { PFColor } from "../../elements/Label";
import "../../elements/buttons/SpinnerButton";
import { getURLParam, updateURLParams } from "../../elements/router/RouteMatch";

@customElement("gravity-tools")
export class ToolPage extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFPage,
            PFContent,
            PFDescriptionList,
            PFForm,
            PFButton,
            PFDataList,
            PFInputGroup,
            PFFormControl,
            PFSidebar,
            PFCard,
            AKElement.GlobalStyle,
        ];
    }

    @state()
    host?: string = getURLParam("host", undefined);

    @state()
    pingOutput?: ApiAPIToolPingOutput;

    @state()
    tracerouteOutput?: ApiAPIToolTracerouteOutput;

    @state()
    portmapOutput?: ApiAPIToolPortmapOutput;

    @state()
    errorOutput?: Error;

    renderPing(): TemplateResult {
        if (!this.pingOutput) return html``;
        return html`<dl class="pf-v6-c-description-list pf-m-horizontal">
            <div class="pf-v6-c-description-list__group">
                <dt class="pf-v6-c-description-list__term">
                    <span class="pf-v6-c-description-list__text">Packets sent</span>
                </dt>
                <dd class="pf-v6-c-description-list__description">
                    <div class="pf-v6-c-description-list__text">${this.pingOutput.packetsSent}</div>
                </dd>
            </div>
            <div class="pf-v6-c-description-list__group">
                <dt class="pf-v6-c-description-list__term">
                    <span class="pf-v6-c-description-list__text">Packets received</span>
                </dt>
                <dd class="pf-v6-c-description-list__description">
                    <div class="pf-v6-c-description-list__text">${this.pingOutput.packetsRecv}</div>
                </dd>
            </div>
            <div class="pf-v6-c-description-list__group">
                <dt class="pf-v6-c-description-list__term">
                    <span class="pf-v6-c-description-list__text"
                        >Packets received (duplicates)</span
                    >
                </dt>
                <dd class="pf-v6-c-description-list__description">
                    <div class="pf-v6-c-description-list__text">
                        ${this.pingOutput.packetsRecvDuplicates}
                    </div>
                </dd>
            </div>
            <div class="pf-v6-c-description-list__group">
                <dt class="pf-v6-c-description-list__term">
                    <span class="pf-v6-c-description-list__text">Packet loss</span>
                </dt>
                <dd class="pf-v6-c-description-list__description">
                    <div class="pf-v6-c-description-list__text">${this.pingOutput.packetLoss}</div>
                </dd>
            </div>
            <div class="pf-v6-c-description-list__group">
                <dt class="pf-v6-c-description-list__term">
                    <span class="pf-v6-c-description-list__text">Average Round-trip-time</span>
                </dt>
                <dd class="pf-v6-c-description-list__description">
                    <div class="pf-v6-c-description-list__text">${this.pingOutput.avgRtt}</div>
                </dd>
            </div>
        </dl>`;
    }

    renderTraceroute(): TemplateResult {
        if (!this.tracerouteOutput) return html``;
        return html`<ul class="pf-v6-c-data-list">
            ${this.tracerouteOutput.hops?.map((hop) => {
                return html`<li
                    class="pf-v6-c-data-list__item"
                    aria-labelledby="data-list-basic-item-1"
                >
                    <div class="pf-v6-c-data-list__item-row">
                        <div class="pf-v6-c-data-list__item-content">
                            <div class="pf-v6-c-data-list__cell">
                                <ak-label
                                    color=${hop.success ? PFColor.Green : PFColor.Orange}
                                ></ak-label>
                            </div>
                            <div class="pf-v6-c-data-list__cell">${hop.address}</div>
                            <div class="pf-v6-c-data-list__cell">${hop.elapsedTime}</div>
                        </div>
                    </div>
                </li>`;
            })}
        </ul>`;
    }

    renderPortmap(): TemplateResult {
        if (!this.portmapOutput) return html``;
        return html`<ul class="pf-v6-c-data-list">
            ${this.portmapOutput.ports?.map((port) => {
                return html`<li class="pf-v6-c-data-list__item">
                    <div class="pf-v6-c-data-list__item-row">
                        <div class="pf-v6-c-data-list__item-content">
                            <div class="pf-v6-c-data-list__cell">${port.reason}</div>
                            <div class="pf-v6-c-data-list__cell">${port.name} (${port.port})</div>
                            <div class="pf-v6-c-data-list__cell">${port.protocol}</div>
                        </div>
                    </div>
                </li>`;
            })}
        </ul>`;
    }

    renderResult(): TemplateResult {
        if (this.errorOutput) return html`${this.errorOutput}`;
        if (this.pingOutput) return this.renderPing();
        if (this.tracerouteOutput) return this.renderTraceroute();
        if (this.portmapOutput) return this.renderPortmap();
        return html`<p>No tool used yet</p>`;
    }

    render(): TemplateResult {
        return html`
            <section class="pf-v6-c-page__main-section pf-m-no-padding-mobile">
                <div class="pf-v6-c-card">
                    <div class="pf-v6-c-card__body">
                        <div class="pf-v6-c-input-group">
                            <input
                                class="pf-v6-c-form-control"
                                type="text"
                                placeholder="Host"
                                @change=${(ev: Event) => {
                                    this.host = (ev.target as HTMLInputElement).value;
                                    updateURLParams({ host: this.host });
                                }}
                                .value=${getURLParam("host", "")}
                            />
                            <ak-spinner-button
                                class="pf-m-control"
                                .callAction=${() => {
                                    return new RolesApiApi(DEFAULT_CONFIG)
                                        .toolsPing({
                                            apiAPIToolPingInput: {
                                                host: this.host,
                                            },
                                        })
                                        .then((out) => {
                                            this.pingOutput = out;
                                            this.tracerouteOutput = undefined;
                                            this.portmapOutput = undefined;
                                        })
                                        .catch((exc) => {
                                            this.errorOutput = exc;
                                        });
                                }}
                            >
                                Ping
                            </ak-spinner-button>
                            <ak-spinner-button
                                class="pf-m-control"
                                .callAction=${() => {
                                    return new RolesApiApi(DEFAULT_CONFIG)
                                        .toolsTraceroute({
                                            apiAPIToolTracerouteInput: {
                                                host: this.host,
                                            },
                                        })
                                        .then((out) => {
                                            this.pingOutput = undefined;
                                            this.tracerouteOutput = out;
                                            this.portmapOutput = undefined;
                                        })
                                        .catch((exc) => {
                                            this.errorOutput = exc;
                                        });
                                }}
                            >
                                Traceroute
                            </ak-spinner-button>
                            <ak-spinner-button
                                class="pf-m-control"
                                .callAction=${() => {
                                    return new RolesApiApi(DEFAULT_CONFIG)
                                        .toolsPortmap({
                                            apiAPIToolPortmapInput: {
                                                host: this.host,
                                            },
                                        })
                                        .then((out) => {
                                            this.pingOutput = undefined;
                                            this.tracerouteOutput = undefined;
                                            this.portmapOutput = out;
                                        });
                                }}
                            >
                                Portmap
                            </ak-spinner-button>
                        </div>
                    </div>
                    <div class="pf-v6-c-card__body">${this.renderResult()}</div>
                </div>
            </section>
        `;
    }
}
