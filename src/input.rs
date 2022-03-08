use crate::unit_of_measurement::{
    pressure::{Pressure, Unit::KiloPascal},
    UnitOfMeasurement,
};

// A row in the inputs table has one of each variation.
#[derive(Clone)]
pub enum InputParam {
    Rpm(u32),      // Revolutions per minute
    Ve(u32),       // Volumetric efficiency
    Map(Pressure), // Manifold absolute pressure
}

impl InputParam {
    /* next() and previous() allow InputParam to act as a circular iterator of
     * sorts. next() will return the next variation as they are defined. When
     * it reaches the end, the first variation will be returned:
     *     RPM->VE->MAP->RPM->etc...
     * previous() simply goes the opposite direction:
     *     MAP->VE->RPM->MAP->etc...
     */
    pub fn next(&self) -> Self {
        match self {
            Self::Rpm(_) => Self::Ve(0),
            Self::Ve(_) => Self::Map(Pressure::default()),
            Self::Map(_) => Self::Rpm(0),
        }
    }

    pub fn previous(&self) -> Self {
        match self {
            Self::Rpm(_) => Self::Map(Pressure::default()),
            Self::Ve(_) => Self::Rpm(0),
            Self::Map(_) => Self::Ve(0),
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
            rpm: InputParam::Rpm(7000),
            ve: InputParam::Ve(95),
            map: InputParam::Map(Pressure::from_unit(KiloPascal, 200)),
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
            iter_state: Some(InputParam::Rpm(0)),
        }
    }
}

impl<'a> Iterator for RowIter<'a> {
    type Item = &'a InputParam;

    fn next(&mut self) -> Option<Self::Item> {
        match self.iter_state {
            Some(InputParam::Rpm(_)) => {
                self.iter_state = Some(InputParam::Ve(0));
                Some(&self.row.rpm)
            }
            Some(InputParam::Ve(_)) => {
                self.iter_state = Some(InputParam::Map(Pressure::default()));
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
