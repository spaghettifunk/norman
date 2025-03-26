const std = @import("std");
const json = std.json;
const fs = std.fs;
const mem = std.mem;

pub const MetadataService = struct {
    filePath: []const u8,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator) !MetadataService {
        const filePath = "metadata.db";

        // Check if the file exists
        fs.cwd().access(filePath, .{}) catch |err| {
            if (err == std.posix.AccessError.FileNotFound) {
                // File does not exist, create it with empty JSON array
                var file = try fs.cwd().createFile(filePath, .{});
                defer file.close();
                try file.writer().writeAll("[]"); // Initialize with an empty JSON array
            }
        };

        return MetadataService{
            .filePath = filePath,
            .allocator = allocator,
        };
    }

    pub fn insertLine(self: MetadataService, newLine: struct { id: i64, name: []const u8, city: []const u8 }) !void {
        var data = try self.readJson();
        defer data.deinit();
        try insertLineImpl(self.allocator, &data, newLine);
        try self.writeJson(data);
    }

    pub fn deleteLine(self: MetadataService, idToDelete: i64) !void {
        var data = try self.readJson();
        defer data.deinit();
        try deleteLineImpl(&data, idToDelete);
        try self.writeJson(data);
    }

    pub fn searchLine(self: MetadataService, field: []const u8, value: []const u8) !?json.Object {
        var data = try self.readJson();
        defer data.deinit();
        return searchLineImpl(data, field, value);
    }

    fn readJson(self: MetadataService) !json.Array {
        var file = try fs.cwd().openFile(self.filePath, .{});
        defer file.close();

        var buffer: [1024]u8 = undefined;
        const bytesRead = try file.reader().readAll(&buffer);
        const fileContent = buffer[0..bytesRead];

        var parser = json.Parser.init(fileContent);
        const root = try parser.parse();
        return try root.array();
    }

    fn writeJson(self: MetadataService, data: json.Array) !void {
        var file = try fs.cwd().createFile(self.filePath, .{ .truncate = true });
        defer file.close();
        const writer = file.writer();
        var stringifier = json.Stringifier.init(writer);
        try stringifier.stringify(data.any());
    }
};

fn insertLineImpl(allocator: std.mem.Allocator, data: *json.Array, newLine: struct { id: i64, name: []const u8, city: []const u8 }) !void {
    var obj = try json.Object.init(allocator);
    try obj.addValue("id", .{ .integer = newLine.id });
    try obj.addValue("name", .{ .string = newLine.name });
    try obj.addValue("city", .{ .string = newLine.city });
    try data.append(.{ .object = obj });
}

fn deleteLineImpl(data: *json.Array, idToDelete: i64) !void {
    var i: usize = 0;
    while (i < data.len) {
        var obj = try data.items[i].object();
        if (try obj.getValue("id").integer() == idToDelete) {
            _ = data.orderedRemove(i);
            return;
        }
        i += 1;
    }
}

fn searchLineImpl(data: json.Array, field: []const u8, value: []const u8) !?json.Object {
    for (data.items) |item| {
        var obj = try item.object();
        if (try mem.eql(u8, obj.getValue(field).string(), value)) {
            return obj;
        }
    }
    return null;
}
