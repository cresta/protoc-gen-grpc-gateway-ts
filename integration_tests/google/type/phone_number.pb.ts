/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);
/**
 * An object representing a short code, which is a phone number that is
 * typically much shorter than regular phone numbers and can be used to
 * address messages in MMS and SMS systems, as well as for abbreviated dialing
 * (For example "Text 611 to see how many minutes you have remaining on your
 * plan.").
 *
 * Short codes are restricted to a region and are not internationally
 * dialable, which means the same short code can exist in different regions,
 * with different usage and pricing, even if those regions share the same
 * country calling code (For example: US and CA).
 **/
export type PhoneNumberShortCode = {
  /**
   * Required. The BCP-47 region code of the location where calls to this
   * short code can be made, such as "US" and "BB".
   *
   * Reference(s):
   *  - http://www.unicode.org/reports/tr35/#unicode_region_subtag
   **/
  regionCode?: string;
  /**
   * Required. The short code digits, without a leading plus ('+') or country
   * calling code. For example "611".
   **/
  number?: string;
};

type BasePhoneNumber = {
  extension?: string;
};

/**
 * An object representing a phone number, suitable as an API wire format.
 *
 * This representation:
 *
 *  - should not be used for locale-specific formatting of a phone number, such
 *    as "+1 (650) 253-0000 ext. 123"
 *
 *  - is not designed for efficient storage
 *  - may not be suitable for dialing - specialized libraries (see references)
 *    should be used to parse the number for that purpose
 *
 * To do something meaningful with this number, such as format it for various
 * use-cases, convert it to an `i18n.phonenumbers.PhoneNumber` object first.
 *
 * For instance, in Java this would be:
 *
 *     com.google.type.PhoneNumber wireProto =
 *         com.google.type.PhoneNumber.newBuilder().build();
 *     com.google.i18n.phonenumbers.Phonenumber.PhoneNumber phoneNumber =
 *         PhoneNumberUtil.getInstance().parse(wireProto.getE164Number(), "ZZ");
 *     if (!wireProto.getExtension().isEmpty()) {
 *       phoneNumber.setExtension(wireProto.getExtension());
 *     }
 *
 *  Reference(s):
 *   - https://github.com/google/libphonenumber
 **/
export type PhoneNumber = BasePhoneNumber &
  OneOf<{
    e164Number: string;
    shortCode: PhoneNumberShortCode;
  }>;