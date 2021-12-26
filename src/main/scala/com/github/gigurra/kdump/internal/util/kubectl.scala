package com.github.gigurra.kdump.internal.util

object kubectl {

  import scala.sys.process._

  def namespaces(): List[String] =
    string.splitLines(runCommand("get", "namespaces", "-o", "name"))
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

  def removeK8sResourcePrefix(in: String): String =
    string.removeUpToAndIncluding(in, "/")

  def resourceTypeNames(): List[String] =
    string.splitLines(runCommand("api-resources", "-o", "name", "--verbs", "get"))
      .map(_.trim)
      .filter(_.nonEmpty)

  def globalResourceTypeNames(): List[String] =
    string.splitLines(runCommand("api-resources", "--namespaced=false", "-o", "name", "--verbs", "get"))
      .map(_.trim)
      .filter(_.nonEmpty)

  def namespacedResourceTypeNames(): List[String] =
    string.splitLines(runCommand("api-resources", "--namespaced=true", "-o", "name", "--verbs", "get"))
      .map(_.trim)
      .filter(_.nonEmpty)
}
