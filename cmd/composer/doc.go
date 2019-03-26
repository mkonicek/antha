// The Composer Command
//
// The composer command is the main command line tool for composing
// workflow snippets together, transpiling all required elements,
// generating a workflow main, compiling it all, and optionally
// invoking the compiled binary.
//
// A workflow and a workflow snippet are both json files. A workflow
// snippet may not contain enough information to be considered a full
// workflow. For example, a workflow snippet may simply define a
// configuration for a robot, or may define a locally checked-out
// repository. Some snippets may be more appropriate to be shared (for
// example element instances, parameters, connections) and some may
// contain configuration relevant only to a particular
// computer. Therefore in general, the composer command accepts as
// inputs a set of snippets, which get merged together before being
// validated and interpreted.
//
// The steps the composer follows are as follows:
//
//  1. Load all workflow snippets provided (see below), parse as JSON,
//     merge together and validate the result.
//
//  2. Checkout all the known repositories under outdir/src at the
//     indicated revisions.
//
//  3. Transpile the transitive closure of all antha element files
//     into Go code.
//
//  4. Use the merged and validated workflow to generate a main.go in
//     outdir/workflow. This is the entry point for the execution of
//     the workflow itself.
//
//  5. Prepare the device instruction plugins as necessary based on
//     their configuration. If files are indicated, then copy into
//     outdir/bin/drivers. If go:// URLs are provided, then build
//     those into outdir/bin/drivers.
//
//  6. Save the merged, validated and tweaked workflow into
//     outdir/workflow/data/workflow.json. The complete workflow is
//     required during execution of the workflow itself, not just
//     generation of the workflow binary (for example, for
//     configuration of device instruction plugins).
//
//  7. Compile the generated main.go into outdir/bin/job-id and build
//     into the generated binary the merged and validated
//     workflow. The generated binary is relatively self contained at
//     this point - the only external dependencies are non-linked-in
//     instruction plugins, and input data files.
//
//  8. If required, execute the compiled workflow.
//
// The composer command accepts the following flags:
//
//  -outdir path/to/directory (optional)
//    An optional path to a directory to use for checking out element
//    sources, storing the generated main, workflow and compiled output.
//
//    If this option is not provided then a fresh temporary direction
//    is created. If this option is provided, the indicated directory
//    must either not exist, or must be empty.
//
//    The layout of this directory is documented in github.com/antha-lang/antha/composer
//
//  -indir path/to/directory (optional)
//    An optional path to a directory. This is used for two purposes:
//      1. If path/to/directory/workflow exists and is a directory
//         then all json files within that directory are read in and
//         parsed as workflow snippets
//      2. If path/to/directory/data exists and is a directory, and
//         the generated workflow binary is executed, then calls to
//         Laboratory.FileManager.ReadAll (or similar read-based
//         calls) will have their paths interpreted as relative to
//         this directory.
//
//  -keep=false
//    By default, after steps 1-7 are successfully run, the contents
//    of outdir/src and outdir/workflow are removed (i.e. the checked
//    out sources, the merged and validated workflow and generated
//    main.go). These are not removed if -keep=true
//
//  -run=true
//    By default, the compiled workflow is executed (step 8). If
//    -run=false then the composer exits after steps 1-7.q
//
//  -linkedDrivers=false
//    By default we do not require a checkout of the
//    github.com/Synthace/instruction-plugins repository, which means
//    that every device within the workflow that requires an
//    instruction plugin must have a valid 'Connection' field. If
//    -linkedDrivers=true then during compilation of the workflow, we
//    attempt to compile and link in all the instruction plugins. If
//    this succeeds, then any device in the workflow configured
//    without a 'Connection' field will instead fall back on the
//    relevant linked instruction plugin.
//
// On the command line, any further arguments are interpreted as paths
// to json files (in addition to -indir, if provided) and are loaded
// and merged as workflow snippets. If you provide - as an argument
// then composer will try to parse a workflow snippet from stdin.
//
// Log messages are produced on stderr, and the composer command exits
// with a code of 0 iff no fatal error in encountered.
package main
