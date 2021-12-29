
pub fn parse_stdout_table(lines: &Vec<String>) {

    let headings_line = &lines[0]; // .split_whitespace().collect::<Vec<&str>>()
    println!("head: {:?}", headings_line);

    let mut begin_indices: Vec<usize> = Vec::<usize>::new();
    let mut prev_is_space = true;

    for (i, c) in headings_line.chars().enumerate() {

        if prev_is_space && !c.is_whitespace() {
            begin_indices.push(i);
        }
        prev_is_space = c.is_whitespace();

    }

    println!("beginIndices: {:?}", begin_indices);
}