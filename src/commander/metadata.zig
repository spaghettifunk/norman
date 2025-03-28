const std = @import("std");
const json = std.json;
const fs = std.fs;
const mem = std.mem;

const models = @import("../models/models.zig");

pub const MetadataDB = struct {
    tables: []models.TableSpec,
    ingestionJobs: []models.IngestionJob,
};

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
