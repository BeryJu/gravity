import YAML from "yaml";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import "../../../elements/CodeMirror";
import "../../../elements/CodeMirror";
import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-cluster-join-compose")
export class ClusterJoinCompose extends WizardFormPage {
    sidebarLabel = () => "Deployment";

    activeCallback = async () => {
        this.host.canBack = false;
        this.host.isValid = true;
    };

    nextCallback = async (): Promise<boolean> => {
        return true;
    };

    renderCompose(): KeyUnknown {
        return {
            services: {
                gravity: {
                    hostname: this.host.state["identifier"],
                    image: "ghcr.io/beryju/gravity:stable",
                    restart: "unless-stopped",
                    network_mode: "host",
                    user: "root",
                    environment: {
                        ETCD_JOIN_CLUSTER: `${this.host.state["join_token"]},http://${this.host.state["node_ip"]}:8008`,
                        BOOTSTRAP_ROLES: this.host.state["roles"],
                    },
                    volumes: ["data:/data"],
                    logging: {
                        driver: "json-file",
                        options: {
                            "max-size": "10m",
                            "max-file": "3",
                        },
                    },
                },
            },
            volumes: {
                data: {
                    driver: "local",
                },
            },
        };
    }

    renderForm(): TemplateResult {
        return html`
            <p>Use the compose file below to deploy the new gravity node.</p>
            <ak-codemirror mode="yaml" value=${YAML.stringify(this.renderCompose())}>
            </ak-codemirror>
        `;
    }
}
