use std::path::Path;
use regex::Regex;

pub fn path_exists(path: &str) -> bool {
    Path::new(path).exists()
}

pub fn delete_all_if_exists(path: &str) -> std::io::Result<()> {
    if path_exists(path) {
        std::fs::remove_dir_all(Path::new(path))
    } else {
        Ok(())
    }
}

pub fn create_dir_all(path: &str) -> std::io::Result<()> {
    std::fs::create_dir_all(Path::new(path))
}

pub fn sanitize(path: &str) -> String {
    let regex: Regex = regex::Regex::new(r"[^a-zA-Z0-9\-_.]+").expect("BUG: file sanitize regex is invalid");
    regex.replace_all(path, "_").to_string()
}
