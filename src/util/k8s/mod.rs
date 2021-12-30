pub mod kubectl;

use serde::Deserialize;
use serde_yaml::Value;
use crate::util; // access all modules between util modules

#[derive(Debug, PartialEq, Clone)]
pub struct ApiVersion {
    pub name: String,
    pub version: String,
}

#[derive(Debug, PartialEq, Clone)]
pub struct ApiResourceType {
    pub name: String,
    pub short_names: Vec<String>,
    pub namespaced: bool,
    pub kind: String,
    pub verbs: Vec<String>,
    pub api_version: ApiVersion,
}

impl ApiResourceType {
    pub fn is_secret(&self) -> bool {
        return self.name.to_lowercase() == "secret" || self.name.to_lowercase() == "secrets";
    }

    pub fn qualified_name(&self) -> String {
        return if self.api_version.name.is_empty() {
            self.name.clone()
        } else {
            self.name.clone() + "." + &self.api_version.name.clone()
        };
    }
}


#[derive(PartialEq, Clone)]
pub struct ApiResource {
    pub raw_source: String,
    pub parsed_fields: ApiResourceParsedFields,
}


impl ApiResource {
    pub fn is_secret(&self) -> bool { return self.parsed_fields.is_secret(); }
    pub fn qualified_type_name(&self) -> String { return self.parsed_fields.qualified_type_name(); }
}

#[derive(Debug, PartialEq, Clone)]
#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ApiResourceParsedFields {
    pub kind: String,
    pub api_version: String,
    pub metadata: ApiResourceParsedFieldsMetaData,
}

impl ApiResourceParsedFields {
    pub fn is_secret(&self) -> bool {
        return self.kind.to_lowercase() == "secret" || self.kind.to_lowercase() == "secrets";
    }

    pub fn qualified_type_name(&self) -> String {
        let parsed_api_version = parse_api_version(&self.api_version);
        return if parsed_api_version.name.is_empty() {
            self.kind.to_lowercase().clone()
        } else {
            self.kind.to_lowercase().clone() + "." + &parsed_api_version.name.clone()
        };
    }
}

#[derive(Debug, PartialEq, Clone)]
#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ApiResourceParsedFieldsMetaData {
    pub name: String,
    pub namespace: Option<String>,
}


pub fn parse_api_version(input: &str) -> ApiVersion {
    let api_version_str_parts = util::string::split_to_vec(input, "/", true);
    return
        if api_version_str_parts.len() > 1 {
            ApiVersion {
                name: api_version_str_parts[0].to_string(),
                version: api_version_str_parts[1].to_string(),
            }
        } else {
            ApiVersion {
                name: "".to_string(),
                version: api_version_str_parts[0].to_string(),
            }
        };
}

pub fn parse_resource_list(data: &str) -> Vec<ApiResource> {
    let deserialized_map: serde_yaml::Value = serde_yaml::from_str(&data).unwrap();

    let root_object = deserialized_map.as_mapping().unwrap();

    let item_list: &Vec<Value> = root_object.get(&Value::from("items")).unwrap().as_sequence().unwrap();

    return item_list.iter().map(parse_resource).collect::<Vec<ApiResource>>();
}

pub fn parse_resource(data: &Value) -> ApiResource {
    let fields: ApiResourceParsedFields = serde_yaml::from_value::<ApiResourceParsedFields>(data.to_owned()).unwrap();

    ApiResource {
        raw_source: serde_yaml::to_string(data).unwrap(),
        parsed_fields: fields,
    }
}
