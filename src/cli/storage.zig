const std = @import("std");
const configuration = @import("../configuration.zig");

pub fn run(config: configuration.Config) !void {
    std.debug.print("Run storage - {s}:{d}", .{ config.storage.host, config.storage.port });
}
