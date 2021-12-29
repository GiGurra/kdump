pub fn run_command(cmd: &mut std::process::Command) -> String {
    let output = cmd.output().unwrap();
    if output.status.code().unwrap() == 0 {
        return String::from_utf8(output.stdout).unwrap();
    }

    panic!("command finished with status={}, stderr={}", output.status, String::from_utf8(output.stderr).unwrap());
}