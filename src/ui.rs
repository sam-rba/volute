use crate::{
    app::App,
    unit_of_measurement::{
        pressure::{self, Pressure},
        UnitOfMeasurement,
    },
};
use std::ptr;
use tui::{
    layout::Constraint,
    style::{Color, Style},
    widgets::{self, Block, Borders, Cell, Paragraph, Table, Widget},
};

pub fn input_table(app: &App) -> impl Widget {
    let rows = app.rows().iter().map(|row| {
        let cells = row.iter().map(|item| {
            if ptr::eq(item, app.selected_input_param()) {
                Cell::from(item.string()).style(Style::default().fg(Color::Yellow))
            } else {
                Cell::from(item.string())
            }
        });
        widgets::Row::new(cells)
    });

    Table::new(rows)
        .header(widgets::Row::new(vec!["rpm", "ve", "map"]))
        .block(Block::default().borders(Borders::ALL).title("inputs"))
        .widths(&[
            Constraint::Length(5), // rpm
            Constraint::Length(3), // ve
            Constraint::Length(3), // map
        ])
}

pub fn output_table(app: &App) -> impl Widget {
    let map = match app.rows()[0].map.string().parse::<i32>() {
        Ok(p) => Pressure::from_unit(pressure::Unit::KiloPascal, p),
        Err(_) => Pressure::default(),
    };
    Paragraph::new(map.as_unit(pressure::Unit::KiloPascal).to_string())
        .block(Block::default().title("map").borders(Borders::ALL))
}
