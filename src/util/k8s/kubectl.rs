use std::str::FromStr;
use regex::Regex;
use std::collections::HashMap;
use crate::util::k8s::*;
use crate::util;

#[derive(Debug, PartialEq, Clone)]
pub struct ApiResourceTypes {
    pub all: Vec<ApiResourceType>,
    pub accessible: AccessibleApiResourceTypes,
}

impl ApiResourceTypes {
    pub fn from(all_values: Vec<ApiResourceType>) -> ApiResourceTypes {
        return ApiResourceTypes {
            all: all_values.to_vec(),
            accessible: AccessibleApiResourceTypes::from(&all_values),
        };
    }
}

#[derive(Debug, PartialEq, Clone)]
pub struct AccessibleApiResourceTypes {
    pub all: Vec<ApiResourceType>,
    pub namespaced: Vec<ApiResourceType>,
    pub global: Vec<ApiResourceType>,
}

impl AccessibleApiResourceTypes {
    pub fn from(all_values: &Vec<ApiResourceType>) -> AccessibleApiResourceTypes {
        let accessible_resources =
            all_values.clone()
                .iter()
                .filter(|x| x.verbs.contains(&"get".to_string()))
                .map(|x| x.clone())
                .collect::<Vec<ApiResourceType>>();

        return AccessibleApiResourceTypes {
            all: accessible_resources.to_vec(),
            namespaced: accessible_resources.iter().filter(|x| x.namespaced).map(|x| x.clone()).collect::<Vec<ApiResourceType>>(),
            global: accessible_resources.iter().filter(|x| !x.namespaced).map(|x| x.clone()).collect::<Vec<ApiResourceType>>(),
        };
    }
}

pub fn api_resource_types() -> ApiResourceTypes {

    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("api-resources").arg("-o").arg("wide")
    );

    let lines: Vec<String> = result.lines().map(|x| String::from(x)).collect::<Vec<String>>();

    let line_maps: Vec<HashMap<String, String>> = util::string::parse_stdout_table(&lines);

    let list: Vec<ApiResourceType> = line_maps.iter().map(map_to_resource_type).collect::<Vec<ApiResourceType>>();

    return ApiResourceTypes::from(list);
}

pub fn download_everything(types_to_download: &Vec<&ApiResourceType>) -> String {
    let qualified_names: Vec<String> = types_to_download.iter().map(|x| x.qualified_name()).collect();
    let qualified_names_joined: String = qualified_names.join(",");
    return util::shell::run_command(std::process::Command::new("kubectl").arg("get").arg(&qualified_names_joined).arg("-o").arg("yaml"));
}

fn map_to_resource_type(map: &HashMap<String, String>) -> ApiResourceType {
    return ApiResourceType {
        name: String::from(&map["NAME"]),
        short_names: util::string::split_to_vec(&map["SHORTNAMES"], ",", true),
        namespaced: bool::from_str(&map["NAMESPACED"]).unwrap(),
        kind: String::from(&map["KIND"]),
        verbs: util::string::split_to_vec_r(util::string::remove_wrap(&map["VERBS"]), &Regex::new(r"\s+").unwrap(), true),
        api_version: parse_api_version(&map["APIVERSION"]),
    };
}