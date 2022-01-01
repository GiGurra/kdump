use std::process::Output;
use std::string::FromUtf8Error;
use crate::util::shell::RunCommandError::*;

#[derive(Debug, PartialEq, Clone)]
pub struct CommandDescription {
    program: String,
    args: Vec<String>,
}

#[derive(Debug)]
pub enum RunCommandError {
    RunCommandFailed(CommandDescription, std::io::Error),
    RunCommandInvalidOutput(CommandDescription, FromUtf8Error),
    RunCommandNoExitCode(CommandDescription),
    RunCommandNonZeroExitCode(CommandDescription, i32, Result<String, FromUtf8Error>),
}

fn get_description(cmd: &std::process::Command) -> CommandDescription {
    CommandDescription {
        program: cmd.get_program().to_str().get_or_insert("<unknown>").to_string(),
        args: cmd.get_args().map(|x| x.to_str().get_or_insert("<unknown>").to_string()).collect(),
    }
}

pub fn run_command(cmd: &mut std::process::Command) -> Result<String, RunCommandError> {
    let output: Output =
        cmd.output().map_err(|err| RunCommandFailed(get_description(cmd), err))?;

    match output.status.code() {
        None => Err(RunCommandNoExitCode(get_description(cmd))),
        Some(0) => String::from_utf8(output.stdout)
            .map_err(|err| RunCommandInvalidOutput(get_description(cmd), err.clone())),
        Some(status_code) => Err(RunCommandNonZeroExitCode(get_description(cmd), status_code, String::from_utf8(output.stderr))),
    }
}
