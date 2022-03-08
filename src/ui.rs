use crate::{app::App, input::InputParam, unit_of_measurement::UnitOfMeasurement};
use std::ptr;
use tui::{
    layout::Constraint,
    style::{Color, Style},
    widgets::{self, Block, Borders, Cell, Paragraph, Table, Widget},
};

pub fn input_table(app: &App) -> impl Widget {
    let rows = app.rows().iter().map(|row| {
        let cells = row.iter().map(|item| {
            let item_str = match item {
                InputParam::Rpm(rpm) => rpm.to_string(),
                InputParam::Ve(ve) => ve.to_string(),
                InputParam::Map(p) => p.as_unit(app.pressure_unit).to_string(),
            };
            if ptr::eq(item, app.selected_input_param()) {
                Cell::from(item_str).style(Style::default().fg(Color::Yellow))
            } else {
                Cell::from(item_str)
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
    if let InputParam::Map(p) = &app.rows()[0].map {
        Paragraph::new(p.as_unit(app.pressure_unit).to_string())
            .block(Block::default().title("map").borders(Borders::ALL))
    } else {
        Paragraph::new("err").block(Block::default().title("map").borders(Borders::ALL))
    }
}
