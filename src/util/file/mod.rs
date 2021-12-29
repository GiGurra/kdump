use std::path::Path;

pub fn path_exists(path: &str) -> bool {
    return Path::new(path).exists();
}
