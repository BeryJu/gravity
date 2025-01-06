import { AuthAPIMeOutput, RolesApiApi } from "gravity-api";

import { css, html, nothing } from "lit";
import { customElement, property, state } from "lit/decorators.js";

import PFAvatar from "@patternfly/patternfly/components/Avatar/avatar.css";
import PFBrand from "@patternfly/patternfly/components/Brand/brand.css";
import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFDrawer from "@patternfly/patternfly/components/Drawer/drawer.css";
import PFDropdown from "@patternfly/patternfly/components/Dropdown/dropdown.css";
import PFNotificationBadge from "@patternfly/patternfly/components/NotificationBadge/notification-badge.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";
import PFDisplay from "@patternfly/patternfly/utilities/Display/display.css";

import { DEFAULT_CONFIG } from "../api/Config";
import { EVENT_API_DRAWER_TOGGLE } from "../common/constants";
import { AKElement } from "./Base";

@customElement("ak-nav-buttons")
export class NavigationButtons extends AKElement {
    @property({ type: Boolean, reflect: true })
    notificationDrawerOpen = false;

    @property({ type: Boolean, reflect: true })
    apiDrawerOpen = false;

    @property({ type: Number })
    notificationsCount = 0;

    @state()
    me?: AuthAPIMeOutput;

    static get styles() {
        return [
            PFBase,
            PFDisplay,
            PFBrand,
            PFPage,
            PFAvatar,
            PFButton,
            PFDrawer,
            PFDropdown,
            PFNotificationBadge,
            AKElement.GlobalStyle,
            css`
                .pf-c-page__header-tools {
                    display: flex;
                }
            `,
        ];
    }

    async firstUpdated() {
        this.me = await new RolesApiApi(DEFAULT_CONFIG).apiUsersMe();
    }

    renderApiDrawerTrigger() {
        return nothing;
        const onClick = (ev: Event) => {
            ev.stopPropagation();
            this.dispatchEvent(
                new Event(EVENT_API_DRAWER_TOGGLE, { bubbles: true, composed: true }),
            );
        };
        return html`<div class="pf-c-page__header-tools-item pf-m-hidden pf-m-visible-on-lg">
            <button class="pf-c-button pf-m-plain" type="button" @click=${onClick}>
                <pf-tooltip position="top" content="Open API drawer">
                    <i class="fas fa-code" aria-hidden="true"></i>
                </pf-tooltip>
            </button>
        </div>`;
    }

    get userDisplayName() {
        return this.me?.username;
    }

    render() {
        return html`<div class="pf-c-page__header-tools">
            <div class="pf-c-page__header-tools-group">
                ${this.renderApiDrawerTrigger()}
                <div class="pf-c-page__header-tools-item">
                    <a href="/auth/logout" class="pf-c-button pf-m-plain">
                        <pf-tooltip position="top" content="Sign out">
                            <i class="fas fa-sign-out-alt" aria-hidden="true"></i>
                        </pf-tooltip>
                    </a>
                </div>
                <slot name="extra"></slot>
            </div>
            ${this.userDisplayName != ""
                ? html`<div class="pf-c-page__header-tools-group">
                      <div class="pf-c-page__header-tools-item pf-m-hidden pf-m-visible-on-md">
                          ${this.userDisplayName}
                      </div>
                  </div>`
                : nothing}
        </div>`;
    }
}
