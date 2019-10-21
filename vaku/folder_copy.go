package vaku

import (
	"github.com/pkg/errors"
)

// folderCopyWorkerInput takes input/output channels for input to the job
type folderCopyWorkerInput struct {
	inputsC  <-chan map[string]*PathInput
	resultsC chan<- error
}

// FolderCopy takes in a source PathInput and target PathInput and copies every path in
// the source to the target. Note that this will copy the input path if it is a secret and
// all paths under the input path that result from calling FolderList() on that path. Also
// note that this will overwrite any existing keys at the target paths.
func (c *Client) FolderCopy(s *PathInput, t *PathInput) error {
	var err error

	// Init both paths to get mount info
	s.opType = "readwrite"
	err = c.InitPathInput(s)
	if err != nil {
		return errors.Wrapf(err, "failed to init path %s", s.Path)
	}
	t.opType = "readwrite"
	err = c.InitPathInput(t)
	if err != nil {
		return errors.Wrapf(err, "failed to init path %s", t.Path)
	}

	// Get the keys to copy
	list, err := c.FolderList(&PathInput{
		Path:           s.Path,
		TrimPathPrefix: true,
	})
	if err != nil {
		return errors.Wrapf(err, "Failed to list %s", s.Path)
	}

	// Concurrency channels for workers
	inputsC := make(chan map[string]*PathInput, len(list))
	resultsC := make(chan error, len(list))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderCopyWorker(&folderCopyWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all paths to copy to the inputs channel
	for _, p := range list {
		inputsC <- map[string]*PathInput{
			"source": {
				Path:          c.PathJoin(s.Path, p),
				mountPath:     s.mountPath,
				mountlessPath: s.mountlessPath,
				mountVersion:  s.mountVersion,
			},
			"target": {
				Path:          c.PathJoin(t.Path, p),
				mountPath:     t.mountPath,
				mountlessPath: t.mountlessPath,
				mountVersion:  t.mountVersion,
			},
		}
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(list); j++ {
		o := <-resultsC
		if o != nil {
			err = errors.Wrap(o, "Failed to copy path")
		}
	}

	return err
}

// folderCopyWorker does the work of copying a single path to a new destination
func (c *Client) folderCopyWorker(i *folderCopyWorkerInput) {
	var err error
	for {
		inputs, more := <-i.inputsC
		if more {
			err = c.PathCopy(inputs["source"], inputs["target"])
			if err != nil {
				i.resultsC <- errors.Wrapf(err, "Failed to copy path %s to %s", inputs["source"].Path, inputs["target"].Path)
				continue
			}
			i.resultsC <- nil
		} else {
			return
		}
	}
}

// FolderCopy takes in a source PathInput and target PathInput and copies every path in
// the source to the target. Note that this will copy the input path if it is a secret and
// all paths under the input path that result from calling FolderList() on that path. Also
// note that this will overwrite any existing keys at the target paths.
func (c *CopyClient) FolderCopy(s *PathInput, t *PathInput) error {
	var err error

	// Init both paths to get mount info
	s.opType = "readwrite"
	err = c.Source.InitPathInput(s)
	if err != nil {
		return errors.Wrapf(err, "failed to init path %s", s.Path)
	}
	t.opType = "readwrite"
	err = c.Target.InitPathInput(t)
	if err != nil {
		return errors.Wrapf(err, "failed to init path %s", t.Path)
	}

	// Get the keys to copy
	list, err := c.Source.FolderList(&PathInput{
		Path:           s.Path,
		TrimPathPrefix: true,
	})
	if err != nil {
		return errors.Wrapf(err, "Failed to list %s", s.Path)
	}

	// Concurrency channels for workers
	inputsC := make(chan map[string]*PathInput, len(list))
	resultsC := make(chan error, len(list))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderCopyWorker(&folderCopyWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all paths to copy to the inputs channel
	for _, p := range list {
		inputsC <- map[string]*PathInput{
			"source": {
				Path:          c.Source.PathJoin(s.Path, p),
				mountPath:     s.mountPath,
				mountlessPath: s.mountlessPath,
				mountVersion:  s.mountVersion,
			},
			"target": {
				Path:          c.Target.PathJoin(t.Path, p),
				mountPath:     t.mountPath,
				mountlessPath: t.mountlessPath,
				mountVersion:  t.mountVersion,
			},
		}
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(list); j++ {
		o := <-resultsC
		if o != nil {
			err = errors.Wrap(o, "Failed to copy path")
		}
	}

	return err
}

// folderCopyWorker does the work of copying a single path to a new destination
func (c *CopyClient) folderCopyWorker(i *folderCopyWorkerInput) {
	var err error
	for {
		inputs, more := <-i.inputsC
		if more {
			err = c.PathCopy(inputs["source"], inputs["target"])
			if err != nil {
				i.resultsC <- errors.Wrapf(err, "Failed to copy path %s to %s", inputs["source"].Path, inputs["target"].Path)
				continue
			}
			i.resultsC <- nil
		} else {
			return
		}
	}
}