package collatinus

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func lemmatise(chunks []string, output io.Writer, language Language) error {
	err := setLanguage(language)
	if err != nil {
		return fmt.Errorf("failed to set language to %s: %w", language, err)
	}

	for _, chunk := range chunks {
		cmd := exec.Command("/collatinus/bin/Client_C11", "-p2", chunk)

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe: %w", err)
		}

		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start command: %w", err)
		}

		scanner := bufio.NewScanner(stdoutPipe)
		scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Fprintln(output, line)
			if err != nil {
				return fmt.Errorf("failed to write to output buffer: %w", err)
			}
		}

		err = scanner.Err()
		if err != nil {
			return fmt.Errorf("scanner error: %w", err)
		}

		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("command failed: %w", err)
		}
	}

	return nil
}

func setLanguage(language Language) error {
	cmd := exec.Command("/collatinus/bin/Client_C11", string(language))

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
