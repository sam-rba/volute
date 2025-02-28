#include "test.h"

int
main(void) {
	test_rad_per_sec();
	test_deg_per_sec();
	test_rpm();
	test_as_rad_per_sec();
	test_as_deg_per_sec();
	test_as_rpm();

	test_pascal();
	test_millibar();
	test_kilopascal();
	test_bar();
	test_psi();
	test_as_pascal();
	test_as_millibar();
	test_as_kilopascal();
	test_as_bar();
	test_as_psi();

	test_cubic_centimetre();
	test_litre();
	test_cubic_metre();
	test_cubic_inch();
	test_as_cubic_centimetre();
	test_as_litre();
	test_as_cubic_metre();
	test_as_cubic_inch();
}
