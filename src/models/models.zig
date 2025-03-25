const StorageType = enum([]u8) { disk = "disk", blobstorage = "blobstorage", kafka = "kafka", pulsar = "pulsar" };

pub const IngestionJob = struct {
    name: []const u8,
    type_: StorageType,
    ingestionConfiguration: struct {
        kafka: struct {
            brokers: []const u8,
            topic: []const u8,
            offset: []const u8,
        },
        pulsar: struct {},
        disk: struct {},
        blobstorage: struct {},
    },
    segmentConfiguration: struct {
        table: []const u8,
        timeColumnName: []const u8,
    },
};

const FieldType = enum([]u8) { int = "int", string = "string", float = "float", long = "long", unixtimestamp = "unixtimestamp" };

pub const TableSpec = struct {
    name: []const u8,
    schema: struct {
        dimensions: []Dimension,
        metrics: []Metric,
        datetime: Datetime,
    },
};

pub const Dimension = struct {
    name: []const u8,
    type_: FieldType,
};

pub const Metric = struct {
    name: []const u8,
    type_: FieldType,
};

pub const Datetime = struct {
    name: []const u8,
    type_: FieldType,
};
