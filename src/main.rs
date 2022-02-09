use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use std::{error::Error, io};
use tui::{
    backend::{Backend, CrosstermBackend},
    layout::{Constraint, Direction, Layout},
    widgets::Paragraph,
    Frame, Terminal,
};
use volute::{
    app::{App, Tab},
    input::InputMode,
    ui,
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
            match app.tab {
                Tab::Const => match key.code {
                    KeyCode::Char('q') => {
                        return Ok(());
                    }
                    KeyCode::Char('L') => app.next_tab(),
                    KeyCode::Char('H') => app.previous_tab(),
                    _ => {}
                },
                Tab::Input => match app.input_mode {
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
                        KeyCode::Char(c) => {
                            if ('0'..'a').contains(&c) {
                                app.selected_input_param_mut().push(c);
                            }
                        }
                        KeyCode::Backspace => app.selected_input_param_mut().pop(),
                        _ => {}
                    },
                },
                Tab::Config => match key.code {
                    KeyCode::Char('q') => {
                        return Ok(());
                    }
                    KeyCode::Char('L') => app.next_tab(),
                    KeyCode::Char('H') => app.previous_tab(),
                    _ => {}
                },
            }
        }
    }
}

fn ui<B: Backend>(f: &mut Frame<B>, app: &App) {
    match app.tab {
        Tab::Const => {
            let layout = Layout::default()
                .direction(Direction::Vertical)
                .constraints(ui::constraints(app).as_ref())
                .split(f.size());
            f.render_widget(ui::tabs(app), layout[0]);
            f.render_widget(Paragraph::new("Const Tab"), layout[1]);
        }
        Tab::Input => {
            let layout = Layout::default()
                .direction(Direction::Vertical)
                .constraints(ui::constraints(app).as_ref())
                .split(f.size());
            f.render_widget(ui::tabs(app), layout[0]);
            f.render_widget(ui::footer(app), layout[3]);

            let table_layout = Layout::default()
                .direction(Direction::Horizontal)
                .constraints([Constraint::Percentage(50), Constraint::Percentage(50)].as_ref())
                .split(layout[1]);
            f.render_widget(ui::input_table(app), table_layout[0]);
            f.render_widget(ui::output_table(app), table_layout[1]);
        }
        Tab::Config => {
            let layout = Layout::default()
                .direction(Direction::Vertical)
                .constraints(ui::constraints(app).as_ref())
                .split(f.size());
            f.render_widget(ui::tabs(app), layout[0]);
            f.render_widget(Paragraph::new("Config tab"), layout[1]);
        }
    }
}
