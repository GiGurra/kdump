package com.github.gigurra.kdump.internal.util

object panic {
  def apply(message: String): Nothing = throw new Error(message)
}
