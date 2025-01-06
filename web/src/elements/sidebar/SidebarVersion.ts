import { ClusterApi, InstanceAPIClusterInfoOutput } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFAvatar from "@patternfly/patternfly/components/Avatar/avatar.css";
import PFNav from "@patternfly/patternfly/components/Nav/nav.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKElement } from "../Base";

@customElement("ak-sidebar-version")
export class SidebarVersion extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFNav,
            PFAvatar,
            AKElement.GlobalStyle,
            css`
                :host {
                    display: flex;
                    width: 100%;
                    flex-direction: column;
                    justify-content: space-between;
                    padding: 1rem !important;
                }
                p {
                    text-align: center;
                    width: 100%;
                    font-size: var(--pf-global--FontSize--xs);
                }
            `,
        ];
    }

    @state()
    version?: InstanceAPIClusterInfoOutput;

    async firstUpdated() {
        this.version = await new ClusterApi(DEFAULT_CONFIG).clusterGetClusterInfo();
    }

    render(): TemplateResult {
        return html`<p class="pf-c-title">Gravity</p>
            <p class="pf-c-title">Version ${this.version?.clusterVersion}</p> `;
    }
}
