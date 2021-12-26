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

  extension (fullString: String)
    def removeUpToAndIncluding(key: String): String =
      fullString.indexOf(key) match
        case -1 => fullString
        case index => fullString.drop(index + key.length)

}
