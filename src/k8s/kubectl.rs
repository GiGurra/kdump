use std::collections::HashMap;
use std::str::FromStr;

use gigurra_rust_util as util;
use util::shell::RunCommandError;

use crate::ApiResourceType;

#[derive(Debug, PartialEq, Clone)]
pub struct ApiResourceTypes {
    pub all: Vec<ApiResourceType>,
    pub accessible: AccessibleApiResourceTypes,
}

impl ApiResourceTypes {
    pub fn from(all_values: &[ApiResourceType]) -> Self {
        Self {
            all: Vec::from(all_values),
            accessible: AccessibleApiResourceTypes::from(all_values),
        }
    }
}

#[derive(Debug, PartialEq, Clone)]
pub struct AccessibleApiResourceTypes {
    pub all: Vec<ApiResourceType>,
    pub namespaced: Vec<ApiResourceType>,
    pub global: Vec<ApiResourceType>,
}

impl AccessibleApiResourceTypes {
    pub fn from(all_values: &[ApiResourceType]) -> Self {
        let accessible_resources =
            all_values
                .iter()
                .filter(|x| x.verbs.contains(&"get".to_string()))
                .cloned()
                .collect::<Vec<ApiResourceType>>();

        Self {
            all: accessible_resources.clone(),
            namespaced: accessible_resources.iter().filter(|x| x.namespaced).cloned().collect::<Vec<ApiResourceType>>(),
            global: accessible_resources.iter().filter(|x| !x.namespaced).cloned().collect::<Vec<ApiResourceType>>(),
        }
    }
}

pub fn api_resource_types() -> Result<ApiResourceTypes, RunCommandError> {
    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("api-resources").arg("-o").arg("wide")
    );

    result.map(|command_result_output| {
        let lines: Vec<String> = command_result_output.lines().map(String::from).collect::<Vec<String>>();
        let line_maps: Vec<HashMap<String, String>> = util::string::parse_stdout_table(&lines);
        let list: Vec<ApiResourceType> = line_maps.iter().map(map_to_resource_type).collect::<Vec<ApiResourceType>>();

        ApiResourceTypes::from(&list)
    })
}

pub fn download_everything(types_to_download: &[&ApiResourceType]) -> Result<String, RunCommandError> {
    let qualified_names: Vec<String> = types_to_download.iter().map(|x| x.qualified_name()).collect();
    let qualified_names_joined: String = qualified_names.join(",");
    util::shell::run_command(
        std::process::Command::new("kubectl")
            .arg("get")
            .arg(&qualified_names_joined)
            .arg("--all-namespaces")
            .arg("-o")
            .arg("yaml")
    )
}

fn map_to_resource_type(map: &HashMap<String, String>) -> ApiResourceType {
    ApiResourceType {
        name: String::from(&map["NAME"]),
        short_names: util::string::split_to_vec(&map["SHORTNAMES"], ",", true),
        namespaced: bool::from_str(&map["NAMESPACED"]).expect("non-bool 'NAMESPACED' in map"),
        kind: String::from(&map["KIND"]),
        verbs: util::string::split_to_vec_r(util::string::remove_wrap(&map["VERBS"]), &util::regex::Regex::new(r"\s+").expect("BUG: invalid regex to split VERBS"), true),
        api_version: super::parse_api_version(&map["APIVERSION"]),
    }
}