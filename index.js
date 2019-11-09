#!/usr/bin/env node

const execSync = require('child_process').execSync;
const fs = require("fs");
const yargs = require("yargs");

const tableParser = require('table-parser');

async function main() {

    const cmdLine = parseCmdLine();

    console.log("Running kube-dump script");

    const contexts = (cmdLine.c || getContexts()).filter( c => {
        return !cmdLine.ec || cmdLine.ec.indexOf(c) < 0
    });

    console.log("contexts: " + contexts);

    for (const context of contexts) {
        console.log(" - processing context '" + context + "'");

        if (fs.existsSync(context)) {
            throw new Error("Output directory '" + context + "' already exists!");
        }
        fs.mkdirSync(context, {recursive: true});

        execSync("kubectl config use-context " + context);

        const allResources =
            tableParser
                .parse(execSync("kubectl api-resources -o wide").toString())
                .map(toResource)
                .filter(isReadableResource)
                .filter(r => !cmdLine.er || cmdLine.er.indexOf(r.name) < 0);

        const namespacedResources =
            allResources
                .filter(r => r.namespaced)
                .filter(r => !cmdLine.nr || cmdLine.nr.indexOf(r.name) >= 0);

        const globalResources =
            allResources
                .filter(r => !r.namespaced)
                .filter(r => !cmdLine.gr || cmdLine.gr.indexOf(r.name) >= 0);

        const namespaces = (cmdLine.n || getNamespaces()).filter(n => {
            return !cmdLine.en || cmdLine.en.indexOf(n) < 0
        });

        for (const namespace of namespaces) {

            console.log("   - processing namespace resources for namespace: " + namespace);

            fs.mkdirSync(context + "/" + namespace, {recursive: true});

            execSync("kubectl config set-context --current --namespace=" + namespace);

            for (const namespacedResource of namespacedResources) {

                const resourceYaml = execSync("kubectl get " + namespacedResource.name + " -o yaml", {maxBuffer: 100*1024*1024}).toString();
                fs.writeFileSync(context + "/" + namespace + "/" + namespacedResource.name + ".yml", resourceYaml)

            }

        }

        if (cmdLine.eg) {
            console.log("   - NOT processing global resources for context, since --eg flag was specified");
        }
        else {

            console.log("   - processing global resources for context");

            for (const globalResource of globalResources) {

                const resourceYaml = execSync("kubectl get " + globalResource.name + " -o yaml").toString();
                fs.writeFileSync(context + "/" + globalResource.name + ".yml", resourceYaml)
            }

        }
    }

    console.log("kube-dump script finished!")
}

function parseCmdLine() {
    return yargs(process.argv.slice(2))
        .option('context', {
            alias: 'c',
            description: 'Specify contexts. If omitted - use all available',
            type: 'array'
        })
        .option('exclude-context', {
            alias: 'ec',
            description: 'Exclude contexts',
            type: 'array'
        })
        .option('namespace', {
            alias: 'n',
            description: 'Specify namespaces. If omitted - use all available',
            type: 'array',
        })
        .option('exclude-namespace', {
            alias: 'en',
            description: 'Exclude namespaces',
            type: 'array',
        })
        .option('exclude-global', {
            alias: 'eg',
            description: 'Exclude global (non-namespaced) data',
            type: 'boolean',
            default: false
        })
        .option('namespaced-resource', {
            alias: 'nr',
            description: 'Specify namespaced resources. If omitted - use all available',
            type: 'array'
        })
        .option('global-resource', {
            alias: 'gr',
            description: 'Specify global resources. If omitted - use all available',
            type: 'array'
        })
        .option('exclude-resource', {
            alias: 'er',
            description: 'Exclude resource',
            type: 'array'
        })
        .help()
        .strict()
        .argv;
}

function isReadableResource(resource) {
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
