mod util;

fn main() {
    let output_dir = "test";

    print!("Checking if output dir '{}' exists... ", output_dir);

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
