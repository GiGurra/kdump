use std::path::Path;

pub fn path_exists(path: &str) -> bool {
    return Path::new(path).exists();
}

pub fn delete_all_if_exists(path: &str) {
    if path_exists(path) {
        std::fs::remove_dir_all(Path::new(path)).unwrap();
    }
}

pub fn create_dir_all(path: &str) {
    std::fs::create_dir_all(Path::new(path)).unwrap();
}
