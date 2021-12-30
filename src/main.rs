use std::collections::HashMap;
use log::LevelFilter;
use crate::util::k8s::ApiResourceType;
use crate::util::k8s::kubectl::ApiResourceTypes;
use simple_logger::SimpleLogger;
use crate::k8s::ApiResource;
use crate::util::k8s;
use itertools::Itertools;

mod util;
mod config;


fn main() {
    SimpleLogger::new().with_level(LevelFilter::Info).init().unwrap();

    log::info!("Checking output dir..");
    let app_config = config::AppConfig::default(); // TODO: Implement cmd line args
    ensure_root_output_dir(&app_config);

    log::info!("Checking what k8s types to download...");

    let all_resource_type_defs: ApiResourceTypes = util::k8s::kubectl::api_resource_types();
    let resource_type_defs_to_download: Vec<&ApiResourceType> = app_config.types_do_download(&all_resource_type_defs);

    log::info!("Downloading all objects...");

    let everything_as_string = util::k8s::kubectl::download_everything(&resource_type_defs_to_download);

    log::info!("Deserializing yaml...");

    let resources: Vec<ApiResource> = k8s::parse_resource_list(&everything_as_string, true);
    let resources_by_namespace: HashMap<Option<String>, Vec<&ApiResource>> = resources.iter().into_group_map_by(|a| a.parsed_fields.metadata.namespace.clone());

    log::info!("Writing yaml files...");

    for (namespace_opt, resources) in resources_by_namespace {
        let output_dir: String = match namespace_opt {
            Some(namespace) => app_config.output_dir.to_string() + "/" + &namespace.to_string(),
            None => app_config.output_dir.to_string(),
        };
        util::file::create_dir_all(&output_dir);
        for resource in resources {
            let file_name = util::file::sanitize(&resource.parsed_fields.metadata.name) + "." + &util::file::sanitize(&resource.qualified_type_name()) + ".yaml";
            let file_path = output_dir.to_string() + "/" + &file_name;
            if resource.is_secret() {
                log::warn!("Secrets not implemented, ignoring {}", file_path);  // TODO: Implement secrets handling/encryption
            } else {
                std::fs::write(&file_path, &resource.raw_source).expect(&format!("Unable to write file {}", file_path));
            }
        }
    }

    log::info!("DONE!");
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
