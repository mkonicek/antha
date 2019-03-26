// The Migrate Command
//
// The migrate command migrates workflows from outdated historical
// schema versions to the current schema version. Currently, it will
// accept version=1.2 workflow files as inputs, and will produce
// SchemaVersion=2.0 workflows.
//
// The following flags are available:
//
//  -from=path/to/file.json
//    The workflow to migrate. If not provided, the command reads from
//    stdin.
//
//  -to=path/to/file.json
//    The file to write the migrated workflow to. If not provided, the
//    command writes to stdout.
//
//  -validate=true (default: true)
//    Whether or not to attempt to validate the migrated workflow. In
//    some cases it may be necessary to disable validation if it is
//    known that the generated workflow is incomplete (e.g. you know
//    you're only producing a workflow snippet).
//
//  -gilson-device=myFirstGilson (optional)
//    In version=1.2 workflows, the only supported liquid handler
//    device is the Gilson PipetMax, and only one such device is
//    supported per workflow. To migrate those configuration values to
//    corresponding entries in the current SchemaVersion=2.0 workflow,
//    a device name must be provided. The default is that no such
//    device name is provided, thus by default PipetMax device
//    configuration settings are not migrated.
//
// Additionally, as normal, workflow snippets in the current
// SchemaVersion=2.0 format may be provided as additional arguments to
// the command. As a minimum, it is necessary to provide sufficient
// Repositories such that every element instance in the input
// historical workflow (version=1.2) can be located within one of the
// Repositories. If this is not possible then migration will fail.
//
// Log messages are produced on stderr.
package main
