@0xb3ca6d2462778bb1;

using Go = import "/go.capnp";

$Go.package("gnss");
$Go.import("github.com/mothergoose31/GNSS-GO/GNSS");

struct BaseEphemeris {
  pseudoRandomNumber @0 :Text;
  epoch @1 :GPSTime;
  ephemerisType @2 :EphemerisType;
  isHealthy @3 :Bool;
  maximumTimeDifference @4 :Float64;
  fileEpoch @5 :GPSTime;
  fileName @6 :Text;
  fileSource @7 :Text;
}

struct GPSEphemeris {
  baseEphemeris @0 :BaseEphemeris;
  ephemerisData @1 :Ephemeris;
  toe @2 :GPSTime;  
  toc @3 :GPSTime;
  squareRootOfSemiMajorAxis @4 :Float64;
}


struct RINEXHeader {
  version @0 :Float64;
  type @1 :Text;
  satelliteSystem @2 :Text;
  programName @3 :Text;
  agency @4 :Text;
  date @5 :Time;
  comments @6 :List(Text);
  leapSeconds @7 :Int32;
}

struct Ephemeris {

  svId @0 :UInt16;
  year @1 :UInt16;
  month @2 :UInt16;
  day @3 :UInt16;
  hour @4 :UInt16;
  minute @5 :UInt16;
  second @6 :Float32;
  af0 @7 :Float64;
  af1 @8 :Float64;
  af2 @9 :Float64;

  iode @10 :Float64;
  crs @11 :Float64;
  deltaN @12 :Float64;
  m0 @13 :Float64;

  cuc @14 :Float64;
  ecc @15 :Float64;
  cus @16 :Float64;
  a @17 :Float64; 

  toe @18 :Float64;
  cic @19 :Float64;
  omega0 @20 :Float64;
  cis @21 :Float64;

  i0 @22 :Float64;
  crc @23 :Float64;
  omega @24 :Float64;
  omegaDot @25 :Float64;

  iDot @26 :Float64;
  codesL2 @27 :Float64;
  l2 @28 :Float64;

  svAcc @29 :Float64;
  svHealth @30 :Float64;
  tgd @31 :Float64;
  iodc @32 :Float64;

  transmissionTime @33 :Float64;
  fitInterval @34 :Float64;

  toc @35 :Float64;

  ionoCoeffsValid @36 :Bool;
  ionoAlpha @37 :List(Float64);
  ionoBeta @38 :List(Float64);

  towCount @39 :UInt32;
  toeWeek @40 :UInt16;
  tocWeek @41 :UInt16;
}


struct RINEXEphemeris {
  header @0 :RINEXHeader;
  satelliteId @1 :Int32;
  epoch @2 :Time;
  clockBias @3 :Float64;
  relativeFrequencyBias @4 :Float64;
  messageFrameTime @5 :Float64;
  positionX @6 :Float64;
  velocityX @7 :Float64;
  accelerationX @8 :Float64;
  positionY @9 :Float64;
  velocityY @10 :Float64;
  accelerationY @11 :Float64;
  positionZ @12 :Float64;
  velocityZ @13 :Float64;
  accelerationZ @14 :Float64;
  health @15 :Float64;
  frequencyChannelOffset @16 :Int32;
  informationAge @17 :Float64;
}

struct GroupedEphemerides {
  satelliteId @0 :Int32;
  sortedEphemerides @1 :List(RINEXEphemeris);
}

struct SP3FormatEphemeris {
  header @0 :SP3Header;
  epochs @1 :List(SP3Epoch);
}

struct SP3Header {
  version @0 :Text;
  start @1 :Time;
  numberOfEpochs @2 :Int32;
  dataUsed @3 :Text;
  coordinateSystem @4 :Text;
  orbitType @5 :Text;
  agency @6 :Text;
  gpsWeek @7 :Int32;
  secondsOfWeek @8 :Float64;
  epochInterval @9 :Float64;
  modifiedJulianDay @10 :Int32;
  fractionalDay @11 :Float64;
}

struct SP3Epoch {
  time @0 :Time;
  entries @1 :List(SP3Entry);
}

struct SP3Entry {
  satelliteVehicleNumber @0 :Text;
  xPosition @1 :Float64;
  yPosition @2 :Float64;
  zPosition @3 :Float64;
  clockBias @4 :Float64;
}

struct GPSTime {
  week @0 :Int32;
  timeOfWeek @1 :Float64;
}

struct Time {
  seconds @0 :Int64;
  nanoseconds @1 :Int32;
}

enum EphemerisType {
  navigation @0;
  finalOrbit @1;
  rapidOrbit @2;
  ultraRapidOrbit @3;
  qcomPoly @4;
}