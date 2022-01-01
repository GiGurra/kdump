use std::process::Output;
use std::string::FromUtf8Error;
use crate::util::shell::RunCommandError::{Failed, InvalidOutput, NoExitCode, NonZeroExitCode};

#[derive(Debug, PartialEq, Clone)]
pub struct CommandDescription {
    program: String,
    args: Vec<String>,
}

#[derive(Debug)]
pub enum RunCommandError {
    Failed(CommandDescription, std::io::Error),
    InvalidOutput(CommandDescription, FromUtf8Error),
    NoExitCode(CommandDescription),
    NonZeroExitCode(CommandDescription, i32, Result<String, FromUtf8Error>),
}

fn get_description(cmd: &std::process::Command) -> CommandDescription {
    CommandDescription {
        program: (**cmd.get_program().to_str().get_or_insert("<unknown>")).to_string(),
        args: cmd.get_args().map(|x| (**x.to_str().get_or_insert("<unknown>")).to_string()).collect(),
    }
}

pub fn run_command(cmd: &mut std::process::Command) -> Result<String, RunCommandError> {
    let output: Output =
        cmd.output().map_err(|err| Failed(get_description(cmd), err))?;

    match output.status.code() {
        None => Err(NoExitCode(get_description(cmd))),
        Some(0) => String::from_utf8(output.stdout)
            .map_err(|err| InvalidOutput(get_description(cmd), err)),
        Some(status_code) => Err(NonZeroExitCode(get_description(cmd), status_code, String::from_utf8(output.stderr))),
    }
}
