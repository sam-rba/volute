use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use std::{error::Error, io, ptr};
use tui::{
    backend::{Backend, CrosstermBackend},
    layout::{Constraint, Direction, Layout},
    style::{Color, Modifier, Style},
    text::Spans,
    widgets::{self, Block, Borders, Cell, Paragraph, Table, Tabs, Widget},
    Frame, Terminal,
};
use volute::{
    app::{App, CONFIG_TAB_INDEX, INPUT_TAB_INDEX},
    input::InputMode,
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
                INPUT_TAB_INDEX => match app.input_mode {
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
                        KeyCode::Char(c) => app.selected_input_param_mut().push(c),
                        KeyCode::Backspace => app.selected_input_param_mut().pop(),
                        _ => {}
                    },
                },
                CONFIG_TAB_INDEX => match key.code {
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
    let titles = app.tab_titles().iter().map(|t| Spans::from(*t)).collect();
    let tabs = Tabs::new(titles)
        .block(Block::default().borders(Borders::ALL).title("Tabs"))
        .select(app.tab_index)
        .highlight_style(
            Style::default()
                .fg(Color::Yellow)
                .add_modifier(Modifier::BOLD),
        );

    match app.tab_index {
        INPUT_TAB_INDEX => {
            let layout = Layout::default()
                .direction(Direction::Vertical)
                .constraints(
                    [
                        Constraint::Length(3),                           // Tabs
                        Constraint::Length(app.rows().len() as u16 + 3), // Input table
                        Constraint::Max(100),                            // Spacer
                        Constraint::Length(1),                           // Footer
                    ]
                    .as_ref(),
                )
                .split(f.size());
            f.render_widget(tabs, layout[0]);
            f.render_widget(input_table(app), layout[1]);
            f.render_widget(footer(app), layout[3]);
        }
        CONFIG_TAB_INDEX => {
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
