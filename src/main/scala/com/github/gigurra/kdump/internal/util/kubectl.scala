package com.github.gigurra.kdump.internal.util

object kubectl {

  import scala.sys.process._

  def namespaces(): List[String] =
    runCommand2Lines("get", "namespaces", "-o", "name")
      .map(removeK8sResourcePrefix)
      .map(_.trim)
      .filter(_.nonEmpty)

  def currentNamespace(): String =
    runCommand("config", "view", "--minify", "--output", "jsonpath={..namespace}")

  def currentContext(): String =
    runCommand("config", "current-context")

  def runCommand(args: String*): String =
    if !shell.commandExists("kubectl") then
      panic("kubectl not on path!")
    ("kubectl" +: args).!!.trim

  def runCommand2Lines(args: String*): List[String] =
    runCommand(args: _*).linesIterator.toList

  def removeK8sResourcePrefix(in: String): String =
    in.removeUpToAndIncluding("/")

  def resourceTypeNames(): List[String] =
    runCommand2Lines("api-resources", "-o", "name", "--verbs", "get")
      .map(_.trim)
      .filter(_.nonEmpty)

  def globalResourceTypeNames(): List[String] =
    runCommand2Lines("api-resources", "--namespaced=false", "-o", "name", "--verbs", "get")
      .map(_.trim)
      .filter(_.nonEmpty)

  def namespacedResourceTypeNames(): List[String] =
    runCommand2Lines("api-resources", "--namespaced=true", "-o", "name", "--verbs", "get")
      .map(_.trim)
      .filter(_.nonEmpty)

  def listNamespacedResourcesOfType(namespace: String, resourceTypeName: String): List[String] =
    runCommand2Lines("-n", namespace, "get", resourceTypeName, "-o", "name")
      .map(_.trim)
      .filter(_.nonEmpty)
      .map(_.removeUpToAndIncluding("/"))

  def downloadAllResources(namespace: String, resourceTypes: List[String], format: String): String =
    runCommand("-n", namespace, "get", resourceTypes.mkString(","), "-o", format).trim

  def downloadNamespacedResource(namespace: String, resourceType: String, resourceName: String, format: String): String =
    runCommand("-n", namespace, "get", resourceType, resourceName, "-o", format).trim

  def downloadGlobalResource(resourceType: String, resourceName: String, format: String): String =
    runCommand("get", resourceType, resourceName, "-o", format).trim

  extension (fullString: String)
    def removeUpToAndIncluding(key: String): String =
      fullString.indexOf(key) match
        case -1 => fullString
        case index => fullString.drop(index + key.length)

}
