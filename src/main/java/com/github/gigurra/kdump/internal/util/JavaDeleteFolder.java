package com.github.gigurra.kdump.internal.util;

import java.io.File;
import java.nio.file.Path;
import java.nio.file.Files;
import java.io.IOException;
import java.util.Comparator;

public class JavaDeleteFolder {
    public static void deleteDirectoryStream(Path path) throws IOException {
        Files.walk(path)
                .sorted(Comparator.reverseOrder())
                .map(Path::toFile)
                .forEach(File::delete);
    }
}
