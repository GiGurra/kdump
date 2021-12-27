package com.github.gigurra.kdump.internal.util

import java.io.File
import java.nio.charset.StandardCharsets
import java.nio.file.{Files, Paths}

object file {

  def resolve(path: String): os.Path = os.Path(Paths.get(path).toAbsolutePath)

  def delete(path: String): Unit = os.remove.all(resolve(path))

  def exists(path: String): Boolean = os.exists(resolve(path))

  def mkDirs(path: String): Unit = os.makeDir.all(resolve(path))

  def string2File(path: String, contents: String): Unit = os.write.over(resolve(path), contents)

  def sanitizeFileName(name: String): String = "[^a-zA-Z0-9\\-_.]+".r.replaceAllIn(name, "_")

  def listFilesInDir(dir: String): Seq[File] = if exists(dir) then os.list(resolve(dir)).map(_.toIO) else Nil
}
