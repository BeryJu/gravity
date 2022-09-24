import { RolesDnsApi } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-dns-zones")
export class DNSScopeCard extends AdminStatusCard<number> {
    header = "DNS Zones";
    headerLink = "#/dns/zones";

    getPrimaryValue(): Promise<number> {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsGetZones().then((zones) => {
            return zones.zones?.length || 0;
        });
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    getStatus(value: number): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }
}
