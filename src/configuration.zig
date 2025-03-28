const std = @import("std");

pub const Config = struct { commander: struct {
    host: []u8,
    port: u16,
}, broker: struct {
    host: []u8,
    port: u16,
}, storage: struct {
    host: []u8,
    port: u16,
}, log: struct {
    level: []u8,
} };

pub fn load(allocator: std.mem.Allocator, path: []const u8) !Config {
    // TODO: validate if 512 is a sufficient value
    const data = try std.fs.cwd().readFileAlloc(allocator, path, 512);
    defer allocator.free(data);

    const result = try std.json.parseFromSlice(Config, allocator, data, .{});
    // defer result.deinit();

    return result.value;
}
