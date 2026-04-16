/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
import { Observable } from 'rxjs';

import * as GoogleProtobufTimestamp from "../protobuf/timestamp.pb";
/**
 * Represents a time interval, encoded as a Timestamp start (inclusive) and a
 * Timestamp end (exclusive).
 *
 * The start must be less than or equal to the end.
 * When the start equals the end, the interval is empty (matches no time).
 * When both start and end are unspecified, the interval matches any time.
 **/
export type Interval = {
  /**
   * Optional. Inclusive start of the interval.
   *
   * If specified, a Timestamp matching this interval will have to be the same
   * or after the start.
   **/
  startTime?: string;
  /**
   * Optional. Exclusive end of the interval.
   *
   * If specified, a Timestamp matching this interval will have to be before the
   * end.
   **/
  endTime?: string;
};