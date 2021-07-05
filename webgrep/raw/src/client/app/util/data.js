import { constants } from 'supercharged/shared';
import { parseJSON } from 'shared/util/data';

/**
 * Parse a Supercharged response payload into its canonical client-side representation.
 *
 * @param {String} resp JSON string response sent by the server.
 * @return {Object} Object containing properties `err` and `data` describing the server response.
 */
export const parseSuperchargedResponse = (resp) => {
  const { data: { success, id, code, message, data } = {}, ok } = parseJSON(resp);

  if (!ok) {
    return {
      err: {
        code: constants.error.CODE_CLIENT_UNDEFINED,
        message: 'Failed to parse incoming response from server.',
      },
      id,
      data: undefined,
    };
  }

  if (!success) {
    return {
      err: { code, message },
      id,
      data,
    };
  }

  return {
    err: null,
    id,
    data,
  };
};

/**
 * Static functions for serializing values for storage in URL state.
 */
export class URLStateSerializer {
  static identity(value) {
    return value;
  }

  static string(value) {
    return encodeURIComponent(value);
  }

  static array(value) {
    return value.map(URLStateSerializer.string).join(',');
  }
}

/**
 * Static functions for deserializing values from URL state.
 */
export class URLStateDeserializer {
  static identity(value) {
    return value;
  }

  static boolean(value) {
    return value.toLowerCase() === 'true';
  }

  static string(value) {
    return value;
  }

  static number(value) {
    return parseFloat(value);
  }

  static array(value) {
    return value.split(',').map(URLStateDeserializer.string);
  }
}
