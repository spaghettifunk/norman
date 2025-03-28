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

const FieldType = enum {
    int,
    string,
    float,
    long,
    unixtimestamp,
};

pub const Table = struct {
    id: u64,
    name: []const u8,
    schema: struct {
        dimensions: []const Dimension,
        metrics: []const Metric,
        datetime: Datetime,
    },

    pub fn jsonStringify(self: Table, jw: anytype) !void {
        try jw.beginObject();

        try jw.objectField("id");
        try jw.write(self.id);
        try jw.objectField("name");
        try jw.write(self.name);
        try jw.objectField("schema");
        try jw.beginObject();
        try jw.objectField("dimensions");
        try jw.beginArray();
        for (self.schema.dimensions) |dimension| {
            try dimension.jsonStringify(jw);
        }
        try jw.endArray();
        try jw.objectField("metrics");
        try jw.beginArray();
        for (self.schema.metrics) |metric| {
            try metric.jsonStringify(jw);
        }
        try jw.endArray();
        try jw.objectField("datetime");
        try self.schema.datetime.jsonStringify(jw);
        try jw.endObject();
        try jw.endObject();
    }
};

pub const Dimension = struct {
    name: []const u8,
    type_: FieldType,

    pub fn jsonStringify(self: Dimension, jw: anytype) !void {
        try jw.beginObject();
        try jw.objectField("name");
        try jw.write(self.name);
        try jw.objectField("type");
        try jw.write(self.type_);
        try jw.endObject();
    }
};

pub const Metric = struct {
    name: []const u8,
    type_: FieldType,

    pub fn jsonStringify(self: Metric, jw: anytype) !void {
        try jw.beginObject();
        try jw.objectField("name");
        try jw.write(self.name);
        try jw.objectField("type");
        try jw.write(self.type_);
        try jw.endObject();
    }
};

pub const Datetime = struct {
    name: []const u8,
    type_: FieldType,

    pub fn jsonStringify(self: Datetime, jw: anytype) !void {
        try jw.beginObject();
        try jw.objectField("name");
        try jw.write(self.name);
        try jw.objectField("type");
        try jw.write(self.type_);
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

test "Table jsonStringify test" {
    const allocator = std.testing.allocator;

    const table = Table{
        .id = 123,
        .name = "my_table",
        .schema = .{
            .dimensions = &[_]Dimension{
                .{ .name = "city", .type_ = .string },
                .{ .name = "year", .type_ = .int },
            },
            .metrics = &[_]Metric{
                .{ .name = "sales", .type_ = .float },
                .{ .name = "count", .type_ = .long },
            },
            .datetime = .{ .name = "timestamp", .type_ = .unixtimestamp },
        },
    };

    const result = try json.stringifyAlloc(allocator, table, .{});
    defer allocator.free(result);

    const expected =
        \\{"id":123,"name":"my_table","schema":{"dimensions":[{"name":"city","type":"string"},{"name":"year","type":"int"}],"metrics":[{"name":"sales","type":"float"},{"name":"count","type":"long"}],"datetime":{"name":"timestamp","type":"unixtimestamp"}}}
    ;

    try testing.expectEqualStrings(expected, result);
}

test "Dimension jsonStringify test" {
    const allocator = std.testing.allocator;
    const dimension = Dimension{ .name = "city", .type_ = .string };

    const result = try json.stringifyAlloc(allocator, dimension, .{});
    defer allocator.free(result);

    const expected = "{\"name\":\"city\",\"type\":\"string\"}";
    try testing.expectEqualStrings(expected, result);
}

test "Metric jsonStringify test" {
    const allocator = std.testing.allocator;
    const metric = Metric{ .name = "sales", .type_ = .float };

    const result = try json.stringifyAlloc(allocator, metric, .{});
    defer allocator.free(result);

    const expected = "{\"name\":\"sales\",\"type\":\"float\"}";
    try testing.expectEqualStrings(expected, result);
}

test "Datetime jsonStringify test" {
    const allocator = std.testing.allocator;
    const datetime = Datetime{ .name = "timestamp", .type_ = .unixtimestamp };

    const result = try json.stringifyAlloc(allocator, datetime, .{});
    defer allocator.free(result);

    const expected = "{\"name\":\"timestamp\",\"type\":\"unixtimestamp\"}";
    try testing.expectEqualStrings(expected, result);
}
