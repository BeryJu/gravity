import { BackupAPIBackupStatusOutput, RolesBackupApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import "../../../elements/forms/ConfirmationForm";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-backup")
export class BackupCard extends AdminStatusCard<BackupAPIBackupStatusOutput> {
    header = "Backup";

    getPrimaryValue(): Promise<BackupAPIBackupStatusOutput> {
        return new RolesBackupApi(DEFAULT_CONFIG).backupStatus();
    }

    getStatus(value: BackupAPIBackupStatusOutput): Promise<AdminStatus> {
        const failed = (value.status || []).filter((v) => v.status !== "success");
        if (failed.length > 0) {
            return Promise.resolve<AdminStatus>({
                icon: "fa fa-exclamation-triangle pf-m-warning",
            });
        }
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }

    renderHeaderLink(): TemplateResult {
        return html` <ak-forms-confirm
            successMessage="Successfully started backup"
            errorMessage="Failed to start backup"
            action="Start"
            .onConfirm=${() => {
                return new RolesBackupApi(DEFAULT_CONFIG).backupStart({
                    wait: false,
                });
            }}
        >
            <span slot="header">Start backup</span>
            <p slot="body">Start a backup using the configured settings.</p>
            <a slot="trigger">
                <i class="fa fa-link"> </i>
            </a>
            <div slot="modal"></div>
        </ak-forms-confirm>`;
    }

    renderValue(): TemplateResult {
        return html`${(this.value?.status || []).map((s) => {
            return html`${s.time?.toLocaleDateString()} ${s.time?.toLocaleTimeString()}`;
        })}`;
    }
}
