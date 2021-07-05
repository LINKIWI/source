import fs from 'fs';
import yaml from 'js-yaml';
import { objLookup } from 'shared/util/data';

/**
 * Client for interacting with a YAML configuration file on disk.
 */
export default class ConfigClient {
  /**
   * Create a config client instance.
   *
   * @param {String} path Path to the config file on disk.
   */
  constructor(path) {
    this.config = yaml.safeLoad(fs.readFileSync(path));
  }

  /**
   * Get a config value by a path descriptor, where keys are delimited by a dot.
   *
   * @param {String} path Path to the config key.
   * @param {String} delimiter Delimiter character of the path descriptor.
   * @returns {*} Corresponding config value as defined in the config file.
   */
  get(path, delimiter = '.') {
    return objLookup(this.config, path.split(delimiter));
  }
}
