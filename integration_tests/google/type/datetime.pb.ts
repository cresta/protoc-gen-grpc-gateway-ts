/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
import { Observable } from 'rxjs';

import * as GoogleProtobufDuration from "../protobuf/duration.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);
type BaseDateTime = {
  year?: number;
  month?: number;
  day?: number;
  hours?: number;
  minutes?: number;
  seconds?: number;
  nanos?: number;
};

/**
 * Represents civil time (or occasionally physical time).
 *
 * This type can represent a civil time in one of a few possible ways:
 *
 *  * When utc_offset is set and time_zone is unset: a civil time on a calendar
 *    day with a particular offset from UTC.
 *  * When time_zone is set and utc_offset is unset: a civil time on a calendar
 *    day in a particular time zone.
 *  * When neither time_zone nor utc_offset is set: a civil time on a calendar
 *    day in local time.
 *
 * The date is relative to the Proleptic Gregorian Calendar.
 *
 * If year, month, or day are 0, the DateTime is considered not to have a
 * specific year, month, or day respectively.
 *
 * This type may also be used to represent a physical time if all the date and
 * time fields are set and either case of the `time_offset` oneof is set.
 * Consider using `Timestamp` message for physical time instead. If your use
 * case also would like to store the user's timezone, that can be done in
 * another field.
 *
 * This type is more flexible than some applications may want. Make sure to
 * document and validate your application's limitations.
 **/
export type DateTime = BaseDateTime &
  OneOf<{
    utcOffset: string;
    timeZone: TimeZone;
  }>;

/**
 * Represents a time zone from the
 * [IANA Time Zone Database](https://www.iana.org/time-zones).
 **/
export type TimeZone = {
  /**
   * IANA Time Zone Database time zone. For example "America/New_York".
   **/
  id?: string;
  /**
   * Optional. IANA Time Zone Database version number. For example "2019a".
   **/
  version?: string;
};