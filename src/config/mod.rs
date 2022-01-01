use crate::{ApiResourceType, util};

use clap::{AppSettings, Parser, Subcommand};

/// Dump all kubernetes resources as yaml files to a dir
#[derive(Parser, Debug, PartialEq, Clone)]
#[clap(about, version, setting(AppSettings::SubcommandRequired))]
pub struct CliArgs {
    #[clap(subcommand)]
    command: Command,
}

/// Doc comment
#[derive(Subcommand, Debug, PartialEq, Clone)]
#[clap(setting(AppSettings::DeriveDisplayOrder | AppSettings::DisableHelpSubcommand))]
enum Command {
    /// Normal usage. Download all resources
    #[clap(setting(AppSettings::DeriveDisplayOrder))]
    Download {
        /// REQUIRED: output directory to create
        #[clap(short, long, required = true)]
        output_dir: String,

        /// if to delete previous output directory (default: false)
        #[clap(long)]
        delete_previous_dir: bool,

        /// symmetric secrets encryption hex key for aes GCM (lower case 64 chars)
        #[clap(long, required = false)]
        secrets_encryption_key: Option<String>,

        /// disable default excluded types
        #[clap(long)]
        no_default_excluded_types: bool,

        /// add additional excluded types
        #[clap(long, required = false)]
        excluded_types: Vec<String>,
    },

    /// Don't download resources - instead show resource types available for download in the cluster
    ClusterResourceTypes,

    /// Don't download resources - instead show default excluded types
    DefaultExcludedTypes,
}

pub struct AppCfg {
    pub output_dir: String,
    pub delete_previous_dir: bool,
    pub secrets_encryption_key: Option<Vec<u8>>,
    pub excluded_types: Vec<String>,
}

impl Default for AppCfg {
    fn default() -> Self {
        Self {
            output_dir: "test".to_string(),
            delete_previous_dir: false,
            secrets_encryption_key: None,
            excluded_types: default_resources_excluded(),
        }
    }
}

impl AppCfg {
    pub fn from_cli_args() -> Self {
        let cli_args: CliArgs = CliArgs::parse();
        match cli_args.command {
            Command::DefaultExcludedTypes => {
                println!("Default excluded types:");
                for tpe in default_resources_excluded() {
                    println!(" - {}", tpe);
                }
                std::process::exit(0);
            }
            Command::ClusterResourceTypes => {
                let types = util::k8s::kubectl::api_resource_types()
                    .expect("Failed to download k8s resource types");
                println!("Cluster types:");
                for tpe in types.accessible.all {
                    println!(" - {}", tpe.qualified_name());
                }
                std::process::exit(0);
            }
            Command::Download {
                output_dir,
                delete_previous_dir,
                secrets_encryption_key,
                no_default_excluded_types,
                excluded_types,
            } => {
                let mut result = Self::default();
                if no_default_excluded_types {
                    result.excluded_types.clear();
                }
                result.excluded_types.extend_from_slice(&excluded_types);
                result.secrets_encryption_key = secrets_encryption_key.map(|x| parse_encryption_key(&x));
                result.delete_previous_dir = delete_previous_dir;
                result.output_dir = output_dir;

                result
            }
        }
    }

    pub const fn include_secrets(&self) -> bool {
        self.secrets_encryption_key.is_some()
    }

    pub fn is_type_included(&self, tpe: &util::k8s::ApiResourceType) -> bool {
        !self.excluded_types.contains(&tpe.name) &&
            !self.excluded_types.contains(&tpe.qualified_name()) &&
            (!tpe.is_secret() || self.include_secrets())
    }

    pub fn types_do_download<'a>(&self, all_resource_type_defs: &'a util::k8s::kubectl::ApiResourceTypes) -> Vec<&'a ApiResourceType> {
        all_resource_type_defs.accessible.all
            .iter()
            .filter(|x| self.is_type_included(x))
            .collect::<Vec<&ApiResourceType>>()
    }
}


pub fn default_resources_excluded() -> Vec<String> {
    vec![
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
    ].iter().map(|x| (*x).to_string()).collect()
}


fn parse_encryption_key(hex_str: &str) -> Vec<u8> {
    log::info!("verifying encryption key length and format...");
    if hex_str.len() != 64 {
        panic!("key string was of size {}, must be exactly 64 hex characters", hex_str.len());
    }
    hex::decode(hex_str).expect("key was not a valid hex string")
}
