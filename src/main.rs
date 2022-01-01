use std::collections::HashMap;
use itertools::Itertools;
use crate::config::AppConfig;
use crate::util::k8s::ApiResourceType;
use crate::util::k8s::kubectl::ApiResourceTypes;
use crate::k8s::ApiResource;
use crate::util::k8s;

mod util;
mod config;

fn main() {
    util::logging::init();

    let app_config: AppConfig = config::AppConfig::from_cli_args();

    log::info!("Checking app configuration..");
    make_root_output_dir(&app_config);

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
        let output_dir: String = make_resource_output_dir(&app_config, &namespace_opt);
        for resource in resources {
            let file_path = make_resource_file_path(&output_dir, resource);
            let output_string = make_resource_file_contents(&app_config, resource);
            std::fs::write(&file_path, &output_string).expect(&format!("Unable to write file {}", file_path));
        }
    }

    log::info!("DONE!");
}

fn make_resource_output_dir(app_config: &config::AppConfig, namespace_opt: &Option<String>) -> String {
    let output_dir: String = match namespace_opt {
        Some(namespace) => app_config.output_dir.to_string() + "/" + &namespace.to_string(),
        None => app_config.output_dir.to_string(),
    };

    util::file::create_dir_all(&output_dir);

    return output_dir;
}

fn make_resource_file_path(output_dir: &str, resource: &k8s::ApiResource) -> String {
    let file_name: String = {
        let file_name_base: String = util::file::sanitize(&resource.parsed_fields.metadata.name) + "." + &util::file::sanitize(&resource.qualified_type_name()) + ".yaml";
        if resource.is_secret() {
            String::from(file_name_base + ".aes")
        } else {
            file_name_base
        }
    };

    return output_dir.to_string() + "/" + &file_name;
}

fn make_resource_file_contents(app_config: &config::AppConfig, resource: &k8s::ApiResource) -> String {
    return if resource.is_secret() {
        let encryption_key = app_config.secrets_encryption_key.as_ref().expect("BUG: encryption key has been removed or was never set");
        let encrypted = util::crypt::encrypt(&resource.raw_source, &encryption_key);
        String::from(encrypted.nonce_hex_string + &encrypted.encrypted_hex_string)
    } else {
        String::from(&resource.raw_source)
    };
}

fn make_root_output_dir(app_config: &config::AppConfig) {
    if app_config.delete_previous_dir {
        util::file::delete_all_if_exists(&app_config.output_dir);
    }

    if util::file::path_exists(&app_config.output_dir) {
        panic!("output path exists!: {}", app_config.output_dir);
    }

    util::file::create_dir_all(&app_config.output_dir);
}
