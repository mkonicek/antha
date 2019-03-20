// There are two runtimes that we must think about:
//
// 1. The composer runtime. This is where we parse and validate the
// workflow, checkout the necessary sources, generate a main.go, and
// compile it all.
//
// 2. The workflow runtime. This is where we run the generated main.go
// and consequently the elements, passing their effects into the
// planner and hopefully generating some instructions for robots.
//
// Both composer and workflow runtimes have some flags in common:
// namely -indir and -outdir. By default, they share the -indir flag,
// but have separate -outdir directories. The filesystem layout is as
// follows:
//
//  inDir/
//    data/
//      # calls to lab.FileManager.ReadAll() must read from in here (workflow runtime)
//    workflow/
//      *.json # treated as workflow fragments (composer runtime)
//
//  outDir/ (for composer runtime)
//    src/
//      # where we check out the element repositories, and do the transpilation of elements
//    workflow/
//      main.go # the generated main.go
//      data/
//        workflow.json # the merged and validated workflow
//    logs.txt # logs from the composer only
//    bin/
//      jobId # the compiled binary for the whole workflow (result of compiling main.go)
//      drivers/
//        driverId # binary instruction plugin (iff go:// or file:// Connection string in driver config)
//
//  outDir/ (for workflow runtime)
//    data/
//      # calls to lab.FileManager.Write* write in here
//    elements/
//      # each element gets a elemName.json file in here which is the serialization of the element after run
//    logs.txt # logs from the workflow runtime
//    errors.txt # if any element gives an error (or panic) it will go in here
//    devices/
//      task-id/
//        device-specific-files # files generated for the device as part of the task
package composer
