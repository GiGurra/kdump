mod util;

struct AppConfig {
    output_dir: String,
    delete_prev_dir: bool,
}

impl Default for AppConfig {
    fn default() -> Self {
        return AppConfig {
            output_dir: String::from("test"),
            delete_prev_dir: false
        };
    }
}

fn main() {
    println!("Checking output dir..");
    let appConfig = AppConfig::default();
    let root_output_dir = appConfig.output_dir;

    print!("Checking if output dir '{}' exists... ", root_output_dir);

    let output_dir_already_exists = util::file::path_exists("test");

    println!("{}", if output_dir_already_exists { "yes" } else { "no" });

    /*
        let command_output = Command::new("kubectl")
            .output()
            .expect("failed to execute process");

        let status_code = command_output.status.code().unwrap();

        let output_str = std::str::from_utf8(&command_output.stdout).unwrap();
        let err_str = std::str::from_utf8(&command_output.stderr).unwrap();

        println!("Hello, world!, cmd line result={}, output={}, err={}", status_code, output_str, err_str);*/
}

fn ensure_root_output_dir(appConfig: AppConfig) {

    if appConfig.delete_prev_dir {
        util::file::delete_all_if_exists(&appConfig.output_dir)
    }
    /*

        if appConfig.DeletePrevDir {
            fileutil.DeleteIfExists(out, fmt.Sprintf("removal of outputdir '%s' failed", out))
        }

        if fileutil.Exists(out, fmt.Sprintf("checking outputdir '%s' failed", out)) {
            log.Fatal("Bailing! output-dir already exists: " + out)
        }

        fileutil.CreateFolderIfNotExists(out, fmt.Sprintf("could not create folder '%s'", out))

        return out*/

}
