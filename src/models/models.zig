const std = @import("std");
const testing = std.testing;
const json = std.json;

const StorageType = enum { disk, blobstorage, kafka, pulsar };

pub const IngestionJob = struct {
    id: u64,
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

    pub fn jsonStringify(self: IngestionJob, jw: anytype) !void {
        // root object
        try jw.beginObject();

        try jw.objectField("id");
        try jw.write(self.id);

        try jw.objectField("name");
        try jw.write(self.name);

        try jw.objectField("type");
        try jw.write(self.type_);

        // ingestion configuration struct
        try jw.objectField("ingestionConfiguration");
        try jw.beginObject();
        try jw.objectField("kafka");
        try jw.beginObject();
        try jw.objectField("brokers");
        try jw.write(self.ingestionConfiguration.kafka.brokers);
        try jw.objectField("topic");
        try jw.write(self.ingestionConfiguration.kafka.topic);
        try jw.objectField("offset");
        try jw.write(self.ingestionConfiguration.kafka.offset);
        try jw.endObject();
        try jw.endObject();

        // segment configuration
        try jw.objectField("segmentConfiguration");
        try jw.beginObject();
        try jw.objectField("table");
        try jw.write(self.segmentConfiguration.table);
        try jw.objectField("timeColumnName");
        try jw.write(self.segmentConfiguration.timeColumnName);
        try jw.endObject();

        try jw.endObject();
    }
};

test "IngestionJob jsonStringify" {
    const allocator = std.testing.allocator;

    const job = IngestionJob{
        .id = 123,
        .name = "test_job",
        .type_ = .kafka,
        .ingestionConfiguration = .{
            .kafka = .{
                .brokers = "broker1,broker2",
                .topic = "test_topic",
                .offset = "latest",
            },
            .pulsar = .{},
            .disk = .{},
            .blobstorage = .{},
        },
        .segmentConfiguration = .{
            .table = "test_table",
            .timeColumnName = "timestamp",
        },
    };

    const result = try json.stringifyAlloc(allocator, job, .{});
    defer allocator.free(result);

    const expected =
        \\{"id":123,"name":"test_job","type":"kafka","ingestionConfiguration":{"kafka":{"brokers":"broker1,broker2","topic":"test_topic","offset":"latest"}},"segmentConfiguration":{"table":"test_table","timeColumnName":"timestamp"}}
    ;

    try testing.expectEqualStrings(expected, result);
}

const FieldType = enum([]u8) { int = "int", string = "string", float = "float", long = "long", unixtimestamp = "unixtimestamp" };

pub const TableSpec = struct {
    id: u64,
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
