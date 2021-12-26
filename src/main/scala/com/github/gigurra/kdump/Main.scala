package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util._

@main def kdump(): Unit =
  val appConfig = AppConfig.default
  dumpCurrentContext(appConfig)

def dumpCurrentContext(appConfig: AppConfig): Unit =

  val currentK8sContext = kubectl.currentContext()
  val outputDir = appConfig.outDir(currentK8sContext)

  file.delete(outputDir)
  file.mkDirs(outputDir)

  println(s"Downloading all resources from current context to dir '$outputDir'")

  val namespaces = kubectl.namespaces()

  println(namespaces)


