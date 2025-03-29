const std = @import("std");
const cli = @import("cli.zig");
const configuration = @import("configuration.zig");

pub fn main() !void {
    if (comptime @import("builtin").os.tag == .windows) {
        std.debug.print("Norman does not run on Windows. Sorry\n", .{});
        return error.PlatformNotSupported;
    }

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    // default config path
    var config_path: []const u8 = "./config.json";

    var i: usize = 1;
    while (i < args.len) : (i += 1) {
        if (std.mem.eql(u8, args[i], "-c") or std.mem.eql(u8, args[i], "--configuration")) {
            i += 1;
            if (i >= args.len) {
                std.debug.print("Error: Configuration path not provided after -c/--configuration.\n", .{});
                std.process.exit(1);
            }
            config_path = args[i];
        } else {
            break; // Stop processing flags when a non-flag argument is encountered
        }
    }

    if (i >= args.len) {
        std.debug.print("Error: Command not provided.\n", .{});
        std.debug.print("Usage: norman [-c/--configuration <config_path>] <command> [command args...]\n", .{});
        std.debug.print("Available commands: storage, commander, broker\n", .{});
        std.process.exit(1);
    }

    const command = args[i];
    const config = try configuration.load(allocator, config_path);

    if (std.mem.eql(u8, command, "storage")) {
        try cli.storage.run(config);
    } else if (std.mem.eql(u8, command, "commander")) {
        try cli.commander.run(allocator, config);
    } else if (std.mem.eql(u8, command, "broker")) {
        try cli.broker.run(config);
    } else {
        std.debug.print("Error: Unknown command: {s}\n", .{command});
        std.debug.print("Available commands: storage, commander, broker\n", .{});
        std.process.exit(1);
    }
}
