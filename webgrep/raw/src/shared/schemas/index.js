/*eslint-disable*/
import * as $protobuf from "protobufjs/minimal";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const IndexSpec = $root.IndexSpec = (() => {

    /**
     * Properties of an IndexSpec.
     * @exports IIndexSpec
     * @interface IIndexSpec
     * @property {string|null} [name] IndexSpec name
     * @property {Array.<IPathSpec>|null} [paths] IndexSpec paths
     * @property {Array.<IRepoSpec>|null} [repos] IndexSpec repos
     */

    /**
     * Constructs a new IndexSpec.
     * @exports IndexSpec
     * @classdesc Represents an IndexSpec.
     * @implements IIndexSpec
     * @constructor
     * @param {IIndexSpec=} [properties] Properties to set
     */
    function IndexSpec(properties) {
        this.paths = [];
        this.repos = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * IndexSpec name.
     * @member {string} name
     * @memberof IndexSpec
     * @instance
     */
    IndexSpec.prototype.name = "";

    /**
     * IndexSpec paths.
     * @member {Array.<IPathSpec>} paths
     * @memberof IndexSpec
     * @instance
     */
    IndexSpec.prototype.paths = $util.emptyArray;

    /**
     * IndexSpec repos.
     * @member {Array.<IRepoSpec>} repos
     * @memberof IndexSpec
     * @instance
     */
    IndexSpec.prototype.repos = $util.emptyArray;

    /**
     * Creates a new IndexSpec instance using the specified properties.
     * @function create
     * @memberof IndexSpec
     * @static
     * @param {IIndexSpec=} [properties] Properties to set
     * @returns {IndexSpec} IndexSpec instance
     */
    IndexSpec.create = function create(properties) {
        return new IndexSpec(properties);
    };

    /**
     * Encodes the specified IndexSpec message. Does not implicitly {@link IndexSpec.verify|verify} messages.
     * @function encode
     * @memberof IndexSpec
     * @static
     * @param {IIndexSpec} message IndexSpec message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    IndexSpec.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
        if (message.paths != null && message.paths.length)
            for (let i = 0; i < message.paths.length; ++i)
                $root.PathSpec.encode(message.paths[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.repos != null && message.repos.length)
            for (let i = 0; i < message.repos.length; ++i)
                $root.RepoSpec.encode(message.repos[i], writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified IndexSpec message, length delimited. Does not implicitly {@link IndexSpec.verify|verify} messages.
     * @function encodeDelimited
     * @memberof IndexSpec
     * @static
     * @param {IIndexSpec} message IndexSpec message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    IndexSpec.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an IndexSpec message from the specified reader or buffer.
     * @function decode
     * @memberof IndexSpec
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {IndexSpec} IndexSpec
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    IndexSpec.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.IndexSpec();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.name = reader.string();
                break;
            case 2:
                if (!(message.paths && message.paths.length))
                    message.paths = [];
                message.paths.push($root.PathSpec.decode(reader, reader.uint32()));
                break;
            case 3:
                if (!(message.repos && message.repos.length))
                    message.repos = [];
                message.repos.push($root.RepoSpec.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an IndexSpec message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof IndexSpec
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {IndexSpec} IndexSpec
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    IndexSpec.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an IndexSpec message.
     * @function verify
     * @memberof IndexSpec
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    IndexSpec.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.paths != null && message.hasOwnProperty("paths")) {
            if (!Array.isArray(message.paths))
                return "paths: array expected";
            for (let i = 0; i < message.paths.length; ++i) {
                let error = $root.PathSpec.verify(message.paths[i]);
                if (error)
                    return "paths." + error;
            }
        }
        if (message.repos != null && message.hasOwnProperty("repos")) {
            if (!Array.isArray(message.repos))
                return "repos: array expected";
            for (let i = 0; i < message.repos.length; ++i) {
                let error = $root.RepoSpec.verify(message.repos[i]);
                if (error)
                    return "repos." + error;
            }
        }
        return null;
    };

    /**
     * Creates an IndexSpec message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof IndexSpec
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {IndexSpec} IndexSpec
     */
    IndexSpec.fromObject = function fromObject(object) {
        if (object instanceof $root.IndexSpec)
            return object;
        let message = new $root.IndexSpec();
        if (object.name != null)
            message.name = String(object.name);
        if (object.paths) {
            if (!Array.isArray(object.paths))
                throw TypeError(".IndexSpec.paths: array expected");
            message.paths = [];
            for (let i = 0; i < object.paths.length; ++i) {
                if (typeof object.paths[i] !== "object")
                    throw TypeError(".IndexSpec.paths: object expected");
                message.paths[i] = $root.PathSpec.fromObject(object.paths[i]);
            }
        }
        if (object.repos) {
            if (!Array.isArray(object.repos))
                throw TypeError(".IndexSpec.repos: array expected");
            message.repos = [];
            for (let i = 0; i < object.repos.length; ++i) {
                if (typeof object.repos[i] !== "object")
                    throw TypeError(".IndexSpec.repos: object expected");
                message.repos[i] = $root.RepoSpec.fromObject(object.repos[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from an IndexSpec message. Also converts values to other types if specified.
     * @function toObject
     * @memberof IndexSpec
     * @static
     * @param {IndexSpec} message IndexSpec
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    IndexSpec.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults) {
            object.paths = [];
            object.repos = [];
        }
        if (options.defaults)
            object.name = "";
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.paths && message.paths.length) {
            object.paths = [];
            for (let j = 0; j < message.paths.length; ++j)
                object.paths[j] = $root.PathSpec.toObject(message.paths[j], options);
        }
        if (message.repos && message.repos.length) {
            object.repos = [];
            for (let j = 0; j < message.repos.length; ++j)
                object.repos[j] = $root.RepoSpec.toObject(message.repos[j], options);
        }
        return object;
    };

    /**
     * Converts this IndexSpec to JSON.
     * @function toJSON
     * @memberof IndexSpec
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    IndexSpec.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return IndexSpec;
})();

export const Metadata = $root.Metadata = (() => {

    /**
     * Properties of a Metadata.
     * @exports IMetadata
     * @interface IMetadata
     * @property {string|null} [urlPattern] Metadata urlPattern
     * @property {string|null} [remote] Metadata remote
     * @property {string|null} [github] Metadata github
     * @property {Array.<string>|null} [labels] Metadata labels
     */

    /**
     * Constructs a new Metadata.
     * @exports Metadata
     * @classdesc Represents a Metadata.
     * @implements IMetadata
     * @constructor
     * @param {IMetadata=} [properties] Properties to set
     */
    function Metadata(properties) {
        this.labels = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Metadata urlPattern.
     * @member {string} urlPattern
     * @memberof Metadata
     * @instance
     */
    Metadata.prototype.urlPattern = "";

    /**
     * Metadata remote.
     * @member {string} remote
     * @memberof Metadata
     * @instance
     */
    Metadata.prototype.remote = "";

    /**
     * Metadata github.
     * @member {string} github
     * @memberof Metadata
     * @instance
     */
    Metadata.prototype.github = "";

    /**
     * Metadata labels.
     * @member {Array.<string>} labels
     * @memberof Metadata
     * @instance
     */
    Metadata.prototype.labels = $util.emptyArray;

    /**
     * Creates a new Metadata instance using the specified properties.
     * @function create
     * @memberof Metadata
     * @static
     * @param {IMetadata=} [properties] Properties to set
     * @returns {Metadata} Metadata instance
     */
    Metadata.create = function create(properties) {
        return new Metadata(properties);
    };

    /**
     * Encodes the specified Metadata message. Does not implicitly {@link Metadata.verify|verify} messages.
     * @function encode
     * @memberof Metadata
     * @static
     * @param {IMetadata} message Metadata message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Metadata.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.urlPattern != null && message.hasOwnProperty("urlPattern"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.urlPattern);
        if (message.remote != null && message.hasOwnProperty("remote"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.remote);
        if (message.github != null && message.hasOwnProperty("github"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.github);
        if (message.labels != null && message.labels.length)
            for (let i = 0; i < message.labels.length; ++i)
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.labels[i]);
        return writer;
    };

    /**
     * Encodes the specified Metadata message, length delimited. Does not implicitly {@link Metadata.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Metadata
     * @static
     * @param {IMetadata} message Metadata message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Metadata.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Metadata message from the specified reader or buffer.
     * @function decode
     * @memberof Metadata
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Metadata} Metadata
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Metadata.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Metadata();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.urlPattern = reader.string();
                break;
            case 2:
                message.remote = reader.string();
                break;
            case 3:
                message.github = reader.string();
                break;
            case 4:
                if (!(message.labels && message.labels.length))
                    message.labels = [];
                message.labels.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Metadata message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Metadata
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Metadata} Metadata
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Metadata.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Metadata message.
     * @function verify
     * @memberof Metadata
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Metadata.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.urlPattern != null && message.hasOwnProperty("urlPattern"))
            if (!$util.isString(message.urlPattern))
                return "urlPattern: string expected";
        if (message.remote != null && message.hasOwnProperty("remote"))
            if (!$util.isString(message.remote))
                return "remote: string expected";
        if (message.github != null && message.hasOwnProperty("github"))
            if (!$util.isString(message.github))
                return "github: string expected";
        if (message.labels != null && message.hasOwnProperty("labels")) {
            if (!Array.isArray(message.labels))
                return "labels: array expected";
            for (let i = 0; i < message.labels.length; ++i)
                if (!$util.isString(message.labels[i]))
                    return "labels: string[] expected";
        }
        return null;
    };

    /**
     * Creates a Metadata message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Metadata
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Metadata} Metadata
     */
    Metadata.fromObject = function fromObject(object) {
        if (object instanceof $root.Metadata)
            return object;
        let message = new $root.Metadata();
        if (object.urlPattern != null)
            message.urlPattern = String(object.urlPattern);
        if (object.remote != null)
            message.remote = String(object.remote);
        if (object.github != null)
            message.github = String(object.github);
        if (object.labels) {
            if (!Array.isArray(object.labels))
                throw TypeError(".Metadata.labels: array expected");
            message.labels = [];
            for (let i = 0; i < object.labels.length; ++i)
                message.labels[i] = String(object.labels[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a Metadata message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Metadata
     * @static
     * @param {Metadata} message Metadata
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Metadata.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.labels = [];
        if (options.defaults) {
            object.urlPattern = "";
            object.remote = "";
            object.github = "";
        }
        if (message.urlPattern != null && message.hasOwnProperty("urlPattern"))
            object.urlPattern = message.urlPattern;
        if (message.remote != null && message.hasOwnProperty("remote"))
            object.remote = message.remote;
        if (message.github != null && message.hasOwnProperty("github"))
            object.github = message.github;
        if (message.labels && message.labels.length) {
            object.labels = [];
            for (let j = 0; j < message.labels.length; ++j)
                object.labels[j] = message.labels[j];
        }
        return object;
    };

    /**
     * Converts this Metadata to JSON.
     * @function toJSON
     * @memberof Metadata
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Metadata.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Metadata;
})();

export const PathSpec = $root.PathSpec = (() => {

    /**
     * Properties of a PathSpec.
     * @exports IPathSpec
     * @interface IPathSpec
     * @property {string|null} [path] PathSpec path
     * @property {string|null} [name] PathSpec name
     * @property {string|null} [orderedContents] PathSpec orderedContents
     * @property {IMetadata|null} [metadata] PathSpec metadata
     */

    /**
     * Constructs a new PathSpec.
     * @exports PathSpec
     * @classdesc Represents a PathSpec.
     * @implements IPathSpec
     * @constructor
     * @param {IPathSpec=} [properties] Properties to set
     */
    function PathSpec(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PathSpec path.
     * @member {string} path
     * @memberof PathSpec
     * @instance
     */
    PathSpec.prototype.path = "";

    /**
     * PathSpec name.
     * @member {string} name
     * @memberof PathSpec
     * @instance
     */
    PathSpec.prototype.name = "";

    /**
     * PathSpec orderedContents.
     * @member {string} orderedContents
     * @memberof PathSpec
     * @instance
     */
    PathSpec.prototype.orderedContents = "";

    /**
     * PathSpec metadata.
     * @member {IMetadata|null|undefined} metadata
     * @memberof PathSpec
     * @instance
     */
    PathSpec.prototype.metadata = null;

    /**
     * Creates a new PathSpec instance using the specified properties.
     * @function create
     * @memberof PathSpec
     * @static
     * @param {IPathSpec=} [properties] Properties to set
     * @returns {PathSpec} PathSpec instance
     */
    PathSpec.create = function create(properties) {
        return new PathSpec(properties);
    };

    /**
     * Encodes the specified PathSpec message. Does not implicitly {@link PathSpec.verify|verify} messages.
     * @function encode
     * @memberof PathSpec
     * @static
     * @param {IPathSpec} message PathSpec message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PathSpec.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.path != null && message.hasOwnProperty("path"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.path);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
        if (message.orderedContents != null && message.hasOwnProperty("orderedContents"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.orderedContents);
        if (message.metadata != null && message.hasOwnProperty("metadata"))
            $root.Metadata.encode(message.metadata, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified PathSpec message, length delimited. Does not implicitly {@link PathSpec.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PathSpec
     * @static
     * @param {IPathSpec} message PathSpec message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PathSpec.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PathSpec message from the specified reader or buffer.
     * @function decode
     * @memberof PathSpec
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PathSpec} PathSpec
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PathSpec.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.PathSpec();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.path = reader.string();
                break;
            case 2:
                message.name = reader.string();
                break;
            case 3:
                message.orderedContents = reader.string();
                break;
            case 4:
                message.metadata = $root.Metadata.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PathSpec message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PathSpec
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PathSpec} PathSpec
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PathSpec.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PathSpec message.
     * @function verify
     * @memberof PathSpec
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PathSpec.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.path != null && message.hasOwnProperty("path"))
            if (!$util.isString(message.path))
                return "path: string expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.orderedContents != null && message.hasOwnProperty("orderedContents"))
            if (!$util.isString(message.orderedContents))
                return "orderedContents: string expected";
        if (message.metadata != null && message.hasOwnProperty("metadata")) {
            let error = $root.Metadata.verify(message.metadata);
            if (error)
                return "metadata." + error;
        }
        return null;
    };

    /**
     * Creates a PathSpec message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PathSpec
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PathSpec} PathSpec
     */
    PathSpec.fromObject = function fromObject(object) {
        if (object instanceof $root.PathSpec)
            return object;
        let message = new $root.PathSpec();
        if (object.path != null)
            message.path = String(object.path);
        if (object.name != null)
            message.name = String(object.name);
        if (object.orderedContents != null)
            message.orderedContents = String(object.orderedContents);
        if (object.metadata != null) {
            if (typeof object.metadata !== "object")
                throw TypeError(".PathSpec.metadata: object expected");
            message.metadata = $root.Metadata.fromObject(object.metadata);
        }
        return message;
    };

    /**
     * Creates a plain object from a PathSpec message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PathSpec
     * @static
     * @param {PathSpec} message PathSpec
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PathSpec.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.path = "";
            object.name = "";
            object.orderedContents = "";
            object.metadata = null;
        }
        if (message.path != null && message.hasOwnProperty("path"))
            object.path = message.path;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.orderedContents != null && message.hasOwnProperty("orderedContents"))
            object.orderedContents = message.orderedContents;
        if (message.metadata != null && message.hasOwnProperty("metadata"))
            object.metadata = $root.Metadata.toObject(message.metadata, options);
        return object;
    };

    /**
     * Converts this PathSpec to JSON.
     * @function toJSON
     * @memberof PathSpec
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PathSpec.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PathSpec;
})();

export const RepoSpec = $root.RepoSpec = (() => {

    /**
     * Properties of a RepoSpec.
     * @exports IRepoSpec
     * @interface IRepoSpec
     * @property {string|null} [path] RepoSpec path
     * @property {string|null} [name] RepoSpec name
     * @property {Array.<string>|null} [revisions] RepoSpec revisions
     * @property {IMetadata|null} [metadata] RepoSpec metadata
     * @property {boolean|null} [walkSubmodules] RepoSpec walkSubmodules
     */

    /**
     * Constructs a new RepoSpec.
     * @exports RepoSpec
     * @classdesc Represents a RepoSpec.
     * @implements IRepoSpec
     * @constructor
     * @param {IRepoSpec=} [properties] Properties to set
     */
    function RepoSpec(properties) {
        this.revisions = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * RepoSpec path.
     * @member {string} path
     * @memberof RepoSpec
     * @instance
     */
    RepoSpec.prototype.path = "";

    /**
     * RepoSpec name.
     * @member {string} name
     * @memberof RepoSpec
     * @instance
     */
    RepoSpec.prototype.name = "";

    /**
     * RepoSpec revisions.
     * @member {Array.<string>} revisions
     * @memberof RepoSpec
     * @instance
     */
    RepoSpec.prototype.revisions = $util.emptyArray;

    /**
     * RepoSpec metadata.
     * @member {IMetadata|null|undefined} metadata
     * @memberof RepoSpec
     * @instance
     */
    RepoSpec.prototype.metadata = null;

    /**
     * RepoSpec walkSubmodules.
     * @member {boolean} walkSubmodules
     * @memberof RepoSpec
     * @instance
     */
    RepoSpec.prototype.walkSubmodules = false;

    /**
     * Creates a new RepoSpec instance using the specified properties.
     * @function create
     * @memberof RepoSpec
     * @static
     * @param {IRepoSpec=} [properties] Properties to set
     * @returns {RepoSpec} RepoSpec instance
     */
    RepoSpec.create = function create(properties) {
        return new RepoSpec(properties);
    };

    /**
     * Encodes the specified RepoSpec message. Does not implicitly {@link RepoSpec.verify|verify} messages.
     * @function encode
     * @memberof RepoSpec
     * @static
     * @param {IRepoSpec} message RepoSpec message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    RepoSpec.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.path != null && message.hasOwnProperty("path"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.path);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
        if (message.revisions != null && message.revisions.length)
            for (let i = 0; i < message.revisions.length; ++i)
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.revisions[i]);
        if (message.metadata != null && message.hasOwnProperty("metadata"))
            $root.Metadata.encode(message.metadata, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.walkSubmodules != null && message.hasOwnProperty("walkSubmodules"))
            writer.uint32(/* id 5, wireType 0 =*/40).bool(message.walkSubmodules);
        return writer;
    };

    /**
     * Encodes the specified RepoSpec message, length delimited. Does not implicitly {@link RepoSpec.verify|verify} messages.
     * @function encodeDelimited
     * @memberof RepoSpec
     * @static
     * @param {IRepoSpec} message RepoSpec message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    RepoSpec.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a RepoSpec message from the specified reader or buffer.
     * @function decode
     * @memberof RepoSpec
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {RepoSpec} RepoSpec
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    RepoSpec.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.RepoSpec();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.path = reader.string();
                break;
            case 2:
                message.name = reader.string();
                break;
            case 3:
                if (!(message.revisions && message.revisions.length))
                    message.revisions = [];
                message.revisions.push(reader.string());
                break;
            case 4:
                message.metadata = $root.Metadata.decode(reader, reader.uint32());
                break;
            case 5:
                message.walkSubmodules = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a RepoSpec message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof RepoSpec
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {RepoSpec} RepoSpec
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    RepoSpec.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a RepoSpec message.
     * @function verify
     * @memberof RepoSpec
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    RepoSpec.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.path != null && message.hasOwnProperty("path"))
            if (!$util.isString(message.path))
                return "path: string expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.revisions != null && message.hasOwnProperty("revisions")) {
            if (!Array.isArray(message.revisions))
                return "revisions: array expected";
            for (let i = 0; i < message.revisions.length; ++i)
                if (!$util.isString(message.revisions[i]))
                    return "revisions: string[] expected";
        }
        if (message.metadata != null && message.hasOwnProperty("metadata")) {
            let error = $root.Metadata.verify(message.metadata);
            if (error)
                return "metadata." + error;
        }
        if (message.walkSubmodules != null && message.hasOwnProperty("walkSubmodules"))
            if (typeof message.walkSubmodules !== "boolean")
                return "walkSubmodules: boolean expected";
        return null;
    };

    /**
     * Creates a RepoSpec message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof RepoSpec
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {RepoSpec} RepoSpec
     */
    RepoSpec.fromObject = function fromObject(object) {
        if (object instanceof $root.RepoSpec)
            return object;
        let message = new $root.RepoSpec();
        if (object.path != null)
            message.path = String(object.path);
        if (object.name != null)
            message.name = String(object.name);
        if (object.revisions) {
            if (!Array.isArray(object.revisions))
                throw TypeError(".RepoSpec.revisions: array expected");
            message.revisions = [];
            for (let i = 0; i < object.revisions.length; ++i)
                message.revisions[i] = String(object.revisions[i]);
        }
        if (object.metadata != null) {
            if (typeof object.metadata !== "object")
                throw TypeError(".RepoSpec.metadata: object expected");
            message.metadata = $root.Metadata.fromObject(object.metadata);
        }
        if (object.walkSubmodules != null)
            message.walkSubmodules = Boolean(object.walkSubmodules);
        return message;
    };

    /**
     * Creates a plain object from a RepoSpec message. Also converts values to other types if specified.
     * @function toObject
     * @memberof RepoSpec
     * @static
     * @param {RepoSpec} message RepoSpec
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    RepoSpec.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.revisions = [];
        if (options.defaults) {
            object.path = "";
            object.name = "";
            object.metadata = null;
            object.walkSubmodules = false;
        }
        if (message.path != null && message.hasOwnProperty("path"))
            object.path = message.path;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.revisions && message.revisions.length) {
            object.revisions = [];
            for (let j = 0; j < message.revisions.length; ++j)
                object.revisions[j] = message.revisions[j];
        }
        if (message.metadata != null && message.hasOwnProperty("metadata"))
            object.metadata = $root.Metadata.toObject(message.metadata, options);
        if (message.walkSubmodules != null && message.hasOwnProperty("walkSubmodules"))
            object.walkSubmodules = message.walkSubmodules;
        return object;
    };

    /**
     * Converts this RepoSpec to JSON.
     * @function toJSON
     * @memberof RepoSpec
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    RepoSpec.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return RepoSpec;
})();

export const Query = $root.Query = (() => {

    /**
     * Properties of a Query.
     * @exports IQuery
     * @interface IQuery
     * @property {string|null} [line] Query line
     * @property {string|null} [file] Query file
     * @property {string|null} [repo] Query repo
     * @property {string|null} [tags] Query tags
     * @property {boolean|null} [foldCase] Query foldCase
     * @property {string|null} [notFile] Query notFile
     * @property {string|null} [notRepo] Query notRepo
     * @property {string|null} [notTags] Query notTags
     * @property {number|null} [maxMatches] Query maxMatches
     * @property {boolean|null} [filenameOnly] Query filenameOnly
     * @property {number|null} [contextLines] Query contextLines
     */

    /**
     * Constructs a new Query.
     * @exports Query
     * @classdesc Represents a Query.
     * @implements IQuery
     * @constructor
     * @param {IQuery=} [properties] Properties to set
     */
    function Query(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Query line.
     * @member {string} line
     * @memberof Query
     * @instance
     */
    Query.prototype.line = "";

    /**
     * Query file.
     * @member {string} file
     * @memberof Query
     * @instance
     */
    Query.prototype.file = "";

    /**
     * Query repo.
     * @member {string} repo
     * @memberof Query
     * @instance
     */
    Query.prototype.repo = "";

    /**
     * Query tags.
     * @member {string} tags
     * @memberof Query
     * @instance
     */
    Query.prototype.tags = "";

    /**
     * Query foldCase.
     * @member {boolean} foldCase
     * @memberof Query
     * @instance
     */
    Query.prototype.foldCase = false;

    /**
     * Query notFile.
     * @member {string} notFile
     * @memberof Query
     * @instance
     */
    Query.prototype.notFile = "";

    /**
     * Query notRepo.
     * @member {string} notRepo
     * @memberof Query
     * @instance
     */
    Query.prototype.notRepo = "";

    /**
     * Query notTags.
     * @member {string} notTags
     * @memberof Query
     * @instance
     */
    Query.prototype.notTags = "";

    /**
     * Query maxMatches.
     * @member {number} maxMatches
     * @memberof Query
     * @instance
     */
    Query.prototype.maxMatches = 0;

    /**
     * Query filenameOnly.
     * @member {boolean} filenameOnly
     * @memberof Query
     * @instance
     */
    Query.prototype.filenameOnly = false;

    /**
     * Query contextLines.
     * @member {number} contextLines
     * @memberof Query
     * @instance
     */
    Query.prototype.contextLines = 0;

    /**
     * Creates a new Query instance using the specified properties.
     * @function create
     * @memberof Query
     * @static
     * @param {IQuery=} [properties] Properties to set
     * @returns {Query} Query instance
     */
    Query.create = function create(properties) {
        return new Query(properties);
    };

    /**
     * Encodes the specified Query message. Does not implicitly {@link Query.verify|verify} messages.
     * @function encode
     * @memberof Query
     * @static
     * @param {IQuery} message Query message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Query.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.line != null && message.hasOwnProperty("line"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.line);
        if (message.file != null && message.hasOwnProperty("file"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.file);
        if (message.repo != null && message.hasOwnProperty("repo"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.repo);
        if (message.tags != null && message.hasOwnProperty("tags"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.tags);
        if (message.foldCase != null && message.hasOwnProperty("foldCase"))
            writer.uint32(/* id 5, wireType 0 =*/40).bool(message.foldCase);
        if (message.notFile != null && message.hasOwnProperty("notFile"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.notFile);
        if (message.notRepo != null && message.hasOwnProperty("notRepo"))
            writer.uint32(/* id 7, wireType 2 =*/58).string(message.notRepo);
        if (message.notTags != null && message.hasOwnProperty("notTags"))
            writer.uint32(/* id 8, wireType 2 =*/66).string(message.notTags);
        if (message.maxMatches != null && message.hasOwnProperty("maxMatches"))
            writer.uint32(/* id 9, wireType 0 =*/72).int32(message.maxMatches);
        if (message.filenameOnly != null && message.hasOwnProperty("filenameOnly"))
            writer.uint32(/* id 10, wireType 0 =*/80).bool(message.filenameOnly);
        if (message.contextLines != null && message.hasOwnProperty("contextLines"))
            writer.uint32(/* id 11, wireType 0 =*/88).int32(message.contextLines);
        return writer;
    };

    /**
     * Encodes the specified Query message, length delimited. Does not implicitly {@link Query.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Query
     * @static
     * @param {IQuery} message Query message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Query.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Query message from the specified reader or buffer.
     * @function decode
     * @memberof Query
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Query} Query
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Query.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Query();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.line = reader.string();
                break;
            case 2:
                message.file = reader.string();
                break;
            case 3:
                message.repo = reader.string();
                break;
            case 4:
                message.tags = reader.string();
                break;
            case 5:
                message.foldCase = reader.bool();
                break;
            case 6:
                message.notFile = reader.string();
                break;
            case 7:
                message.notRepo = reader.string();
                break;
            case 8:
                message.notTags = reader.string();
                break;
            case 9:
                message.maxMatches = reader.int32();
                break;
            case 10:
                message.filenameOnly = reader.bool();
                break;
            case 11:
                message.contextLines = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Query message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Query
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Query} Query
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Query.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Query message.
     * @function verify
     * @memberof Query
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Query.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.line != null && message.hasOwnProperty("line"))
            if (!$util.isString(message.line))
                return "line: string expected";
        if (message.file != null && message.hasOwnProperty("file"))
            if (!$util.isString(message.file))
                return "file: string expected";
        if (message.repo != null && message.hasOwnProperty("repo"))
            if (!$util.isString(message.repo))
                return "repo: string expected";
        if (message.tags != null && message.hasOwnProperty("tags"))
            if (!$util.isString(message.tags))
                return "tags: string expected";
        if (message.foldCase != null && message.hasOwnProperty("foldCase"))
            if (typeof message.foldCase !== "boolean")
                return "foldCase: boolean expected";
        if (message.notFile != null && message.hasOwnProperty("notFile"))
            if (!$util.isString(message.notFile))
                return "notFile: string expected";
        if (message.notRepo != null && message.hasOwnProperty("notRepo"))
            if (!$util.isString(message.notRepo))
                return "notRepo: string expected";
        if (message.notTags != null && message.hasOwnProperty("notTags"))
            if (!$util.isString(message.notTags))
                return "notTags: string expected";
        if (message.maxMatches != null && message.hasOwnProperty("maxMatches"))
            if (!$util.isInteger(message.maxMatches))
                return "maxMatches: integer expected";
        if (message.filenameOnly != null && message.hasOwnProperty("filenameOnly"))
            if (typeof message.filenameOnly !== "boolean")
                return "filenameOnly: boolean expected";
        if (message.contextLines != null && message.hasOwnProperty("contextLines"))
            if (!$util.isInteger(message.contextLines))
                return "contextLines: integer expected";
        return null;
    };

    /**
     * Creates a Query message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Query
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Query} Query
     */
    Query.fromObject = function fromObject(object) {
        if (object instanceof $root.Query)
            return object;
        let message = new $root.Query();
        if (object.line != null)
            message.line = String(object.line);
        if (object.file != null)
            message.file = String(object.file);
        if (object.repo != null)
            message.repo = String(object.repo);
        if (object.tags != null)
            message.tags = String(object.tags);
        if (object.foldCase != null)
            message.foldCase = Boolean(object.foldCase);
        if (object.notFile != null)
            message.notFile = String(object.notFile);
        if (object.notRepo != null)
            message.notRepo = String(object.notRepo);
        if (object.notTags != null)
            message.notTags = String(object.notTags);
        if (object.maxMatches != null)
            message.maxMatches = object.maxMatches | 0;
        if (object.filenameOnly != null)
            message.filenameOnly = Boolean(object.filenameOnly);
        if (object.contextLines != null)
            message.contextLines = object.contextLines | 0;
        return message;
    };

    /**
     * Creates a plain object from a Query message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Query
     * @static
     * @param {Query} message Query
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Query.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.line = "";
            object.file = "";
            object.repo = "";
            object.tags = "";
            object.foldCase = false;
            object.notFile = "";
            object.notRepo = "";
            object.notTags = "";
            object.maxMatches = 0;
            object.filenameOnly = false;
            object.contextLines = 0;
        }
        if (message.line != null && message.hasOwnProperty("line"))
            object.line = message.line;
        if (message.file != null && message.hasOwnProperty("file"))
            object.file = message.file;
        if (message.repo != null && message.hasOwnProperty("repo"))
            object.repo = message.repo;
        if (message.tags != null && message.hasOwnProperty("tags"))
            object.tags = message.tags;
        if (message.foldCase != null && message.hasOwnProperty("foldCase"))
            object.foldCase = message.foldCase;
        if (message.notFile != null && message.hasOwnProperty("notFile"))
            object.notFile = message.notFile;
        if (message.notRepo != null && message.hasOwnProperty("notRepo"))
            object.notRepo = message.notRepo;
        if (message.notTags != null && message.hasOwnProperty("notTags"))
            object.notTags = message.notTags;
        if (message.maxMatches != null && message.hasOwnProperty("maxMatches"))
            object.maxMatches = message.maxMatches;
        if (message.filenameOnly != null && message.hasOwnProperty("filenameOnly"))
            object.filenameOnly = message.filenameOnly;
        if (message.contextLines != null && message.hasOwnProperty("contextLines"))
            object.contextLines = message.contextLines;
        return object;
    };

    /**
     * Converts this Query to JSON.
     * @function toJSON
     * @memberof Query
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Query.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Query;
})();

export const Bounds = $root.Bounds = (() => {

    /**
     * Properties of a Bounds.
     * @exports IBounds
     * @interface IBounds
     * @property {number|null} [left] Bounds left
     * @property {number|null} [right] Bounds right
     */

    /**
     * Constructs a new Bounds.
     * @exports Bounds
     * @classdesc Represents a Bounds.
     * @implements IBounds
     * @constructor
     * @param {IBounds=} [properties] Properties to set
     */
    function Bounds(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Bounds left.
     * @member {number} left
     * @memberof Bounds
     * @instance
     */
    Bounds.prototype.left = 0;

    /**
     * Bounds right.
     * @member {number} right
     * @memberof Bounds
     * @instance
     */
    Bounds.prototype.right = 0;

    /**
     * Creates a new Bounds instance using the specified properties.
     * @function create
     * @memberof Bounds
     * @static
     * @param {IBounds=} [properties] Properties to set
     * @returns {Bounds} Bounds instance
     */
    Bounds.create = function create(properties) {
        return new Bounds(properties);
    };

    /**
     * Encodes the specified Bounds message. Does not implicitly {@link Bounds.verify|verify} messages.
     * @function encode
     * @memberof Bounds
     * @static
     * @param {IBounds} message Bounds message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Bounds.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.left != null && message.hasOwnProperty("left"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.left);
        if (message.right != null && message.hasOwnProperty("right"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.right);
        return writer;
    };

    /**
     * Encodes the specified Bounds message, length delimited. Does not implicitly {@link Bounds.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Bounds
     * @static
     * @param {IBounds} message Bounds message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Bounds.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Bounds message from the specified reader or buffer.
     * @function decode
     * @memberof Bounds
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Bounds} Bounds
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Bounds.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Bounds();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.left = reader.int32();
                break;
            case 2:
                message.right = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Bounds message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Bounds
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Bounds} Bounds
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Bounds.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Bounds message.
     * @function verify
     * @memberof Bounds
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Bounds.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.left != null && message.hasOwnProperty("left"))
            if (!$util.isInteger(message.left))
                return "left: integer expected";
        if (message.right != null && message.hasOwnProperty("right"))
            if (!$util.isInteger(message.right))
                return "right: integer expected";
        return null;
    };

    /**
     * Creates a Bounds message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Bounds
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Bounds} Bounds
     */
    Bounds.fromObject = function fromObject(object) {
        if (object instanceof $root.Bounds)
            return object;
        let message = new $root.Bounds();
        if (object.left != null)
            message.left = object.left | 0;
        if (object.right != null)
            message.right = object.right | 0;
        return message;
    };

    /**
     * Creates a plain object from a Bounds message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Bounds
     * @static
     * @param {Bounds} message Bounds
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Bounds.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.left = 0;
            object.right = 0;
        }
        if (message.left != null && message.hasOwnProperty("left"))
            object.left = message.left;
        if (message.right != null && message.hasOwnProperty("right"))
            object.right = message.right;
        return object;
    };

    /**
     * Converts this Bounds to JSON.
     * @function toJSON
     * @memberof Bounds
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Bounds.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Bounds;
})();

export const SearchResult = $root.SearchResult = (() => {

    /**
     * Properties of a SearchResult.
     * @exports ISearchResult
     * @interface ISearchResult
     * @property {string|null} [tree] SearchResult tree
     * @property {string|null} [version] SearchResult version
     * @property {string|null} [path] SearchResult path
     * @property {number|Long|null} [lineNumber] SearchResult lineNumber
     * @property {Array.<string>|null} [contextBefore] SearchResult contextBefore
     * @property {Array.<string>|null} [contextAfter] SearchResult contextAfter
     * @property {IBounds|null} [bounds] SearchResult bounds
     * @property {string|null} [line] SearchResult line
     */

    /**
     * Constructs a new SearchResult.
     * @exports SearchResult
     * @classdesc Represents a SearchResult.
     * @implements ISearchResult
     * @constructor
     * @param {ISearchResult=} [properties] Properties to set
     */
    function SearchResult(properties) {
        this.contextBefore = [];
        this.contextAfter = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * SearchResult tree.
     * @member {string} tree
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.tree = "";

    /**
     * SearchResult version.
     * @member {string} version
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.version = "";

    /**
     * SearchResult path.
     * @member {string} path
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.path = "";

    /**
     * SearchResult lineNumber.
     * @member {number|Long} lineNumber
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.lineNumber = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchResult contextBefore.
     * @member {Array.<string>} contextBefore
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.contextBefore = $util.emptyArray;

    /**
     * SearchResult contextAfter.
     * @member {Array.<string>} contextAfter
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.contextAfter = $util.emptyArray;

    /**
     * SearchResult bounds.
     * @member {IBounds|null|undefined} bounds
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.bounds = null;

    /**
     * SearchResult line.
     * @member {string} line
     * @memberof SearchResult
     * @instance
     */
    SearchResult.prototype.line = "";

    /**
     * Creates a new SearchResult instance using the specified properties.
     * @function create
     * @memberof SearchResult
     * @static
     * @param {ISearchResult=} [properties] Properties to set
     * @returns {SearchResult} SearchResult instance
     */
    SearchResult.create = function create(properties) {
        return new SearchResult(properties);
    };

    /**
     * Encodes the specified SearchResult message. Does not implicitly {@link SearchResult.verify|verify} messages.
     * @function encode
     * @memberof SearchResult
     * @static
     * @param {ISearchResult} message SearchResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SearchResult.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.tree != null && message.hasOwnProperty("tree"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.tree);
        if (message.version != null && message.hasOwnProperty("version"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.version);
        if (message.path != null && message.hasOwnProperty("path"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.path);
        if (message.lineNumber != null && message.hasOwnProperty("lineNumber"))
            writer.uint32(/* id 4, wireType 0 =*/32).int64(message.lineNumber);
        if (message.contextBefore != null && message.contextBefore.length)
            for (let i = 0; i < message.contextBefore.length; ++i)
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.contextBefore[i]);
        if (message.contextAfter != null && message.contextAfter.length)
            for (let i = 0; i < message.contextAfter.length; ++i)
                writer.uint32(/* id 6, wireType 2 =*/50).string(message.contextAfter[i]);
        if (message.bounds != null && message.hasOwnProperty("bounds"))
            $root.Bounds.encode(message.bounds, writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
        if (message.line != null && message.hasOwnProperty("line"))
            writer.uint32(/* id 8, wireType 2 =*/66).string(message.line);
        return writer;
    };

    /**
     * Encodes the specified SearchResult message, length delimited. Does not implicitly {@link SearchResult.verify|verify} messages.
     * @function encodeDelimited
     * @memberof SearchResult
     * @static
     * @param {ISearchResult} message SearchResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SearchResult.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a SearchResult message from the specified reader or buffer.
     * @function decode
     * @memberof SearchResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {SearchResult} SearchResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SearchResult.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.SearchResult();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.tree = reader.string();
                break;
            case 2:
                message.version = reader.string();
                break;
            case 3:
                message.path = reader.string();
                break;
            case 4:
                message.lineNumber = reader.int64();
                break;
            case 5:
                if (!(message.contextBefore && message.contextBefore.length))
                    message.contextBefore = [];
                message.contextBefore.push(reader.string());
                break;
            case 6:
                if (!(message.contextAfter && message.contextAfter.length))
                    message.contextAfter = [];
                message.contextAfter.push(reader.string());
                break;
            case 7:
                message.bounds = $root.Bounds.decode(reader, reader.uint32());
                break;
            case 8:
                message.line = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a SearchResult message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof SearchResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {SearchResult} SearchResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SearchResult.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a SearchResult message.
     * @function verify
     * @memberof SearchResult
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    SearchResult.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.tree != null && message.hasOwnProperty("tree"))
            if (!$util.isString(message.tree))
                return "tree: string expected";
        if (message.version != null && message.hasOwnProperty("version"))
            if (!$util.isString(message.version))
                return "version: string expected";
        if (message.path != null && message.hasOwnProperty("path"))
            if (!$util.isString(message.path))
                return "path: string expected";
        if (message.lineNumber != null && message.hasOwnProperty("lineNumber"))
            if (!$util.isInteger(message.lineNumber) && !(message.lineNumber && $util.isInteger(message.lineNumber.low) && $util.isInteger(message.lineNumber.high)))
                return "lineNumber: integer|Long expected";
        if (message.contextBefore != null && message.hasOwnProperty("contextBefore")) {
            if (!Array.isArray(message.contextBefore))
                return "contextBefore: array expected";
            for (let i = 0; i < message.contextBefore.length; ++i)
                if (!$util.isString(message.contextBefore[i]))
                    return "contextBefore: string[] expected";
        }
        if (message.contextAfter != null && message.hasOwnProperty("contextAfter")) {
            if (!Array.isArray(message.contextAfter))
                return "contextAfter: array expected";
            for (let i = 0; i < message.contextAfter.length; ++i)
                if (!$util.isString(message.contextAfter[i]))
                    return "contextAfter: string[] expected";
        }
        if (message.bounds != null && message.hasOwnProperty("bounds")) {
            let error = $root.Bounds.verify(message.bounds);
            if (error)
                return "bounds." + error;
        }
        if (message.line != null && message.hasOwnProperty("line"))
            if (!$util.isString(message.line))
                return "line: string expected";
        return null;
    };

    /**
     * Creates a SearchResult message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof SearchResult
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {SearchResult} SearchResult
     */
    SearchResult.fromObject = function fromObject(object) {
        if (object instanceof $root.SearchResult)
            return object;
        let message = new $root.SearchResult();
        if (object.tree != null)
            message.tree = String(object.tree);
        if (object.version != null)
            message.version = String(object.version);
        if (object.path != null)
            message.path = String(object.path);
        if (object.lineNumber != null)
            if ($util.Long)
                (message.lineNumber = $util.Long.fromValue(object.lineNumber)).unsigned = false;
            else if (typeof object.lineNumber === "string")
                message.lineNumber = parseInt(object.lineNumber, 10);
            else if (typeof object.lineNumber === "number")
                message.lineNumber = object.lineNumber;
            else if (typeof object.lineNumber === "object")
                message.lineNumber = new $util.LongBits(object.lineNumber.low >>> 0, object.lineNumber.high >>> 0).toNumber();
        if (object.contextBefore) {
            if (!Array.isArray(object.contextBefore))
                throw TypeError(".SearchResult.contextBefore: array expected");
            message.contextBefore = [];
            for (let i = 0; i < object.contextBefore.length; ++i)
                message.contextBefore[i] = String(object.contextBefore[i]);
        }
        if (object.contextAfter) {
            if (!Array.isArray(object.contextAfter))
                throw TypeError(".SearchResult.contextAfter: array expected");
            message.contextAfter = [];
            for (let i = 0; i < object.contextAfter.length; ++i)
                message.contextAfter[i] = String(object.contextAfter[i]);
        }
        if (object.bounds != null) {
            if (typeof object.bounds !== "object")
                throw TypeError(".SearchResult.bounds: object expected");
            message.bounds = $root.Bounds.fromObject(object.bounds);
        }
        if (object.line != null)
            message.line = String(object.line);
        return message;
    };

    /**
     * Creates a plain object from a SearchResult message. Also converts values to other types if specified.
     * @function toObject
     * @memberof SearchResult
     * @static
     * @param {SearchResult} message SearchResult
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    SearchResult.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults) {
            object.contextBefore = [];
            object.contextAfter = [];
        }
        if (options.defaults) {
            object.tree = "";
            object.version = "";
            object.path = "";
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.lineNumber = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.lineNumber = options.longs === String ? "0" : 0;
            object.bounds = null;
            object.line = "";
        }
        if (message.tree != null && message.hasOwnProperty("tree"))
            object.tree = message.tree;
        if (message.version != null && message.hasOwnProperty("version"))
            object.version = message.version;
        if (message.path != null && message.hasOwnProperty("path"))
            object.path = message.path;
        if (message.lineNumber != null && message.hasOwnProperty("lineNumber"))
            if (typeof message.lineNumber === "number")
                object.lineNumber = options.longs === String ? String(message.lineNumber) : message.lineNumber;
            else
                object.lineNumber = options.longs === String ? $util.Long.prototype.toString.call(message.lineNumber) : options.longs === Number ? new $util.LongBits(message.lineNumber.low >>> 0, message.lineNumber.high >>> 0).toNumber() : message.lineNumber;
        if (message.contextBefore && message.contextBefore.length) {
            object.contextBefore = [];
            for (let j = 0; j < message.contextBefore.length; ++j)
                object.contextBefore[j] = message.contextBefore[j];
        }
        if (message.contextAfter && message.contextAfter.length) {
            object.contextAfter = [];
            for (let j = 0; j < message.contextAfter.length; ++j)
                object.contextAfter[j] = message.contextAfter[j];
        }
        if (message.bounds != null && message.hasOwnProperty("bounds"))
            object.bounds = $root.Bounds.toObject(message.bounds, options);
        if (message.line != null && message.hasOwnProperty("line"))
            object.line = message.line;
        return object;
    };

    /**
     * Converts this SearchResult to JSON.
     * @function toJSON
     * @memberof SearchResult
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    SearchResult.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return SearchResult;
})();

export const FileResult = $root.FileResult = (() => {

    /**
     * Properties of a FileResult.
     * @exports IFileResult
     * @interface IFileResult
     * @property {string|null} [tree] FileResult tree
     * @property {string|null} [version] FileResult version
     * @property {string|null} [path] FileResult path
     * @property {IBounds|null} [bounds] FileResult bounds
     */

    /**
     * Constructs a new FileResult.
     * @exports FileResult
     * @classdesc Represents a FileResult.
     * @implements IFileResult
     * @constructor
     * @param {IFileResult=} [properties] Properties to set
     */
    function FileResult(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FileResult tree.
     * @member {string} tree
     * @memberof FileResult
     * @instance
     */
    FileResult.prototype.tree = "";

    /**
     * FileResult version.
     * @member {string} version
     * @memberof FileResult
     * @instance
     */
    FileResult.prototype.version = "";

    /**
     * FileResult path.
     * @member {string} path
     * @memberof FileResult
     * @instance
     */
    FileResult.prototype.path = "";

    /**
     * FileResult bounds.
     * @member {IBounds|null|undefined} bounds
     * @memberof FileResult
     * @instance
     */
    FileResult.prototype.bounds = null;

    /**
     * Creates a new FileResult instance using the specified properties.
     * @function create
     * @memberof FileResult
     * @static
     * @param {IFileResult=} [properties] Properties to set
     * @returns {FileResult} FileResult instance
     */
    FileResult.create = function create(properties) {
        return new FileResult(properties);
    };

    /**
     * Encodes the specified FileResult message. Does not implicitly {@link FileResult.verify|verify} messages.
     * @function encode
     * @memberof FileResult
     * @static
     * @param {IFileResult} message FileResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FileResult.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.tree != null && message.hasOwnProperty("tree"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.tree);
        if (message.version != null && message.hasOwnProperty("version"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.version);
        if (message.path != null && message.hasOwnProperty("path"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.path);
        if (message.bounds != null && message.hasOwnProperty("bounds"))
            $root.Bounds.encode(message.bounds, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified FileResult message, length delimited. Does not implicitly {@link FileResult.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FileResult
     * @static
     * @param {IFileResult} message FileResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FileResult.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FileResult message from the specified reader or buffer.
     * @function decode
     * @memberof FileResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FileResult} FileResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FileResult.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.FileResult();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.tree = reader.string();
                break;
            case 2:
                message.version = reader.string();
                break;
            case 3:
                message.path = reader.string();
                break;
            case 4:
                message.bounds = $root.Bounds.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FileResult message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FileResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FileResult} FileResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FileResult.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FileResult message.
     * @function verify
     * @memberof FileResult
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FileResult.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.tree != null && message.hasOwnProperty("tree"))
            if (!$util.isString(message.tree))
                return "tree: string expected";
        if (message.version != null && message.hasOwnProperty("version"))
            if (!$util.isString(message.version))
                return "version: string expected";
        if (message.path != null && message.hasOwnProperty("path"))
            if (!$util.isString(message.path))
                return "path: string expected";
        if (message.bounds != null && message.hasOwnProperty("bounds")) {
            let error = $root.Bounds.verify(message.bounds);
            if (error)
                return "bounds." + error;
        }
        return null;
    };

    /**
     * Creates a FileResult message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FileResult
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FileResult} FileResult
     */
    FileResult.fromObject = function fromObject(object) {
        if (object instanceof $root.FileResult)
            return object;
        let message = new $root.FileResult();
        if (object.tree != null)
            message.tree = String(object.tree);
        if (object.version != null)
            message.version = String(object.version);
        if (object.path != null)
            message.path = String(object.path);
        if (object.bounds != null) {
            if (typeof object.bounds !== "object")
                throw TypeError(".FileResult.bounds: object expected");
            message.bounds = $root.Bounds.fromObject(object.bounds);
        }
        return message;
    };

    /**
     * Creates a plain object from a FileResult message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FileResult
     * @static
     * @param {FileResult} message FileResult
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FileResult.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.tree = "";
            object.version = "";
            object.path = "";
            object.bounds = null;
        }
        if (message.tree != null && message.hasOwnProperty("tree"))
            object.tree = message.tree;
        if (message.version != null && message.hasOwnProperty("version"))
            object.version = message.version;
        if (message.path != null && message.hasOwnProperty("path"))
            object.path = message.path;
        if (message.bounds != null && message.hasOwnProperty("bounds"))
            object.bounds = $root.Bounds.toObject(message.bounds, options);
        return object;
    };

    /**
     * Converts this FileResult to JSON.
     * @function toJSON
     * @memberof FileResult
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FileResult.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FileResult;
})();

export const SearchStats = $root.SearchStats = (() => {

    /**
     * Properties of a SearchStats.
     * @exports ISearchStats
     * @interface ISearchStats
     * @property {number|Long|null} [re2Time] SearchStats re2Time
     * @property {number|Long|null} [gitTime] SearchStats gitTime
     * @property {number|Long|null} [sortTime] SearchStats sortTime
     * @property {number|Long|null} [indexTime] SearchStats indexTime
     * @property {number|Long|null} [analyzeTime] SearchStats analyzeTime
     * @property {number|Long|null} [totalTime] SearchStats totalTime
     * @property {SearchStats.ExitReason|null} [exitReason] SearchStats exitReason
     */

    /**
     * Constructs a new SearchStats.
     * @exports SearchStats
     * @classdesc Represents a SearchStats.
     * @implements ISearchStats
     * @constructor
     * @param {ISearchStats=} [properties] Properties to set
     */
    function SearchStats(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * SearchStats re2Time.
     * @member {number|Long} re2Time
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.re2Time = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchStats gitTime.
     * @member {number|Long} gitTime
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.gitTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchStats sortTime.
     * @member {number|Long} sortTime
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.sortTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchStats indexTime.
     * @member {number|Long} indexTime
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.indexTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchStats analyzeTime.
     * @member {number|Long} analyzeTime
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.analyzeTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchStats totalTime.
     * @member {number|Long} totalTime
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.totalTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * SearchStats exitReason.
     * @member {SearchStats.ExitReason} exitReason
     * @memberof SearchStats
     * @instance
     */
    SearchStats.prototype.exitReason = 0;

    /**
     * Creates a new SearchStats instance using the specified properties.
     * @function create
     * @memberof SearchStats
     * @static
     * @param {ISearchStats=} [properties] Properties to set
     * @returns {SearchStats} SearchStats instance
     */
    SearchStats.create = function create(properties) {
        return new SearchStats(properties);
    };

    /**
     * Encodes the specified SearchStats message. Does not implicitly {@link SearchStats.verify|verify} messages.
     * @function encode
     * @memberof SearchStats
     * @static
     * @param {ISearchStats} message SearchStats message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SearchStats.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.re2Time != null && message.hasOwnProperty("re2Time"))
            writer.uint32(/* id 1, wireType 0 =*/8).int64(message.re2Time);
        if (message.gitTime != null && message.hasOwnProperty("gitTime"))
            writer.uint32(/* id 2, wireType 0 =*/16).int64(message.gitTime);
        if (message.sortTime != null && message.hasOwnProperty("sortTime"))
            writer.uint32(/* id 3, wireType 0 =*/24).int64(message.sortTime);
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            writer.uint32(/* id 4, wireType 0 =*/32).int64(message.indexTime);
        if (message.analyzeTime != null && message.hasOwnProperty("analyzeTime"))
            writer.uint32(/* id 5, wireType 0 =*/40).int64(message.analyzeTime);
        if (message.exitReason != null && message.hasOwnProperty("exitReason"))
            writer.uint32(/* id 6, wireType 0 =*/48).int32(message.exitReason);
        if (message.totalTime != null && message.hasOwnProperty("totalTime"))
            writer.uint32(/* id 7, wireType 0 =*/56).int64(message.totalTime);
        return writer;
    };

    /**
     * Encodes the specified SearchStats message, length delimited. Does not implicitly {@link SearchStats.verify|verify} messages.
     * @function encodeDelimited
     * @memberof SearchStats
     * @static
     * @param {ISearchStats} message SearchStats message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SearchStats.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a SearchStats message from the specified reader or buffer.
     * @function decode
     * @memberof SearchStats
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {SearchStats} SearchStats
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SearchStats.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.SearchStats();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.re2Time = reader.int64();
                break;
            case 2:
                message.gitTime = reader.int64();
                break;
            case 3:
                message.sortTime = reader.int64();
                break;
            case 4:
                message.indexTime = reader.int64();
                break;
            case 5:
                message.analyzeTime = reader.int64();
                break;
            case 7:
                message.totalTime = reader.int64();
                break;
            case 6:
                message.exitReason = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a SearchStats message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof SearchStats
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {SearchStats} SearchStats
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SearchStats.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a SearchStats message.
     * @function verify
     * @memberof SearchStats
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    SearchStats.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.re2Time != null && message.hasOwnProperty("re2Time"))
            if (!$util.isInteger(message.re2Time) && !(message.re2Time && $util.isInteger(message.re2Time.low) && $util.isInteger(message.re2Time.high)))
                return "re2Time: integer|Long expected";
        if (message.gitTime != null && message.hasOwnProperty("gitTime"))
            if (!$util.isInteger(message.gitTime) && !(message.gitTime && $util.isInteger(message.gitTime.low) && $util.isInteger(message.gitTime.high)))
                return "gitTime: integer|Long expected";
        if (message.sortTime != null && message.hasOwnProperty("sortTime"))
            if (!$util.isInteger(message.sortTime) && !(message.sortTime && $util.isInteger(message.sortTime.low) && $util.isInteger(message.sortTime.high)))
                return "sortTime: integer|Long expected";
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            if (!$util.isInteger(message.indexTime) && !(message.indexTime && $util.isInteger(message.indexTime.low) && $util.isInteger(message.indexTime.high)))
                return "indexTime: integer|Long expected";
        if (message.analyzeTime != null && message.hasOwnProperty("analyzeTime"))
            if (!$util.isInteger(message.analyzeTime) && !(message.analyzeTime && $util.isInteger(message.analyzeTime.low) && $util.isInteger(message.analyzeTime.high)))
                return "analyzeTime: integer|Long expected";
        if (message.totalTime != null && message.hasOwnProperty("totalTime"))
            if (!$util.isInteger(message.totalTime) && !(message.totalTime && $util.isInteger(message.totalTime.low) && $util.isInteger(message.totalTime.high)))
                return "totalTime: integer|Long expected";
        if (message.exitReason != null && message.hasOwnProperty("exitReason"))
            switch (message.exitReason) {
            default:
                return "exitReason: enum value expected";
            case 0:
            case 1:
            case 2:
                break;
            }
        return null;
    };

    /**
     * Creates a SearchStats message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof SearchStats
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {SearchStats} SearchStats
     */
    SearchStats.fromObject = function fromObject(object) {
        if (object instanceof $root.SearchStats)
            return object;
        let message = new $root.SearchStats();
        if (object.re2Time != null)
            if ($util.Long)
                (message.re2Time = $util.Long.fromValue(object.re2Time)).unsigned = false;
            else if (typeof object.re2Time === "string")
                message.re2Time = parseInt(object.re2Time, 10);
            else if (typeof object.re2Time === "number")
                message.re2Time = object.re2Time;
            else if (typeof object.re2Time === "object")
                message.re2Time = new $util.LongBits(object.re2Time.low >>> 0, object.re2Time.high >>> 0).toNumber();
        if (object.gitTime != null)
            if ($util.Long)
                (message.gitTime = $util.Long.fromValue(object.gitTime)).unsigned = false;
            else if (typeof object.gitTime === "string")
                message.gitTime = parseInt(object.gitTime, 10);
            else if (typeof object.gitTime === "number")
                message.gitTime = object.gitTime;
            else if (typeof object.gitTime === "object")
                message.gitTime = new $util.LongBits(object.gitTime.low >>> 0, object.gitTime.high >>> 0).toNumber();
        if (object.sortTime != null)
            if ($util.Long)
                (message.sortTime = $util.Long.fromValue(object.sortTime)).unsigned = false;
            else if (typeof object.sortTime === "string")
                message.sortTime = parseInt(object.sortTime, 10);
            else if (typeof object.sortTime === "number")
                message.sortTime = object.sortTime;
            else if (typeof object.sortTime === "object")
                message.sortTime = new $util.LongBits(object.sortTime.low >>> 0, object.sortTime.high >>> 0).toNumber();
        if (object.indexTime != null)
            if ($util.Long)
                (message.indexTime = $util.Long.fromValue(object.indexTime)).unsigned = false;
            else if (typeof object.indexTime === "string")
                message.indexTime = parseInt(object.indexTime, 10);
            else if (typeof object.indexTime === "number")
                message.indexTime = object.indexTime;
            else if (typeof object.indexTime === "object")
                message.indexTime = new $util.LongBits(object.indexTime.low >>> 0, object.indexTime.high >>> 0).toNumber();
        if (object.analyzeTime != null)
            if ($util.Long)
                (message.analyzeTime = $util.Long.fromValue(object.analyzeTime)).unsigned = false;
            else if (typeof object.analyzeTime === "string")
                message.analyzeTime = parseInt(object.analyzeTime, 10);
            else if (typeof object.analyzeTime === "number")
                message.analyzeTime = object.analyzeTime;
            else if (typeof object.analyzeTime === "object")
                message.analyzeTime = new $util.LongBits(object.analyzeTime.low >>> 0, object.analyzeTime.high >>> 0).toNumber();
        if (object.totalTime != null)
            if ($util.Long)
                (message.totalTime = $util.Long.fromValue(object.totalTime)).unsigned = false;
            else if (typeof object.totalTime === "string")
                message.totalTime = parseInt(object.totalTime, 10);
            else if (typeof object.totalTime === "number")
                message.totalTime = object.totalTime;
            else if (typeof object.totalTime === "object")
                message.totalTime = new $util.LongBits(object.totalTime.low >>> 0, object.totalTime.high >>> 0).toNumber();
        switch (object.exitReason) {
        case "NONE":
        case 0:
            message.exitReason = 0;
            break;
        case "TIMEOUT":
        case 1:
            message.exitReason = 1;
            break;
        case "MATCH_LIMIT":
        case 2:
            message.exitReason = 2;
            break;
        }
        return message;
    };

    /**
     * Creates a plain object from a SearchStats message. Also converts values to other types if specified.
     * @function toObject
     * @memberof SearchStats
     * @static
     * @param {SearchStats} message SearchStats
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    SearchStats.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.re2Time = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.re2Time = options.longs === String ? "0" : 0;
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.gitTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.gitTime = options.longs === String ? "0" : 0;
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.sortTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.sortTime = options.longs === String ? "0" : 0;
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.indexTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.indexTime = options.longs === String ? "0" : 0;
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.analyzeTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.analyzeTime = options.longs === String ? "0" : 0;
            object.exitReason = options.enums === String ? "NONE" : 0;
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.totalTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.totalTime = options.longs === String ? "0" : 0;
        }
        if (message.re2Time != null && message.hasOwnProperty("re2Time"))
            if (typeof message.re2Time === "number")
                object.re2Time = options.longs === String ? String(message.re2Time) : message.re2Time;
            else
                object.re2Time = options.longs === String ? $util.Long.prototype.toString.call(message.re2Time) : options.longs === Number ? new $util.LongBits(message.re2Time.low >>> 0, message.re2Time.high >>> 0).toNumber() : message.re2Time;
        if (message.gitTime != null && message.hasOwnProperty("gitTime"))
            if (typeof message.gitTime === "number")
                object.gitTime = options.longs === String ? String(message.gitTime) : message.gitTime;
            else
                object.gitTime = options.longs === String ? $util.Long.prototype.toString.call(message.gitTime) : options.longs === Number ? new $util.LongBits(message.gitTime.low >>> 0, message.gitTime.high >>> 0).toNumber() : message.gitTime;
        if (message.sortTime != null && message.hasOwnProperty("sortTime"))
            if (typeof message.sortTime === "number")
                object.sortTime = options.longs === String ? String(message.sortTime) : message.sortTime;
            else
                object.sortTime = options.longs === String ? $util.Long.prototype.toString.call(message.sortTime) : options.longs === Number ? new $util.LongBits(message.sortTime.low >>> 0, message.sortTime.high >>> 0).toNumber() : message.sortTime;
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            if (typeof message.indexTime === "number")
                object.indexTime = options.longs === String ? String(message.indexTime) : message.indexTime;
            else
                object.indexTime = options.longs === String ? $util.Long.prototype.toString.call(message.indexTime) : options.longs === Number ? new $util.LongBits(message.indexTime.low >>> 0, message.indexTime.high >>> 0).toNumber() : message.indexTime;
        if (message.analyzeTime != null && message.hasOwnProperty("analyzeTime"))
            if (typeof message.analyzeTime === "number")
                object.analyzeTime = options.longs === String ? String(message.analyzeTime) : message.analyzeTime;
            else
                object.analyzeTime = options.longs === String ? $util.Long.prototype.toString.call(message.analyzeTime) : options.longs === Number ? new $util.LongBits(message.analyzeTime.low >>> 0, message.analyzeTime.high >>> 0).toNumber() : message.analyzeTime;
        if (message.exitReason != null && message.hasOwnProperty("exitReason"))
            object.exitReason = options.enums === String ? $root.SearchStats.ExitReason[message.exitReason] : message.exitReason;
        if (message.totalTime != null && message.hasOwnProperty("totalTime"))
            if (typeof message.totalTime === "number")
                object.totalTime = options.longs === String ? String(message.totalTime) : message.totalTime;
            else
                object.totalTime = options.longs === String ? $util.Long.prototype.toString.call(message.totalTime) : options.longs === Number ? new $util.LongBits(message.totalTime.low >>> 0, message.totalTime.high >>> 0).toNumber() : message.totalTime;
        return object;
    };

    /**
     * Converts this SearchStats to JSON.
     * @function toJSON
     * @memberof SearchStats
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    SearchStats.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * ExitReason enum.
     * @name SearchStats.ExitReason
     * @enum {string}
     * @property {number} NONE=0 NONE value
     * @property {number} TIMEOUT=1 TIMEOUT value
     * @property {number} MATCH_LIMIT=2 MATCH_LIMIT value
     */
    SearchStats.ExitReason = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "NONE"] = 0;
        values[valuesById[1] = "TIMEOUT"] = 1;
        values[valuesById[2] = "MATCH_LIMIT"] = 2;
        return values;
    })();

    return SearchStats;
})();

export const ServerInfo = $root.ServerInfo = (() => {

    /**
     * Properties of a ServerInfo.
     * @exports IServerInfo
     * @interface IServerInfo
     * @property {string|null} [name] ServerInfo name
     * @property {Array.<ServerInfo.ITree>|null} [trees] ServerInfo trees
     * @property {boolean|null} [hasTags] ServerInfo hasTags
     * @property {number|Long|null} [indexTime] ServerInfo indexTime
     */

    /**
     * Constructs a new ServerInfo.
     * @exports ServerInfo
     * @classdesc Represents a ServerInfo.
     * @implements IServerInfo
     * @constructor
     * @param {IServerInfo=} [properties] Properties to set
     */
    function ServerInfo(properties) {
        this.trees = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ServerInfo name.
     * @member {string} name
     * @memberof ServerInfo
     * @instance
     */
    ServerInfo.prototype.name = "";

    /**
     * ServerInfo trees.
     * @member {Array.<ServerInfo.ITree>} trees
     * @memberof ServerInfo
     * @instance
     */
    ServerInfo.prototype.trees = $util.emptyArray;

    /**
     * ServerInfo hasTags.
     * @member {boolean} hasTags
     * @memberof ServerInfo
     * @instance
     */
    ServerInfo.prototype.hasTags = false;

    /**
     * ServerInfo indexTime.
     * @member {number|Long} indexTime
     * @memberof ServerInfo
     * @instance
     */
    ServerInfo.prototype.indexTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * Creates a new ServerInfo instance using the specified properties.
     * @function create
     * @memberof ServerInfo
     * @static
     * @param {IServerInfo=} [properties] Properties to set
     * @returns {ServerInfo} ServerInfo instance
     */
    ServerInfo.create = function create(properties) {
        return new ServerInfo(properties);
    };

    /**
     * Encodes the specified ServerInfo message. Does not implicitly {@link ServerInfo.verify|verify} messages.
     * @function encode
     * @memberof ServerInfo
     * @static
     * @param {IServerInfo} message ServerInfo message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ServerInfo.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
        if (message.trees != null && message.trees.length)
            for (let i = 0; i < message.trees.length; ++i)
                $root.ServerInfo.Tree.encode(message.trees[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.hasTags != null && message.hasOwnProperty("hasTags"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.hasTags);
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            writer.uint32(/* id 4, wireType 0 =*/32).int64(message.indexTime);
        return writer;
    };

    /**
     * Encodes the specified ServerInfo message, length delimited. Does not implicitly {@link ServerInfo.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ServerInfo
     * @static
     * @param {IServerInfo} message ServerInfo message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ServerInfo.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ServerInfo message from the specified reader or buffer.
     * @function decode
     * @memberof ServerInfo
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ServerInfo} ServerInfo
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ServerInfo.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ServerInfo();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.name = reader.string();
                break;
            case 2:
                if (!(message.trees && message.trees.length))
                    message.trees = [];
                message.trees.push($root.ServerInfo.Tree.decode(reader, reader.uint32()));
                break;
            case 3:
                message.hasTags = reader.bool();
                break;
            case 4:
                message.indexTime = reader.int64();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ServerInfo message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ServerInfo
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ServerInfo} ServerInfo
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ServerInfo.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ServerInfo message.
     * @function verify
     * @memberof ServerInfo
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ServerInfo.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.trees != null && message.hasOwnProperty("trees")) {
            if (!Array.isArray(message.trees))
                return "trees: array expected";
            for (let i = 0; i < message.trees.length; ++i) {
                let error = $root.ServerInfo.Tree.verify(message.trees[i]);
                if (error)
                    return "trees." + error;
            }
        }
        if (message.hasTags != null && message.hasOwnProperty("hasTags"))
            if (typeof message.hasTags !== "boolean")
                return "hasTags: boolean expected";
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            if (!$util.isInteger(message.indexTime) && !(message.indexTime && $util.isInteger(message.indexTime.low) && $util.isInteger(message.indexTime.high)))
                return "indexTime: integer|Long expected";
        return null;
    };

    /**
     * Creates a ServerInfo message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ServerInfo
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ServerInfo} ServerInfo
     */
    ServerInfo.fromObject = function fromObject(object) {
        if (object instanceof $root.ServerInfo)
            return object;
        let message = new $root.ServerInfo();
        if (object.name != null)
            message.name = String(object.name);
        if (object.trees) {
            if (!Array.isArray(object.trees))
                throw TypeError(".ServerInfo.trees: array expected");
            message.trees = [];
            for (let i = 0; i < object.trees.length; ++i) {
                if (typeof object.trees[i] !== "object")
                    throw TypeError(".ServerInfo.trees: object expected");
                message.trees[i] = $root.ServerInfo.Tree.fromObject(object.trees[i]);
            }
        }
        if (object.hasTags != null)
            message.hasTags = Boolean(object.hasTags);
        if (object.indexTime != null)
            if ($util.Long)
                (message.indexTime = $util.Long.fromValue(object.indexTime)).unsigned = false;
            else if (typeof object.indexTime === "string")
                message.indexTime = parseInt(object.indexTime, 10);
            else if (typeof object.indexTime === "number")
                message.indexTime = object.indexTime;
            else if (typeof object.indexTime === "object")
                message.indexTime = new $util.LongBits(object.indexTime.low >>> 0, object.indexTime.high >>> 0).toNumber();
        return message;
    };

    /**
     * Creates a plain object from a ServerInfo message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ServerInfo
     * @static
     * @param {ServerInfo} message ServerInfo
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ServerInfo.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.trees = [];
        if (options.defaults) {
            object.name = "";
            object.hasTags = false;
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.indexTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.indexTime = options.longs === String ? "0" : 0;
        }
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.trees && message.trees.length) {
            object.trees = [];
            for (let j = 0; j < message.trees.length; ++j)
                object.trees[j] = $root.ServerInfo.Tree.toObject(message.trees[j], options);
        }
        if (message.hasTags != null && message.hasOwnProperty("hasTags"))
            object.hasTags = message.hasTags;
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            if (typeof message.indexTime === "number")
                object.indexTime = options.longs === String ? String(message.indexTime) : message.indexTime;
            else
                object.indexTime = options.longs === String ? $util.Long.prototype.toString.call(message.indexTime) : options.longs === Number ? new $util.LongBits(message.indexTime.low >>> 0, message.indexTime.high >>> 0).toNumber() : message.indexTime;
        return object;
    };

    /**
     * Converts this ServerInfo to JSON.
     * @function toJSON
     * @memberof ServerInfo
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ServerInfo.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    ServerInfo.Tree = (function() {

        /**
         * Properties of a Tree.
         * @memberof ServerInfo
         * @interface ITree
         * @property {string|null} [name] Tree name
         * @property {string|null} [version] Tree version
         * @property {IMetadata|null} [metadata] Tree metadata
         */

        /**
         * Constructs a new Tree.
         * @memberof ServerInfo
         * @classdesc Represents a Tree.
         * @implements ITree
         * @constructor
         * @param {ServerInfo.ITree=} [properties] Properties to set
         */
        function Tree(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Tree name.
         * @member {string} name
         * @memberof ServerInfo.Tree
         * @instance
         */
        Tree.prototype.name = "";

        /**
         * Tree version.
         * @member {string} version
         * @memberof ServerInfo.Tree
         * @instance
         */
        Tree.prototype.version = "";

        /**
         * Tree metadata.
         * @member {IMetadata|null|undefined} metadata
         * @memberof ServerInfo.Tree
         * @instance
         */
        Tree.prototype.metadata = null;

        /**
         * Creates a new Tree instance using the specified properties.
         * @function create
         * @memberof ServerInfo.Tree
         * @static
         * @param {ServerInfo.ITree=} [properties] Properties to set
         * @returns {ServerInfo.Tree} Tree instance
         */
        Tree.create = function create(properties) {
            return new Tree(properties);
        };

        /**
         * Encodes the specified Tree message. Does not implicitly {@link ServerInfo.Tree.verify|verify} messages.
         * @function encode
         * @memberof ServerInfo.Tree
         * @static
         * @param {ServerInfo.ITree} message Tree message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Tree.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.name != null && message.hasOwnProperty("name"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
            if (message.version != null && message.hasOwnProperty("version"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.version);
            if (message.metadata != null && message.hasOwnProperty("metadata"))
                $root.Metadata.encode(message.metadata, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified Tree message, length delimited. Does not implicitly {@link ServerInfo.Tree.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ServerInfo.Tree
         * @static
         * @param {ServerInfo.ITree} message Tree message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Tree.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a Tree message from the specified reader or buffer.
         * @function decode
         * @memberof ServerInfo.Tree
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ServerInfo.Tree} Tree
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Tree.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ServerInfo.Tree();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.version = reader.string();
                    break;
                case 3:
                    message.metadata = $root.Metadata.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a Tree message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ServerInfo.Tree
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ServerInfo.Tree} Tree
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Tree.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a Tree message.
         * @function verify
         * @memberof ServerInfo.Tree
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        Tree.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            if (message.version != null && message.hasOwnProperty("version"))
                if (!$util.isString(message.version))
                    return "version: string expected";
            if (message.metadata != null && message.hasOwnProperty("metadata")) {
                let error = $root.Metadata.verify(message.metadata);
                if (error)
                    return "metadata." + error;
            }
            return null;
        };

        /**
         * Creates a Tree message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ServerInfo.Tree
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ServerInfo.Tree} Tree
         */
        Tree.fromObject = function fromObject(object) {
            if (object instanceof $root.ServerInfo.Tree)
                return object;
            let message = new $root.ServerInfo.Tree();
            if (object.name != null)
                message.name = String(object.name);
            if (object.version != null)
                message.version = String(object.version);
            if (object.metadata != null) {
                if (typeof object.metadata !== "object")
                    throw TypeError(".ServerInfo.Tree.metadata: object expected");
                message.metadata = $root.Metadata.fromObject(object.metadata);
            }
            return message;
        };

        /**
         * Creates a plain object from a Tree message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ServerInfo.Tree
         * @static
         * @param {ServerInfo.Tree} message Tree
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Tree.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.name = "";
                object.version = "";
                object.metadata = null;
            }
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            if (message.version != null && message.hasOwnProperty("version"))
                object.version = message.version;
            if (message.metadata != null && message.hasOwnProperty("metadata"))
                object.metadata = $root.Metadata.toObject(message.metadata, options);
            return object;
        };

        /**
         * Converts this Tree to JSON.
         * @function toJSON
         * @memberof ServerInfo.Tree
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Tree.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return Tree;
    })();

    return ServerInfo;
})();

export const CodeSearchResult = $root.CodeSearchResult = (() => {

    /**
     * Properties of a CodeSearchResult.
     * @exports ICodeSearchResult
     * @interface ICodeSearchResult
     * @property {ISearchStats|null} [stats] CodeSearchResult stats
     * @property {Array.<ISearchResult>|null} [results] CodeSearchResult results
     * @property {Array.<IFileResult>|null} [fileResults] CodeSearchResult fileResults
     * @property {string|null} [indexName] CodeSearchResult indexName
     * @property {number|Long|null} [indexTime] CodeSearchResult indexTime
     */

    /**
     * Constructs a new CodeSearchResult.
     * @exports CodeSearchResult
     * @classdesc Represents a CodeSearchResult.
     * @implements ICodeSearchResult
     * @constructor
     * @param {ICodeSearchResult=} [properties] Properties to set
     */
    function CodeSearchResult(properties) {
        this.results = [];
        this.fileResults = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CodeSearchResult stats.
     * @member {ISearchStats|null|undefined} stats
     * @memberof CodeSearchResult
     * @instance
     */
    CodeSearchResult.prototype.stats = null;

    /**
     * CodeSearchResult results.
     * @member {Array.<ISearchResult>} results
     * @memberof CodeSearchResult
     * @instance
     */
    CodeSearchResult.prototype.results = $util.emptyArray;

    /**
     * CodeSearchResult fileResults.
     * @member {Array.<IFileResult>} fileResults
     * @memberof CodeSearchResult
     * @instance
     */
    CodeSearchResult.prototype.fileResults = $util.emptyArray;

    /**
     * CodeSearchResult indexName.
     * @member {string} indexName
     * @memberof CodeSearchResult
     * @instance
     */
    CodeSearchResult.prototype.indexName = "";

    /**
     * CodeSearchResult indexTime.
     * @member {number|Long} indexTime
     * @memberof CodeSearchResult
     * @instance
     */
    CodeSearchResult.prototype.indexTime = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * Creates a new CodeSearchResult instance using the specified properties.
     * @function create
     * @memberof CodeSearchResult
     * @static
     * @param {ICodeSearchResult=} [properties] Properties to set
     * @returns {CodeSearchResult} CodeSearchResult instance
     */
    CodeSearchResult.create = function create(properties) {
        return new CodeSearchResult(properties);
    };

    /**
     * Encodes the specified CodeSearchResult message. Does not implicitly {@link CodeSearchResult.verify|verify} messages.
     * @function encode
     * @memberof CodeSearchResult
     * @static
     * @param {ICodeSearchResult} message CodeSearchResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CodeSearchResult.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.stats != null && message.hasOwnProperty("stats"))
            $root.SearchStats.encode(message.stats, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.results != null && message.results.length)
            for (let i = 0; i < message.results.length; ++i)
                $root.SearchResult.encode(message.results[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.fileResults != null && message.fileResults.length)
            for (let i = 0; i < message.fileResults.length; ++i)
                $root.FileResult.encode(message.fileResults[i], writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.indexName != null && message.hasOwnProperty("indexName"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.indexName);
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            writer.uint32(/* id 5, wireType 0 =*/40).int64(message.indexTime);
        return writer;
    };

    /**
     * Encodes the specified CodeSearchResult message, length delimited. Does not implicitly {@link CodeSearchResult.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CodeSearchResult
     * @static
     * @param {ICodeSearchResult} message CodeSearchResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CodeSearchResult.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CodeSearchResult message from the specified reader or buffer.
     * @function decode
     * @memberof CodeSearchResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CodeSearchResult} CodeSearchResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CodeSearchResult.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CodeSearchResult();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.stats = $root.SearchStats.decode(reader, reader.uint32());
                break;
            case 2:
                if (!(message.results && message.results.length))
                    message.results = [];
                message.results.push($root.SearchResult.decode(reader, reader.uint32()));
                break;
            case 3:
                if (!(message.fileResults && message.fileResults.length))
                    message.fileResults = [];
                message.fileResults.push($root.FileResult.decode(reader, reader.uint32()));
                break;
            case 4:
                message.indexName = reader.string();
                break;
            case 5:
                message.indexTime = reader.int64();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CodeSearchResult message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CodeSearchResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CodeSearchResult} CodeSearchResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CodeSearchResult.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CodeSearchResult message.
     * @function verify
     * @memberof CodeSearchResult
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CodeSearchResult.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.stats != null && message.hasOwnProperty("stats")) {
            let error = $root.SearchStats.verify(message.stats);
            if (error)
                return "stats." + error;
        }
        if (message.results != null && message.hasOwnProperty("results")) {
            if (!Array.isArray(message.results))
                return "results: array expected";
            for (let i = 0; i < message.results.length; ++i) {
                let error = $root.SearchResult.verify(message.results[i]);
                if (error)
                    return "results." + error;
            }
        }
        if (message.fileResults != null && message.hasOwnProperty("fileResults")) {
            if (!Array.isArray(message.fileResults))
                return "fileResults: array expected";
            for (let i = 0; i < message.fileResults.length; ++i) {
                let error = $root.FileResult.verify(message.fileResults[i]);
                if (error)
                    return "fileResults." + error;
            }
        }
        if (message.indexName != null && message.hasOwnProperty("indexName"))
            if (!$util.isString(message.indexName))
                return "indexName: string expected";
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            if (!$util.isInteger(message.indexTime) && !(message.indexTime && $util.isInteger(message.indexTime.low) && $util.isInteger(message.indexTime.high)))
                return "indexTime: integer|Long expected";
        return null;
    };

    /**
     * Creates a CodeSearchResult message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CodeSearchResult
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CodeSearchResult} CodeSearchResult
     */
    CodeSearchResult.fromObject = function fromObject(object) {
        if (object instanceof $root.CodeSearchResult)
            return object;
        let message = new $root.CodeSearchResult();
        if (object.stats != null) {
            if (typeof object.stats !== "object")
                throw TypeError(".CodeSearchResult.stats: object expected");
            message.stats = $root.SearchStats.fromObject(object.stats);
        }
        if (object.results) {
            if (!Array.isArray(object.results))
                throw TypeError(".CodeSearchResult.results: array expected");
            message.results = [];
            for (let i = 0; i < object.results.length; ++i) {
                if (typeof object.results[i] !== "object")
                    throw TypeError(".CodeSearchResult.results: object expected");
                message.results[i] = $root.SearchResult.fromObject(object.results[i]);
            }
        }
        if (object.fileResults) {
            if (!Array.isArray(object.fileResults))
                throw TypeError(".CodeSearchResult.fileResults: array expected");
            message.fileResults = [];
            for (let i = 0; i < object.fileResults.length; ++i) {
                if (typeof object.fileResults[i] !== "object")
                    throw TypeError(".CodeSearchResult.fileResults: object expected");
                message.fileResults[i] = $root.FileResult.fromObject(object.fileResults[i]);
            }
        }
        if (object.indexName != null)
            message.indexName = String(object.indexName);
        if (object.indexTime != null)
            if ($util.Long)
                (message.indexTime = $util.Long.fromValue(object.indexTime)).unsigned = false;
            else if (typeof object.indexTime === "string")
                message.indexTime = parseInt(object.indexTime, 10);
            else if (typeof object.indexTime === "number")
                message.indexTime = object.indexTime;
            else if (typeof object.indexTime === "object")
                message.indexTime = new $util.LongBits(object.indexTime.low >>> 0, object.indexTime.high >>> 0).toNumber();
        return message;
    };

    /**
     * Creates a plain object from a CodeSearchResult message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CodeSearchResult
     * @static
     * @param {CodeSearchResult} message CodeSearchResult
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CodeSearchResult.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults) {
            object.results = [];
            object.fileResults = [];
        }
        if (options.defaults) {
            object.stats = null;
            object.indexName = "";
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.indexTime = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.indexTime = options.longs === String ? "0" : 0;
        }
        if (message.stats != null && message.hasOwnProperty("stats"))
            object.stats = $root.SearchStats.toObject(message.stats, options);
        if (message.results && message.results.length) {
            object.results = [];
            for (let j = 0; j < message.results.length; ++j)
                object.results[j] = $root.SearchResult.toObject(message.results[j], options);
        }
        if (message.fileResults && message.fileResults.length) {
            object.fileResults = [];
            for (let j = 0; j < message.fileResults.length; ++j)
                object.fileResults[j] = $root.FileResult.toObject(message.fileResults[j], options);
        }
        if (message.indexName != null && message.hasOwnProperty("indexName"))
            object.indexName = message.indexName;
        if (message.indexTime != null && message.hasOwnProperty("indexTime"))
            if (typeof message.indexTime === "number")
                object.indexTime = options.longs === String ? String(message.indexTime) : message.indexTime;
            else
                object.indexTime = options.longs === String ? $util.Long.prototype.toString.call(message.indexTime) : options.longs === Number ? new $util.LongBits(message.indexTime.low >>> 0, message.indexTime.high >>> 0).toNumber() : message.indexTime;
        return object;
    };

    /**
     * Converts this CodeSearchResult to JSON.
     * @function toJSON
     * @memberof CodeSearchResult
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CodeSearchResult.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CodeSearchResult;
})();

export const InfoRequest = $root.InfoRequest = (() => {

    /**
     * Properties of an InfoRequest.
     * @exports IInfoRequest
     * @interface IInfoRequest
     */

    /**
     * Constructs a new InfoRequest.
     * @exports InfoRequest
     * @classdesc Represents an InfoRequest.
     * @implements IInfoRequest
     * @constructor
     * @param {IInfoRequest=} [properties] Properties to set
     */
    function InfoRequest(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new InfoRequest instance using the specified properties.
     * @function create
     * @memberof InfoRequest
     * @static
     * @param {IInfoRequest=} [properties] Properties to set
     * @returns {InfoRequest} InfoRequest instance
     */
    InfoRequest.create = function create(properties) {
        return new InfoRequest(properties);
    };

    /**
     * Encodes the specified InfoRequest message. Does not implicitly {@link InfoRequest.verify|verify} messages.
     * @function encode
     * @memberof InfoRequest
     * @static
     * @param {IInfoRequest} message InfoRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InfoRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified InfoRequest message, length delimited. Does not implicitly {@link InfoRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof InfoRequest
     * @static
     * @param {IInfoRequest} message InfoRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InfoRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an InfoRequest message from the specified reader or buffer.
     * @function decode
     * @memberof InfoRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {InfoRequest} InfoRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InfoRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.InfoRequest();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an InfoRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof InfoRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {InfoRequest} InfoRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InfoRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an InfoRequest message.
     * @function verify
     * @memberof InfoRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    InfoRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates an InfoRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof InfoRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {InfoRequest} InfoRequest
     */
    InfoRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.InfoRequest)
            return object;
        return new $root.InfoRequest();
    };

    /**
     * Creates a plain object from an InfoRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof InfoRequest
     * @static
     * @param {InfoRequest} message InfoRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    InfoRequest.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this InfoRequest to JSON.
     * @function toJSON
     * @memberof InfoRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    InfoRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return InfoRequest;
})();

export const Empty = $root.Empty = (() => {

    /**
     * Properties of an Empty.
     * @exports IEmpty
     * @interface IEmpty
     */

    /**
     * Constructs a new Empty.
     * @exports Empty
     * @classdesc Represents an Empty.
     * @implements IEmpty
     * @constructor
     * @param {IEmpty=} [properties] Properties to set
     */
    function Empty(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new Empty instance using the specified properties.
     * @function create
     * @memberof Empty
     * @static
     * @param {IEmpty=} [properties] Properties to set
     * @returns {Empty} Empty instance
     */
    Empty.create = function create(properties) {
        return new Empty(properties);
    };

    /**
     * Encodes the specified Empty message. Does not implicitly {@link Empty.verify|verify} messages.
     * @function encode
     * @memberof Empty
     * @static
     * @param {IEmpty} message Empty message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Empty.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified Empty message, length delimited. Does not implicitly {@link Empty.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Empty
     * @static
     * @param {IEmpty} message Empty message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Empty.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Empty message from the specified reader or buffer.
     * @function decode
     * @memberof Empty
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Empty} Empty
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Empty.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Empty();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Empty message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Empty
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Empty} Empty
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Empty.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Empty message.
     * @function verify
     * @memberof Empty
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Empty.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates an Empty message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Empty
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Empty} Empty
     */
    Empty.fromObject = function fromObject(object) {
        if (object instanceof $root.Empty)
            return object;
        return new $root.Empty();
    };

    /**
     * Creates a plain object from an Empty message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Empty
     * @static
     * @param {Empty} message Empty
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Empty.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this Empty to JSON.
     * @function toJSON
     * @memberof Empty
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Empty.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Empty;
})();

export const CodeSearch = $root.CodeSearch = (() => {

    /**
     * Constructs a new CodeSearch service.
     * @exports CodeSearch
     * @classdesc Represents a CodeSearch
     * @extends $protobuf.rpc.Service
     * @constructor
     * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
     * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
     * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
     */
    function CodeSearch(rpcImpl, requestDelimited, responseDelimited) {
        $protobuf.rpc.Service.call(this, rpcImpl, requestDelimited, responseDelimited);
    }

    (CodeSearch.prototype = Object.create($protobuf.rpc.Service.prototype)).constructor = CodeSearch;

    /**
     * Creates new CodeSearch service using the specified rpc implementation.
     * @function create
     * @memberof CodeSearch
     * @static
     * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
     * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
     * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
     * @returns {CodeSearch} RPC service. Useful where requests and/or responses are streamed.
     */
    CodeSearch.create = function create(rpcImpl, requestDelimited, responseDelimited) {
        return new this(rpcImpl, requestDelimited, responseDelimited);
    };

    /**
     * Callback as used by {@link CodeSearch#info}.
     * @memberof CodeSearch
     * @typedef InfoCallback
     * @type {function}
     * @param {Error|null} error Error, if any
     * @param {ServerInfo} [response] ServerInfo
     */

    /**
     * Calls Info.
     * @function info
     * @memberof CodeSearch
     * @instance
     * @param {IInfoRequest} request InfoRequest message or plain object
     * @param {CodeSearch.InfoCallback} callback Node-style callback called with the error, if any, and ServerInfo
     * @returns {undefined}
     * @variation 1
     */
    Object.defineProperty(CodeSearch.prototype.info = function info(request, callback) {
        return this.rpcCall(info, $root.InfoRequest, $root.ServerInfo, request, callback);
    }, "name", { value: "Info" });

    /**
     * Calls Info.
     * @function info
     * @memberof CodeSearch
     * @instance
     * @param {IInfoRequest} request InfoRequest message or plain object
     * @returns {Promise<ServerInfo>} Promise
     * @variation 2
     */

    /**
     * Callback as used by {@link CodeSearch#search}.
     * @memberof CodeSearch
     * @typedef SearchCallback
     * @type {function}
     * @param {Error|null} error Error, if any
     * @param {CodeSearchResult} [response] CodeSearchResult
     */

    /**
     * Calls Search.
     * @function search
     * @memberof CodeSearch
     * @instance
     * @param {IQuery} request Query message or plain object
     * @param {CodeSearch.SearchCallback} callback Node-style callback called with the error, if any, and CodeSearchResult
     * @returns {undefined}
     * @variation 1
     */
    Object.defineProperty(CodeSearch.prototype.search = function search(request, callback) {
        return this.rpcCall(search, $root.Query, $root.CodeSearchResult, request, callback);
    }, "name", { value: "Search" });

    /**
     * Calls Search.
     * @function search
     * @memberof CodeSearch
     * @instance
     * @param {IQuery} request Query message or plain object
     * @returns {Promise<CodeSearchResult>} Promise
     * @variation 2
     */

    /**
     * Callback as used by {@link CodeSearch#reload}.
     * @memberof CodeSearch
     * @typedef ReloadCallback
     * @type {function}
     * @param {Error|null} error Error, if any
     * @param {Empty} [response] Empty
     */

    /**
     * Calls Reload.
     * @function reload
     * @memberof CodeSearch
     * @instance
     * @param {IEmpty} request Empty message or plain object
     * @param {CodeSearch.ReloadCallback} callback Node-style callback called with the error, if any, and Empty
     * @returns {undefined}
     * @variation 1
     */
    Object.defineProperty(CodeSearch.prototype.reload = function reload(request, callback) {
        return this.rpcCall(reload, $root.Empty, $root.Empty, request, callback);
    }, "name", { value: "Reload" });

    /**
     * Calls Reload.
     * @function reload
     * @memberof CodeSearch
     * @instance
     * @param {IEmpty} request Empty message or plain object
     * @returns {Promise<Empty>} Promise
     * @variation 2
     */

    return CodeSearch;
})();

export { $root as default };
