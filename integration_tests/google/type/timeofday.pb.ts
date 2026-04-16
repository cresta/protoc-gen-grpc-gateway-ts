/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*//**
 * Represents a time of day. The date and time zone are either not significant
 * or are specified elsewhere. An API may choose to allow leap seconds. Related
 * types are [google.type.Date][google.type.Date] and
 * `google.protobuf.Timestamp`.
 **/
export type TimeOfDay = {
  /**
   * Hours of a day in 24 hour format. Must be greater than or equal to 0 and
   * typically must be less than or equal to 23. An API may choose to allow the
   * value "24:00:00" for scenarios like business closing time.
   **/
  hours?: number;
  /**
   * Minutes of an hour. Must be greater than or equal to 0 and less than or
   * equal to 59.
   **/
  minutes?: number;
  /**
   * Seconds of a minute. Must be greater than or equal to 0 and typically must
   * be less than or equal to 59. An API may allow the value 60 if it allows
   * leap-seconds.
   **/
  seconds?: number;
  /**
   * Fractions of seconds, in nanoseconds. Must be greater than or equal to 0
   * and less than or equal to 999,999,999.
   **/
  nanos?: number;
};