use std::path::Path;
use regex::Regex;

pub fn path_exists(path: &str) -> bool {
    return Path::new(path).exists();
}

pub fn delete_all_if_exists(path: &str) {
    if path_exists(path) {
        std::fs::remove_dir_all(Path::new(path))
            .expect(&format!("could not delete dir: '{}'", path));
    }
}

pub fn create_dir_all(path: &str) {
    std::fs::create_dir_all(Path::new(path))
        .expect(&format!("could not crate dir: '{}'", path));
}

pub fn sanitize(path: &str) -> String {
    let regex: Regex = regex::Regex::new(r"[^a-zA-Z0-9\-_.]+").expect("BUG: file sanitize regex is invalid");
    return regex.replace_all(path, "_").to_owned().to_string();
}
