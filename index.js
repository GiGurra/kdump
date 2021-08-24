#!/usr/bin/env node

const execSync = require('child_process').execSync;
const fs = require("fs");
const yaml = require('js-yaml');
const yargs = require("yargs");
const crypto = require('crypto');

const tableParser = require('table-parser');

async function main() {

    const cmdLine = parseCmdLine();

    // Make this configurable
    const ignoredTypes = [
        "events",
        "jobs",
        "pods",
        "componentstatuses",
        "endpoints",
        "replicasets",
        "clusterauthtokens",
        "clusteruserattributes",
        "controllerrevisions",
        "apiservices",
        "clusterinformations",
        //"customresourcedefinitions",
        "felixconfigurations",
        "ippools",
        "nodes",
        "priorityclasses",
        "leases"
    ];

    console.log("Running kube-dump script");

    const contexts = (cmdLine.c || getContexts()).filter( c => {
        return !cmdLine.ec || cmdLine.ec.indexOf(c) < 0
    });

    const rootDir = cmdLine.o ? cmdLine.o + "/" : "";

    console.log("contexts: " + contexts);

    for (const context of contexts) {

        console.log(" - processing context '" + context + "'");

        if (fs.existsSync(rootDir + context)) {
            throw new Error("Output directory '" + context + "' already exists!");
        }
        fs.mkdirSync(rootDir + context, {recursive: true});

        execSync("kubectl config use-context " + context);

        const tablesStr = execSync("kubectl api-resources -o wide").toString();

        const allResources =
            tableParser
                .parse(tablesStr)
                .map(toResource)
                .filter(isReadableResource)
                .filter(x => !isIgnoredResource(x, ignoredTypes))
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

            fs.mkdirSync(rootDir + context + "/" + namespace, {recursive: true});

            execSync("kubectl config set-context --current --namespace=" + namespace);

            for (const namespacedResource of namespacedResources) {

                const humanReadableTable = getItems(namespacedResource.name);

                if (humanReadableTable.length === 0 && !cmdLine['include-empty-resources']) {
                    continue;
                }

                const resourceYaml = cleanUpKubernetesItemsYaml(execSync("kubectl get " + namespacedResource.name + " -o yaml", {maxBuffer: 100*1024*1024}).toString());

                if (namespacedResource.name === "secrets" && !cmdLine['include-secrets']) {
                    console.log("WARNING: Skipping resource 'secrets' as include-secrets was not set");
                    continue;
                }

                if (namespacedResource.name === "secrets" && cmdLine['encrypt-secrets']) {

                    if (!cmdLine['encrypt-password']) {
                        throw new Error("Failed encrypting secrets: cmd line flag --encrypt-secrets set, but no --encrypt-password/-p value provided");
                    }

                    if (!cmdLine['encrypt-algorithm']) {
                        throw new Error("Failed encrypting secrets: cmd line flag --encrypt-algorithm empty");
                    }

                    const algorithm = cmdLine['encrypt-algorithm'];
                    const key = Buffer.from(cmdLine['encrypt-password'], 'hex');

                    const prevDumpDir = cmdLine['prev-dump-dir'];

                    if (prevDumpDir) {
                        if (fs.existsSync(prevDumpDir)) {
                            if (fs.existsSync(prevDumpDir + "/" + context + "/" + namespace)) {
                                const filesInOldDir = fs.readdirSync(prevDumpDir + "/" + context + "/" + namespace);
                                const oldFileIndex = filesInOldDir.findIndex(f => f.startsWith(namespacedResource.name + ".iv=") && f.endsWith(".yml"));
                                if (oldFileIndex >= 0) {

                                    const oldFileKey = key;
                                    const oldFileAlgorithm = algorithm;

                                    const oldFileName = filesInOldDir[oldFileIndex];
                                    const oldFileIv = oldFileName.split('.iv=')[1].split('.yml')[0];

                                    const oldFileContents = fs.readFileSync(prevDumpDir + "/" + context + "/" + namespace + "/" + oldFileName);
                                    const decryptedOldFile = decrypt(oldFileContents, oldFileKey, Buffer.from(oldFileIv, 'hex'), oldFileAlgorithm).toString();

                                    if (decryptedOldFile === resourceYaml) {
                                        const fromFile = prevDumpDir + "/" + context + "/" + namespace + "/" + oldFileName;
                                        const toFile = rootDir + context + "/" + namespace + "/" + oldFileName;
                                        fs.copyFileSync(fromFile, toFile);
                                        continue;
                                    }
                                }
                            }
                        }
                        else {
                            console.error("WARNING: --prev-dump-dir given, but specified directory does not exist!");
                        }
                    }

                    const iv = crypto.randomBytes(16);
                    const encryptedData = encrypt(resourceYaml, key, iv, algorithm);

                    fs.writeFileSync(rootDir + context + "/" + namespace + "/" + namespacedResource.name + ".iv=" + iv.toString("hex") + ".yml", encryptedData);
                }
                else {
                    fs.writeFileSync(rootDir + context + "/" + namespace + "/" + namespacedResource.name + ".yml", resourceYaml);
                }

            }

        }

        if (cmdLine.eg) {
            console.log("   - NOT processing global resources for context, since --eg flag was specified");
        }
        else {

            console.log("   - processing global resources for context");

            for (const globalResource of globalResources) {

                const humanReadableTable = getItems(globalResource.name);

                if (humanReadableTable.length === 0 && !cmdLine['include-empty-resources']) {
                    continue;
                }

                const resourceYaml = cleanUpKubernetesItemsYaml(execSync("kubectl get " + globalResource.name + " -o yaml", {maxBuffer: 100*1024*1024}).toString());
                fs.writeFileSync(rootDir + context + "/" + globalResource.name + ".yml", resourceYaml)
            }

        }
    }

    console.log("kube-dump script finished!")
}

function cleanUpKubernetesItemsYaml(itemsString) {

    // TODO: Make this configurable

    const data = yaml.load(itemsString);

    for (const item of data.items) {

        // remove things we simply dont want
        delete item['lastRefresh'];
        delete item['status'];

        // clean up metadata
        if (item['metadata']) {
            delete item['metadata']['generation'];
            delete item['metadata']['resourceVersion'];
            if (item['metadata']['annotations']) {
                delete item['metadata']['annotations']['control-plane.alpha.kubernetes.io/leader'];
                delete item['metadata']['annotations']['deployment.kubernetes.io/revision'];
                delete item['metadata']['annotations']['cattle.io/timestamp'];
            }
        }

        // clean up spec
        if (item['spec']) {

            delete item['spec']['renewTime'];

            if (item['spec']['template']) {
                if (item['spec']['template']['metadata']) {
                    if (item['spec']['template']['metadata']['annotations']) {
                        delete item['spec']['template']['metadata']['annotations']['deployment.kubernetes.io/revision'];
                        delete item['spec']['template']['metadata']['annotations']['cattle.io/timestamp'];
                    }
                }
            }

            if (item['spec']['jobTemplate']) {
                if (item['spec']['jobTemplate']['spec']) {
                    if (item['spec']['jobTemplate']['spec']['template']) {
                        if (item['spec']['jobTemplate']['spec']['template']['metadata']) {
                            if (item['spec']['jobTemplate']['spec']['template']['metadata']['annotations']) {
                                delete item['spec']['jobTemplate']['spec']['template']['metadata']['annotations']['deployment.kubernetes.io/revision'];
                                delete item['spec']['jobTemplate']['spec']['template']['metadata']['annotations']['cattle.io/timestamp'];
                            }
                        }
                    }
                }
            }
        }
    }

    return yaml.dump(data);
}

function decrypt(encryptedData, key, iv, algorithm) {
    const decipher = crypto.createDecipheriv(algorithm, Buffer.from(key), iv);
    const decrypted = decipher.update(encryptedData);
    return Buffer.concat([decrypted, decipher.final()]);
}

function encrypt(text, key, iv, algorithm) {
    const cipher = crypto.createCipheriv(algorithm, Buffer.from(key), iv);
    const encrypted = cipher.update(text);
    return Buffer.concat([encrypted, cipher.final()]);
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
        .option('include-secrets', {
            description: 'If to include secrets, default false. you will need to include encrypt-password or set encrypt-secrets false',
            type: 'boolean',
            default: false
        })
        .option('encrypt-secrets', {
            description: 'If to encrypt the secrets resource. Default and recommended. ' +
                'To decrypt: \n' +
                'openssl enc -d -aes-256-cbc -iv hexIV -K hexKey',
            type: 'boolean',
            default: true
        })
        .option('encrypt-password', {
            alias: 'p',
            description: 'Password for aes-256-cbc encryption of secrets resource. This must be 32 bytes hex (64 characters).' +
                'You can generate one using: \n' +
                'openssl rand -hex 32',
            type: 'string'
        })
        .option('encrypt-algorithm', {
            description: 'Encryption algorithm to use for secrets',
            type: 'string',
            default: 'aes-256-cbc'
        })
        .option('prev-dump-dir', {
            description: 'Directory with contents of previous dump. Useful to compare encrypted secrets to only replace file if something actually changed. (otherwise you will get a git diff every time because encryption IV changes)',
            type: 'string'
        })
        .option('output-dir', {
            alias: 'o',
            description: 'Output directory',
            type: 'string'
        })
        .option('include-empty-resources', {
            description: 'If to write yaml files for resources with no entries',
            type: 'boolean',
            default: false
        })
        .help()
        .strict()
        .argv;
}

function isReadableResource(resource) {
    return resource.verbs.indexOf('get') >= 0;
}

function isIgnoredResource(resource, ignoreList) {
    return ignoreList.indexOf(resource.name) >= 0;
}

function toResource(tableResource) {
    const out = {};
    out.name = tableResource.NAME.toString();
    out.shortNames = tableResource.SHORTNAMES.toString().split(',').map(x => x.trim()).filter(x => x.length > 0);
    //out.apiGroup = tableResource.APIGROUP.toString();
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

function getItems(resourceName) {
    return execSync("kubectl get " + resourceName + " -o name", {maxBuffer: 100*1024*1024})
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
