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

// Intended to emulate vim modes.
enum InputMode {
    Normal, // Navigating the ui.
    Insert, // Editing a parameter.
}

// A row in the inputs table has one of each variation.
#[derive(Clone)]
enum InputParam {
    Rpm(String), // Revolutions per minute
    Ve(String),  // Volumetric efficiency
    Map(String), // Manifold absolute pressure
}
impl InputParam {
    /* Acts like the push() method of a Vec.
     * Appends the given char to the end of the string contained by the
     * InputParam.
     */
    fn push(&mut self, c: char) {
        match self {
            Self::Rpm(rpm) => {
                rpm.push(c);
                *self = Self::Rpm(rpm.to_string());
            }
            Self::Ve(ve) => {
                ve.push(c);
                *self = Self::Ve(ve.to_string());
            }
            Self::Map(map) => {
                map.push(c);
                *self = Self::Map(map.to_string());
            }
        }
    }
    /* Acts like the pop() method of a Vec.
     * Removes the last char from the string contained by the InputParam.
     */
    fn pop(&mut self) {
        match self {
            Self::Rpm(rpm) => {
                rpm.pop();
                *self = Self::Rpm(rpm.to_string());
            }
            Self::Ve(ve) => {
                ve.pop();
                *self = Self::Rpm(ve.to_string());
            }
            Self::Map(map) => {
                map.pop();
                *self = Self::Map(map.to_string());
            }
        }
    }
    // Return a copy of the string contained by the InputParam.
    fn string(&self) -> String {
        match self {
            Self::Rpm(rpm) => rpm.to_string(),
            Self::Ve(ve) => ve.to_string(),
            Self::Map(map) => map.to_string(),
        }
    }
    /* next() and previous() allow InputParam to act as a circular iterator of
     * sorts. next() will return the next variation as they are defined. When
     * it reaches the end, the first variation will be returned:
     *     RPM->VE->MAP->RPM->etc...
     * previous() simply goes the opposite direction:
     *     MAP->VE->RPM->MAP->etc...
     */
    fn next(&self) -> Self {
        match self {
            Self::Rpm(_) => Self::Ve(String::new()),
            Self::Ve(_) => Self::Map(String::new()),
            Self::Map(_) => Self::Rpm(String::new()),
        }
    }
    fn previous(&self) -> Self {
        match self {
            Self::Rpm(_) => Self::Map(String::new()),
            Self::Ve(_) => Self::Rpm(String::new()),
            Self::Map(_) => Self::Ve(String::new()),
        }
    }
}

// A row in the inputs table. Contains one of each variation of InputParam.
#[derive(Clone)]
struct Row {
    rpm: InputParam,
    ve: InputParam,
    map: InputParam,
}
impl Default for Row {
    fn default() -> Self {
        Self {
            rpm: InputParam::Rpm(String::from("7000")),
            ve: InputParam::Ve(String::from("95")),
            map: InputParam::Map(String::from("200")),
        }
    }
}

// Holds the state of the application.
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
            selected_column: InputParam::Rpm(String::new()),
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

// Input handling
fn run_app<B: Backend>(terminal: &mut Terminal<B>, mut app: App) -> io::Result<()> {
    loop {
        terminal.draw(|f| ui(f, &app))?;

        if let Event::Key(key) = event::read()? {
            match app.input_mode {
                // Navigating
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
                    // Append a caracter to the currently selected parameter.
                    KeyCode::Char(c) => match app.selected_column {
                        InputParam::Rpm(_) => {
                            app.rows[app.selected_row].rpm.push(c);
                        }
                        InputParam::Ve(_) => {
                            app.rows[app.selected_row].ve.push(c);
                        }
                        InputParam::Map(_) => {
                            app.rows[app.selected_row].map.push(c);
                        }
                    },
                    // Remove a character from the currently selected parameter.
                    KeyCode::Backspace => match app.selected_column {
                        InputParam::Rpm(_) => {
                            app.rows[app.selected_row].rpm.pop();
                        }
                        InputParam::Ve(_) => {
                            app.rows[app.selected_row].ve.pop();
                        }
                        InputParam::Map(_) => {
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

    /* This is used so I can have named fields instead of indexing a vector of
     * Cells when styling the selected input parameter.
     */
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
