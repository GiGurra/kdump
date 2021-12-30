use std::str::FromStr;
use regex::Regex;
use std::collections::HashMap;
use crate::util::k8s::*;
use crate::util;

pub struct ApiResourceTypes<'a> {
    pub all: Vec<ApiResourceType>,
    accessible: AccessibleApiResourceTypes<'a>,
}

impl ApiResourceTypes<'_> {
    pub fn default<'a>() -> ApiResourceTypes<'a> {
        ApiResourceTypes {
            all: vec![],
            accessible: AccessibleApiResourceTypes::default(),
        }
    }

    pub fn from<'a>(all_values: Vec<ApiResourceType>) -> ApiResourceTypes<'a> {
        ApiResourceTypes {
            all: all_values,
            accessible: AccessibleApiResourceTypes::default(),
        }
    }
}

pub struct AccessibleApiResourceTypes<'a> {
    all: Vec<&'a ApiResourceType>,
    namespaced: Vec<&'a ApiResourceType>,
    global: Vec<&'a ApiResourceType>,
}

impl AccessibleApiResourceTypes<'_> {
    pub fn default<'a>() -> AccessibleApiResourceTypes<'a> {
        AccessibleApiResourceTypes {
            all: vec![],
            namespaced: vec![],
            global: vec![],
        }
    }
}

pub fn api_resource_types<'a>() -> ApiResourceTypes<'a> {
    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("api-resources").arg("-o").arg("wide")
    );

    let lines: Vec<String> = result.lines().map(|x| String::from(x)).collect::<Vec<String>>();

    let line_maps: Vec<HashMap<String, String>> = util::string::parse_stdout_table(&lines);

    let list: Vec<ApiResourceType> = line_maps.iter().map(map_to_resource_type).collect::<Vec<ApiResourceType>>();

    return ApiResourceTypes::from(list);
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