package com.github.gigurra.kdump.internal.util

object string {

  def splitLines(s: String): List[String] = s.linesIterator.toList

  def trimSpaces(s: List[String]): List[String] = s.map(_.trim)

  def removeEmpty(s: List[String]): List[String] = trimSpaces(s).filter(_.nonEmpty)

  def removeUpToAndIncluding(fullString: String, key: String): String =
    fullString.indexOf(key) match
      case -1 => fullString
      case index => fullString.drop(index + key.length)


}
