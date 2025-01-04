import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFList from "@patternfly/patternfly/components/List/list.css";

import { EVENT_REFRESH } from "../../common/constants";
import { MessageLevel } from "../../common/messages";
import { ModalButton } from "../buttons/ModalButton";
import "../buttons/SpinnerButton";
import { showMessage } from "../messages/MessageContainer";
import { PaginatedResponse } from "../table/Table";
import { Table, TableColumn } from "../table/Table";

type BulkDeleteMetadata = { key: string; value: string }[];

@customElement("ak-delete-objects-table")
export class DeleteObjectsTable<T extends object> extends Table<T> {
    @property({ attribute: false })
    objects: T[] = [];

    @property({ attribute: false })
    metadata!: (item: T) => BulkDeleteMetadata;

    static get styles(): CSSResult[] {
        return super.styles.concat(PFList);
    }

    constructor() {
        super();
        this.paginated = false;
    }

    async apiEndpoint(): Promise<PaginatedResponse<T>> {
        return Promise.resolve({
            pagination: {
                count: this.objects.length,
                current: 1,
                totalPages: 1,
                startIndex: 1,
                endIndex: this.objects.length,
            },
            results: this.objects,
        });
    }

    columns(): TableColumn[] {
        return this.metadata(this.objects[0]).map((element) => {
            return new TableColumn(element.key);
        });
    }

    row(item: T): TemplateResult[] {
        return this.metadata(item).map((element) => {
            return html`${element.value}`;
        });
    }

    renderToolbarContainer(): TemplateResult {
        return html``;
    }
}

@customElement("ak-forms-delete-bulk")
export class DeleteBulkForm extends ModalButton {
    @property({ attribute: false })
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    objects: any[] = [];

    @property()
    objectLabel: string | undefined;

    @property({ attribute: false })
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    metadata: (item: any) => BulkDeleteMetadata = (item: any) => {
        const rec = item as Record<string, unknown>;
        const meta = [];
        if (Object.prototype.hasOwnProperty.call(rec, "name")) {
            meta.push({ key: "Name", value: rec.name as string });
        }
        if (Object.prototype.hasOwnProperty.call(rec, "pk")) {
            meta.push({ key: "ID", value: rec.pk as string });
        }
        return meta;
    };

    @property({ attribute: false })
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    delete!: (item: any, extraData: any) => Promise<any>;

    @property({ attribute: false })
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    preDelete: () => Promise<any> = () => {
        return Promise.resolve();
    };

    async confirm(): Promise<void> {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        let extraData: any = undefined;
        try {
            extraData = await this.preDelete();
        } catch (exc) {
            showMessage({
                message: (exc as Error).toString(),
                level: MessageLevel.error,
            });
            return Promise.reject();
        }
        return Promise.all(
            this.objects.map((item) => {
                return this.delete(item, extraData);
            }),
        )
            .then(() => {
                this.onSuccess();
                this.open = false;
                this.dispatchEvent(
                    new CustomEvent(EVENT_REFRESH, {
                        bubbles: true,
                        composed: true,
                    }),
                );
            })
            .catch((e) => {
                this.onError(e);
                throw e;
            });
    }

    onSuccess(): void {
        showMessage({
            message: `Successfully deleted ${this.objects.length} ${this.objectLabel}`,
            level: MessageLevel.success,
        });
    }

    onError(e: Error): void {
        showMessage({
            message: `Failed to delete ${this.objectLabel}: ${e.toString()}`,
            level: MessageLevel.error,
        });
    }

    renderModalInner(): TemplateResult {
        return html`<section
                class="pf-v6-c-modal-box__header pf-m-light"
            >
                <div class="pf-v6-c-content">
                    <h1 class="pf-v6-c-title pf-m-2xl">${`Delete ${this.objectLabel}`}</h1>
                </div>
            </section>
            <section class="pf-v6-c-modal-box__body pf-v6-c-page__main-section pf-m-light">
                <form class="pf-v6-c-form pf-m-horizontal">
                    <p class="pf-v6-c-title">
                        ${`Are you sure you want to delete ${this.objects.length} ${this.objectLabel}?`}
                    </p>
                    <slot name="notice"></slot>
                </form>
            </section>
            <section class="pf-v6-c-modal-box__body pf-v6-c-page__main-section pf-m-light">
                <ak-delete-objects-table .objects=${this.objects} .metadata=${this.metadata}>
                </ak-delete-objects-table>
            </section>
            <footer class="pf-v6-c-modal-box__footer">
                <ak-spinner-button
                    .callAction=${() => {
                        return this.confirm();
                    }}
                    class="pf-m-danger"
                >
                    ${"Delete"} </ak-spinner-button
                >&nbsp;
                <ak-spinner-button
                    .callAction=${async () => {
                        this.open = false;
                    }}
                    class="pf-m-secondary"
                >
                    ${"Cancel"}
                </ak-spinner-button>
            </footer>`;
    }
}
