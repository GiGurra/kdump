mod util;

struct AppConfig {
    output_dir: String,
    delete_prev_dir: bool,
}

impl Default for AppConfig {
    fn default() -> Self {
        return AppConfig {
            output_dir: String::from("test"),  // TODO: Change to default empty when implementing cli args
            delete_prev_dir: true, // TODO: Change to default false when implementing cli args
        };
    }
}

fn main() {
    println!("Checking output dir..");
    let app_config = AppConfig::default();
    ensure_root_output_dir(app_config);

    println!("Downloading all resources from current context");


    let result = util::shell::run_command(
        std::process::Command::new("kubectl").arg("get").arg("all")
    );

    println!("result: {}", result);

    /*
            .output()
            .expect("failed to execute process");

        let status_code = command_output.status.code().unwrap();

        let output_str = std::str::from_utf8(&command_output.stdout).unwrap();
        let err_str = std::str::from_utf8(&command_output.stderr).unwrap();

        println!("Hello, world!, cmd line result={}, output={}, err={}", status_code, output_str, err_str);*/
}

fn ensure_root_output_dir(app_config: AppConfig) {
    if app_config.delete_prev_dir {
        util::file::delete_all_if_exists(&app_config.output_dir);
    }

    if util::file::path_exists(&app_config.output_dir) {
        panic!("output path exists!: {}", app_config.output_dir);
    }

    util::file::create_dir_all(&app_config.output_dir);
}
