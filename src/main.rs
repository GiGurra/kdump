use std::collections::HashMap;
use itertools::Itertools;
use crate::config::AppCfg;
use crate::util::k8s::ApiResourceType;
use crate::util::k8s::kubectl::ApiResourceTypes;
use crate::k8s::ApiResource;
use crate::util::{crypt, k8s};

mod util;
mod config;

fn main() {
    util::logging::init();

    let app_config: AppCfg = config::AppCfg::from_cli_args();

    log::info!("Checking app configuration..");
    check_root_output_dir(&app_config);

    log::info!("Checking what k8s types to download...");

    let all_resource_type_defs: ApiResourceTypes = util::k8s::kubectl::api_resource_types().expect("Failed to download k8s resource types");
    let resource_type_defs_to_download: Vec<&ApiResourceType> = app_config.types_do_download(&all_resource_type_defs);

    log::info!("Downloading all objects...");

    let everything_as_string: String = util::k8s::kubectl::download_everything(&resource_type_defs_to_download)
        .expect("Failed to download all k8s resources");

    log::info!("Deserializing yaml...");

    let resources: Vec<ApiResource> = k8s::parse_resource_list(&everything_as_string, true).expect("could not parse 'everything yaml' from kubectl");
    let resources_by_namespace: HashMap<Option<String>, Vec<&ApiResource>> = resources.iter().into_group_map_by(|a| a.parsed_fields.metadata.namespace.clone());

    log::info!("Writing yaml files...");

    for (namespace_opt, resources) in resources_by_namespace {
        let output_dir: String = make_resource_output_dir(&app_config, &namespace_opt).expect("Unable to create output dir");
        for resource in resources {
            let file_path = make_resource_file_path(&output_dir, resource);
            let output_string = make_resource_file_contents(&app_config, resource).unwrap_or_else(|err| panic!("unable to make resource file contents for: {}, due to {:?}", file_path, err));
            std::fs::write(&file_path, &output_string).unwrap_or_else(|err| panic!("Unable to write file {}, due to {:?}", file_path, err));
        }
    }

    log::info!("DONE!");
}

fn make_resource_output_dir(app_config: &config::AppCfg, namespace_opt: &Option<String>) -> std::io::Result<String> {
    let output_dir: String = match namespace_opt {
        Some(namespace) => app_config.output_dir.to_string() + "/" + &namespace.to_string(),
        None => app_config.output_dir.to_string(),
    };

    util::file::create_dir_all(&output_dir)?;

    Ok(output_dir)
}

fn make_resource_file_path(output_dir: &str, resource: &k8s::ApiResource) -> String {
    let file_name: String = {
        let file_name_base: String = util::file::sanitize(&resource.parsed_fields.metadata.name) + "." + &util::file::sanitize(&resource.qualified_type_name()) + ".yaml";
        if resource.is_secret() {
            file_name_base + ".aes"
        } else {
            file_name_base
        }
    };

    output_dir.to_string() + "/" + &file_name
}

fn make_resource_file_contents(app_config: &config::AppCfg, resource: &k8s::ApiResource) -> Result<String, crypt::EncryptError> {
    if resource.is_secret() {
        let encryption_key = app_config.secrets_encryption_key.as_ref().expect("BUG: encryption key has been removed or was never set");
        let encrypted = util::crypt::encrypt(&resource.raw_source, encryption_key)?;
        Ok(encrypted.nonce_hex_string + &encrypted.encrypted_hex_string)
    } else {
        Ok(resource.raw_source.to_string())
    }
}

fn check_root_output_dir(app_config: &config::AppCfg) {
    if app_config.delete_previous_dir {
        util::file::delete_all_if_exists(&app_config.output_dir)
            .unwrap_or_else(|err| panic!("unable to delete output dir: {}, due to {:?}", &app_config.output_dir, err));
    }

    if util::file::path_exists(&app_config.output_dir) {
        panic!("output path exists!: {}", app_config.output_dir);
    }

    util::file::create_dir_all(&app_config.output_dir)
        .unwrap_or_else(|err| panic!("unable to create output dir: {}, due to {:?}", &app_config.output_dir, err));
}
