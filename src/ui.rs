use crate::{
    app::{App, Tab},
    input::InputMode,
    unit_of_measurement::pressure::*,
};
use std::ptr;
use tui::{
    layout::Constraint,
    style::{Color, Modifier, Style},
    text::Spans,
    widgets::{self, Block, Borders, Cell, Paragraph, Table, Tabs, Widget},
};

pub fn constraints(app: &App) -> Vec<Constraint> {
    match app.tab {
        Tab::Const | Tab::Config => {
            vec![Constraint::Length(3), Constraint::Length(3)]
        }
        Tab::Input => {
            vec![
                Constraint::Length(3),                           // Tabs
                Constraint::Length(app.rows().len() as u16 + 3), // tables
                Constraint::Max(100),                            // Spacer
                Constraint::Length(1),                           // Footer
            ]
        }
    }
}

pub fn tabs(app: &App) -> impl Widget + '_ {
    let titles = app
        .tab_titles()
        .iter()
        .map(|t| Spans::from(t.as_str()))
        .collect();
    Tabs::new(titles)
        .block(Block::default().borders(Borders::ALL).title("Tabs"))
        .select(app.tab as usize)
        .highlight_style(
            Style::default()
                .fg(Color::Yellow)
                .add_modifier(Modifier::BOLD),
        )
}

pub fn input_table(app: &App) -> impl Widget {
    let rows = app.rows().iter().map(|row| {
        let cells = row.iter().map(|item| {
            if ptr::eq(item, app.selected_input_param()) {
                Cell::from(item.string()).style(match app.input_mode {
                    InputMode::Normal => Style::default().fg(Color::Yellow),
                    InputMode::Insert => Style::default()
                        .fg(Color::Blue)
                        .add_modifier(Modifier::ITALIC),
                })
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
        Ok(p) => Pressure::from_unit(PressureUnit::KiloPascal, p),
        Err(_) => Pressure::default(),
    };
    Paragraph::new(map.as_unit(PressureUnit::KiloPascal).to_string())
        .block(Block::default().title("map").borders(Borders::ALL))
}

pub fn footer(app: &App) -> impl Widget {
    match app.input_mode {
        InputMode::Normal => {
            Paragraph::new("Normal").style(Style::default().fg(Color::Black).bg(Color::Yellow))
        }
        InputMode::Insert => {
            Paragraph::new("Insert").style(Style::default().fg(Color::Black).bg(Color::Blue))
        }
    }
}
