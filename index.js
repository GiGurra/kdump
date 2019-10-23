#!/usr/bin/env node

const execSync = require('child_process').execSync;
const fs = require("fs");

const tableParser = require('table-parser');

async function main() {

    console.log("Running kube-dump script");

    const contexts = getContexts();
    console.log("contexts: " + contexts);

    for (const context of contexts) {
        console.log(" - processing context '" + context + "'");

        fs.mkdirSync(context, {recursive: true});

        execSync("kubectl config use-context " + context);

        const allResources =
            tableParser
                .parse(execSync("kubectl api-resources -o wide").toString())
                .map(toResource)
                .filter(allowedResource);

        const namespacedResources = allResources.filter(r => r.namespaced);
        const globalResources = allResources.filter(r => !r.namespaced);

        const namespaces = getNamespaces();

        for (const namespace of namespaces) {

            console.log("   - processing namespace resources for namespace: " + namespace);

            fs.mkdirSync(context + "/" + namespace, {recursive: true});

            execSync("kubectl config set-context --current --namespace=" + namespace);

            for (const namespacedResource of namespacedResources) {

                const resourceYaml = execSync("kubectl get " + namespacedResource.name + " -o yaml", {maxBuffer: 100*1024*1024}).toString();
                fs.writeFileSync(context + "/" + namespace + "/" + namespacedResource.name + ".yml", resourceYaml)

            }

        }

        console.log("   - processing global resources for context");

        for (const globalResource of globalResources) {

            const resourceYaml = execSync("kubectl get " + globalResource.name + " -o yaml").toString();
            fs.writeFileSync(context + "/" + globalResource.name + ".yml", resourceYaml)
        }

    }

    console.log("kube-dump script finished!")
}

function allowedResource(resource) {
    return resource.verbs.indexOf('get') >= 0;
}

function toResource(tableResource) {
    const out = {};
    out.name = tableResource.NAME.toString();
    out.shortNames = tableResource.SHORTNAMES.toString().split(',').map(x => x.trim()).filter(x => x.length > 0);
    out.apiGroup = tableResource.APIGROUP.toString();
    out.namespaced = tableResource.NAMESPACED.toString() === "true";
    out.kind = tableResource.NAMESPACED.toString();
    out.verbs = tableResource.VERBS.toString().slice(1, tableResource.VERBS.toString().length-1).split(',');
    return out
}

function getContexts() {
    return execSync("kubectl config get-contexts -o name")
        .toString()
        .trim()
        .split(/\r?\n/)
        .map(x => x.trim())
        .filter(x => x.length > 0);
}

function getNamespaces() {
    return execSync("kubectl get namespaces -o name")
        .toString()
        .trim()
        .split(/\r?\n/)
        .map(x => x.trim())
        .filter(x => x.length > 0)
        .map(x => x.split("/").slice(-1)[0]);
}

main().catch(error => {
    console.error(error);
    process.exit(1)
});
