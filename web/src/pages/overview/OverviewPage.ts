import { AuthUserMeOutput, RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFContent from "@patternfly/patternfly/components/Content/content.css";
import PFList from "@patternfly/patternfly/components/List/list.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFGrid from "@patternfly/patternfly/layouts/Grid/grid.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKElement } from "../../elements/Base";
import "../../elements/PageHeader";
import "./cards/BuildHashCard";
import "./cards/CurrentInstanceCard";
import "./cards/DHCPScopeCard";
import "./cards/DNSZoneCard";
import "./cards/VersionCard";

@customElement("gravity-overview")
export class OverviewPage extends AKElement {
    @state()
    me?: AuthUserMeOutput;

    static get styles(): CSSResult[] {
        return [PFGrid, PFPage, PFContent, PFList, AKElement.GlobalStyle];
    }

    firstUpdated(): void {
        new RolesApiApi(DEFAULT_CONFIG).apiUsersMe().then((me) => (this.me = me));
    }

    render(): TemplateResult {
        return html` <ak-page-header>
                <span slot="header"> ${this.me ? html`Hello, ${this.me.username}` : html``} </span>
            </ak-page-header>
            <section class="pf-c-page__main-section">
                <div class="pf-l-grid pf-m-gutter">
                    <div class="pf-l-grid__item pf-m-6-col pf-m-3-col-on-2xl">
                        <gravity-overview-card-dhcp-scopes></gravity-overview-card-dhcp-scopes>
                    </div>
                    <div class="pf-l-grid__item pf-m-6-col pf-m-3-col-on-2xl">
                        <gravity-overview-card-dns-zones></gravity-overview-card-dns-zones>
                    </div>
                    <div class="pf-l-grid__item pf-m-6-col pf-m-3-col-on-2xl">
                        <gravity-overview-card-version></gravity-overview-card-version>
                    </div>
                    <div class="pf-l-grid__item pf-m-6-col pf-m-3-col-on-2xl">
                        <gravity-overview-card-build-hash></gravity-overview-card-build-hash>
                    </div>
                    <div class="pf-l-grid__item pf-m-12-col pf-m-4-col-on-2xl">
                        <gravity-overview-card-current-instance></gravity-overview-card-current-instance>
                    </div>
                </div>
            </section>`;
    }
}
