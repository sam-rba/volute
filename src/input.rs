// A row in the inputs table has one of each variation.
#[derive(Clone)]
pub enum InputParam {
    Rpm(String), // Revolutions per minute
    Ve(String),  // Volumetric efficiency
    Map(String), // Manifold absolute pressure
}

impl InputParam {
    /* Acts like the push() method of a Vec.
     * Appends the given char to the end of the string contained by the
     * InputParam.
     */
    pub fn push(&mut self, c: char) {
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
    pub fn pop(&mut self) {
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
    pub fn string(&self) -> String {
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
    pub fn next(&self) -> Self {
        match self {
            Self::Rpm(_) => Self::Ve(String::new()),
            Self::Ve(_) => Self::Map(String::new()),
            Self::Map(_) => Self::Rpm(String::new()),
        }
    }

    pub fn previous(&self) -> Self {
        match self {
            Self::Rpm(_) => Self::Map(String::new()),
            Self::Ve(_) => Self::Rpm(String::new()),
            Self::Map(_) => Self::Ve(String::new()),
        }
    }
}

// A row in the inputs table. Contains one of each variation of InputParam.
#[derive(Clone)]
pub struct Row {
    pub rpm: InputParam,
    pub ve: InputParam,
    pub map: InputParam,
}

impl Row {
    pub fn iter(&self) -> RowIter {
        RowIter::from_row(&self)
    }
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

pub struct RowIter<'a> {
    row: &'a Row,
    iter_state: Option<InputParam>,
}

impl<'a> RowIter<'a> {
    fn from_row(row: &'a Row) -> Self {
        Self {
            row: row,
            iter_state: Some(InputParam::Rpm(String::new())),
        }
    }
}

impl<'a> Iterator for RowIter<'a> {
    type Item = &'a InputParam;

    fn next(&mut self) -> Option<Self::Item> {
        match self.iter_state {
            Some(InputParam::Rpm(_)) => {
                self.iter_state = Some(InputParam::Ve(String::new()));
                Some(&self.row.rpm)
            }
            Some(InputParam::Ve(_)) => {
                self.iter_state = Some(InputParam::Map(String::new()));
                Some(&self.row.ve)
            }
            Some(InputParam::Map(_)) => {
                self.iter_state = None;
                Some(&self.row.map)
            }
            None => None,
        }
    }
}
