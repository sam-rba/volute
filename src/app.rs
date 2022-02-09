use crate::input::{InputMode, InputParam, Row};

#[derive(Copy, Clone)]
pub enum Tab {
    Const = 0,
    Input = 1,
    Config = 2,
}

impl Tab {
    fn next(&self) -> Self {
        match self {
            Self::Const => Self::Input,
            Self::Input => Self::Config,
            Self::Config => Self::Const,
        }
    }

    fn previous(&self) -> Self {
        match self {
            Self::Const => Self::Config,
            Self::Input => Self::Const,
            Self::Config => Self::Input,
        }
    }

    fn string(&self) -> String {
        match self {
            Self::Const => "Const".to_string(),
            Self::Input => "Input".to_string(),
            Self::Config => "Config".to_string(),
        }
    }
}

impl IntoIterator for Tab {
    type Item = Tab;
    type IntoIter = TabIter;

    fn into_iter(self) -> Self::IntoIter {
        TabIter { tab: Some(self) }
    }
}

pub struct TabIter {
    tab: Option<Tab>,
}

impl Iterator for TabIter {
    type Item = Tab;

    fn next(&mut self) -> Option<Self::Item> {
        match self.tab {
            Some(Tab::Const) => {
                self.tab = Some(Tab::Input);
                Some(Tab::Const)
            }
            Some(Tab::Input) => {
                self.tab = Some(Tab::Config);
                Some(Tab::Input)
            }
            Some(Tab::Config) => {
                self.tab = None;
                Some(Tab::Config)
            }
            None => None,
        }
    }
}

pub struct App {
    pub tab: Tab,
    tab_titles: Vec<String>,

    rows: Vec<Row>,
    selected_row: usize,
    selected_column: InputParam,

    pub input_mode: InputMode,
}

impl App {
    pub fn rows(&self) -> &Vec<Row> {
        &self.rows
    }

    pub fn tab_titles(&self) -> &Vec<String> {
        &self.tab_titles
    }

    pub fn next_tab(&mut self) {
        self.tab = self.tab.next();
    }

    pub fn previous_tab(&mut self) {
        self.tab = self.tab.previous();
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
            // If we remove the last row, the selected row will be out of range.
            if self.selected_row >= self.rows.len() {
                self.selected_row = self.rows.len() - 1;
            }
        }
    }

    pub fn selected_input_param(&self) -> &InputParam {
        match self.selected_column {
            InputParam::Rpm(_) => &self.rows[self.selected_row].rpm,
            InputParam::Ve(_) => &self.rows[self.selected_row].ve,
            InputParam::Map(_) => &self.rows[self.selected_row].map,
        }
    }

    pub fn selected_input_param_mut(&mut self) -> &mut InputParam {
        match self.selected_column {
            InputParam::Rpm(_) => &mut self.rows[self.selected_row].rpm,
            InputParam::Ve(_) => &mut self.rows[self.selected_row].ve,
            InputParam::Map(_) => &mut self.rows[self.selected_row].map,
        }
    }
}

impl Default for App {
    fn default() -> App {
        App {
            tab: Tab::Const,
            tab_titles: Tab::Const.into_iter().map(|t| t.string()).collect(),
            rows: vec![Row::default()],
            selected_row: 0,
            selected_column: InputParam::Rpm(String::new()),
            input_mode: InputMode::Normal,
        }
    }
}
