typedef double AngularSpeed;

AngularSpeed rad_per_sec(double x);
AngularSpeed deg_per_sec(double x);
AngularSpeed rpm(double x);
double as_rad_per_sec(AngularSpeed x);
double as_deg_per_sec(AngularSpeed x);
double as_rpm(AngularSpeed x);


typedef double Fraction;

Fraction percent(double x);
double as_percent(double x);


typedef double Pressure;

Pressure pascal(double x);
Pressure millibar(double x);
Pressure kilopascal(double x);
Pressure bar(double x);
Pressure psi(double x);
double as_pascal(Pressure x);
double as_millibar(Pressure x);
double as_kilopascal(Pressure x);
double as_bar(Pressure x);
double as_psi(Pressure x);


typedef double Volume;

Volume cubic_centimetre(double x);
Volume litre(double x);
Volume cubic_metre(double x);
Volume cubic_inch(double x);
double as_cubic_centimetre(Volume x);
double as_litre(Volume x);
double as_cubic_metre(double x);
double as_cubic_inch(double x);
