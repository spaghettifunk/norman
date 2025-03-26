const std = @import("std");
const httpz = @import("httpz");
const configuration = @import("../configuration.zig");
const metadata = @import("metadata.zig");
const models = @import("../models/models.zig");
const version = @import("../const.zig").normanVersion;

var server_instance: ?*httpz.Server(Commander) = null;

pub const Commander = struct {
    allocator: std.mem.Allocator,
    config: configuration.Config,
    md: metadata.MetadataService,

    pub fn init(allocator: std.mem.Allocator, config: configuration.Config) !Commander {
        const md = try metadata.MetadataService.init(allocator);

        return Commander{
            .allocator = allocator,
            .config = config,
            .md = md,
        };
    }

    pub fn start(self: Commander) !void {
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

        var server = try httpz.Server(Commander).init(self.allocator, .{ .port = self.config.commander.port }, self);
        defer server.deinit();

        var router = try server.router(.{});
        router.get("/", index, .{});
        router.get("/api/tables", getTables, .{});
        router.post("/api/tables", createTable, .{});
        router.get("/api/tables/:name", getTableByName, .{});
        router.get("/api/ingestions", getIngestions, .{});
        router.post("/api/ingestions", createIngestion, .{});
        router.get("/api/ingestions/:id", getIngestionById, .{});
        router.get("/metrics", metrics, .{});

        std.debug.print("Listening on http://{s}:{d}/\n", .{ self.config.commander.host, self.config.commander.port });

        // blocks
        server_instance = &server;
        try server.listen();
    }
};

fn getTables(_: Commander, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get tables" }, .{});
}

fn getTableByName(_: Commander, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get table by name" }, .{});
}

fn createTable(_: Commander, req: *httpz.Request, res: *httpz.Response) !void {
    _ = req.body();

    res.status = 200;
    try res.json(.{ .result = "create table" }, .{});
}

fn getIngestions(_: Commander, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get ingestions" }, .{});
}

fn getIngestionById(_: Commander, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    try res.json(.{ .result = "get ingestion by Id" }, .{});
}

fn createIngestion(commander: Commander, req: *httpz.Request, res: *httpz.Response) !void {
    if (req.body()) |body| {
        const maybe_ingestion: ?std.json.Parsed(models.IngestionJob) = std.json.parseFromSlice(models.IngestionJob, commander.allocator, body, .{}) catch |err| {
            std.debug.print("error parsing json: {any}\n", .{err});
            res.status = 400;
            return;
        };
        if (maybe_ingestion) |ingestion| {
            defer ingestion.deinit();
        }
    }

    res.status = 200;
    try res.json(.{ .result = "create ingestion" }, .{});
}

fn index(_: Commander, _: *httpz.Request, res: *httpz.Response) !void {
    // the last parameter to res.json is an std.json.StringifyOptions
    try res.json(.{ .name = "Norman Commander API", .version = version }, .{});
}

fn metrics(_: Commander, _: *httpz.Request, res: *httpz.Response) !void {
    // httpz exposes some prometheus-style metrics
    return httpz.writeMetrics(res.writer());
}

fn shutdown(_: c_int) callconv(.C) void {
    if (server_instance) |server| {
        server_instance = null;
        server.stop();
    }
}
