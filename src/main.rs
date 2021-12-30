use std::collections::BTreeMap;
use log::LevelFilter;
use crate::util::k8s::ApiResourceType;
use crate::util::k8s::kubectl::ApiResourceTypes;
use simple_logger::SimpleLogger;
use crate::k8s::ApiResource;
use crate::util::k8s;

mod util;
mod config;

fn main() {
    SimpleLogger::new().with_level(LevelFilter::Info).init().unwrap();

    log::info!("Checking output dir..");
    let app_config = config::AppConfig::default();
    ensure_root_output_dir(&app_config);

    log::info!("Checking what k8s types to download...");

    let all_resource_type_defs: ApiResourceTypes = util::k8s::kubectl::api_resource_types();
    let resource_type_defs_to_download: Vec<&ApiResourceType> = app_config.types_do_download(&all_resource_type_defs);

    log::info!("Downloading all objects...");

    let everything_as_string = util::k8s::kubectl::download_everything(&resource_type_defs_to_download);

    log::info!("Deserializing yaml...");

    let resources: Vec<ApiResource> = k8s::parse_resource_list(&everything_as_string);

    for resource in resources {
       println!("resource: {}", resource.parsed_fields.qualified_type_name());
    }

    //println!("everything: \n{}", everything_as_string);


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
