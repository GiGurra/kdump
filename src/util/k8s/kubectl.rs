use std::str::FromStr;
use regex::Regex;
use crate::util::k8s::{ApiResourceType, ApiVersion};
use super::super::super::*; // access all modules between util modules

pub fn api_resource_types() -> Vec<super::ApiResourceType> {
    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("api-resources").arg("-o").arg("wide")
    );

    let lines: Vec<String> = result.lines().map(|x| String::from(x)).collect::<Vec<String>>();

    let line_maps: Vec<HashMap<String, String>> = util::string::parse_stdout_table(&lines);

    let result = line_maps.iter().map(map_to_resource_type).collect::<Vec<ApiResourceType>>();

    return result;
}

fn map_to_resource_type(map: &HashMap<String, String>) -> ApiResourceType {
    let mut result = ApiResourceType {
        name: String::from(&map["NAME"]),
        short_names: util::string::split_to_vec(&map["SHORTNAMES"], ",", true),
        namespaced: bool::from_str(&map["NAMESPACED"]).unwrap(),
        kind: String::from(&map["KIND"]),
        verbs: util::string::split_to_vec_r(util::string::remove_wrap(&map["VERBS"]), &Regex::new(r"\s+").unwrap(), true),
        api_version: ApiVersion {
            name: "".to_string(),
            version: "".to_string(),
        },
        qualified_name: "".to_string(),
    };

    let api_version_str_parts = util::string::split_to_vec(&map["APIVERSION"], "/", true);
    if api_version_str_parts.len() > 1 {
        result.api_version.name = api_version_str_parts[0].to_string();
        result.api_version.version = api_version_str_parts[1].to_string();
        result.qualified_name = result.name.clone() + "." + &result.api_version.name.clone();
    } else {
        result.api_version.version = api_version_str_parts[0].to_string();
        result.qualified_name = result.name.clone();
    };

    return result;
}