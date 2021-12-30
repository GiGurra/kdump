use crate::{ApiResourceType, util};

#[derive(Debug, PartialEq, Clone)]
pub struct AppConfig {
    pub output_dir: String,
    pub delete_prev_dir: bool,
    pub excluded_types: Vec<String>,
}

impl AppConfig {
    pub fn is_type_included(&self, tpe: &util::k8s::ApiResourceType) -> bool {
        return !self.excluded_types.contains(&tpe.name) &&
            !self.excluded_types.contains(&tpe.qualified_name());
    }

    pub fn types_do_download<'a>(&self, all_resource_type_defs: &'a util::k8s::kubectl::ApiResourceTypes) -> Vec<&'a ApiResourceType> {
        return all_resource_type_defs.accessible.all
            .iter()
            .filter(|x| self.is_type_included(x))
            .collect::<Vec<&ApiResourceType>>();
    }
}

impl Default for AppConfig {
    fn default() -> Self {
        return AppConfig {
            output_dir: String::from("test"),  // TODO: Change to default empty when implementing cli args
            delete_prev_dir: true, // TODO: Change to default false when implementing cli args
            excluded_types: vec![],
        };
    }
}

