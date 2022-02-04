use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use std::{error::Error, io};
use tui::{
    backend::{Backend, CrosstermBackend},
    layout::{Constraint, Direction, Layout},
    style::{Color, Modifier, Style},
    text::Spans,
    widgets::{self, Block, Borders, Cell, Paragraph, Table, Tabs, Widget},
    Frame, Terminal,
};
use volute::{
    app::App,
    input::{InputMode, InputParam},
};

fn main() -> Result<(), Box<dyn Error>> {
    // setup terminal
    enable_raw_mode()?;
    let mut stdout = io::stdout();
    execute!(stdout, EnterAlternateScreen, EnableMouseCapture)?;
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    // create app and run it
    let app = App::default();
    let res = run_app(&mut terminal, app);

    // restore terminal
    disable_raw_mode()?;
    execute!(
        terminal.backend_mut(),
        LeaveAlternateScreen,
        DisableMouseCapture
    )?;
    terminal.show_cursor()?;

    if let Err(err) = res {
        println!("{:?}", err)
    }

    Ok(())
}

// Input handling
fn run_app<B: Backend>(terminal: &mut Terminal<B>, mut app: App) -> io::Result<()> {
    loop {
        terminal.draw(|f| ui(f, &app))?;

        if let Event::Key(key) = event::read()? {
            match app.tab_index {
                // Input tab
                0 => match app.input_mode {
                    InputMode::Normal => match key.code {
                        KeyCode::Char('q') => {
                            return Ok(());
                        }
                        KeyCode::Char('L') => app.next_tab(),
                        KeyCode::Char('H') => app.previous_tab(),
                        KeyCode::Char('j') => app.next_row(),
                        KeyCode::Char('k') => app.previous_row(),
                        KeyCode::Char('l') => app.next_column(),
                        KeyCode::Char('h') => app.previous_column(),
                        KeyCode::Char('i') => {
                            app.input_mode = InputMode::Insert;
                        }
                        KeyCode::Char('p') => app.insert_row(),
                        KeyCode::Char('d') => app.remove_row(),
                        _ => {}
                    },
                    InputMode::Insert => match key.code {
                        KeyCode::Esc | KeyCode::Enter => {
                            app.input_mode = InputMode::Normal;
                        }
                        KeyCode::Char(c) => match app.selected_column {
                            InputParam::Rpm(_) => app.rows[app.selected_row].rpm.push(c),
                            InputParam::Ve(_) => app.rows[app.selected_row].ve.push(c),
                            InputParam::Map(_) => app.rows[app.selected_row].map.push(c),
                        },
                        KeyCode::Backspace => match app.selected_column {
                            InputParam::Rpm(_) => app.rows[app.selected_row].rpm.pop(),
                            InputParam::Ve(_) => app.rows[app.selected_row].ve.pop(),
                            InputParam::Map(_) => app.rows[app.selected_row].map.pop(),
                        },
                        _ => {}
                    },
                },
                // Config tab
                1 => match key.code {
                    KeyCode::Char('q') => {
                        return Ok(());
                    }
                    KeyCode::Char('L') => app.next_tab(),
                    KeyCode::Char('H') => app.previous_tab(),
                    _ => {}
                },
                _ => unreachable!(),
            }
        }
    }
}

fn ui<B: Backend>(f: &mut Frame<B>, app: &App) {
    let titles = app.tab_titles.iter().map(|t| Spans::from(*t)).collect();
    let tabs = Tabs::new(titles)
        .block(Block::default().borders(Borders::ALL).title("Tabs"))
        .select(app.tab_index)
        .highlight_style(Style::default().add_modifier(Modifier::BOLD));

    match app.tab_index {
        // Input tab
        0 => {
            let layout = Layout::default()
                .direction(Direction::Vertical)
                .constraints(
                    [
                        Constraint::Length(3),
                        Constraint::Min(app.rows.len() as u16 + 2),
                        Constraint::Length(1),
                    ]
                    .as_ref(),
                )
                .split(f.size());
            f.render_widget(tabs, layout[0]);
            f.render_widget(input_table(app), layout[1]);
            f.render_widget(footer(app), layout[2]);
        }
        // Config tab
        1 => {
            let layout = Layout::default()
                .direction(Direction::Vertical)
                .constraints([Constraint::Length(3), Constraint::Min(3)].as_ref())
                .split(f.size());
            f.render_widget(tabs, layout[0]);
            f.render_widget(Paragraph::new("Config tab"), layout[1]);
        }
        _ => unreachable!(),
    }
}

fn input_table(app: &App) -> impl Widget {
    // This is used so I can have named fields instead of indexing a vector of
    // Cells when styling the selected input parameter.
    struct VirtualRow<'a> {
        rpm: Cell<'a>,
        ve: Cell<'a>,
        map: Cell<'a>,
    }

    let mut rows: Vec<VirtualRow> = app
        .rows
        .iter()
        .map(|row| VirtualRow {
            rpm: Cell::from(row.rpm.string()),
            ve: Cell::from(row.ve.string()),
            map: Cell::from(row.map.string()),
        })
        .collect();

    // Highlight the selected parameter
    let selected_parameter = match app.selected_column {
        InputParam::Rpm(_) => &mut rows[app.selected_row].rpm,
        InputParam::Ve(_) => &mut rows[app.selected_row].ve,
        InputParam::Map(_) => &mut rows[app.selected_row].map,
    };
    *selected_parameter = selected_parameter.clone().style(match app.input_mode {
        InputMode::Normal => Style::default().fg(Color::Yellow),
        InputMode::Insert => Style::default()
            .fg(Color::Yellow)
            .add_modifier(Modifier::ITALIC),
    });

    Table::new(
        rows.iter()
            .map(|row| widgets::Row::new(vec![row.rpm.clone(), row.ve.clone(), row.map.clone()]))
            .collect::<Vec<widgets::Row>>(),
    )
    .header(widgets::Row::new(vec!["rpm", "ve", "map"]))
    .block(Block::default().borders(Borders::ALL).title("inputs"))
    .widths(&[
        Constraint::Length(5),
        Constraint::Length(3),
        Constraint::Length(3),
    ])
}

fn footer(app: &App) -> impl Widget {
    match app.input_mode {
        InputMode::Normal => {
            Paragraph::new("Normal").style(Style::default().fg(Color::Black).bg(Color::Yellow))
        }
        InputMode::Insert => {
            Paragraph::new("Insert").style(Style::default().fg(Color::Black).bg(Color::Blue))
        }
    }
}
