use crate::util::k8s::ApiResourceType;
use super::super::super::*; // access all modules between util modules

pub fn api_resource_types() -> Vec<super::ApiResourceType> {
    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("api-resources")
    );

    let lines: Vec<String> = result.lines().map(|x| String::from(x)).collect::<Vec<String>>();

    let line_maps: Vec<HashMap<String, String>> = util::string::parse_stdout_table(&lines);

    let result = line_maps.iter().map(map_to_resource_type).collect::<Vec<ApiResourceType>>();

    return result;
}

fn map_to_resource_type(map: &HashMap<String, String>) -> ApiResourceType {
    ApiResourceType {
        name: String::from(&map["NAME"]),
        short_names: vec![],
        namespaced: false,
        kind: "".to_string(),
        verbs: vec![],
        api_version: super::ApiVersion { version: "".to_string(), name: "".to_string() },
        qualified_name: "".to_string(),
    }
}