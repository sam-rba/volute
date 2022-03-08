use crate::{
    input::{InputParam, Row},
    unit_of_measurement::pressure,
};

pub struct App {
    rows: Vec<Row>,
    selected_row: usize,
    selected_column: InputParam,

    pub pressure_unit: pressure::Unit,
}

impl App {
    pub fn rows(&self) -> &Vec<Row> {
        &self.rows
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
            rows: vec![Row::default()],
            selected_row: 0,
            selected_column: InputParam::Rpm(0),
            pressure_unit: pressure::Unit::KiloPascal,
        }
    }
}
