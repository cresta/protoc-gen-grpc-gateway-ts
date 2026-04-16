/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*//**
 * An object that represents a latitude/longitude pair. This is expressed as a
 * pair of doubles to represent degrees latitude and degrees longitude. Unless
 * specified otherwise, this object must conform to the
 * <a href="https://en.wikipedia.org/wiki/World_Geodetic_System#1984_version">
 * WGS84 standard</a>. Values must be within normalized ranges.
 **/
export type LatLng = {
  /**
   * The latitude in degrees. It must be in the range [-90.0, +90.0].
   **/
  latitude?: number;
  /**
   * The longitude in degrees. It must be in the range [-180.0, +180.0].
   **/
  longitude?: number;
};