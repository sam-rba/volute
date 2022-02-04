use crate::input::{InputMode, InputParam, Row};

pub const INPUT_TAB_INDEX: usize = 0;
pub const CONFIG_TAB_INDEX: usize = 1;

pub struct App {
    pub tab_index: usize,
    tab_titles: [&'static str; 2],

    pub rows: Vec<Row>,
    pub selected_row: usize,
    pub selected_column: InputParam,

    pub input_mode: InputMode,
}

impl App {
    pub fn tab_titles(&self) -> &[&str] {
        &self.tab_titles
    }

    pub fn next_tab(&mut self) {
        self.tab_index = (self.tab_index + 1) % self.tab_titles.len();
    }

    pub fn previous_tab(&mut self) {
        if self.tab_index > 0 {
            self.tab_index -= 1;
        } else {
            self.tab_index = self.tab_titles.len() - 1;
        }
    }

    pub fn next_row(&mut self) {
        if self.selected_row < self.rows.len() - 1 {
            self.selected_row += 1;
        } else {
            self.selected_row = 0;
        }
    }

    pub fn previous_row(&mut self) {
        if self.selected_row > 0 {
            self.selected_row -= 1;
        } else {
            self.selected_row = self.rows.len() - 1;
        }
    }

    pub fn next_column(&mut self) {
        self.selected_column = self.selected_column.next();
    }

    pub fn previous_column(&mut self) {
        self.selected_column = self.selected_column.previous();
    }

    pub fn insert_row(&mut self) {
        let index = self.selected_row;
        self.rows.insert(index, self.rows[index].clone());
    }

    pub fn remove_row(&mut self) {
        if self.rows.len() > 1 {
            self.rows.remove(self.selected_row);
            if self.selected_row > 0 {
                self.selected_row -= 1;
            }
        }
    }
}

impl Default for App {
    fn default() -> App {
        App {
            tab_index: INPUT_TAB_INDEX,
            tab_titles: ["Input", "Config"],
            rows: vec![Row::default()],
            selected_row: 0,
            selected_column: InputParam::Rpm(String::new()),
            input_mode: InputMode::Normal,
        }
    }
}
