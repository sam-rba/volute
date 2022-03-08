use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use std::{error::Error, io};
use tui::{
    backend::{Backend, CrosstermBackend},
    layout::{Constraint, Direction, Layout},
    Frame, Terminal,
};
use volute::{
    app::App,
    input::InputParam,
    ui,
    unit_of_measurement::{pressure::Pressure, UnitOfMeasurement},
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
            match key.code {
                KeyCode::Char('q') => {
                    return Ok(());
                }
                KeyCode::Char('j') => app.next_row(),
                KeyCode::Char('k') => app.previous_row(),
                KeyCode::Char('l') => app.next_column(),
                KeyCode::Char('h') => app.previous_column(),
                KeyCode::Char('y') => app.insert_row(),
                KeyCode::Char('d') => app.remove_row(),
                KeyCode::Char('p') => app.pressure_unit.next(),
                KeyCode::Char(c) => {
                    if ('0'..':').contains(&c) {
                        // 0 to 9 inclusive
                        let digit = c.to_digit(10).unwrap();
                        match app.selected_input_param() {
                            InputParam::Rpm(rpm) => {
                                *app.selected_input_param_mut() =
                                    InputParam::Rpm(*rpm * 10 + digit);
                            }
                            InputParam::Ve(ve) => {
                                *app.selected_input_param_mut() = InputParam::Ve(*ve * 10 + digit);
                            }
                            InputParam::Map(p) => {
                                *app.selected_input_param_mut() =
                                    InputParam::Map(Pressure::from_unit(
                                        app.pressure_unit,
                                        p.as_unit(app.pressure_unit) * 10 + digit as i32,
                                    ))
                            }
                        }
                    }
                }
                KeyCode::Backspace => match app.selected_input_param() {
                    InputParam::Rpm(rpm) => {
                        *app.selected_input_param_mut() = InputParam::Rpm(*rpm / 10);
                    }
                    InputParam::Ve(ve) => {
                        *app.selected_input_param_mut() = InputParam::Ve(*ve / 10);
                    }
                    InputParam::Map(p) => {
                        *app.selected_input_param_mut() = InputParam::Map(Pressure::from_unit(
                            app.pressure_unit,
                            p.as_unit(app.pressure_unit) / 10,
                        ))
                    }
                },
                _ => {}
            }
        }
    }
}

fn ui<B: Backend>(f: &mut Frame<B>, app: &App) {
    let layout = Layout::default()
        .direction(Direction::Horizontal)
        .constraints([Constraint::Percentage(50), Constraint::Percentage(50)].as_ref())
        .split(f.size());
    f.render_widget(ui::input_table(app), layout[0]);
    f.render_widget(ui::output_table(app), layout[1]);
}
