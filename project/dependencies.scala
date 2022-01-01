import sbt._

object dependencies {

  object circe extends scalaLibSource("io.circe", "0.14.1")
  object osLib extends scalaLibSource("com.lihaoyi", "0.8.0")
  object logback extends scalaLibSource("ch.qos.logback", "1.2.6")

  val prod = Seq(
    circe("circe-yaml"),
    circe("circe-parser"),
    circe("circe-generic"),
    osLib("os-lib"),
    logback.java("logback-classic"),
  )
}

abstract class scalaLibSource(org: String, version: String) {
  def apply(component: String): ModuleID = org %% component % version
  def java(component: String): ModuleID = org % component % version
}
