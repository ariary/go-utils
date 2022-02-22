package clipboard

import "os/exec"

func Copy(command string) error {
	// Check for xclip installation
	copyCmdArgs := []string{"xclip", "-in", "-selection", "clipboard"}

	_, err := exec.LookPath("xclip")
	Check(err, "xclip was not found")

	// copy
	copyCmd := exec.Command(copyCmdArgs[0], copyCmdArgs[1:]...)
	in, err := copyCmd.StdinPipe()
	Check(err, "Failed redirecting stdin")

	if err := copyCmd.Start(); err != nil {
		return err
	}

	if _, err := in.Write([]byte(command)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
