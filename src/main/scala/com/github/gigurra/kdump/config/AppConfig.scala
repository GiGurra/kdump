package com.github.gigurra.kdump.config

import com.github.gigurra.kdump.config.AppConfig.SecretsHandling


case class AppConfig(outputDir: String,
                     excludedResourceTypes: Set[String],
                     secretsHandling: AppConfig.SecretsHandling) {

  def isResourceIncluded(resourceTypeName: String): Boolean =
    !excludedResourceTypes.contains(resourceTypeName) &&
      (resourceTypeName != "secrets" || secretsHandling != SecretsHandling.DontStore)
}

object AppConfig {

  def default: AppConfig = AppConfig(
    outputDir = "test",
    excludedResourceTypes = defaultExcludedResourceTypes,
    secretsHandling = SecretsHandling.DontStore,
  )

  def defaultExcludedResourceTypes: Set[String] = Set(
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
  )

  sealed trait SecretsHandling

  object SecretsHandling {
    object DontStore extends SecretsHandling
    object StoreAsPlainText extends SecretsHandling
    case class StoreEncrypted(encryption: String /*todo: impl*/) extends SecretsHandling
  }
}