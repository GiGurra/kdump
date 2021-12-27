import sbt._

object dependencies {

  object circe extends scalaLibSource("io.circe", "0.14.1")

  val prod = Seq(
    circe("circe-yaml"),
    circe("circe-parser"),
    circe("circe-generic"),
  )
}

abstract class scalaLibSource(org: String, version: String) {
  def apply(component: String): ModuleID = org %% component % version
}
