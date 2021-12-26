package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util.kubectl

@main def kdump(): Unit =
  val appConfig = AppConfig.default
  dumpCurrentContext(appConfig)

def dumpCurrentContext(config: AppConfig): Unit =

  val currentK8sContext = kubectl.currentContext()

  println(s"currentK8sContext: $currentK8sContext")


