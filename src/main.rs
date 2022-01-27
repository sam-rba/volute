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
    widgets::{self, Block, Borders, Cell, Paragraph, Table},
    Frame, Terminal,
};

enum InputMode {
    Normal,
    Insert,
}

#[derive(Clone)]
enum InputParam {
    RPM(String), // Revolutions per minute
    VE(String),  // Volumetric efficiency
    MAP(String), // Manifold absolute pressure
}
impl InputParam {
    fn push(&mut self, c: char) {
        match self {
            Self::RPM(rpm) => {
                rpm.push(c);
                *self = Self::RPM(rpm.to_string());
            }
            Self::VE(ve) => {
                ve.push(c);
                *self = Self::VE(ve.to_string());
            }
            Self::MAP(map) => {
                map.push(c);
                *self = Self::MAP(map.to_string());
            }
        }
    }
    fn pop(&mut self) {
        match self {
            Self::RPM(rpm) => {
                rpm.pop();
                *self = Self::RPM(rpm.to_string());
            }
            Self::VE(ve) => {
                ve.pop();
                *self = Self::RPM(ve.to_string());
            }
            Self::MAP(map) => {
                map.pop();
                *self = Self::MAP(map.to_string());
            }
        }
    }
    fn string(&self) -> String {
        match self {
            Self::RPM(rpm) => rpm.to_string(),
            Self::VE(ve) => ve.to_string(),
            Self::MAP(map) => map.to_string(),
        }
    }
    fn next(&self) -> Self {
        match self {
            Self::RPM(_) => Self::VE(String::new()),
            Self::VE(_) => Self::MAP(String::new()),
            Self::MAP(_) => Self::RPM(String::new()),
        }
    }
    fn previous(&self) -> Self {
        match self {
            Self::RPM(_) => Self::MAP(String::new()),
            Self::VE(_) => Self::RPM(String::new()),
            Self::MAP(_) => Self::VE(String::new()),
        }
    }
}

#[derive(Clone)]
struct Row {
    rpm: InputParam,
    ve: InputParam,
    map: InputParam,
}
impl Default for Row {
    fn default() -> Self {
        Self {
            rpm: InputParam::RPM(String::from("7000")),
            ve: InputParam::VE(String::from("95")),
            map: InputParam::MAP(String::from("200")),
        }
    }
}

struct App {
    rows: Vec<Row>,

    selected_row: usize,
    selected_column: InputParam,

    input_mode: InputMode,
}
impl Default for App {
    fn default() -> App {
        App {
            rows: vec![Row::default()],
            selected_row: 0,
            selected_column: InputParam::RPM(String::new()),
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
                    // Enter insert mode
                    KeyCode::Char('i') => {
                        app.input_mode = InputMode::Insert;
                    }
                    // Quit
                    KeyCode::Char('q') => {
                        return Ok(());
                    }
                    // Navigate up
                    KeyCode::Char('k') => {
                        if app.selected_row > 0 {
                            app.selected_row -= 1;
                        } else {
                            app.selected_row = app.rows.len() - 1;
                        }
                    }
                    // Navigate down
                    KeyCode::Char('j') => {
                        if app.selected_row < app.rows.len() - 1 {
                            app.selected_row += 1;
                        } else {
                            app.selected_row = 0;
                        }
                    }
                    // Navigate right
                    KeyCode::Char('l') => {
                        app.selected_column = app.selected_column.next();
                    }
                    // Navigate left
                    KeyCode::Char('h') => {
                        app.selected_column = app.selected_column.previous();
                    }
                    // Add row
                    KeyCode::Char('p') => {
                        app.rows
                            .insert(app.selected_row, app.rows[app.selected_row].clone());
                    }
                    // Remove row
                    KeyCode::Char('d') => {
                        if app.rows.len() > 1 {
                            app.rows.remove(app.selected_row);
                            if app.selected_row > 0 {
                                app.selected_row -= 1;
                            }
                        }
                    }
                    _ => {}
                },
                InputMode::Insert => match key.code {
                    // Exit insert mode
                    KeyCode::Esc => {
                        app.input_mode = InputMode::Normal;
                    }
                    // Exit insert mode
                    KeyCode::Enter => {
                        app.input_mode = InputMode::Normal;
                    }
                    KeyCode::Char(c) => match app.selected_column {
                        InputParam::RPM(_) => {
                            app.rows[app.selected_row].rpm.push(c);
                        }
                        InputParam::VE(_) => {
                            app.rows[app.selected_row].ve.push(c);
                        }
                        InputParam::MAP(_) => {
                            app.rows[app.selected_row].map.push(c);
                        }
                    },
                    KeyCode::Backspace => match app.selected_column {
                        InputParam::RPM(_) => {
                            app.rows[app.selected_row].rpm.pop();
                        }
                        InputParam::VE(_) => {
                            app.rows[app.selected_row].ve.pop();
                        }
                        InputParam::MAP(_) => {
                            app.rows[app.selected_row].map.pop();
                        }
                    },

                    _ => {}
                },
            }
        }
    }
}

fn ui<B: Backend>(f: &mut Frame<B>, app: &App) {
    let layout = Layout::default()
        .direction(Direction::Vertical)
        .constraints(
            [
                Constraint::Min(app.rows.len() as u16 + 2),
                Constraint::Length(1),
            ]
            .as_ref(),
        )
        .split(f.size());

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
        InputParam::RPM(_) => &mut rows[app.selected_row].rpm,
        InputParam::VE(_) => &mut rows[app.selected_row].ve,
        InputParam::MAP(_) => &mut rows[app.selected_row].map,
    };
    *selected_parameter = selected_parameter.clone().style(match app.input_mode {
        InputMode::Normal => Style::default().fg(Color::Yellow),
        InputMode::Insert => Style::default()
            .fg(Color::Yellow)
            .add_modifier(Modifier::ITALIC),
    });

    let table = Table::new(
        rows.iter()
            .map(|row| widgets::Row::new(vec![row.rpm.clone(), row.ve.clone(), row.map.clone()]))
            .collect::<Vec<widgets::Row>>(),
    )
    .header(widgets::Row::new(vec!["rpm", "ve", "map"]))
    .block(Block::default().borders(Borders::ALL).title("Table"))
    .widths(&[
        Constraint::Length(5),
        Constraint::Length(3),
        Constraint::Length(3),
    ]);
    f.render_widget(table, layout[0]);

    let footer = match app.input_mode {
        InputMode::Normal => {
            Paragraph::new("Normal").style(Style::default().fg(Color::Black).bg(Color::Yellow))
        }
        InputMode::Insert => {
            Paragraph::new("Insert").style(Style::default().fg(Color::Black).bg(Color::Blue))
        }
    };
    f.render_widget(footer, layout[1]);
}
