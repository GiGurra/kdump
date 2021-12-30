use crate::{ApiResourceType, util};

use clap::Parser;

/// Dump all kubernetes resources as yaml files to a dir
#[derive(Parser, Debug, PartialEq, Clone)]
#[clap(about, version, author)]
pub struct CliArgs {
    /// output directory to create
    #[clap(short, long)]
    pub output_dir: String,

    /// if to delete previous output directory (default: false)
    #[clap(long)]
    pub delete_previous_dir: bool,

    /// symmetric secrets encryption hex key for aes GCM (lower case 64 chars)
    #[clap(long)]
    pub secrets_encryption_key: Option<String>,
}

#[derive(Debug, PartialEq, Clone)]
pub struct AppConfig {
    pub output_dir: String,
    pub delete_prev_dir: bool,
    pub excluded_types: Vec<String>,
    pub secrets_encryption_key: Option<String>,
}

impl Default for AppConfig {
    fn default() -> Self {
        return AppConfig {
            output_dir: String::from("test"),  // TODO: Change to default empty when implementing cli args
            delete_prev_dir: true, // TODO: Change to default false when implementing cli args
            excluded_types: default_resources_excluded(),
            secrets_encryption_key: None,
        };
    }
}

impl AppConfig {
    pub fn from_cli_args() -> AppConfig {
        let mut result = AppConfig::default();
        let cli_args: CliArgs = CliArgs::parse();
        result.output_dir = cli_args.output_dir;
        result.delete_prev_dir = cli_args.delete_previous_dir;
        result.secrets_encryption_key = cli_args.secrets_encryption_key;
        return result;
    }

    pub fn include_secrets(&self) -> bool {
        return self.secrets_encryption_key.is_some();
    }

    pub fn is_type_included(&self, tpe: &util::k8s::ApiResourceType) -> bool {
        return !self.excluded_types.contains(&tpe.name) &&
            !self.excluded_types.contains(&tpe.qualified_name()) &&
            (!tpe.is_secret() || self.include_secrets());
    }

    pub fn types_do_download<'a>(&self, all_resource_type_defs: &'a util::k8s::kubectl::ApiResourceTypes) -> Vec<&'a ApiResourceType> {
        return all_resource_type_defs.accessible.all
            .iter()
            .filter(|x| self.is_type_included(x))
            .collect::<Vec<&ApiResourceType>>();
    }
}


pub fn default_resources_excluded() -> Vec<String> {
    return vec![
        "limitranges",
        "podtemplates",
        "replicationcontrollers",
        "resourcequotas",
        "events",
        "jobs",
        "jobs.batch",
        "pods",
        "componentstatuses",
        "endpoints",
        "endpointslices.discovery.k8s.io",
        "replicasets.apps",
        "clusterauthtokens",
        "clusteruserattributes",
        "controllerrevisions.apps",
        "apiservices.apiregistration.k8s.io",
        "clusterinformations",
        "felixconfigurations",
        "ippools",
        "nodes",
        "csinodes.storage.k8s.io",
        "csidrivers.storage.k8s.io",
        "priorityclasses.scheduling.k8s.io",
        "ciliumendpoints.cilium.io",
        "ciliumlocalredirectpolicies.cilium.io",
        "ciliumnetworkpolicies.cilium.io",
        "ciliumclusterwidenetworkpolicies.cilium.io",
        "ciliumegressnatpolicies.cilium.io",
        "ciliumexternalworkloads.cilium.io",
        "ciliumidentities.cilium.io",
        "flowschemas.flowcontrol.apiserver.k8s.io",
        "prioritylevelconfigurations.flowcontrol.apiserver.k8s.io",
        "horizontalpodautoscalers.autoscaling",
        "runtimeclasses.node.k8s.io",
        "nodes.metrics.k8s.io",
        "ciliumnodes.cilium.io",
        "events.events.k8s.io",
        "leases.coordination.k8s.io",
        "certificaterequests.cert-manager.io",
        "orders.acme.cert-manager.io",
        "challenges.acme.cert-manager.io",
        "mutatingwebhookconfigurations.admissionregistration.k8s.io",
        "validatingwebhookconfigurations.admissionregistration.k8s.io",
        "certificatesigningrequests.certificates.k8s.io",
        "ingresses.extensions",
        "pods.metrics.k8s.io",
    ].iter().map(|x| x.to_string()).collect();
}

