use std::collections::HashMap;

#[derive(Debug)]
pub struct StdOutTableColumn {
    name: String,
    byte_offset: usize,
}

pub fn parse_stdout_table(lines: &Vec<String>) -> Vec<HashMap<String, String>> {
    let headings_line = &lines[0];
    let data_lines = Vec::from(&lines[1..]);
    let headings_offsets = find_headings_offsets(headings_line);
    let headings = find_headings(headings_line, &headings_offsets);
    let line_values = find_line_values(&data_lines, headings);
    return line_values;
}

fn find_line_values(data_lines: &Vec<String>, headings: Vec<StdOutTableColumn>) -> Vec<HashMap<String, String>> {
    let mut line_values = Vec::<HashMap<String, String>>::new();
    for data_line in data_lines {
        let mut line_value = HashMap::<String, String>::new();
        for (i_heading, heading) in headings.iter().enumerate() {
            let end_index: usize =
                if i_heading + 1 < headings.len() {
                    headings[i_heading + 1].byte_offset
                } else {
                    data_line.len()
                };
            let str_value = &data_line[heading.byte_offset..end_index];
            let heading_name = String::from(&heading.name);
            line_value.insert(heading_name, String::from(str_value.trim()));
        }
        line_values.push(line_value);
    }
    line_values
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