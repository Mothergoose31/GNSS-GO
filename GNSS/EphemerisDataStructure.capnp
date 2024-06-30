using Go = import "/go.capnp";
Go.Package("gnss");

@0xb3ca6d2462778bb1;
struct Ephemeris {
  # GPS ephemeris data according to the RINEX format
  satelliteId @0 :UInt16;
  year @1 :UInt16;
  month @2 :UInt16;
  day @3 :UInt16;
  hour @4 :UInt16;
  minute @5 :UInt16;
  second @6 :Float32;
  af0 @7 :Float64;  # Clock bias coefficient 0
  af1 @8 :Float64;  # Clock drift coefficient 1
  af2 @9 :Float64;  # Clock drift rate coefficient 2

  iode @10 :Float64;
  crs @11 :Float64;
  deltaN @12 :Float64;
  m0 @13 :Float64;

  cuc @14 :Float64;
  eccentricity @15 :Float64;
  cus @16 :Float64;
  semiMajorAxis @17 :Float64;

  toe @18 :Float64;
  cic @19 :Float64;
  omega0 @20 :Float64;
  cis @21 :Float64;

  inclination @22 :Float64;
  crc @23 :Float64;
  perigeeArgument @24 :Float64;
  rateOfRightAscension @25 :Float64;

  rateOfInclination @26 :Float64;
  codesL2 @27 :Float64;
  gpsWeekDEPRECATED @28 :Float64;
  l2 @29 :Float64;

  signalAccuracy @30 :Float64;
  satelliteHealth @31 :Float64;
  tgd @32 :Float64;
  iodc @33 :Float64;

  transmissionTime @34 :Float64;
  fitInterval @35 :Float64;

  timeOfClock @36 :Float64;

  ionosphereCoefficientsValid @37 :Bool;
  ionosphereAlpha @38 :List(Float64);
  ionosphereBeta @39 :List(Float64);

  timeOfWeekCount @40 :UInt32;
  toeWeek @41 :UInt16;
  timeOfClockWeek @42 :UInt16;
}

struct GlonassEphemeris {
  satelliteId @0 :UInt16;
  year @1 :UInt16;
  dayOfYear @2 :UInt16;
  hour @3 :UInt16;
  minute @4 :UInt16;
  second @5 :Float32;

  xPosition @6 :Float64;
  xVelocity @7 :Float64;
  xAcceleration @8 :Float64;
  yPosition @9 :Float64;
  yVelocity @10 :Float64;
  yAcceleration @11 :Float64;
  zPosition @12 :Float64;
  zVelocity @13 :Float64;
  zAcceleration @14 :Float64;

  satelliteType @15 :UInt8;
  userRangeAccuracy @16 :Float32;
  ageOfOperation @17 :UInt8;

  satelliteHealth @18 :UInt8;
  timeCorrectionDEPRECATED @19 :UInt16;
  timeGroupDelay @20 :UInt16;

  tauN @21 :Float64;
  deltaTauN @22 :Float64;
  gammaN @23 :Float64;

  frequencyNumber1 @24 :UInt8;
  frequencyNumber2 @25 :UInt8;
  frequencyNumber3 @26 :UInt8;
  frequencyNumber4 @27 :UInt8;

  frequencyNumberDEPRECATED @28 :UInt32;

  N4 @29 :UInt8;
  NT @30 :UInt16;
  frequencyNumber @31 :Int16;
  timeCorrectionSeconds @32 :UInt32;
}

struct EphemerisCache {
  gpsEphemerides @0 :List(Ephemeris);
  glonassEphemerides @1 :List(GlonassEphemeris);
}