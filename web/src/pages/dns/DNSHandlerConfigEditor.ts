import { CSSResult, TemplateResult, css, html, nothing } from "lit";
import { customElement, property, state } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFDataList from "@patternfly/patternfly/components/DataList/data-list.css";
import PFForm from "@patternfly/patternfly/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly/components/FormControl/form-control.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../elements/Base";

interface HandlerConfig {
    type: string;
    to?: string[];
    cache_ttl?: number;
    net?: string;
    allowlists?: string[];
    blocklists?: string[];
    [key: string]: unknown;
}

const HANDLER_TYPES = [
    { value: "memory", label: "Memory", description: "In-process cache for DNS records" },
    { value: "etcd", label: "etcd", description: "Distributed storage via etcd" },
    { value: "forward_ip", label: "Forward (IP)", description: "Forward queries to an upstream resolver" },
    {
        value: "forward_blocky",
        label: "Forward (Blocky)",
        description: "Forward through Blocky for ad/privacy blocking",
    },
];

function toList(v: string[] | string | undefined): string[] {
    if (!v) return [];
    if (Array.isArray(v)) return v;
    return v.split(";").map((s) => s.trim()).filter(Boolean);
}

function parseTextarea(raw: string): string[] {
    return raw.split("\n").map((s) => s.trim()).filter(Boolean);
}

@customElement("gravity-dns-handler-config-editor")
export class DNSHandlerConfigEditor extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFButton,
            PFDataList,
            PFForm,
            PFFormControl,
            AKElement.GlobalStyle,
            css`
                :host {
                    display: block;
                }
                .handler-list {
                    list-style: none;
                    margin: 0 0 0.75rem 0;
                    padding: 0;
                }
                .handler-item {
                    border: 1px solid var(--pf-global--BorderColor--100);
                    border-radius: var(--pf-global--BorderRadius--sm);
                    margin-bottom: 0.5rem;
                    background: var(--pf-global--BackgroundColor--100);
                    transition: box-shadow 0.15s;
                }
                .handler-item.drag-over {
                    border-color: var(--pf-global--primary-color--100);
                    box-shadow: 0 0 0 2px var(--pf-global--primary-color--100);
                }
                .handler-item.dragging {
                    opacity: 0.45;
                }
                .handler-header {
                    display: flex;
                    align-items: center;
                    padding: 0.5rem 0.75rem;
                    gap: 0.5rem;
                    cursor: default;
                }
                .drag-handle {
                    cursor: grab;
                    color: var(--pf-global--Color--200);
                    font-size: 1.1rem;
                    line-height: 1;
                    flex-shrink: 0;
                    user-select: none;
                    padding: 0 0.25rem;
                }
                .drag-handle:active {
                    cursor: grabbing;
                }
                .handler-type-info {
                    flex: 1;
                    min-width: 0;
                }
                .handler-type-name {
                    font-weight: 600;
                    font-size: var(--pf-global--FontSize--sm);
                    color: var(--pf-global--primary-color--100);
                }
                .handler-type-desc {
                    font-size: var(--pf-global--FontSize--xs);
                    color: var(--pf-global--Color--200);
                    margin-top: 1px;
                }
                .handler-body {
                    padding: 0.75rem 1rem 0.875rem 2.75rem;
                    border-top: 1px solid var(--pf-global--BorderColor--100);
                }
                .handler-field {
                    margin-bottom: 0.75rem;
                }
                .handler-field:last-child {
                    margin-bottom: 0;
                }
                .handler-field-label {
                    display: block;
                    font-size: var(--pf-global--FontSize--sm);
                    font-weight: 600;
                    margin-bottom: 0.3rem;
                    color: var(--pf-global--Color--100);
                }
                .handler-field-hint {
                    font-size: var(--pf-global--FontSize--xs);
                    color: var(--pf-global--Color--200);
                    margin-top: 0.2rem;
                }
                .no-body {
                    padding: 0.5rem 1rem 0.5rem 2.75rem;
                    border-top: 1px solid var(--pf-global--BorderColor--100);
                    font-size: var(--pf-global--FontSize--sm);
                    color: var(--pf-global--Color--200);
                    font-style: italic;
                }
                .empty-state {
                    padding: 1.25rem;
                    text-align: center;
                    color: var(--pf-global--Color--200);
                    border: 2px dashed var(--pf-global--BorderColor--100);
                    border-radius: var(--pf-global--BorderRadius--sm);
                    margin-bottom: 0.75rem;
                    font-size: var(--pf-global--FontSize--sm);
                }
                .add-row {
                    display: flex;
                    align-items: center;
                    gap: 0.5rem;
                }
                .add-row select {
                    flex: 1;
                }
            `,
        ];
    }

    @property({ reflect: true })
    name: string | undefined;

    @state()
    private configs: HandlerConfig[] = [];

    @state()
    private dragIndex: number | null = null;

    @state()
    private dragOverIndex: number | null = null;

    @state()
    private newHandlerType = "forward_ip";

    get value(): HandlerConfig[] {
        return this.configs;
    }

    set value(v: HandlerConfig[] | null | undefined) {
        if (!v) return;
        this.configs = v.map((c) => ({ ...c }));
    }

    private updateConfig(index: number, updates: Partial<HandlerConfig>) {
        const next = [...this.configs];
        next[index] = { ...next[index], ...updates };
        // Remove undefined keys to keep YAML clean
        Object.keys(next[index]).forEach((k) => {
            if (next[index][k] === undefined) delete next[index][k];
        });
        this.configs = next;
        this.dispatchEvent(new Event("change", { bubbles: true }));
    }

    private removeHandler(index: number) {
        this.configs = this.configs.filter((_, i) => i !== index);
        this.dispatchEvent(new Event("change", { bubbles: true }));
    }

    private addHandler() {
        const base: HandlerConfig = { type: this.newHandlerType };
        if (this.newHandlerType.startsWith("forward_")) {
            base.to = ["8.8.8.8:53"];
        }
        this.configs = [...this.configs, base];
        this.dispatchEvent(new Event("change", { bubbles: true }));
    }

    private onDragStart(e: DragEvent, index: number) {
        this.dragIndex = index;
        if (e.dataTransfer) {
            e.dataTransfer.effectAllowed = "move";
            e.dataTransfer.setData("text/plain", String(index));
        }
    }

    private onDragOver(e: DragEvent, index: number) {
        e.preventDefault();
        if (e.dataTransfer) e.dataTransfer.dropEffect = "move";
        if (this.dragOverIndex !== index) this.dragOverIndex = index;
    }

    private onDrop(e: DragEvent, dropIndex: number) {
        e.preventDefault();
        const from = this.dragIndex;
        if (from === null || from === dropIndex) {
            this.dragIndex = null;
            this.dragOverIndex = null;
            return;
        }
        const next = [...this.configs];
        const [moved] = next.splice(from, 1);
        next.splice(dropIndex, 0, moved);
        this.configs = next;
        this.dragIndex = null;
        this.dragOverIndex = null;
        this.dispatchEvent(new Event("change", { bubbles: true }));
    }

    private onDragEnd() {
        this.dragIndex = null;
        this.dragOverIndex = null;
    }

    private handlerMeta(type: string) {
        return HANDLER_TYPES.find((t) => t.value === type) ?? { label: type, description: "" };
    }

    private renderForwardIPBody(config: HandlerConfig, index: number): TemplateResult {
        const toVal = toList(config.to as string[] | string | undefined).join("\n");
        return html`
            <div class="handler-body">
                <div class="handler-field">
                    <label class="handler-field-label">Upstream resolvers</label>
                    <textarea
                        class="pf-c-form-control"
                        rows="3"
                        placeholder="8.8.8.8:53&#10;1.1.1.1:53"
                        .value=${toVal}
                        @change=${(e: Event) =>
                            this.updateConfig(index, {
                                to: parseTextarea((e.target as HTMLTextAreaElement).value),
                            })}
                    ></textarea>
                    <p class="handler-field-hint">One resolver per line (host:port)</p>
                </div>
                <div class="handler-field">
                    <label class="handler-field-label">Cache TTL (seconds)</label>
                    <input
                        type="number"
                        class="pf-c-form-control"
                        placeholder="0 — disabled"
                        .value=${config.cache_ttl !== undefined ? String(config.cache_ttl) : ""}
                        @change=${(e: Event) => {
                            const raw = (e.target as HTMLInputElement).value;
                            const v = raw === "" ? undefined : parseInt(raw, 10);
                            this.updateConfig(index, { cache_ttl: isNaN(v as number) ? undefined : v });
                        }}
                    />
                    <p class="handler-field-hint">0 = disabled · −1 = never cache · −2 = cache forever (import mode)</p>
                </div>
                <div class="handler-field">
                    <label class="handler-field-label">Network protocol</label>
                    <select
                        class="pf-c-form-control"
                        @change=${(e: Event) => {
                            const v = (e.target as HTMLSelectElement).value;
                            this.updateConfig(index, { net: v || undefined });
                        }}
                    >
                        <option value="" ?selected=${!config.net}>Default (UDP)</option>
                        <option value="tcp" ?selected=${config.net === "tcp"}>TCP</option>
                        <option value="udp" ?selected=${config.net === "udp"}>UDP (explicit)</option>
                    </select>
                </div>
            </div>
        `;
    }

    private renderForwardBlockyBody(config: HandlerConfig, index: number): TemplateResult {
        const toVal = toList(config.to as string[] | string | undefined).join("\n");
        const allowVal = toList(config.allowlists as string[] | string | undefined).join("\n");
        const blockVal = toList(config.blocklists as string[] | string | undefined).join("\n");
        return html`
            <div class="handler-body">
                <div class="handler-field">
                    <label class="handler-field-label">Upstream resolvers</label>
                    <textarea
                        class="pf-c-form-control"
                        rows="3"
                        placeholder="8.8.8.8:53&#10;1.1.1.1:53"
                        .value=${toVal}
                        @change=${(e: Event) =>
                            this.updateConfig(index, {
                                to: parseTextarea((e.target as HTMLTextAreaElement).value),
                            })}
                    ></textarea>
                    <p class="handler-field-hint">One resolver per line (host:port)</p>
                </div>
                <div class="handler-field">
                    <label class="handler-field-label">Cache TTL (seconds)</label>
                    <input
                        type="number"
                        class="pf-c-form-control"
                        placeholder="0 — disabled"
                        .value=${config.cache_ttl !== undefined ? String(config.cache_ttl) : ""}
                        @change=${(e: Event) => {
                            const raw = (e.target as HTMLInputElement).value;
                            const v = raw === "" ? undefined : parseInt(raw, 10);
                            this.updateConfig(index, { cache_ttl: isNaN(v as number) ? undefined : v });
                        }}
                    />
                    <p class="handler-field-hint">0 = disabled · −1 = never cache · −2 = cache forever (import mode)</p>
                </div>
                <div class="handler-field">
                    <label class="handler-field-label">Allowlists</label>
                    <textarea
                        class="pf-c-form-control"
                        rows="3"
                        placeholder="https://example.com/allowlist.txt&#10;s.youtube.com"
                        .value=${allowVal}
                        @change=${(e: Event) => {
                            const lines = parseTextarea((e.target as HTMLTextAreaElement).value);
                            this.updateConfig(index, { allowlists: lines.length ? lines : undefined });
                        }}
                    ></textarea>
                    <p class="handler-field-hint">URLs or domain names (one per line). Domains in allowlists bypass blocking.</p>
                </div>
                <div class="handler-field">
                    <label class="handler-field-label">Blocklists</label>
                    <textarea
                        class="pf-c-form-control"
                        rows="3"
                        placeholder="https://example.com/blocklist.txt"
                        .value=${blockVal}
                        @change=${(e: Event) => {
                            const lines = parseTextarea((e.target as HTMLTextAreaElement).value);
                            this.updateConfig(index, { blocklists: lines.length ? lines : undefined });
                        }}
                    ></textarea>
                    <p class="handler-field-hint">URLs to blocklist files (one per line). Leave empty to use Blocky defaults.</p>
                </div>
            </div>
        `;
    }

    private renderHandlerBody(config: HandlerConfig, index: number): TemplateResult | typeof nothing {
        switch (config.type) {
            case "memory":
            case "etcd":
                return html`<div class="no-body">No additional configuration required.</div>`;
            case "forward_ip":
                return this.renderForwardIPBody(config, index);
            case "forward_blocky":
                return this.renderForwardBlockyBody(config, index);
            default:
                return nothing;
        }
    }

    render(): TemplateResult {
        return html`
            ${this.configs.length === 0
                ? html`<div class="empty-state">No handlers configured. Add one below.</div>`
                : html`
                      <ul class="handler-list">
                          ${this.configs.map((config, index) => {
                              const meta = this.handlerMeta(config.type);
                              const isDragging = this.dragIndex === index;
                              const isDragOver = this.dragOverIndex === index && !isDragging;
                              return html`
                                  <li
                                      class="handler-item${isDragging
                                          ? " dragging"
                                          : ""}${isDragOver ? " drag-over" : ""}"
                                      draggable="true"
                                      @dragstart=${(e: DragEvent) => this.onDragStart(e, index)}
                                      @dragover=${(e: DragEvent) => this.onDragOver(e, index)}
                                      @dragleave=${() => {
                                          if (this.dragOverIndex === index)
                                              this.dragOverIndex = null;
                                      }}
                                      @drop=${(e: DragEvent) => this.onDrop(e, index)}
                                      @dragend=${() => this.onDragEnd()}
                                  >
                                      <div class="handler-header">
                                          <span class="drag-handle" title="Drag to reorder">⠿</span>
                                          <div class="handler-type-info">
                                              <div class="handler-type-name">${meta.label}</div>
                                              ${meta.description
                                                  ? html`<div class="handler-type-desc">${meta.description}</div>`
                                                  : nothing}
                                          </div>
                                          <button
                                              class="pf-c-button pf-m-plain"
                                              type="button"
                                              title="Remove handler"
                                              @click=${() => this.removeHandler(index)}
                                              aria-label="Remove ${meta.label} handler"
                                          >
                                              <i class="fas fa-times" aria-hidden="true"></i>
                                          </button>
                                      </div>
                                      ${this.renderHandlerBody(config, index)}
                                  </li>
                              `;
                          })}
                      </ul>
                  `}
            <div class="add-row">
                <select
                    class="pf-c-form-control"
                    @change=${(e: Event) => {
                        this.newHandlerType = (e.target as HTMLSelectElement).value;
                    }}
                >
                    ${HANDLER_TYPES.map(
                        (t) =>
                            html`<option
                                value=${t.value}
                                ?selected=${t.value === this.newHandlerType}
                            >
                                ${t.label} — ${t.description}
                            </option>`,
                    )}
                </select>
                <button
                    class="pf-c-button pf-m-primary"
                    type="button"
                    @click=${() => this.addHandler()}
                >
                    <i class="fas fa-plus" aria-hidden="true"></i>
                    Add handler
                </button>
            </div>
        `;
    }
}
