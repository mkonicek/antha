// The Migrate Command
//
// The migrate command migrates workflows from outdated historical
// schema versions to the current schema version. Currently, it will
// accept version=1.2 workflow files as inputs, and will produce
// SchemaVersion=2.0 workflows.
//
// The version=1.2 format is deficient in that it does not contain
// enough information for workflow to be repeatedly executed. Thus
// when migrating, you are required to provide additional information
// (particularly Repositories) in the SchemaVersion=2.0 format which
// is then combined with the old workflow to create a workflow.
//
// The following flags are available:
//
//  -from=path/to/file.json
//    The workflow to migrate. If not provided, the command reads from
//    stdin.
//
//  -outdir=path/to/directory
//    A directory to write the results to. In SchemaVersion=2.0, the
//    workflow.json itself cannot contain file content. Thus a
//    directory must be provided so that file content that was within
//    the old workflow can be extracted and written out. The new
//    workflow is written to directory/workflow/workflow.json, and any
//    file contents are written to directory/data/. The new workflow
//    will contain references to any extracted files, and the
//    directory layout matches the requirements of the composer -indir
//    flag. If not provided, a fresh temporary directory is used.
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
// Example:
//   migrate -from=path/to/old.json myRepositories.json
//
// Log messages are produced on stderr.
package main
