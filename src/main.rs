use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use std::{error::Error, io};
use tui::{
    backend::{Backend, CrosstermBackend},
    layout::{Constraint, Direction, Layout},
    style::{Color, Style},
    widgets::{Block, Borders, Paragraph},
    Frame, Terminal,
};
use unicode_width::UnicodeWidthStr;

enum InputMode {
    Normal,
    Insert,
}

enum Row {
    Rpm(String),
    Ve(String),
    Map(String),
}

impl Row {
    fn push(&mut self, c: char) {
        match self {
            Row::Rpm(rpm) => {
                rpm.push(c);
                *self = Row::Rpm(rpm.to_string());
            }
            Row::Ve(ve) => {
                ve.push(c);
                *self = Row::Ve(ve.to_string());
            }
            Row::Map(map) => {
                map.push(c);
                *self = Row::Map(map.to_string());
            }
        }
    }

    fn pop(&mut self) {
        match self {
            Row::Rpm(rpm) => {
                rpm.pop();
                *self = Row::Rpm(rpm.to_string());
            }
            Row::Ve(ve) => {
                ve.pop();
                *self = Row::Rpm(ve.to_string());
            }
            Row::Map(map) => {
                map.pop();
                *self = Row::Map(map.to_string());
            }
        }
    }

    fn string(&self) -> String {
        match self {
            Row::Rpm(rpm) => rpm.to_string(),
            Row::Ve(ve) => ve.to_string(),
            Row::Map(map) => map.to_string(),
        }
    }

    fn next(&self) -> Self {
        match self {
            Row::Rpm(_) => Row::Ve(String::new()),
            Row::Ve(_) => Row::Map(String::new()),
            Row::Map(_) => Row::Rpm(String::new()),
        }
    }

    fn previous(&self) -> Self {
        match self {
            Row::Rpm(_) => Row::Map(String::new()),
            Row::Ve(_) => Row::Rpm(String::new()),
            Row::Map(_) => Row::Ve(String::new()),
        }
    }
}

struct Column {
    rpm: Row,
    ve: Row,
    map: Row,
}

/// App holds the state of the application
struct App {
    column: Column,

    selected_row: Row,

    input_mode: InputMode,
}

impl Default for App {
    fn default() -> App {
        App {
            column: Column {
                rpm: Row::Rpm(String::from("7000")),
                ve: Row::Ve(String::from("95")),
                map: Row::Map(String::from("150")),
            },
            selected_row: Row::Rpm(String::new()),
            input_mode: InputMode::Normal,
        }
    }
}

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

fn run_app<B: Backend>(terminal: &mut Terminal<B>, mut app: App) -> io::Result<()> {
    loop {
        terminal.draw(|f| ui(f, &app))?;

        if let Event::Key(key) = event::read()? {
            match app.input_mode {
                InputMode::Normal => match key.code {
                    KeyCode::Char('i') => {
                        app.input_mode = InputMode::Insert;
                    }
                    KeyCode::Char('q') => {
                        return Ok(());
                    }
                    KeyCode::Char('j') => {
                        app.selected_row = app.selected_row.next();
                    }
                    KeyCode::Char('k') => {
                        app.selected_row = app.selected_row.previous();
                    }
                    _ => {}
                },
                InputMode::Insert => match key.code {
                    KeyCode::Enter => {
                        app.input_mode = InputMode::Normal;
                    }
                    KeyCode::Char(c) => match app.selected_row {
                        Row::Rpm(_) => {
                            app.column.rpm.push(c);
                        }
                        Row::Ve(_) => {
                            app.column.ve.push(c);
                        }
                        Row::Map(_) => {
                            app.column.map.push(c);
                        }
                    },
                    KeyCode::Backspace => match app.selected_row {
                        Row::Rpm(_) => {
                            app.column.rpm.pop();
                        }
                        Row::Ve(_) => {
                            app.column.ve.pop();
                        }
                        Row::Map(_) => {
                            app.column.map.pop();
                        }
                    },
                    KeyCode::Esc => {
                        app.input_mode = InputMode::Normal;
                    }
                    _ => {}
                },
            }
        }
    }
}

fn ui<B: Backend>(f: &mut Frame<B>, app: &App) {
    let chunks = Layout::default()
        .direction(Direction::Vertical)
        .margin(2)
        .constraints([Constraint::Length(3), Constraint::Length(3), Constraint::Length(3)].as_ref())
        .split(f.size());

    let rpm = Paragraph::new(app.column.rpm.string())
        .style(match app.selected_row {
            Row::Rpm(_) => Style::default().fg(Color::Yellow),
            _ => Style::default(),
        })
        .block(Block::default().borders(Borders::ALL).title("rpm"));
    f.render_widget(rpm, chunks[0]);

    let ve = Paragraph::new(app.column.ve.string())
        .style(match app.selected_row {
            Row::Ve(_) => Style::default().fg(Color::Yellow),
            _ => Style::default(),
        })
        .block(Block::default().borders(Borders::ALL).title("ve"));
    f.render_widget(ve, chunks[1]);

    let map = Paragraph::new(app.column.map.string())
        .style(match app.selected_row {
            Row::Map(_) => Style::default().fg(Color::Yellow),
            _ => Style::default(),
        })
        .block(Block::default().borders(Borders::ALL).title("map"));
    f.render_widget(map, chunks[2]);
}
