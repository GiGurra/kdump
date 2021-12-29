#[derive(Debug)]
pub struct StdOutTableColumn {
    name: String,
    byte_offset: usize,
}

pub fn parse_stdout_table(lines: &Vec<String>) {
    let headings_line = &lines[0]; // .split_whitespace().collect::<Vec<&str>>()
    println!("head: {:?}", headings_line);

    let headings_offsets = find_headings_offsets(headings_line);
    let headings = find_headings(headings_line, &headings_offsets);

    let mut lineValues = Vec::<std::collections::HashMap<String, String>>::new();

    println!("headings: {:?}", headings);
}

fn find_headings_offsets(headings_line: &String) -> Vec<usize> {
    let mut begin_offsets = Vec::<usize>::new();
    let mut prev_is_space = true;

    for (i, c) in headings_line.bytes().enumerate() {
        if prev_is_space && !c.is_ascii_whitespace() {
            begin_offsets.push(i);
        }
        prev_is_space = c.is_ascii_whitespace();
    }

    return begin_offsets;
}

fn find_headings(headings_line: &String, headings_offsets: &Vec<usize>) -> Vec<StdOutTableColumn> {

    let mut headings = Vec::<StdOutTableColumn>::new();
    for (vec_index, byte_offset) in headings_offsets.iter().enumerate() {
        let end_offset: usize =
            if vec_index + 1 < headings_offsets.len() {
                headings_offsets[vec_index + 1]
            } else {
                headings_line.len()
            };
        let name_bytes = &headings_line.as_bytes()[*byte_offset..end_offset];
        let name = String::from(String::from_utf8(Vec::from(name_bytes)).unwrap().trim());

        headings.push(StdOutTableColumn { name, byte_offset: *byte_offset });
    }
    return headings;
}