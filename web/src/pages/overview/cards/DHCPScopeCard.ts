import { RolesDhcpApi } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-dhcp-scopes")
export class DHCPScopeCard extends AdminStatusCard<number> {
    header = "DHCP Scopes";
    headerLink = "#/dhcp/scopes";

    getPrimaryValue(): Promise<number> {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes().then((scopes) => {
            return scopes.scopes?.length || 0;
        });
    }

    getStatus(value: number): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }
}
