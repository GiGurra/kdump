use super::super::super::*; // access all modules between util modules

pub fn api_resource_types() -> Vec<HashMap<String, String>> {

    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("api-resources")
    );

    let lines = result.lines().map(|x| String::from(x)).collect::<Vec<String>>();

    return util::string::parse_stdout_table(&lines);
}