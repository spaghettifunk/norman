const std = @import("std");
const httpz = @import("httpz");
const configuration = @import("../configuration.zig");
const metadata = @import("metadata.zig");
const version = @import("../const.zig").normanVersion;

var server_instance: ?*httpz.Server(void) = null;

pub const Commander = struct {
    pub fn start(allocator: std.mem.Allocator, config: configuration.Config) !void {
        std.posix.sigaction(std.posix.SIG.INT, &.{
            .handler = .{ .handler = shutdown },
            .mask = std.posix.empty_sigset,
            .flags = 0,
        }, null);
        std.posix.sigaction(std.posix.SIG.TERM, &.{
            .handler = .{ .handler = shutdown },
            .mask = std.posix.empty_sigset,
            .flags = 0,
        }, null);

        var server = try httpz.Server(void).init(allocator, .{ .port = config.commander.port }, {});
        defer server.deinit();

        var router = try server.router(.{});
        router.get("/", index, .{});
        router.post("/api/tables", getTables, .{});
        router.post("/api/tables", createTable, .{});
        router.post("/api/tables/:name", getTableByName, .{});
        router.post("/api/ingestions", getIngestions, .{});
        router.post("/api/ingestions", createIngestion, .{});
        router.post("/api/ingestions/:id", getIngestionById, .{});
        router.get("/metrics", metrics, .{});

        std.debug.print("Listening on http://{s}:{d}/\n", .{ config.commander.host, config.commander.port });

        // blocks
        server_instance = &server;
        try server.listen();
    }

    fn initializeMetadataService(allocator: std.mem.Allocator) !void {
        var service = metadata.MetadataService.init(allocator);
        defer service.deinit();

        try service.insertLine(.{ .id = 4, .name = "David", .city = "Paris" });
        try service.deleteLine(2);

        const searchResult = try service.searchLine("city", "Tokyo");
        if (searchResult) |result| {
            std.debug.print("Found: {}\n", .{result});
        } else {
            std.debug.print("Not found.\n", .{});
        }
    }
};

fn getTables(_: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get tables" }, .{});
}

fn getTableByName(_: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get table by name" }, .{});
}

fn createTable(_: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "create table" }, .{});
}

fn getIngestions(_: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get ingestions" }, .{});
}

fn getIngestionById(_: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get ingestion by Id" }, .{});
}

fn createIngestion(_: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "create ingestion" }, .{});
}

fn index(_: *httpz.Request, res: *httpz.Response) !void {
    // the last parameter to res.json is an std.json.StringifyOptions
    try res.json(.{ .name = "Norman Commander API", .version = version }, .{});
}

fn metrics(_: *httpz.Request, res: *httpz.Response) !void {
    // httpz exposes some prometheus-style metrics
    return httpz.writeMetrics(res.writer());
}

fn shutdown(_: c_int) callconv(.C) void {
    if (server_instance) |server| {
        server_instance = null;
        server.stop();
    }
}
