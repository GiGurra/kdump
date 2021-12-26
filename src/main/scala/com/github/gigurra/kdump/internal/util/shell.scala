package com.github.gigurra.kdump.internal.util

import java.io.File

object shell {

  def commandExists(command: String): Boolean =

    // This is basically copy pasta of how it is done in go stdlib
    // Apparently there doesnt seem to exist a std solution in java/scala :S
    // (prob because nobody writes CLI apps in java/scala :D)

    val entries: List[String] = sys.env("PATH").split(File.pathSeparatorChar).toList

    def commandExistsInDir(dir: String): Boolean = new File(dir).listFiles().toList.exists(isMatch)

    def isMatch(file: File): Boolean =

      def isEmptyOrValidExecutableExtension(str: String): Boolean =
        str.isEmpty || Set(".exe", ".cmd", ".com", ".bat", ".app").contains(str.toLowerCase)

      file.isFile &&
        !file.isDirectory &&
        file.canExecute &&
        file.getName.startsWith(command) &&
        isEmptyOrValidExecutableExtension(file.getName.drop(command.length))

    entries.exists(commandExistsInDir)

}
