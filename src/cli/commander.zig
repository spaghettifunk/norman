const std = @import("std");
const configuration = @import("../configuration.zig");
const commander = @import("../commander/service.zig");

pub fn run(allocator: std.mem.Allocator, config: configuration.Config) !void {
    std.debug.print("Run commander - {s}:{d}\n", .{ config.commander.host, config.commander.port });

    try commander.Commander.start(allocator, config);
}
