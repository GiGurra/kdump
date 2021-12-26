package com.github.gigurra.kdump.internal.util

import java.io.File
import java.nio.charset.StandardCharsets
import java.nio.file.{Files, Paths}

object file {

  def delete(path: String): Unit = Files.deleteIfExists(Paths.get(path))

  def exists(path: String): Boolean = Files.exists(Paths.get(path))

  def mkDirs(path: String): Unit = Files.createDirectories(Paths.get(path))

  def string2File(path: String, contents: String): Unit = Files.writeString(Paths.get(path), contents, StandardCharsets.UTF_8)

  def sanitizeFileName(name: String): String = "[^a-zA-Z0-9\\-_.]+".r.replaceAllIn(name, "_")

  def listFilesInDir(dir: String): List[File] =
    Option(new File(dir).listFiles()).toList.flatMap(_.toList)
}
