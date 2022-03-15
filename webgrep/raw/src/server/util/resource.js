import fs from 'fs';
import path from 'path';

/**
 * Read a dist resource by path.
 *
 * @param {String} resource Path to the resource, relative to the build dist directory.
 * @param {Function} cb Callback invoked with (err, file) after attempting to open the file.
 */
const distResource = (resource, cb) =>
  fs.readFile(path.join(__dirname, '../../../dist', resource), cb);

/**
 * Read a client resource by path.
 *
 * @param {String} resource Path to the client resource, relative to the build dist directory.
 * @param {Function} cb Callback invoked with (err, file) after attempting to open the file.
 */
export const clientResource = (resource, cb) => distResource(path.join('client', resource), cb);

/**
 * Read a server resource by path.
 *
 * @param {String} resource Path to the server resource, relative to the build dist directory.
 * @param {Function} cb Callback invoked with (err, file) after attempting to open the file.
 */
export const serverResource = (resource, cb) => distResource(path.join('server', resource), cb);
