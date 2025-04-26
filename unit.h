#define STANDARD_PRESSURE millibar(1013)
#define STANDARD_TEMPERATURE celsius(20)


typedef double AngularSpeed;
typedef AngularSpeed (*AngularSpeedMaker)(double);
typedef double (*AngularSpeedReader)(AngularSpeed);

AngularSpeed rad_per_sec(double x);
AngularSpeed deg_per_sec(double x);
AngularSpeed rpm(double x);
double as_rad_per_sec(AngularSpeed x);
double as_deg_per_sec(AngularSpeed x);
double as_rpm(AngularSpeed x);


typedef double Fraction;
typedef Fraction (*FractionMaker)(double);
typedef double (*FractionReader)(Fraction);

Fraction percent(double x);
double as_percent(double x);


typedef double Pressure;
typedef Pressure (*PressureMaker)(double);
typedef double (*PressureReader)(Pressure);

Pressure pascal(double x);
Pressure millibar(double x);
Pressure kilopascal(double x);
Pressure bar(double x);
Pressure psi(double x);
Pressure inch_mercury(double x);
double as_pascal(Pressure x);
double as_millibar(Pressure x);
double as_kilopascal(Pressure x);
double as_bar(Pressure x);
double as_psi(Pressure x);
double as_inch_mercury(Pressure x);

extern const size_t n_pressure_units;
extern const char *const pressure_units[];
extern const PressureMaker pressure_makers[];
extern const PressureReader pressure_readers[];


typedef double Temperature;
typedef Temperature (*TemperatureMaker)(double);
typedef double (*TemperatureReader)(Temperature);

Temperature kelvin(double x);
Temperature celsius(double x);
Temperature fahrenheit(double x);
Temperature rankine(double x);
double as_kelvin(Temperature t);
double as_celsius(Temperature t);
double as_fahrenheit(Temperature t);
double as_rankine(Temperature t);

extern const size_t n_temperature_units;
extern const char *const temperature_units[];
extern const TemperatureMaker temperature_makers[];
extern const TemperatureReader temperature_readers[];


typedef double Volume;
typedef Volume (*VolumeMaker)(double);
typedef double (*VolumeReader)(Volume);

Volume cubic_centimetre(double x);
Volume litre(double x);
Volume cubic_metre(double x);
Volume cubic_inch(double x);
double as_cubic_centimetre(Volume x);
double as_litre(Volume x);
double as_cubic_metre(double x);
double as_cubic_inch(double x);

extern const size_t n_volume_units;
extern const char *const volume_units[];
extern const VolumeMaker volume_makers[];
extern const VolumeReader volume_readers[];


typedef double VolumeFlowRate;
typedef VolumeFlowRate (*VolumeFlowRateMaker)(double);
typedef double (*VolumeFlowRateReader)(VolumeFlowRate);

VolumeFlowRate cubic_metre_per_sec(double x);
VolumeFlowRate cubic_metre_per_min(double x);
VolumeFlowRate cubic_foot_per_min(double x);
double as_cubic_metre_per_sec(VolumeFlowRate x);
double as_cubic_metre_per_min(VolumeFlowRate x);
double as_cubic_foot_per_min(VolumeFlowRate x);

extern const size_t n_volume_flow_rate_units;
extern const char *const volume_flow_rate_units[];
extern const VolumeFlowRateMaker volume_flow_rate_makers[];
extern const VolumeFlowRateReader volume_flow_rate_readers[];


typedef double MassFlowRate;
typedef MassFlowRate (*MassFlowRateMaker)(double);
typedef double (*MassFlowRateReader)(MassFlowRate);

MassFlowRate kilo_per_sec(double x);
MassFlowRate pound_per_min(double x);
double as_kilo_per_sec(MassFlowRate m);
double as_pound_per_min(MassFlowRate m);

extern const size_t n_mass_flow_rate_units;
extern const char *const mass_flow_rate_units[];
extern const MassFlowRateMaker mass_flow_rate_makers[];
extern const MassFlowRateReader mass_flow_rate_readers[];
