package com.github.gigurra.kdump.internal.util

import io.circe.*

import scala.annotation.tailrec
import scala.util.Try

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

  def downloadAllNamespacedResources(namespace: String, resourceTypes: List[String]): Seq[K8sResource] =
    resourceTypes.foreach(println)
    System.exit(0)
    val strVal = runCommand("-n", namespace, "get", resourceTypes.mkString(","), "-o", "yaml").trim
    parseK8sResourceList(strVal)

  def downloadAllResources(resourceTypes: List[String]): Seq[K8sResource] =
    val strVal = runCommand("get", resourceTypes.mkString(","), "-o", "yaml", "--all-namespaces").trim
    parseK8sResourceList(strVal)

  case class K8sResource(apiVersion: String,
                         kind: String,
                         name: String,
                         namespace: Option[String],
                         sourceYaml: String) {

    lazy val isSecret: Boolean =
      Set("secret", "secrets").contains(kind.toLowerCase)

    lazy val typeSuffix: String =
      apiVersion.indexOf("/") match
        case -1 => ""
        case index => s".${apiVersion.take(index)}"

    lazy val qualifiedKind: String = s"$kind$typeSuffix".toLowerCase
    lazy val qualifiedName: String = s"$name.$qualifiedKind"

    override def toString: String = s"K8sResource(qualifiedKind=$qualifiedKind, name=$name, namespace=$namespace)"
  }

  private def parseK8sResourceList(yamlString: String): Seq[K8sResource] =

    val json: JsonObject =
      yaml.parser.parse(yamlString)
        .getOrElse(panic("failed parsing root input as yaml!"))
        .asObject.getOrElse(panic("input wasn't an object"))

    val resourceList: Seq[JsonObject] = json("items")
      .getOrElse(panic("input yaml missing 'items' list"))
      .asArray.getOrElse(panic("input yaml 'items' wasn't a list"))
      .map(_.asObject.getOrElse(panic("items in yaml 'items' list weren't objects")))

    resourceList.map(parseK8sResource)

  private def parseK8sResource(obj: JsonObject): K8sResource =
    K8sResource(
      apiVersion = extractString(obj, "apiVersion"),
      kind = extractString(obj, "kind"),
      name = extractString(obj, "metadata", "name"),
      namespace = Try(extractString(obj, "metadata", "namespace")).toOption,
      sourceYaml = yaml.printer.print(Json.fromJsonObject(obj)),
    )

  @tailrec private def extractString(src: Json, path: String*): String =
    if (path.isEmpty)
      src.asString.getOrElse(panic(s"wasnt a string: $src"))
    else
      val newRoot: Json =src.asObject
        .getOrElse(panic(s"wasnt an object: $src"))
        .apply(path.head).getOrElse(panic(s"didnt have key ${path.head} in object: $src"))
      extractString(newRoot, path.tail: _*)

  private def extractString(src: JsonObject, path: String*): String =
    extractString(Json.fromJsonObject(src), path: _*)

  extension (fullString: String)
    def removeUpToAndIncluding(key: String): String =
      fullString.indexOf(key) match
        case -1 => fullString
        case index => fullString.drop(index + key.length)

}
