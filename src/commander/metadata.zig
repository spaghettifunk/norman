const std = @import("std");
const uuid = @import("uuid");
const json = std.json;
const fs = std.fs;
const mem = std.mem;
const assert = std.debug.assert;

const models = @import("../models/models.zig");

pub const MetadataDB = struct {
    tables: std.ArrayList(models.Table),
    ingestionJobs: std.ArrayList(models.IngestionJob),

    pub fn jsonStringify(self: MetadataDB, jw: anytype) !void {
        // root object
        try jw.beginObject();

        try jw.objectField("tables");
        try jw.beginObject();

        for (self.tables.items) |table| {
            try table.jsonStringify(jw);
        }

        try jw.endObject();

        try jw.objectField("ingestionJobs");
        try jw.beginObject();

        for (self.ingestionJobs.items) |ij| {
            try ij.jsonStringify(jw);
        }

        try jw.endObject();

        try jw.endObject();
    }
};

pub const MetadataService = struct {
    filePath: []const u8 = "metadata.db",
    allocator: std.mem.Allocator,
    db: MetadataDB,

    pub fn init(allocator: std.mem.Allocator) !MetadataService {
        var md = MetadataService{
            .allocator = allocator,
            .db = undefined,
        };

        // Check if the file exists
        fs.cwd().access(md.filePath, .{}) catch |err| {
            switch (err) {
                std.posix.AccessError.FileNotFound => {
                    // File does not exist, create it with empty JSON array
                    var file = try fs.cwd().createFile(md.filePath, .{});
                    defer file.close();
                    try file.writer().writeAll("{\"tables\":[], \"ingestionJobs\":[]}");
                },
                else => {
                    // something went wrong
                    std.debug.print("{any}\n", .{err});
                    std.process.exit(1);
                },
            }
        };

        // load db in memory
        try md.readDBFile();

        return md;
    }

    pub fn store(self: MetadataService) !void {
        // 1. Read the original file.
        const originalFile = try std.fs.openFileAbsolute(self.filePath, .{ .mode = .read_write });
        defer originalFile.close();

        const originalSize = try originalFile.getEndPos();
        const originalBuffer = try std.heap.page_allocator.alloc(u8, originalSize);
        defer std.heap.page_allocator.free(originalBuffer);

        const bytes_read = try originalFile.read(originalBuffer);
        if (bytes_read != originalSize) {
            return error.UnexpectedEndOfFile;
        }

        // 2. Create a temporary copy.
        const tempPath = "metadata.temp.db";
        const temp_file = try std.fs.createFileAbsolute(tempPath, .{});
        defer temp_file.close();

        _ = try temp_file.write(originalBuffer);

        // 3. Overwrite the original file with new data.
        try originalFile.seekTo(0); // Reset the file pointer to the beginning.
        try originalFile.writeAll(""); // Clear the contents of the file.

        const result = try json.stringifyAlloc(self.allocator, self.db, .{});
        defer self.allocator.free(result);

        _ = try originalFile.write(result);

        // 4. Delete the temporary copy.
        try std.fs.deleteFileAbsolute(tempPath);
    }

    pub fn addTable(self: MetadataService, table: models.Table) !void {
        const id = uuid.v7.new();
        table.id = id;
        try self.db.tables.append(table);
    }

    pub fn addInjestionJob(self: MetadataService, ij: models.IngestionJob) !void {
        const id = uuid.v7.new();
        ij.id = id;
        try self.db.ingestionJobs.append(ij);
    }

    pub fn deleteTable(self: MetadataService, table: models.Table) void {
        var i: u8 = 0;
        while (i < self.db.tables.items.len) : (i += 1) {
            if (self.db.tables.items[i].id == table.id) {
                break;
            }
        }
        const t = self.db.tables.swapRemove(i);
        assert(t.id == table.id);
    }

    pub fn deleteInjestionJob(self: MetadataService, ij: models.IngestionJob) void {
        var i: u8 = 0;
        while (i < self.db.ingestionJobs.items.len) : (i += 1) {
            if (self.db.ingestionJobs.items[i].id == ij.id) {
                break;
            }
        }
        const t = self.db.ingestionJobs.swapRemove(i);
        assert(t.id == ij.id);
    }

    fn readDBFile(self: *MetadataService) !void {
        // TODO: validate if 512 is a sufficient value
        const data = try std.fs.cwd().readFileAlloc(self.allocator, self.filePath, 512);
        defer self.allocator.free(data);

        const result = try std.json.parseFromSlice(MetadataDB, self.allocator, data, .{});
        const db = result.value;
        self.db = db;
    }
};
