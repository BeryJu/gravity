import { BackupAPIBackupStatusOutput, RolesBackupApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-backup")
export class BackupCard extends AdminStatusCard<BackupAPIBackupStatusOutput> {
    header = "Backup";

    getPrimaryValue(): Promise<BackupAPIBackupStatusOutput> {
        return new RolesBackupApi(DEFAULT_CONFIG).backupStatus();
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    getStatus(value: BackupAPIBackupStatusOutput): Promise<AdminStatus> {
        const isSuccess = (value.status || []).filter((v) => v.status !== "success");
        if (!isSuccess) {
            return Promise.resolve<AdminStatus>({
                icon: "fa fa-exclamation-triangle pf-m-warning"
            });
        }
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }

    renderValue(): TemplateResult {
        return html`${(this.value?.status || []).map(s => {
            return html`${s.time?.toLocaleDateString()} ${s.time?.toLocaleTimeString()}`;
        })}`;
    }
}
