import { AuthUserMeOutput, RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFContent from "@patternfly/patternfly/components/Content/content.css";
import PFList from "@patternfly/patternfly/components/List/list.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFGrid from "@patternfly/patternfly/layouts/Grid/grid.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKElement } from "../../elements/Base";
import "../../elements/PageHeader";
import "./cards/DHCPScopeCard";
import "./cards/DNSZoneCard";

@customElement("gravity-overview")
export class OverviewPage extends AKElement {
    @property()
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
                    <!-- row 1 -->
                    <div
                        class="pf-l-grid__item pf-m-6-col pf-m-4-col-on-xl pf-m-2-col-on-2xl graph-container"
                    >
                        <gravity-overview-card-dhcp-scopes></gravity-overview-card-dhcp-scopes>
                    </div>
                    <div
                        class="pf-l-grid__item pf-m-6-col pf-m-4-col-on-xl pf-m-2-col-on-2xl graph-container"
                    >
                        <gravity-overview-card-dns-zones></gravity-overview-card-dns-zones>
                    </div>
                </div>
            </section>`;
    }
}
