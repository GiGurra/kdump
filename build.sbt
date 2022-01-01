name := "kdump"

version := "0.1.0-SNAPSHOT"

scalaVersion := "3.1.0"

libraryDependencies ++= dependencies.prod

//enablePlugins(PackPlugin)

enablePlugins(JlinkPlugin)

jlinkIgnoreMissingDependency := JlinkIgnore.everything

jlinkModules := {
  jlinkModules.value :+ "jdk.unsupported"
}