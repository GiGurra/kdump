pub mod kubectl;

#[derive(Debug, PartialEq)]
pub struct ApiVersion {
    name: String,
    version: String,
}

#[derive(Debug, PartialEq)]
pub struct ApiResourceType {
    name: String,
    short_names: Vec<String>,
    namespaced: bool,
    kind: String,
    verbs: Vec<String>,
    api_version: ApiVersion,
    qualified_name: String,
}