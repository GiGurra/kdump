pub mod kubectl;

use serde::Deserialize;
use serde_yaml::{Mapping, Value};
use crate::util;
// access all modules between util modules

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
        self.name.to_lowercase() == "secret" || self.name.to_lowercase() == "secrets"
    }

    pub fn qualified_name(&self) -> String {
        if self.api_version.name.is_empty() {
            self.name.clone()
        } else {
            self.name.clone() + "." + &self.api_version.name.clone()
        }
    }
}


#[derive(PartialEq, Clone)]
pub struct ApiResource {
    pub raw_source: String,
    pub parsed_fields: ApiResourceParsedFields,
}


impl ApiResource {
    pub fn is_secret(&self) -> bool { self.parsed_fields.is_secret() }
    pub fn qualified_type_name(&self) -> String { self.parsed_fields.qualified_type_name() }
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
        self.kind.to_lowercase() == "secret" || self.kind.to_lowercase() == "secrets"
    }

    pub fn qualified_type_name(&self) -> String {
        let parsed_api_version = parse_api_version(&self.api_version);
        if parsed_api_version.name.is_empty() {
            self.kind.to_lowercase()
        } else {
            self.kind.to_lowercase() + "." + &parsed_api_version.name
        }
    }
}

#[derive(Debug, PartialEq, Clone)]
#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ApiResourceParsedFieldsMetaData {
    pub name: String,
    pub namespace: Option<String>,
}

#[derive(Debug, PartialEq, Clone)]
#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ApiResourceList {
    pub api_version: String,
    pub kind: String,
    pub items: Vec<serde_yaml::Mapping>,
}

pub fn parse_api_version(input: &str) -> ApiVersion {
    let api_version_str_parts = util::string::split_to_vec(input, "/", true);

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
    }
}

pub fn parse_resource_list(data: &str, remove_status_fields: bool) -> serde_yaml::Result<Vec<ApiResource>> {
    let deserialized_resource_list: ApiResourceList = serde_yaml::from_str(data)?;

    let item_list = &deserialized_resource_list.items;

    item_list.iter()
        .map(|x| parse_resource(x, remove_status_fields))
        .collect()
}

pub fn parse_resource(data: &Mapping, remove_status_fields: bool) -> serde_yaml::Result<ApiResource> {
    let upcast = Value::from(data.to_owned());
    let fields: ApiResourceParsedFields = serde_yaml::from_value(upcast)?;

    let mut data_copy = data.clone();

    if remove_status_fields {
        data_copy.remove(&Value::from("status"));
        data_copy.remove(&Value::from("lastRefresh"));
    }

    Ok(ApiResource {
        raw_source: serde_yaml::to_string(&data_copy)?,
        parsed_fields: fields,
    })
}
