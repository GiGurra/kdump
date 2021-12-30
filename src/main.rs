use crate::util::k8s::ApiResourceType;
use crate::util::k8s::kubectl::ApiResourceTypes;

mod util;
mod config;

fn main() {
    println!("Checking output dir..");
    let app_config = config::AppConfig::default();
    ensure_root_output_dir(&app_config);

    println!("Downloading all resources from current context");

    let all_resource_type_defs: ApiResourceTypes = util::k8s::kubectl::api_resource_types();
    let resource_type_defs_to_download: Vec<&ApiResourceType> = app_config.types_do_download(&all_resource_type_defs);

    for resource in &resource_type_defs_to_download {
        println!("resource: {:?}", resource);
    }
}

fn ensure_root_output_dir(app_config: &config::AppConfig) {
    if app_config.delete_prev_dir {
        util::file::delete_all_if_exists(&app_config.output_dir);
    }

    if util::file::path_exists(&app_config.output_dir) {
        panic!("output path exists!: {}", app_config.output_dir);
    }

    util::file::create_dir_all(&app_config.output_dir);
}
