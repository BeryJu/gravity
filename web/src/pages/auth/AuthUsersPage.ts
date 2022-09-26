import { AuthUser, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./AuthUserForm";

@customElement("gravity-auth-users")
export class AuthUsersPage extends TablePage<AuthUser> {
    pageTitle(): string {
        return "Users";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    apiEndpoint(page: number): Promise<PaginatedResponse<AuthUser>> {
        return new RolesApiApi(DEFAULT_CONFIG).apiGetUsers().then((users) => {
            const data = (users.users || []).filter((l) =>
                l.username.toLowerCase().includes(this.search.toLowerCase()),
            );
            data.sort((a, b) => {
                if (a.username > b.username) return 1;
                if (a.username < b.username) return -1;
                return 0;
            });
            return PaginationWrapper(data);
        });
    }

    columns(): TableColumn[] {
        return [new TableColumn("Username"), new TableColumn("Actions")];
    }

    row(item: AuthUser): TemplateResult<1 | 2>[] {
        return [
            html`${item.username}`,
            html`<ak-forms-modal>
                <span slot="submit"> ${"Update"} </span>
                <span slot="header"> ${"Update user"} </span>
                <gravity-auth-user-form slot="form" .instancePk=${item.username}>
                </gravity-auth-user-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"User(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: AuthUser) => {
                return [{ key: "Username", value: item.username }];
            }}
            .delete=${(item: AuthUser) => {
                return new RolesApiApi(DEFAULT_CONFIG).apiDeleteUsers({
                    username: item.username,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${"Create"} </span>
                <span slot="header"> ${"Create User"} </span>
                <gravity-auth-user-form slot="form"> </gravity-auth-user-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${"Create"}</button>
            </ak-forms-modal>
        `;
    }
}
