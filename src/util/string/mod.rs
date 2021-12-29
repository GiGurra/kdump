use std::fmt::{Debug, Formatter};

#[derive(Debug)]
pub struct StdOutTableColumn {
    name: String,
    byte_offset: usize,
}

pub fn parse_stdout_table(lines: &Vec<String>) {
    let headings_line = &lines[0]; // .split_whitespace().collect::<Vec<&str>>()
    println!("head: {:?}", headings_line);

    let mut begin_offsets = Vec::<usize>::new();
    let mut prev_is_space = true;

    for (i, c) in headings_line.bytes().enumerate() {
        if prev_is_space && !c.is_ascii_whitespace() {
            begin_offsets.push(i);
        }
        prev_is_space = c.is_ascii_whitespace();
    }

    let mut headings = Vec::<StdOutTableColumn>::new();
    for (vec_index, byte_offset) in begin_offsets.iter().enumerate() {
        let end_offset: usize =
            if vec_index + 1 < begin_offsets.len() {
                begin_offsets[vec_index + 1]
            } else {
                headings_line.len()
            };
        let name_bytes = &headings_line.as_bytes()[*byte_offset..end_offset];
        let name = String::from(String::from_utf8(Vec::from(name_bytes)).unwrap().trim());

        headings.push(StdOutTableColumn { name, byte_offset: *byte_offset });
    }
    println!("headings: {:?}", headings);
}