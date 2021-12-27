package com.github.gigurra.kdump.internal.util

import io.circe.*

import scala.annotation.tailrec
import scala.util.Try

object k8sYaml {

  case class K8sResource(apiVersion: String,
                         kind: String,
                         name: String,
                         namespace: Option[String],
                         source: String) {

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

  def parseList(yamlString: String): Seq[K8sResource] =

    val json: JsonObject =
      yaml.parser.parse(yamlString)
        .getOrElse(panic("failed parsing root input as yaml!"))
        .asObject.getOrElse(panic("input wasn't an object"))

    val itemsList: Seq[JsonObject] = json("items")
      .getOrElse(panic("input yaml missing 'items' list"))
      .asArray.getOrElse(panic("input yaml 'items' wasn't a list"))
      .map(_.asObject.getOrElse(panic("items in yaml 'items' list weren't objects")))

    itemsList.map { obj =>

      K8sResource(
        apiVersion = extractString(obj, "apiVersion"),
        kind = extractString(obj, "kind"),
        name = extractString(obj, "metadata", "name"),
        namespace = Try(extractString(obj, "metadata", "namespace")).toOption,
        source = yaml.printer.print(Json.fromJsonObject(obj)),
      )
    }
}

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



