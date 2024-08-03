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
  timeOfEphemeris @2 :GPSTime;
  timeOfClock @3 :GPSTime;
  squareRootOfSemiMajorAxis @4 :Float64;
}

struct Ephemeris {
  satelliteId @0 :UInt16;
  year @1 :UInt16;
  month @2 :UInt16;
  day @3 :UInt16;
  hour @4 :UInt16;
  minute @5 :UInt16;
  second @6 :Float32;
  clockBiasCoefficient @7 :Float64;
  clockDriftCoefficient @8 :Float64;
  clockDriftRateCoefficient @9 :Float64;
  issueOfDataEphemeris @10 :Float64;
  radiusSineCorrectionTerm @11 :Float64;
  meanMotionDifference @12 :Float64;
  meanAnomaly @13 :Float64;
  latitudeCosineCorrectionTerm @14 :Float64;
  eccentricity @15 :Float64;
  latitudeSineCorrectionTerm @16 :Float64;
  semiMajorAxis @17 :Float64;
  timeOfEphemeris @18 :Float64;
  inclinationCosineCorrectionTerm @19 :Float64;
  rightAscensionOfAscendingNode @20 :Float64;
  inclinationSineCorrectionTerm @21 :Float64;
  inclination @22 :Float64;
  radiusCosineCorrectionTerm @23 :Float64;
  argumentOfPerigee @24 :Float64;
  rateOfRightAscension @25 :Float64;
  rateOfInclination @26 :Float64;
  l2CodeFlags @27 :Float64;
  gpsWeekDeprecated @28 :Float64;
  l2PDataFlag @29 :Float64;
  signalAccuracy @30 :Float64;
  satelliteHealth @31 :Float64;
  totalGroupDelay @32 :Float64;
  issueOfDataClock @33 :Float64;
  transmissionTime @34 :Float64;
  fitInterval @35 :Float64;
  timeOfClock @36 :Float64;
  ionosphereCoefficientsValid @37 :Bool;
  ionosphereAlpha @38 :List(Float64);
  ionosphereBeta @39 :List(Float64);
  timeOfWeekCount @40 :UInt32;
  timeOfEphemerisWeek @41 :UInt16;
  timeOfClockWeek @42 :UInt16;
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