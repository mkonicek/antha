package main

import (
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func findElements(l *logger.Logger, paths []string, consumer func(*workflow.Repository, *workflow.ElementType) error) error {
	if rs, err := workflow.ReadersFromPaths(paths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if repoToEts, err := wf.Repositories.FindAllElementTypes(); err != nil {
		return err
	} else {
		for repo, ets := range repoToEts {
			for _, et := range ets {
				if err := consumer(repo, &et); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
