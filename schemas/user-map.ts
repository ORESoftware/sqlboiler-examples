/* tslint:disable */


/**
 * AUTO-GENERATED FILE @ 2020-02-23 23:03:55 - DO NOT EDIT!
 *
 * This file was automatically generated by schemats v.3.0.3
 * $ schemats generate -c postgres://username:password@localhost/postgres -t user_map_table -s public
 *
 */


export namespace user_map_tableFields {
    export type id = number;
    export type user_id = number;
    export type key = string;
    export type value_type = string;
    export type value = object;
    export type added = boolean | null;

}

export interface user_map_table {
    id: user_map_tableFields.id;
    user_id: user_map_tableFields.user_id;
    key: user_map_tableFields.key;
    value_type: user_map_tableFields.value_type;
    value: user_map_tableFields.value;
    added: user_map_tableFields.added;

}
