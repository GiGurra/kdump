pub mod kubectl;

#[derive(Debug, PartialEq)]
pub struct ApiVersion {
    version: String,
    name: String,
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