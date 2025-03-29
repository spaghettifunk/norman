const std = @import("std");
const configuration = @import("../configuration.zig");

pub fn run(config: configuration.Config) !void {
    std.debug.print("Run broker - {s}:{d}", .{ config.broker.host, config.broker.port });
}
