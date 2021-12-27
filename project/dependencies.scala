import sbt._

object dependencies {

  object circe extends scalaLibSource("io.circe", "0.14.1")
  object osLib extends scalaLibSource("com.lihaoyi", "0.8.0")

  val prod = Seq(
    circe("circe-yaml"),
    circe("circe-parser"),
    circe("circe-generic"),
    osLib("os-lib"),
  )
}

abstract class scalaLibSource(org: String, version: String) {
  def apply(component: String): ModuleID = org %% component % version
}
