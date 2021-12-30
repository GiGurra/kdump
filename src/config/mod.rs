use crate::{ApiResourceType, util};

#[derive(Debug, PartialEq, Clone)]
pub struct AppConfig {
    pub output_dir: String,
    pub delete_prev_dir: bool,
    pub excluded_types: Vec<String>,
}

impl AppConfig {
    pub fn is_type_included(&self, tpe: &util::k8s::ApiResourceType) -> bool {
        return !self.excluded_types.contains(&tpe.name) &&
            !self.excluded_types.contains(&tpe.qualified_name());
    }

    pub fn types_do_download<'a>(&self, all_resource_type_defs: &'a util::k8s::kubectl::ApiResourceTypes) -> Vec<&'a ApiResourceType> {
        return all_resource_type_defs.accessible.all
            .iter()
            .filter(|x| self.is_type_included(x))
            .collect::<Vec<&ApiResourceType>>();
    }
}

impl Default for AppConfig {
    fn default() -> Self {
        return AppConfig {
            output_dir: String::from("test"),  // TODO: Change to default empty when implementing cli args
            delete_prev_dir: true, // TODO: Change to default false when implementing cli args
            excluded_types: default_resources_excluded(),
        };
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

